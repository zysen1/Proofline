package diffx

import "testing"

func TestParseUnified(t *testing.T) {
	diff := "" +
		"diff --git a/billing.py b/billing.py\n" +
		"--- a/billing.py\n" +
		"+++ b/billing.py\n" +
		"@@ -4 +4 @@\n" +
		"-    raise NotImplementedError(\"x\")\n" +
		"+    return price * (100 - percent) // 100\n" +
		"diff --git a/tests/test_billing.py b/tests/test_billing.py\n" +
		"--- a/tests/test_billing.py\n" +
		"+++ b/tests/test_billing.py\n" +
		"@@ -0,0 +8,2 @@\n" +
		"+def test_new():\n" +
		"+    assert True\n"
	got := ParseUnified(diff)
	if len(got["billing.py"]) != 1 || got["billing.py"][0] != 4 {
		t.Fatalf("billing lines = %#v", got["billing.py"])
	}
	if len(got["tests/test_billing.py"]) != 2 {
		t.Fatalf("test lines = %#v", got["tests/test_billing.py"])
	}
}
