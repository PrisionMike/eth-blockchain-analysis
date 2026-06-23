# Architecture & How It Works

## Data Flow

```
┌─────────────────────────────────────┐
│  Config File                        │
│  (JSON/YAML/TOML/CSV/TXT)          │
│  - Contract addresses               │
│  - Function selectors               │
│  - Time range (days_back)           │
└──────────────────┬──────────────────┘
                   │
                   ▼
        ┌──────────────────────┐
        │  Config Parser       │
        │  (main.go)          │
        │ - JSON unmarshaler   │
        │ - CSV reader         │
        │ - YAML/TOML support  │
        └──────────────────────┘
                   │
                   ▼
        ┌──────────────────────────────┐
        │  Etherscan API Fetcher       │
        │  (analyzer.go)              │
        │ - Query transactions         │
        │ - Filter by selector         │
        │ - Extract calldata size      │
        └──────────────────────┬───────┘
                               │
                ┌──────────────┴──────────────┐
                │                             │
    Per Address: ~13s query time         Results grouped
    Rate limit:  ~5 calls/sec            by function
                │
                ▼
    ┌──────────────────────────┐
    │  Blockchain Transaction  │
    │  Data (per contract)     │
    │                          │
    │ Address: 0x1111...       │
    │ Transactions: [...]      │
    │ Filtered by selector     │
    └────────────┬─────────────┘
                 │
                 ▼
        ┌──────────────────────────┐
        │  Statistics Calculator    │
        │  (analyzer.go)           │
        │                          │
        │  For each function:      │
        │  - Collect calldata sizes│
        │  - Sort sizes            │
        │  - Calculate:            │
        │    • Min/Max             │
        │    • Mean (avg)          │
        │    • Median (middle)     │
        │    • Mode (most common)  │
        └────────────┬─────────────┘
                     │
                     ▼
        ┌────────────────────────┐
        │  Results               │
        │  (per contract/func)   │
        └────────────┬───────────┘
                     │
        ┌────────────┴──────────────────┬───────────────┐
        │                               │               │
        ▼                               ▼               ▼
    JSON Output                    CSV Output      Text Output
    (parseable)              (spreadsheet)      (human-readable)
