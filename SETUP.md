# Ethereum Calldata Analysis Tool - Setup & Usage

## What You Have

A complete Go-based tool that analyzes Ethereum smart contract function calls and computes calldata statistics (mean, median, mode) over a configurable time period.

### Files in this directory:

```
eth-analysis                 # Compiled binary (ready to run)
main.go                      # Config parsing + CLI logic
analyzer.go                  # Etherscan API + statistics
config.example.{json,yml,toml,csv,txt}  # Example configs
README.md                    # Full documentation
Makefile                     # Build shortcuts
QUICKSTART.sh               # Setup guide
```

---

## 5-Minute Setup

### 1. Get Your Etherscan API Key (2 min)

1. Go to **https://etherscan.io/apis**
2. Create a free account or sign in
3. Click "Create App Token"
4. Copy your API key

### 2. Create config.json (2 min)

Copy `config.example.json` to `config.json` and replace:
- `YOUR_ETHERSCAN_API_KEY` → your actual API key
- Add your contracts and functions

Example:
```json
{
  "etherscan_api_key": "abc123xyz...",
  "network": "mainnet",
  "days_back": 7,
  "contracts": [
    {
      "address": "0x1111111254fb6c44bac0bed2854e76f90643097d",
      "name": "1inch Router",
      "functions": [
        {"name": "swap", "selector": "0x7c025200"}
      ]
    }
  ]
}
```

### 3. Run It (1 min)

```bash
./eth-analysis -config config.json -output text
```

Done! See your results instantly.

---

## How to Find Function Selectors

**Option 1: Use 4byte.directory** (easiest)
1. Go to https://www.4byte.directory/
2. Search for your function: `swap(address[],uint256,uint256)`
3. Copy the 4-byte selector: `0x7c025200`

**Option 2: Calculate manually**
- Selector = first 4 bytes of `keccak256("function_signature")`
- Python: `from eth_utils import function_signature_to_4byte_selector`
- Online: https://emn178.github.io/online-tools/keccak_256.html

---

## Example Usage

### Single contract, past week:
```json
{
  "etherscan_api_key": "YOUR_KEY",
  "days_back": 7,
  "contracts": [{
    "address": "0xE592427A0AEce92De3Edee1F18E0157C05861564",
    "name": "Uniswap V3",
    "functions": [
      {"name": "exactInputSingle", "selector": "0x414bf389"}
    ]
  }]
}
```

```bash
./eth-analysis -config config.json -output text
```

Output:
```
=== Uniswap V3 ===
Function                       Count      Min          Max        Mean       Median       Mode       Mode Cnt
----------------------------------------------------
exactInputSingle              4521       68           2048       512.30     512.00       512        892
```

### Multiple contracts, past 2 weeks, CSV output:
```bash
./eth-analysis -config config.json -output csv > results.csv
```

Then open `results.csv` in Excel/Google Sheets.

### Past month, JSON output:
```json
{
  "etherscan_api_key": "YOUR_KEY",
  "days_back": 30,
  ...
}
```

```bash
./eth-analysis -config config.json > analysis.json
```

---

## Config Format Examples

All these formats work identically:

**JSON:**
```json
{"etherscan_api_key": "...", "days_back": 7, "contracts": [...]}
```

**YAML:**
```yaml
etherscan_api_key: ...
days_back: 7
contracts:
  - address: "0x..."
    name: "Contract"
    functions:
      - name: "swap"
        selector: "0x..."
```

**TOML:**
```toml
etherscan_api_key = "..."
days_back = 7
[[contracts]]
address = "0x..."
name = "Contract"
[[contracts.functions]]
name = "swap"
selector = "0x..."
```

**CSV:**
```csv
address,name,function_name,selector
0x1111...,1inch Router,swap,0x7c025200
```

**TXT:**
```txt
0x1111... | 1inch Router | swap | 0x7c025200
```

---

## Understanding the Output

Each function shows:

| Column | Meaning |
|--------|---------|
| **Count** | Total function calls in the time period |
| **Min Size** | Smallest calldata (in bytes) |
| **Max Size** | Largest calldata (in bytes) |
| **Mean Size** | Average calldata size |
| **Median Size** | Middle value (50th percentile) |
| **Mode Size** | Most common calldata size |
| **Mode Count** | How many times the mode appeared |

### Example Interpretation:
```
swap: Count=1250, Min=196, Max=4096, Mean=852.50, Median=768, Mode=512 (Count=183)
```

This means:
- 1,250 swap calls in the past period
- Typically around **512 bytes** (the mode)
- Average closer to **852 bytes**
- Range from 196 to 4,096 bytes
- 768 bytes is the middle value

---

## Performance & Limits

**Speed:**
- Single contract: 5-15 seconds
- Multiple contracts (3-5): 30-60 seconds
- Etherscan API: ~5 calls/second (free tier)

**Time ranges:**
- 1 week: ✓ (fast)
- 2-3 weeks: ✓ (fast)
- 1 month: ✓ (40-60 sec)
- 3 months: ⚠️ (slower, may hit rate limits)
- 1 year: ✗ (too many transactions)

**If you hit rate limits:**
- Upgrade your Etherscan API key (paid tier)
- Analyze one contract at a time
- Reduce `days_back`

---

## Common Issues

**"No records found"**
- The contract didn't have calls in the time period
- Try increasing `days_back`
- Verify the address is correct (check Etherscan)

**"Unsupported file format"**
- File extension must be `.json`, `.yml`, `.yaml`, `.toml`, `.csv`, or `.txt`

**"etherscan error"**
- Check your API key is correct
- Verify you're on mainnet (not testnet)
- Wait a moment if you hit rate limits

---

## Next Steps

1. **Add your contracts** - Copy example config, add your contracts
2. **Run analysis** - `./eth-analysis -config config.json`
3. **Export results** - Use `-output csv` for spreadsheet, `-output json` for programmatic access
4. **Iterate** - Change `days_back` to see trends over different periods

Questions? Check the full README.md or visit https://etherscan.io/apis

Happy analyzing! 🚀
