package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
	"strings"
	"strconv"	
	"time"
)

type Parser interface {
    GetCurrentBlock() int
    Subscribe(address string) bool
}

type EthParser struct {
    currentBlock int
    addresses    map[string]bool
    mu           sync.Mutex
}

func (p *EthParser) GetCurrentBlock() int {
    p.mu.Lock()
    defer p.mu.Unlock()
    return p.currentBlock
}

func (p *EthParser) Subscribe(address string) bool {
    p.mu.Lock()
    defer p.mu.Unlock()
    if _, exists := p.addresses[address]; exists {
        return false
    }
    p.addresses[address] = true
    return true
}

func NewEthParser() *EthParser {
    return &EthParser{
        currentBlock: 0,
        addresses:    make(map[string]bool),
    }
}
func getLatestBlockNumber() (int, error) {
    resp, err := http.Post("http://localhost:8545", "application/json", strings.NewReader(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`))
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    var result struct {
        Result string `json:"result"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return 0, err
    }

    blockNumber, err := strconv.ParseInt(result.Result[2:], 16, 64)
    if err != nil {
        return 0, err
    }
    return int(blockNumber), nil
}
func (p *EthParser) UpdateBlockNumber() {
    for {
        blockNumber, err := getLatestBlockNumber()
        if err != nil {
            fmt.Println("Error fetching block number:", err)
            continue
        }
        p.mu.Lock()
        p.currentBlock = blockNumber
        p.mu.Unlock()
        time.Sleep(10 * time.Second)  // Poll every 10 seconds
    }
}
func main() {
    parser := NewEthParser()

    go parser.UpdateBlockNumber()

    http.HandleFunc("/currentBlock", func(w http.ResponseWriter, r *http.Request) {
        block := parser.GetCurrentBlock()
        fmt.Fprintf(w, "Current Block: %d", block)
    })

    http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
        address := r.URL.Query().Get("address")
        if address == "" {
            http.Error(w, "Address parameter is required", http.StatusBadRequest)
            return
        }
        success := parser.Subscribe(address)
        if success {
            fmt.Fprintf(w, "Subscribed to address: %s", address)
        } else {
            fmt.Fprintf(w, "Address already subscribed: %s", address)
        }
    })

    fmt.Println("Server running on port 8080")
    http.ListenAndServe(":8080", nil)
}
