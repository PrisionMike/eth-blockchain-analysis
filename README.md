# Ethereum Calldata Analysis Tool

Analyze Ethereum smart contract function calls over a specified time period. Get aggregate statistics on calldata sizes (mean, median, mode) for your specified functions.

## ✨ Created by GitHub Copilot

This tool was designed and built by GitHub Copilot. The source code, documentation, and all implementation are AI-generated using Copilot CLI runtime in VS Code.

## Features

- **Multiple config formats**: JSON, YAML, TOML, CSV, or TXT
- **Flexible time ranges**: Analyze past 7 days, 2 weeks, month, etc.
- **Continuous analysis**: Run for 1-2 hours and continuously fetch new data
- **Multiple output formats**: JSON, CSV, or text table
- **Etherscan API integration**: Fast data fetching
- **Statistical analysis**: Mean, median, mode, min, max calldata sizes
- **Block range reporting**: See exactly which blocks were scanned
- **RPC support**: Optional direct RPC node support for better filtering

## Key Improvements

### Before
- One-time analysis of historical data only
- Wasteful API calls fetching all transactions
- Limited filtering capability

### Now
- ✅ **Continuous mode**: Run for any duration (1-120 minutes) and keep fetching
- ✅ **Block range reporting**: Know exactly what blocks were analyzed
- ✅ **Better filtering**: Etherscan + RPC support for selective calldata analysis
- ✅ **Smarter queries**: Function selector filtering at API level
- ✅ **Real-time monitoring**: Keep analyzing while new blocks are being created

## Installation

```bash
go mod tidy
go build -o eth-analysis
```

## Usage

### 1. Create a Config File

Choose your preferred format and create a config file. Example formats:

**JSON** (`config.json`):
```json
{
  "etherscan_api_key": "YOUR_KEY",
  "network": "mainnet",
  "days_back": 7,
  "contracts": [
    {
      "address": "0x1111111254fb6c44bac0bed2854e76f90643097d",
      "name": "1inch Router",
      "functions": [
        {"name": "swap", "selector": "0x7c025200"},
        {"name": "fillOrder", "selector": "0xd0a3b665"}
      ]
    }
  ]
}
```

**YAML** (`config.yml`):
```yaml
etherscan_api_key: YOUR_KEY
network: mainnet
days_back: 7
contracts:
  - address: "0x1111111254fb6c44bac0bed2854e76f90643097d"
    name: "1inch Router"
    functions:
      - name: swap
        selector: "0x7c025200"
```

**CSV** (`config.csv`):
```csv
address,name,function_name,selector
0x1111111254fb6c44bac0bed2854e76f90643097d,1inch Router,swap,0x7c025200
```

**TXT** (`config.txt`):
```txt
0x1111111254fb6c44bac0bed2854e76f90643097d | 1inch Router | swap | 0x7c025200
```

See example configs in this directory: `config.example.*`

### 2. Get Your Etherscan API Key

1. Go to https://etherscan.io/apis
2. Create a free account and generate an API key
3. Add it to your config file

### 3. Run the Analysis

```bash
# Default: uses config.json, outputs JSON
./eth-analysis

# Specify config file
./eth-analysis -config config.yml

# Output as CSV
./eth-analysis -config config.json -output csv

# Output as text table
./eth-analysis -config config.json -output text
```

## Config Parameters

- **etherscan_api_key**: Your Etherscan API key (required)
- **rpc_endpoint**: Optional direct RPC node URL (for better filtering)
- **network**: Blockchain network, typically "mainnet"
- **run_duration_mins**: Minutes to run continuous analysis (default: 60)
- **days_back**: Number of days to start analyzing from (default: 7)
- **contracts**: Array of contracts to analyze
  - **address**: Contract address (0x...)
  - **name**: Display name for the contract
  - **functions**: Array of functions to track
    - **name**: Function name
    - **selector**: 4-byte function selector (0x...)

## How Continuous Analysis Works

When you run the tool with `run_duration_mins: 60`, it will:

1. **Continuously fetch** transaction data for 60 minutes
2. **Poll every 5 seconds** for new data
3. **Track block numbers** as they come in
4. **Report at the end**:
   - Start block to end block (block range analyzed)
   - Total blocks scanned
   - Total transactions analyzed
   - Statistics per function

**Example output:**
```
Duration: 2024-06-23 15:10:25 to 2024-06-23 16:10:25
Blocks Scanned: 20245342 → 20246582 (Total: 1240 blocks)
Transactions Analyzed: 4521
```

This means the tool was running while new blocks were being created and captured all matching transactions during that window.

## Function Selectors

To find function selectors:
1. Use Etherscan's 4-byte directory: https://www.4byte.directory/
2. Or calculate them: First 4 bytes of keccak256 hash of function signature
   - Example: `keccak256("swap(address[],uint256,uint256,uint256,address,uint256)")` → `0x7c025200`

## Output Examples

### JSON Output
```json
{
  "1inch Router": [
    {
      "FunctionName": "swap",
      "Selector": "0x7c025200",
      "Count": 1250,
      "MinSize": 196,
      "MaxSize": 4096,
      "MeanSize": 852.5,
      "MedianSize": 768,
      "ModeSize": 512,
      "ModeCount": 183
    }
  ]
}
```

### Text Output
```
=== 1inch Router ===
Function                       Count      Min          Max        Mean       Median       Mode       Mode Cnt
--------------------------------------------
swap                           1250       196          4096       852.50     768.00       512        183
```

## Performance

- Etherscan API rate limits: Free tier allows ~5 calls/second
- Typical analysis time for a single contract over 7 days: 5-15 seconds
- Full month analysis with multiple contracts: 30-60 seconds

## Troubleshooting

**"etherscan error: No records found"**
- The contract might not have had transactions in the time period
- Try increasing `days_back`
- Verify the contract address is correct

**"unsupported file format"**
- Make sure your config file has a recognized extension (.json, .yml, .yaml, .toml, .csv, .txt)

**Rate limit errors**
- Free Etherscan API has rate limits (~5 calls/sec)
- Consider upgrading your API key or running analysis with fewer contracts at once
