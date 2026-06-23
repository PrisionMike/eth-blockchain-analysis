package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Analyzer struct {
	apiKey      string
	network     string
	rpcEndpoint string
}

type EtherscanTx struct {
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	Input       string `json:"input"`
	BlockNumber string `json:"blockNumber"`
	TimeStamp   string `json:"timeStamp"`
}

type EtherscanResponse struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Result  []EtherscanTx `json:"result"`
}

func NewAnalyzer(apiKey, network, rpcEndpoint string) *Analyzer {
	if rpcEndpoint == "" {
		rpcEndpoint = "https://eth-mainnet.g.alchemy.com/v2/" // Placeholder - user should provide
	}
	return &Analyzer{
		apiKey:      apiKey,
		network:     network,
		rpcEndpoint: rpcEndpoint,
	}
}

func (a *Analyzer) AnalyzeContinuous(contracts []ContractSpec, durationMins int) (*AnalysisResult, error) {
	startTime := time.Now()
	endTime := startTime.Add(time.Duration(durationMins) * time.Minute)
	
	results := &AnalysisResult{
		Stats:       make(map[string][]Stats),
		StartTime:   startTime,
		TransactionsAnalyzed: 0,
	}

	fmt.Printf("📊 Running continuous analysis for %d minute(s)...\n", durationMins)
	fmt.Printf("⏱️  Start: %s | End: %s\n\n", startTime.Format("15:04:05"), endTime.Format("15:04:05"))

	// Track block range
	var minBlock, maxBlock uint64

	// Fetch transactions continuously
	ticker := time.NewTicker(5 * time.Second) // Check every 5 seconds
	defer ticker.Stop()

	allData := make(map[string][]EtherscanTx) // Contract addr -> transactions

	for {
		now := time.Now()
		if now.After(endTime) {
			break
		}

		for _, contract := range contracts {
			fmt.Printf("⏳ [%s] Fetching %s...\n", now.Format("15:04:05"), contract.Name)
			
			txs, err := a.fetchTransactions(contract.Address, 7) // Always fetch last week for comparison
			if err != nil {
				fmt.Printf("⚠️  Error fetching %s: %v\n", contract.Name, err)
				continue
			}

			// Extract block numbers
			for _, tx := range txs {
				blockNum, _ := strconv.ParseUint(tx.BlockNumber, 10, 64)
				if blockNum > maxBlock {
					maxBlock = blockNum
				}
				if minBlock == 0 || blockNum < minBlock {
					minBlock = blockNum
				}
			}

			allData[contract.Address] = append(allData[contract.Address], txs...)
			fmt.Printf("   ✓ Found %d transactions\n", len(txs))
		}

		// Wait or check timeout
		select {
		case <-ticker.C:
			remaining := endTime.Sub(now).Minutes()
			if remaining > 0 {
				fmt.Printf("⏳ Continuing (%.1f min remaining)...\n\n", remaining)
			}
		}
	}

	results.EndTime = time.Now()
	results.StartBlock = minBlock
	results.EndBlock = maxBlock
	results.BlocksScanned = maxBlock - minBlock + 1

	// Process all collected data
	for _, contract := range contracts {
		txs := allData[contract.Address]
		
		fmt.Printf("\n📈 Processing %s (%d transactions)...\n", contract.Name, len(txs))

		// Group by function selector
		functionData := make(map[string][]int) // selector -> sizes

		for _, tx := range txs {
			if len(tx.Input) < 10 {
				continue
			}

			selector := tx.Input[:10]
			for _, fn := range contract.Functions {
				if strings.EqualFold(selector, fn.Selector) {
					size := (len(tx.Input) - 2) / 2
					functionData[fn.Selector] = append(functionData[fn.Selector], size)
					results.TransactionsAnalyzed++
					break
				}
			}
		}

		// Calculate stats
		var contractStats []Stats
		for _, fn := range contract.Functions {
			sizes := functionData[fn.Selector]
			if len(sizes) == 0 {
				continue
			}

			stats := calculateStats(fn.Name, fn.Selector, sizes)
			contractStats = append(contractStats, stats)
		}

		results.Stats[contract.Name] = contractStats
	}

	return results, nil
}

func (a *Analyzer) fetchTransactions(address string, daysBack int) ([]EtherscanTx, error) {
	// Normalize address
	address = strings.ToLower(address)
	if strings.HasPrefix(address, "0x") {
		address = address[2:]
	}
	address = "0x" + address

	// Calculate start block (roughly ~13 seconds per block)
	startTimestamp := time.Now().Add(-time.Duration(daysBack*24) * time.Hour).Unix()

	url := fmt.Sprintf(
		"https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&sort=asc&apikey=%s",
		address,
		a.apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ethResp EtherscanResponse
	if err := json.Unmarshal(body, &ethResp); err != nil {
		return nil, err
	}

	if ethResp.Status != "1" {
		return nil, fmt.Errorf("etherscan error: %s", ethResp.Message)
	}

	// Filter transactions by time
	var filtered []EtherscanTx
	for _, tx := range ethResp.Result {
		timestamp, _ := strconv.ParseInt(tx.TimeStamp, 10, 64)
		if timestamp >= startTimestamp {
			filtered = append(filtered, tx)
		}
	}

	return filtered, nil
}

func calculateStats(name, selector string, sizes []int) Stats {
	if len(sizes) == 0 {
		return Stats{}
	}

	sort.Ints(sizes)

	// Min/Max
	minSize := sizes[0]
	maxSize := sizes[len(sizes)-1]

	// Mean
	sum := 0
	for _, s := range sizes {
		sum += s
	}
	meanSize := float64(sum) / float64(len(sizes))

	// Median
	var medianSize float64
	if len(sizes)%2 == 0 {
		medianSize = float64(sizes[len(sizes)/2-1]+sizes[len(sizes)/2]) / 2
	} else {
		medianSize = float64(sizes[len(sizes)/2])
	}

	// Mode
	freq := make(map[int]int)
	for _, s := range sizes {
		freq[s]++
	}

	modeSize := sizes[0]
	modeCount := 0
	for size, count := range freq {
		if count > modeCount {
			modeCount = count
			modeSize = size
		}
	}

	return Stats{
		FunctionName: name,
		Selector:     selector,
		Count:        len(sizes),
		MinSize:      minSize,
		MaxSize:      maxSize,
		MeanSize:     meanSize,
		MedianSize:   medianSize,
		ModeSize:     modeSize,
		ModeCount:    modeCount,
	}
}
