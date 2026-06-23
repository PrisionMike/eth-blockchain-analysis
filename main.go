package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type FunctionSpec struct {
	Name     string `json:"name" yaml:"name" toml:"name"`
	Selector string `json:"selector" yaml:"selector" toml:"selector"`
}

type ContractSpec struct {
	Address   string           `json:"address" yaml:"address" toml:"address"`
	Name      string `json:"name" yaml:"name" toml:"name"`
	Functions []FunctionSpec `json:"functions" yaml:"functions" toml:"functions"`
}

type Config struct {
	EtherscanAPIKey string         `json:"etherscan_api_key" yaml:"etherscan_api_key" toml:"etherscan_api_key"`
	Network         string         `json:"network" yaml:"network" toml:"network"`
	Contracts       []ContractSpec `json:"contracts" yaml:"contracts" toml:"contracts"`
	DaysBack        int            `json:"days_back" yaml:"days_back" toml:"days_back"`
	RPCEndpoint     string         `json:"rpc_endpoint" yaml:"rpc_endpoint" toml:"rpc_endpoint"`
	RunDurationMins int            `json:"run_duration_mins" yaml:"run_duration_mins" toml:"run_duration_mins"`
}

type Stats struct {
	FunctionName string
	Selector     string
	Count        int
	MinSize      int
	MaxSize      int
	MeanSize     float64
	MedianSize   float64
	ModeSize     int
	ModeCount    int
}

type AnalysisResult struct {
	Stats            map[string][]Stats
	StartBlock       uint64
	EndBlock         uint64
	StartTime        time.Time
	EndTime          time.Time
	BlocksScanned    uint64
	TransactionsAnalyzed int
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config

	// Detect format by file extension
	switch {
	case strings.HasSuffix(filename, ".json"):
		err = json.Unmarshal(data, &config)
	case strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml"):
		err = yaml.Unmarshal(data, &config)
	case strings.HasSuffix(filename, ".toml"):
		err = toml.Unmarshal(data, &config)
	case strings.HasSuffix(filename, ".csv"):
		config, err = loadFromCSV(filename)
	case strings.HasSuffix(filename, ".txt"):
		config, err = loadFromTXT(filename)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", filename)
	}

	if err != nil {
		return nil, err
	}

	if config.DaysBack == 0 {
		config.DaysBack = 7
	}

	return &config, nil
}

func loadFromCSV(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	config := Config{
		Network:   "mainnet",
		Contracts: []ContractSpec{},
	}

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return config, err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return config, err
		}

		// Expected columns: address, name, function_name, selector
		if len(record) >= 4 {
			address := record[0]
			name := record[1]
			funcName := record[2]
			selector := record[3]

			// Find or create contract
			var contract *ContractSpec
			for i := range config.Contracts {
				if config.Contracts[i].Address == address {
					contract = &config.Contracts[i]
					break
				}
			}

			if contract == nil {
				config.Contracts = append(config.Contracts, ContractSpec{
					Address:   address,
					Name:      name,
					Functions: []FunctionSpec{},
				})
				contract = &config.Contracts[len(config.Contracts)-1]
			}

			contract.Functions = append(contract.Functions, FunctionSpec{
				Name:     funcName,
				Selector: selector,
			})
		}
	}

	return config, nil
}

func loadFromTXT(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	config := Config{
		Network:   "mainnet",
		Contracts: []ContractSpec{},
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Format: address|name|function_name|selector
		parts := strings.Split(line, "|")
		if len(parts) >= 4 {
			address := strings.TrimSpace(parts[0])
			name := strings.TrimSpace(parts[1])
			funcName := strings.TrimSpace(parts[2])
			selector := strings.TrimSpace(parts[3])

			// Find or create contract
			var contract *ContractSpec
			for i := range config.Contracts {
				if config.Contracts[i].Address == address {
					contract = &config.Contracts[i]
					break
				}
			}

			if contract == nil {
				config.Contracts = append(config.Contracts, ContractSpec{
					Address:   address,
					Name:      name,
					Functions: []FunctionSpec{},
				})
				contract = &config.Contracts[len(config.Contracts)-1]
			}

			contract.Functions = append(contract.Functions, FunctionSpec{
				Name:     funcName,
				Selector: selector,
			})
		}
	}

	return config, nil
}

func main() {
	configFile := flag.String("config", "config.json", "Path to config file")
	outputFormat := flag.String("output", "json", "Output format: json, csv, or text")
	flag.Parse()

	// Load config
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if config.EtherscanAPIKey == "" || config.EtherscanAPIKey == "YOUR_ETHERSCAN_API_KEY" {
		log.Fatal("Please set etherscan_api_key in your config file")
	}

	// Set default run duration
	if config.RunDurationMins == 0 {
		config.RunDurationMins = 5 // Default 5 minutes
	}

	log.Printf("Starting continuous analysis for %d minute(s)...", config.RunDurationMins)

	// Fetch and analyze
	analyzer := NewAnalyzer(config.EtherscanAPIKey, config.Network, config.RPCEndpoint)
	result, err := analyzer.AnalyzeContinuous(config.Contracts, config.RunDurationMins)
	if err != nil {
		log.Fatalf("Analysis failed: %v", err)
	}

	// Output results
	switch *outputFormat {
	case "json":
		outputJSON(result)
	case "csv":
		outputCSV(result)
	case "text":
		outputText(result)
	default:
		log.Fatalf("Unknown output format: %s", *outputFormat)
	}
}

func outputJSON(result *AnalysisResult) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(data))
}

func outputCSV(result *AnalysisResult) {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	writer.Write([]string{"Contract", "Function", "Selector", "Count", "Min Size", "Max Size", "Mean Size", "Median Size", "Mode Size", "Mode Count"})

	for contractName, funcs := range result.Stats {
		for _, stat := range funcs {
			writer.Write([]string{
				contractName,
				stat.FunctionName,
				stat.Selector,
				fmt.Sprintf("%d", stat.Count),
				fmt.Sprintf("%d", stat.MinSize),
				fmt.Sprintf("%d", stat.MaxSize),
				fmt.Sprintf("%.2f", stat.MeanSize),
				fmt.Sprintf("%.2f", stat.MedianSize),
				fmt.Sprintf("%d", stat.ModeSize),
				fmt.Sprintf("%d", stat.ModeCount),
			})
		}
	}
}

func outputText(result *AnalysisResult) {
	fmt.Printf("\n╔════════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                    ANALYSIS RESULTS                               ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════════════════╝\n\n")

	fmt.Printf("Duration: %s to %s\n", result.StartTime.Format("2006-01-02 15:04:05"), result.EndTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Blocks Scanned: %d → %d (Total: %d blocks)\n", result.StartBlock, result.EndBlock, result.BlocksScanned)
	fmt.Printf("Transactions Analyzed: %d\n\n", result.TransactionsAnalyzed)

	for contractName, funcs := range result.Stats {
		fmt.Printf("=== %s ===\n", contractName)
		fmt.Printf("%-30s %-10s %-12s %-10s %-10s %-12s %-12s %-10s\n",
			"Function", "Count", "Min", "Max", "Mean", "Median", "Mode", "Mode Cnt")
		fmt.Println(strings.Repeat("-", 120))

		for _, stat := range funcs {
			fmt.Printf("%-30s %-10d %-12d %-10d %-10.2f %-12.2f %-10d %-8d\n",
				stat.FunctionName, stat.Count, stat.MinSize, stat.MaxSize,
				stat.MeanSize, stat.MedianSize, stat.ModeSize, stat.ModeCount)
		}
		fmt.Println()
	}
}
