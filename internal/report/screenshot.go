package report

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Screenshot writes a PNG of the HTML report. Chrome is used when it is on PATH.
// Otherwise a bitmap fallback is written so the demo still produces a PNG offline.
func Screenshot(htmlPath, pngPath string, lines []string) error {
	if err := chrome(htmlPath, pngPath); err == nil {
		if st, statErr := os.Stat(pngPath); statErr == nil && st.Size() > 1000 {
			return nil
		}
	}
	return bitmap(pngPath, lines)
}

func chrome(htmlPath, pngPath string) error {
	absHTML, err := filepath.Abs(htmlPath)
	if err != nil {
		return err
	}
	absPNG, err := filepath.Abs(pngPath)
	if err != nil {
		return err
	}
	bin, err := exec.LookPath("google-chrome")
	if err != nil {
		bin, err = exec.LookPath("chromium")
		if err != nil {
			return err
		}
	}
	profile, err := os.MkdirTemp("", "proofline-chrome-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(profile)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin,
		"--headless=new",
		"--disable-gpu",
		"--no-sandbox",
		"--disable-dev-shm-usage",
		"--hide-scrollbars",
		"--no-first-run",
		"--user-data-dir="+profile,
		"--force-device-scale-factor=1",
		"--virtual-time-budget=4000",
		"--timeout=8000",
		"--window-size=1200,980",
		"--screenshot="+absPNG,
		"file://"+absHTML,
	)
	cmd.Env = append(os.Environ(), "HOME=/tmp")
	if err := cmd.Start(); err != nil {
		return err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	// Headless Chrome in this environment writes the PNG and then fails to exit.
	// Keep the file once its size settles, and do not wait forever.
	deadline := time.Now().Add(15 * time.Second)
	var last int64
	stable := 0
	for time.Now().Before(deadline) {
		select {
		case err := <-done:
			if err != nil && !pngReady(absPNG) {
				return fmt.Errorf("chrome: %w", err)
			}
			return nil
		default:
		}
		if st, err := os.Stat(absPNG); err == nil && st.Size() > 1000 {
			if st.Size() == last {
				stable++
				if stable >= 3 {
					_ = cmd.Process.Kill()
					<-done
					return nil
				}
			} else {
				stable = 0
				last = st.Size()
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
	<-done
	if pngReady(absPNG) {
		return nil
	}
	return fmt.Errorf("chrome did not write a screenshot")
}

func pngReady(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.Size() > 1000
}

func bitmap(path string, lines []string) error {
	const scale = 2
	width, height := 920*scale, 640*scale
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	bg := color.RGBA{244, 240, 230, 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bg)
		}
	}
	face := basicfont.Face7x13
	y := 28
	for _, line := range lines {
		drawScaled(img, line, 16, y, scale, face)
		y += 18
		if y > 620 {
			break
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func drawScaled(dst *image.RGBA, text string, x, y, scale int, face font.Face) {
	tmp := image.NewRGBA(image.Rect(0, 0, 900, 20))
	d := &font.Drawer{
		Dst:  tmp,
		Src:  image.NewUniform(color.RGBA{28, 25, 21, 255}),
		Face: face,
		Dot:  fixed.P(0, 13),
	}
	d.DrawString(text)
	b := tmp.Bounds()
	for yy := 0; yy < b.Dy(); yy++ {
		for xx := 0; xx < b.Dx(); xx++ {
			c := tmp.RGBAAt(xx, yy)
			if c.A == 0 {
				continue
			}
			for sy := 0; sy < scale; sy++ {
				for sx := 0; sx < scale; sx++ {
					dst.Set(x+xx*scale+sx, y+yy*scale+sy, c)
				}
			}
		}
	}
}
