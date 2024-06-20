package ethereumRPC

import (
	"context"
	"ethereum-parser/internal/services/models"
	"log"
	"os"
)

type Parser interface {
	GetCurrentBlock(ctx context.Context) int
	Subscribe(ctx context.Context, address string) bool
	GetTransactions(ctx context.Context, address string) []models.Transaction
}

type parser struct {
	EthereumRpcURL string // Just for simplicity; better to do configuration as dependency injection
	Storage        StorageService
}

type StorageService interface {
	SubscribeToAddress(ctx context.Context, address string) bool
	GetLastParsedBlock(ctx context.Context) int
	SetLastParsedBlock(bctx context.Context, lockNum int)
	IsSubscribed(ctx context.Context, address string) bool
	AddTXtoAddressRealTime(ctx context.Context, blockNum int, tx models.Transaction, address string)
	AddTXtoAddress(ctx context.Context, tx models.Transaction, address string)
	GetTransactionsByAddress(ctx context.Context, address string) []models.Transaction

	GetSubscriptionsStorage(ctx context.Context) map[string]bool
	GetTXsPerAddressOfLatestBlockStorage(ctx context.Context) map[int]map[string][]models.Transaction
	GetTXsPerAddressTotalStorage(ctx context.Context) map[string][]models.Transaction
}

func NewParser(storage StorageService) Parser {
	return &parser{
		EthereumRpcURL: os.Getenv("ETHEREUM_RPC_ENDPOINT_URL"),
		Storage:        storage,
	}
}

// GetCurrentBlock last parsed block
func (p *parser) GetCurrentBlock(ctx context.Context) int {
	return p.Storage.GetLastParsedBlock(ctx)
}

// Subscribe add address to observer
func (p *parser) Subscribe(ctx context.Context, address string) bool {
	return p.Storage.SubscribeToAddress(ctx, address)
}

// GetTransactions list of inbound or outbound transactions for an address
func (p *parser) GetTransactions(ctx context.Context, address string) []models.Transaction {
	// Based on the assumption that notifications are available for subscribed addresses
	if !p.Storage.IsSubscribed(ctx, address) {
		log.Println("the address is not subscribed")
		return nil
	}

	return p.Storage.GetTransactionsByAddress(ctx, address)
}
