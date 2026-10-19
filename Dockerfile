# GitHub Action image: one binary plus the test runners it shells out to.
FROM golang:1.22-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o /out/proofline ./cmd/proofline

FROM debian:bookworm-slim
RUN apt-get update \
  && apt-get install -y --no-install-recommends \
    ca-certificates \
    git \
    python3 \
    python3-pytest \
    nodejs \
    npm \
  && rm -rf /var/lib/apt/lists/*
COPY --from=build /out/proofline /usr/local/bin/proofline
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
WORKDIR /github/workspace
ENTRYPOINT ["/entrypoint.sh"]
