# Deployment & Distribution

## For Your Own Use

1. **Keep the binary**
   ```bash
   ./eth-analysis -config config.json -output text
   ```

2. **Or rebuild anytime**
   ```bash
   go mod tidy
   go build -o eth-analysis
   ```

## For Others (Open Source)

### Build for multiple platforms:
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o eth-analysis-linux-amd64

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o eth-analysis-darwin-amd64

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o eth-analysis-darwin-arm64

# Windows
GOOS=windows GOARCH=amd64 go build -o eth-analysis-windows-amd64.exe
```

### Create a release:
```bash
mkdir -p release
cp eth-analysis-* release/
cp README.md SETUP.md config.example.json release/
cd release && tar czf eth-analysis.tar.gz * && cd ..
```

## Docker (Optional)

If you want to containerize:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o eth-analysis

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /build/eth-analysis /usr/local/bin/
ENTRYPOINT ["eth-analysis"]
```

Build: `docker build -t eth-analysis .`
Use: `docker run -v $(pwd):/data eth-analysis -config /data/config.json`

## GitHub Actions (CI/CD)

Auto-build and release on new tags:

```yaml
name: Release
on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: go build -o eth-analysis
      - uses: softprops/action-gh-release@v1
        with:
          files: eth-analysis
```

## Distribution Tips

- Make config examples very clear
- Emphasize: "NEVER commit config.json with API keys"
- Provide selector lookup instructions
- Include troubleshooting section
- Keep binary small (~10MB is reasonable)

