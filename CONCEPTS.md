# Blockchain Concepts Explained

This document explains key blockchain concepts used in this tool.

## Events vs Function Calls

### Smart Contracts Have Two Parts

When someone interacts with a smart contract, the contract can:

1. **Execute Code** (function call)
   - Process logic
   - Update state
   - Send ETH
   - Call other contracts

2. **Emit Events** (optional)
   - Log structured data
   - Easy to search/filter
   - NOT stored in contract storage
   - Like a "receipt" of what happened

### Function Calls (What We Analyze)

```solidity
function swap(
    address[] tokens,
    uint256 minAmountOut
) external {
    // Execute logic
    // NO EVENT EMITTED
}
```

**In Ethereum:**
- Transaction calldata: `0x7c025200` + encoded parameters
- Stored in transaction `input` field
- Must be decoded to see what happened
- Harder to filter at blockchain level

**Our Tool:**
- Downloads calldata
- Checks first 4 bytes (selector)
- Calculates size
- Works for ALL functions

### Events

```solidity
event Swap(
    address indexed token0,
    address indexed token1,
    uint256 amount,
    address indexed sender
);

function swap(...) external {
    // Execute logic
    emit Swap(token0, token1, amount, msg.sender);
}
```

**In Ethereum:**
- Stored as structured "logs"
- Topic0: Event signature hash (keccak256)
- Topic1-3: Indexed parameters
- Data: Non-indexed parameters
- Easy to query with `getLogs`

**Etherscan API:**
- `/getLogs` filters by Topic0
- Returns only matching events
- Faster, fewer API calls
- **Only works if event is emitted**

---

## Function Selectors Explained

### What is a Selector?

**Function Selector** = First 4 bytes of function signature hash

```
Function signature: "swap(address[],uint256,uint256)"
↓
keccak256 hash: 0x7c025200aabbcc...
↓
First 4 bytes (selector): 0x7c025200
```

### How Ethereum Uses It

```
Calldata Structure:
┌─────────┬──────────────────┐
│ Selector│ Encoded Parameters
│ 4 bytes │ Variable length
│ 0x7c... │ 0xabab...
└─────────┴──────────────────┘
```

When a transaction is sent:
1. First 4 bytes tell contract WHICH function to call
2. Remaining bytes are the function parameters
3. Contract decodes and executes

### Finding Selectors

**Method 1: Use 4byte.directory**
- Website: https://www.4byte.directory/
- Search: "swap(address[],uint256)"
- Returns selector instantly

**Method 2: Calculate Manually**
```python
from eth_utils import function_signature_to_4byte_selector
selector = function_signature_to_4byte_selector("swap(address[],uint256)")
# Returns b'\x7c\x02\x52\x00' or '0x7c025200'
```

**Method 3: Use ethers.js**
```javascript
const abi = [
  {
    name: "swap",
    inputs: [{type: "address[]"}, {type: "uint256"}]
  }
];
const iface = new ethers.Interface(abi);
console.log(iface.getSighash("swap")); // 0x7c025200
```

---

## Calldata Size

### What is Calldata Size?

**Calldata** = Function selector + encoded parameters

```
swap([0x123, 0x456], 1000)
↓
Calldata hex: 0x7c025200 
              + 0000000000000000000000000000000000000000000000000000000000000002
              + 0000000000000000000000000000000000000000000000000000000000000123
              + ...
↓
Total bytes: 196 (hex chars ÷ 2)
```

### Why Size Matters

1. **Gas Cost**: Calldata costs ~16 gas per byte
2. **Optimization**: Smaller calldata = cheaper transactions
3. **Usage Patterns**: Popular param counts have modes
4. **Efficiency**: Median shows typical usage

### Size Example

```
Mean: 852.30 bytes
Median: 768 bytes
Mode: 512 bytes (appears 892 times)

Interpretation:
• Most common: 512 bytes (typical swap)
• Average: 852 bytes (some complex calls)
• Middle value: 768 bytes (balanced usage)
```

---

## Block Numbers & Ranges

### What is a Block?

**Block** = Container for transactions on Ethereum

```
Block #20245342
├─ Timestamp: 2024-06-23 15:10:25
├─ 120 transactions
├─ Previous block: #20245341
└─ Next block: #20245343
```

### New Block Every ~12 Seconds

