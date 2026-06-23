# Blockchain Data Query Methods & Filtering Guide

This document explains different ways to query Ethereum data and the differences between **events** and **function calls**.

## Events vs Function Calls

### Function Calls (What We Analyze)

**Function Call** = When someone invokes a smart contract function

```
User → Transaction with calldata → Smart Contract Function
        └─ calldata hex: "0x7c025200aabbcc..."
           └─ First 4 bytes (selector): 0x7c025200 = "swap()"
           └─ Rest (arguments): aabbcc... = encoded parameters
```

**Characteristics:**
- Stored in transaction `input` field
- Hard to filter at API level (need to download full txs)
- What we use in this tool
- Can be any function (doesn't have to emit events)

**Example:**
```
Function: swap(address[] tokens, uint256 amount)
Calldata: 0x7c025200 + encoded(tokens, amount)
Size: 196-4096 bytes depending on parameter count
```

### Events (Emitted Logs)

**Event** = Structured data emitted by a smart contract

```
Smart Contract → emit Swap(address indexed from, address indexed to, uint256 amount)
                 └─ Stored as Log/Event
                 └─ Topic0 (signature): keccak256("Swap(...)")
                 └─ Easy to query!
```

**Characteristics:**
- Indexed separately in blockchain
- Easy to filter via Etherscan `getLogs` API
- **Only works if the function emits an event**
- Limited to data that's actually emitted
- Cleaner filtering but limited scope

**Example:**
```
Event: event Swap(address indexed token0, address indexed token1, uint256 amount)
Topics:
  Topic0: 0x1234... (function signature hash)
  Topic1: 0xaaaa... (token0 address)
  Topic2: 0xbbbb... (token1 address)
```

---

## Query Methods

### Method 1: Etherscan `txlist` (Current Implementation)

**What it does:** Fetches ALL transactions to/from an address

```
Request:
https://api.etherscan.io/api?module=account&action=txlist&address=0x1111...&apikey=KEY

Returns: [
  {tx1: from, to, input: "0x7c025200...", blockNumber: 20245342},
  {tx2: ...},
  {tx3: ...},
  ...
]

Then we filter by:
1. Check input starts with our selector (0x7c025200)
2. Extract input length
3. Calculate size
```

**Pros:**
- ✅ Simple to use
- ✅ No extra setup needed
- ✅ Works for ALL functions (even non-emitting)

**Cons:**
- ❌ Returns ALL transactions (wasting API calls)
- ❌ No server-side filtering
- ❌ Rate limited (5 calls/sec free tier)
- ❌ Slow for popular contracts

**Cost:** ~100 API calls for 1 contract over 1 month

---

### Method 2: Etherscan `getLogs` (Better for Events)

**What it does:** Filters by contract + event signature at API level

```
Request:
https://api.etherscan.io/api?module=logs&action=getLogs&address=0x1111...&topic0=0x1234...&apikey=KEY

Returns: [
  {log1: topics: [0x1234...], data: "0xaaaa...", blockNumber: 20245342},
  {log2: ...},
  ...
]
```

**Pros:**
- ✅ Server-side filtering (fewer results)
- ✅ Much faster
- ✅ Smarter queries

**Cons:**
- ❌ **Only works if function emits events**
- ❌ Can't get raw calldata (only event data)
- ❌ Missing functions that don't emit

**When to use:**
- Contract has events you want to track
- Want to monitor token transfers (Transfer event)
- Need structured event data

**Example:**
```json
{
  "address": "0x1111111254fb6c44bac0bed2854e76f90643097d",
  "functions": [
    {
      "name": "swap",
      "selector": "0x7c025200",
      "event_topic": "0x1234567890abcdef..."  // Event signature
    }
  ]
}
```

---

### Method 3: Direct RPC Node (Most Powerful)

**What it does:** Query the blockchain directly via RPC

```
POST https://eth-mainnet.g.alchemy.com/v2/YOUR_KEY

Request body:
{
  "jsonrpc": "2.0",
  "method": "eth_getLogs",
  "params": [{
    "address": "0x1111...",
    "fromBlock": "0x1234567",
    "toBlock": "0x1234999",
    "topics": ["0x1234..."]
  }],
  "id": 1
}

Response: [
  {address: "0x1111...", data: "0x...", topics: [...], ...},
  ...
]
```

**Pros:**
- ✅ Maximum flexibility
- ✅ Advanced filtering (multiple conditions)
- ✅ Higher rate limits (10-200 req/sec)
- ✅ No rate limiting issues

**Cons:**
- ⚠️ Need API key (but free tier available)
- ⚠️ Requires RPC URL knowledge
- ⚠️ More setup

**Free options:**
- Alchemy: https://alchemy.com (300M calls/month free)
- Infura: https://infura.io (100k/day free)
- Ankr: https://ankr.com (200 req/sec free)

**Setup:**
```json
{
  "rpc_endpoint": "https://eth-mainnet.g.alchemy.com/v2/YOUR_ALCHEMY_KEY",
  ...
}
```

---

### Method 4: The Graph (Best for Complex Analysis)

**What it does:** Pre-indexed blockchain data with powerful queries

```
GraphQL Query:
query {
  transactions(
    where: {
      to: "0x1111111254fb6c44bac0bed2854e76f90643097d"
      from: "0x..."
    }
    orderBy: blockNumber
    orderDirection: desc
  ) {
    id
    from
    to
    input
    blockNumber
    timestamp
  }
}
```

**Pros:**
- ✅ Pre-indexed (fast)
- ✅ Complex filtering
- ✅ No rate limits
- ✅ Subgraphs for specific protocols

**Cons:**
- ⚠️ Requires subgraph setup
- ⚠️ Data lags real-time

**When to use:**
- Complex multi-contract analysis
- Want sophisticated filtering
- Need protocol-specific data

**Popular subgraphs:**
- Uniswap: `uniswap/uniswap-v3`
- 1inch: `1inch-protocol/1inch-subgraph`
- Aave: `aave/protocol-v2`

Website: https://thegraph.com

---

## Comparison Table

| Feature | txlist | getLogs | RPC | The Graph |
|---------|--------|---------|-----|-----------|
| **Setup** | ⭐ Easy | ⭐ Easy | ⭐⭐ Medium | ⭐⭐⭐ Hard |
| **Speed** | Slow | Fast | Very Fast | Instant |
| **Rate Limits** | 5/sec | 5/sec | 10-200/sec | ∞ |
| **Filter by function** | ❌ | ✅ (events only) | ✅ | ✅ |
| **Works without events** | ✅ | ❌ | ✅ | ✅ |
| **API cost** | Free | Free | Free tier | Free |
| **Best for** | Simple analysis | Event tracking | High-volume | Complex queries |

---

## Current Tool Implementation

### Current (v1.0)
- Uses: **Etherscan `txlist`**
- Filtering: Client-side (after download)
- API calls: ~100 per contract per month
- Rate limited: Yes (5/sec)

### Recommended Next Steps

**For better performance, implement:**

1. **Hybrid approach** (Recommended)
   - Use `getLogs` if function has events
   - Fall back to `txlist` for pure functions
   - Smart routing based on contract type

2. **Add RPC support**
   ```go
   if config.RpcEndpoint != "" {
     use Direct RPC with getLogs
   } else {
     use Etherscan txlist
   }
   ```

3. **Cache known events**
   - Build a map of "contract → emitted events"
   - Prioritize event-based queries

4. **The Graph integration** (Advanced)
   - For multi-protocol analysis
   - Pre-indexed data, instant results

---

## Example: Filtering Strategies

### Scenario 1: Track Uniswap Swaps (Has Events)

```json
{
  "address": "0x1111111254fb6c44bac0bed2854e76f90643097d",
  "functions": [
    {
      "name": "swap",
      "selector": "0x7c025200",
      "has_event": true,
      "event_signature": "event Swap(address indexed sender, address indexed recipient, int256 amount0, int256 amount1)"
    }
  ]
}
```

**Best query:** `getLogs` with event topic
**Why:** Much faster, exact results

---

### Scenario 2: Track Fallback Function (No Event)

```json
{
  "address": "0xdeadbeef...",
  "functions": [
    {
      "name": "fallback",
      "selector": null,
      "has_event": false
    }
  ]
}
```

**Best query:** Direct RPC or `txlist`
**Why:** Can't use events (none emitted)

---

### Scenario 3: High-Volume Analysis (1000+ calls/sec)

```json
{
  "rpc_endpoint": "https://eth-mainnet.g.alchemy.com/v2/KEY",
  "use_the_graph": true,
  "contracts": [...]
}
```

**Best query:** The Graph + RPC parallel
**Why:** No rate limits, fast

---

## Setup Your Preferred Method

### Option A: Stick with Etherscan (Default)
```json
{
  "etherscan_api_key": "YOUR_KEY",
  "network": "mainnet"
}
```
- Simple
- Works well for <5 contracts
- OK for 1-week analysis

### Option B: Use Direct RPC (Recommended)
```json
{
  "etherscan_api_key": "YOUR_KEY",
  "rpc_endpoint": "https://eth-mainnet.g.alchemy.com/v2/YOUR_ALCHEMY_KEY",
  "network": "mainnet"
}
```
- Better filtering
- Higher rate limits
- Good for 10+ contracts

### Option C: Add The Graph Support (Advanced)
```json
{
  "use_the_graph": true,
  "subgraph": "uniswap/uniswap-v3",
  "contracts": [...]
}
```
- Pre-indexed
- Instant results
- Best for multi-contract analysis

---

## Troubleshooting

**"Only works for events" - What does this mean?**
- Some smart contract functions don't emit events
- The `getLogs` API only sees events
- For functions without events, use `txlist` or RPC

**How do I know if my function emits events?**
1. Check the ABI (Application Binary Interface)
2. Look for `event` declarations
3. Use Etherscan contract page → Code tab

**Should I use events or calldata?**
| Use Events If | Use Calldata If |
|---|---|
| Function emits events | Function has no events |
| Want structured data | Need raw parameters |
| Fewer API calls | Want complete info |
| Don't care about params | Need all details |

**My function doesn't emit events. How do I filter?**
- Use `txlist` (current method)
- Upgrade to direct RPC (better)
- Use The Graph (if available)

---

## Future Improvements

1. ✨ Auto-detect if contract has events
2. ✨ Smart routing (events → `getLogs`, pure functions → RPC)
3. ✨ The Graph integration
4. ✨ Event caching for known contracts
5. ✨ Multi-provider fallback (if one is down)

---

**Questions?** Check the main README.md or run `./eth-analysis -help`
