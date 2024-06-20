package inmemory

import (
	"ethereum-parser/internal/services/models"
	"sync"
)

type Storage interface {
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

func (s *storage) SubscribeToAddress(address string) bool {
	s.mu.Lock()
	s.subscriptions[address] = true
	s.mu.Unlock()

	return true
}

func (s *storage) GetLastParsedBlock() int {
	s.mu.RLock()
	currentBlockNum := s.currentParsedBlock
	s.mu.RUnlock()

	return currentBlockNum
}

func (s *storage) SetLastParsedBlock(blockNum int) {
	s.mu.Lock()
	s.currentParsedBlock = blockNum
	s.mu.Unlock()
}

func (s *storage) IsSubscribed(address string) bool {
	s.mu.RLock()
	isSubscribed := s.subscriptions[address]
	s.mu.RUnlock()

	return isSubscribed
}

func (s *storage) AddTXtoAddressRealTime(blockNum int, tx models.Transaction, address string) {
	if len(s.tXsPerAddressOfLatestBlock[blockNum]) == 0 {
		s.tXsPerAddressOfLatestBlock[blockNum] = make(map[string][]models.Transaction)
	}

	if s.IsSubscribed(address) {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.tXsPerAddressOfLatestBlock[blockNum][address] = append(s.tXsPerAddressOfLatestBlock[blockNum][address], tx)
	}
}

func (s *storage) AddTXtoAddress(tx models.Transaction, address string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tXsPerAddressTotal[address]; !ok {
		s.tXsPerAddressTotal[address] = make([]models.Transaction, 0)
	}

	if s.IsSubscribed(address) {
		s.tXsPerAddressTotal[address] = append(s.tXsPerAddressTotal[address], tx)
	}
}

func (s *storage) GetTransactionsByAddress(address string) []models.Transaction {
	s.mu.RLock()
	transactions := s.tXsPerAddressTotal[address]
	s.mu.RUnlock()

	return transactions
}

func (s *storage) GetSubscriptionsStorage() map[string]bool {
	s.mu.RLock()
	subscriptions := s.subscriptions
	s.mu.RUnlock()

	return subscriptions
}

func (s *storage) GetTXsPerAddressOfLatestBlockStorage() map[int]map[string][]models.Transaction {
	s.mu.RLock()
	tXs := s.tXsPerAddressOfLatestBlock
	s.mu.RUnlock()

	return tXs
}

func (s *storage) GetTXsPerAddressTotalStorage() map[string][]models.Transaction {
	s.mu.RLock()
	tXs := s.tXsPerAddressTotal
	s.mu.RUnlock()

	return tXs
}
