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
	mu                         sync.Mutex
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
	defer s.mu.Unlock()

	s.subscriptions[address] = true
	if _, ok := s.subscriptions[address]; !ok {
		return false
	}

	return true
}

func (s *storage) GetLastParsedBlock() int {
	return s.currentParsedBlock
}

func (s *storage) SetLastParsedBlock(blockNum int) {
	s.currentParsedBlock = blockNum
}

func (s *storage) IsSubscribed(address string) bool {
	isSubscribed, _ := s.subscriptions[address]
	if isSubscribed {
		return true
	}

	return false
}

func (s *storage) AddTXtoAddressRealTime(blockNum int, tx models.Transaction, address string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.tXsPerAddressOfLatestBlock[blockNum]) == 0 {
		s.tXsPerAddressOfLatestBlock[blockNum] = make(map[string][]models.Transaction)
	}
	// for easier testing let's consider that all transactions, including related to unsubscribed addresses, will be stored
	s.tXsPerAddressOfLatestBlock[blockNum][address] = append(s.tXsPerAddressOfLatestBlock[blockNum][address], tx)
	/*
		if s.IsSubscribed(address) {
			s.mu.Lock()
			defer s.mu.Unlock()

			s.tXsPerAddressOfLatestBlock[blockNum][address] = append(s.tXsPerAddressOfLatestBlock[blockNum][address], tx)
		}
	*/
}

func (s *storage) AddTXtoAddress(tx models.Transaction, address string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// for easier testing let's consider that all transactions, including related to unsubscribed addresses, will be stored
	s.tXsPerAddressTotal[address] = append(s.tXsPerAddressTotal[address], tx)
	/*
		if s.IsSubscribed(address) {
			s.mu.Lock()
			defer s.mu.Unlock()

			s.tXsPerAddressTotal[address] = append(s.tXsPerAddressTotal[address], tx)
		}
	*/
}

func (s *storage) GetTransactionsByAddress(address string) []models.Transaction {
	// for simplicity only returns the TXs in the latest block(CURRENT BLOCK);
	// Otherwise, tXsPerAddressTotal can be used to store all the TXs bound to an address.

	return s.tXsPerAddressOfLatestBlock[s.GetLastParsedBlock()][address]
	//return s.tXsPerAddressTotal[address]
}

func (s *storage) GetSubscriptionsStorage() map[string]bool {
	return s.subscriptions
}

func (s *storage) GetTXsPerAddressOfLatestBlockStorage() map[int]map[string][]models.Transaction {
	return s.tXsPerAddressOfLatestBlock
}

func (s *storage) GetTXsPerAddressTotalStorage() map[string][]models.Transaction {
	return s.tXsPerAddressTotal
}
