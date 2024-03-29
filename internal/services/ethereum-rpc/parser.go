package ethereumRPC

import (
	"ethereum-parser/internal/services/models"
	inmemory "ethereum-parser/internal/storage/in-memory"
	"log"
	"os"
)

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []models.Transaction
}

type parser struct {
	EthereumRpcURL string // Just for simplicity; better to do configuration as dependency injection
	Storage        inmemory.Storage
}

func NewParser(storage inmemory.Storage) Parser {
	return &parser{
		EthereumRpcURL: os.Getenv("ETHEREUM_RPC_ENDPOINT_URL"),
		Storage:        storage,
	}
}

// GetCurrentBlock last parsed block
func (p *parser) GetCurrentBlock() int {
	return p.Storage.GetLastParsedBlock()
}

// Subscribe add address to observer
func (p *parser) Subscribe(address string) bool {
	return p.Storage.SubscribeToAddress(address)
}

// GetTransactions list of inbound or outbound transactions for an address
func (p *parser) GetTransactions(address string) []models.Transaction {
	// Based on the assumption that notifications are available for subscribed addresses
	if !p.Storage.IsSubscribed(address) {
		log.Println("the address is not subscribed")
		return nil
	}

	return p.Storage.GetTransactionsByAddress(address)
}