```
Per minute:     ~5 blocks
Per hour:       ~300 blocks
Per day:        ~7,200 blocks
Per week:       ~50,000 blocks
Per month:      ~200,000 blocks
```

### Why Block Range Matters

When you run continuous analysis for 60 minutes:

```
Start: Block 20245342 (15:10:25)
End:   Block 20246582 (16:10:25)
Total: 1,240 blocks scanned

Why it matters:
• Reproducible (can verify data again)
• Auditable (exactly what was analyzed)
• Compliance (proof of analysis range)
• Performance (can estimate blocks/min)
```

---

## Transaction Calldata vs Event Data

### Transaction Calldata

```json
{
  "hash": "0xabcd...",
  "to": "0x1111...",
  "input": "0x7c025200...",
  "value": "0",
  "gas": "150000"
}
```

**Characteristics:**
- What: Raw function call data
- Size: 196 to 4096+ bytes
- Available: In all transactions
- Filtering: Client-side only
- Cost: 16 gas/byte

### Event Data

```json
{
  "address": "0x1111...",
  "topics": [
    "0x1234...",  // Event signature
    "0xaaaa...",  // Indexed param 1
    "0xbbbb..."   // Indexed param 2
  ],
  "data": "0x...",
  "blockNumber": "0x1345678"
}
```

**Characteristics:**
- What: Structured log entry
- Size: Much smaller than calldata
- Available: If function emits event
- Filtering: Server-side (getLogs)
- Cost: 8 gas/byte

### When to Use Each

| Situation | Use |
|-----------|-----|
| Function emits events | Query events (getLogs) |
| Function has no events | Query calldata (txlist) |
| Need parameter details | Query calldata |
| Only need event data | Query events |
| Analyzing token transfers | Query Transfer event |
| Analyzing custom logic | Query calldata |

---

## Rate Limits & API Choices

### Etherscan (Current Default)

```
Free tier:   5 calls/second
Paid tier:   Up to 100 calls/second
Price:       ~$0.0001 per call
Limit:       1 month data
```

For 1 contract over 1 month:
- Blocks: ~200,000
- Transactions: ~1,000,000s
- API calls: ~5-10 (batched)
- Time: ~2 seconds
- Cost: Free

### Direct RPC (Alchemy, Infura)

```
Alchemy free:    10 calls/second
Alchemy paid:    10,000+ calls/second
Infura free:     100 calls/second
Infura paid:     Unlimited

Setup: Get API key (free)
Cost:  Free tier available
```

For same analysis:
- API calls: ~2-5 (filtered better)
- Time: ~1 second
- Cost: Free (no paid tier needed)

### The Graph

```
Query speed:  Instant (pre-indexed)
Rate limits:  None
Cost:         Free
Setup:        Visit thegraph.com

For same analysis:
- API calls: 1 (single query)
- Time: <100ms
- Cost: Free
```

---

## Continuous Analysis

### What It Does

```
Start: 15:10:25 (Block 20245342)
│
├─ Poll every 5 seconds
├─ Fetch new transactions
├─ Analyze matching functions
├─ Update running statistics
│
End:   16:10:25 (Block 20246582)

Result: All transactions during this hour
```

### Why Run Continuously?

1. **Real-time monitoring**: See activity as it happens
2. **Complete data**: Don't miss any transactions
3. **Block range proof**: Know exactly what was analyzed
4. **Auditable**: Can reproduce with same block range
5. **Compliance**: Proof of analysis period

### Example Output

```
Duration: 2024-06-23 15:10:25 to 16:10:25
Blocks: 20245342 → 20246582 (1,240 blocks)
Transactions: 4,521 analyzed

Function stats:
  swap: Count=3,200, Mean=512.3, Median=512
  fillOrder: Count=1,321, Mean=768.2, Median=768
```

---

## Key Takeaways

✅ **Events** = Structured logs, easy to filter, not all functions emit them

✅ **Function calls** = Raw calldata, harder to filter, work for any function

✅ **Selectors** = First 4 bytes of function signature hash

✅ **Calldata size** = Function selector + parameters

✅ **Blocks** = ~1 every 12 seconds, track them for reproducibility

✅ **Rate limits** = Etherscan (5/sec), RPC (10-100/sec), The Graph (∞)

✅ **Continuous mode** = Proves audit trail of analysis period

---

**Need more help?** Check API_FILTERING_GUIDE.md or visit https://ethereum.org
