# Ethereum-Blockchain-Parser
# Ethereum Blockchain Parser

## Overview

This project implements an Ethereum blockchain parser in Go, which allows users to query transactions for subscribed addresses. The parser interacts with the Ethereum blockchain using JSONRPC and exposes a public interface for external usage via HTTP.

## Features

- **GetCurrentBlock**: Retrieves the last parsed block.
- **Subscribe**: Allows subscribing to specific Ethereum addresses to observe their transactions.
- **UpdateBlockNumber**: Continuously polls for new blocks and updates the current block number.

## Limitations

- Implemented in Go language.
- Avoids usage of external libraries.
- Uses Ethereum JSONRPC to interact with the Ethereum Blockchain.
- Utilizes memory storage for storing data, with easy extendability for future storage support.
- Exposes a public interface via HTTP.

## How to Use

1. **Clone the Repository:**
    ```bash
    git clone https://github.com/your-username/ethereum-parser.git
    cd ethereum-parser
    ```

2. **Run the Server:**
    ```bash
    go run main.go
    ```

3. **Endpoints:**
    - **/currentBlock**: Get the current block number.
        ```bash
        curl http://localhost:8080/currentBlock
        ```
    - **/subscribe**: Subscribe to an address.
        ```bash
        curl "http://localhost:8080/subscribe?address=<ethereum_address>"
        ```

## Code Explanation

```go
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
