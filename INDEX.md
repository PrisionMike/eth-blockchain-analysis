# Ethereum Calldata Analysis Tool - File Index

## 🚀 Start Here

1. **SETUP.md** - Quick 5-minute setup guide (⭐ READ THIS FIRST)
2. **QUICKSTART.sh** - Interactive walkthrough
3. **README.md** - Complete reference documentation

## 🔧 Using the Tool

- **eth-analysis** - The compiled binary (ready to run!)
- **config.example.*** - Example config files in different formats:
  - `config.example.json` - JSON format (recommended)
  - `config.example.yml` - YAML format
  - `config.example.toml` - TOML format
  - `config.example.csv` - CSV format
  - `config.example.txt` - Text format

## 📝 Configuration

Copy any example config and edit it:
```bash
cp config.example.json config.json
# Edit config.json with your Etherscan API key and contracts
./eth-analysis -config config.json -output text
```

## 📚 Documentation

- **README.md** - Full API reference, examples, troubleshooting
- **SETUP.md** - Step-by-step setup instructions
- **ARCHITECTURE.md** - How the tool works internally
- **DEPLOYMENT.md** - Building for different platforms, Docker, CI/CD
- **QUICKSTART.sh** - Interactive setup guide
- **Makefile** - Build and development commands

## 💻 Source Code

- **main.go** - CLI parsing, config loading (JSON, YAML, TOML, CSV, TXT)
- **analyzer.go** - Etherscan API integration, statistics calculation
- **go.mod / go.sum** - Go dependencies

## ⚙️ Build & Development

- **Makefile** - Run `make help` to see available commands:
  - `make build` - Rebuild the binary
  - `make clean` - Remove build artifacts

- **.gitignore** - Prevents accidentally committing API keys or results

## 📊 Quick Command Reference

### Basic usage (text output):
```bash
./eth-analysis -config config.json -output text
```

### JSON output (programmatic):
```bash
./eth-analysis -config config.json -output json > results.json
```

### CSV output (spreadsheet):
```bash
./eth-analysis -config config.json -output csv > results.csv
```

## 🎯 What the Tool Does

1. **Reads** your config file (contract addresses, function selectors, time range)
2. **Fetches** transaction data from Etherscan API
3. **Filters** by function selector (first 4 bytes of calldata)
4. **Calculates** statistics on calldata sizes:
   - Count (how many calls)
   - Min/Max (smallest/largest)
   - Mean (average)
   - Median (middle value)
   - Mode (most common value)
5. **Outputs** results in your chosen format

## ⚡ Performance

- Single contract (7 days): ~5-15 seconds
- Multiple contracts (3-5): ~30-60 seconds  
- Past month: ~40-60 seconds
- Etherscan API: ~5 calls/second (free tier)

## 📖 Tutorial

### Step 1: Get API Key (2 min)
1. Visit https://etherscan.io/apis
2. Create account / sign in
3. Generate API key
4. Copy it

### Step 2: Configure (2 min)
1. Copy: `cp config.example.json config.json`
2. Edit and paste your API key
3. Add your contract addresses and function selectors

### Step 3: Find Function Selectors (5 min)
Visit https://www.4byte.directory/ and search for your functions

### Step 4: Run (1 min)
```bash
./eth-analysis -config config.json -output text
```

Done! ✨

## 🆘 Troubleshooting

**See README.md for common issues and solutions**

## 📞 Questions?

All documentation is in this directory:
- General questions → README.md
- Setup help → SETUP.md
- How it works → ARCHITECTURE.md
- Advanced usage → DEPLOYMENT.md

---

**Total Time to First Run: ~15 minutes**

Start with SETUP.md! 🚀
