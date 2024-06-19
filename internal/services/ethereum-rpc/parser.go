package ethereumRPC

import (
	"ethereum-parser/internal/services/models"
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
	Storage        StorageService
}

type StorageService interface {
	SubscribeToAddress(address string) bool
	GetLastParsedBlock() int
	SetLastParsedBlock(blockNum int)
	IsSubscribed(address string) bool
	AddTXtoAddressRealTime(blockNum int, tx models.Transaction, address string)
	AddTXtoAddress(tx models.Transaction, address string)
	GetTransactionsByAddress(address string) []models.Transaction

	GetSubscriptionsStorage() map[string]bool
	GetTXsPerAddressOfLatestBlockStorage() map[int]map[string][]models.Transaction
	GetTXsPerAddressTotalStorage() map[string][]models.Transaction
}

func NewParser(storage StorageService) Parser {
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
