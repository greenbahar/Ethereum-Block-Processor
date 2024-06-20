package inmemory

import (
	"context"
	"ethereum-parser/internal/services/models"
	"sync"
)

type Storage interface {
	SubscribeToAddress(ctx context.Context, address string) bool
	GetLastParsedBlock(ctx context.Context) int
	SetLastParsedBlock(ctx context.Context, blockNum int)
	IsSubscribed(ctx context.Context, address string) bool
	AddTXtoAddressRealTime(ctx context.Context, blockNum int, tx models.Transaction, address string)
	AddTXtoAddress(ctx context.Context, tx models.Transaction, address string)
	GetTransactionsByAddress(ctx context.Context, address string) []models.Transaction

	GetSubscriptionsStorage(ctx context.Context) map[string]bool
	GetTXsPerAddressOfLatestBlockStorage(ctx context.Context) map[int]map[string][]models.Transaction
	GetTXsPerAddressTotalStorage(ctx context.Context) map[string][]models.Transaction
}

type storage struct {
	currentParsedBlock         int
	subscriptions              map[string]bool
	tXsPerAddressOfLatestBlock map[int]map[string][]models.Transaction // blockNum: {address: transactionsOfCurrentBlock}
	tXsPerAddressTotal         map[string][]models.Transaction         // {address: allTransactions}
	mu                         sync.RWMutex
}

func NewStorage() Storage {
	return &storage{
		currentParsedBlock:         0,
		subscriptions:              make(map[string]bool),
		tXsPerAddressOfLatestBlock: make(map[int]map[string][]models.Transaction),
		tXsPerAddressTotal:         make(map[string][]models.Transaction),
	}
}

func (s *storage) SubscribeToAddress(ctx context.Context, address string) bool {
	s.mu.Lock()
	s.subscriptions[address] = true
	s.mu.Unlock()

	return true
}

func (s *storage) GetLastParsedBlock(ctx context.Context) int {
	s.mu.RLock()
	currentBlockNum := s.currentParsedBlock
	s.mu.RUnlock()

	return currentBlockNum
}

func (s *storage) SetLastParsedBlock(ctx context.Context, blockNum int) {
	s.mu.Lock()
	s.currentParsedBlock = blockNum
	s.mu.Unlock()
}

func (s *storage) IsSubscribed(ctx context.Context, address string) bool {
	s.mu.RLock()
	isSubscribed := s.subscriptions[address]
	s.mu.RUnlock()

	return isSubscribed
}

func (s *storage) AddTXtoAddressRealTime(ctx context.Context, blockNum int, tx models.Transaction, address string) {
	if len(s.tXsPerAddressOfLatestBlock[blockNum]) == 0 {
		s.tXsPerAddressOfLatestBlock[blockNum] = make(map[string][]models.Transaction)
	}

	if s.IsSubscribed(ctx, address) {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.tXsPerAddressOfLatestBlock[blockNum][address] = append(s.tXsPerAddressOfLatestBlock[blockNum][address], tx)
	}
}

func (s *storage) AddTXtoAddress(ctx context.Context, tx models.Transaction, address string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tXsPerAddressTotal[address]; !ok {
		s.tXsPerAddressTotal[address] = make([]models.Transaction, 0)
	}

	if s.IsSubscribed(ctx, address) {
		s.tXsPerAddressTotal[address] = append(s.tXsPerAddressTotal[address], tx)
	}
}

func (s *storage) GetTransactionsByAddress(ctx context.Context, address string) []models.Transaction {
	s.mu.RLock()
	transactions := s.tXsPerAddressTotal[address]
	s.mu.RUnlock()

	return transactions
}

func (s *storage) GetSubscriptionsStorage(ctx context.Context) map[string]bool {
	s.mu.RLock()
	subscriptions := s.subscriptions
	s.mu.RUnlock()

	return subscriptions
}

func (s *storage) GetTXsPerAddressOfLatestBlockStorage(ctx context.Context) map[int]map[string][]models.Transaction {
	s.mu.RLock()
	tXs := s.tXsPerAddressOfLatestBlock
	s.mu.RUnlock()

	return tXs
}

func (s *storage) GetTXsPerAddressTotalStorage(ctx context.Context) map[string][]models.Transaction {
	s.mu.RLock()
	tXs := s.tXsPerAddressTotal
	s.mu.RUnlock()

	return tXs
}
