package syncwithmainnet

import (
	"bytes"
	"encoding/json"
	"ethereum-parser/internal/services/models"
	"ethereum-parser/pkg/utils"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Synchronizer interface {
	SyncWithMainNetViaRPC()
}

type synchronizer struct {
	EthereumRpcURL string
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

func NewSynchronizer(storage StorageService) Synchronizer {
	return &synchronizer{
		EthereumRpcURL: os.Getenv("ETHEREUM_RPC_ENDPOINT_URL"),
		Storage:        storage,
	}
}

func (s *synchronizer) SyncWithMainNetViaRPC() {
	// todo error handling
	blocks, _ := s.getBlocks()
	blockNumber, err := utils.ConvertHexToInt(blocks.Result.Number)
	if err != nil {
		log.Println("error in conversion hex to integer")
	}

	s.Storage.SetLastParsedBlock(blockNumber)

	// todo enhancement: worker pool
	for _, tx := range blocks.Result.Transactions {
		// AddTXtoAddress: FOR NON-REALTIME NOTIFICATION
		s.Storage.AddTXtoAddress(tx, tx.From)
		s.Storage.AddTXtoAddress(tx, tx.To)

		// AddTXtoAddressRealTime: FOR REALTIME NOTIFICATION for the inbound/outbound TXs according to CURRENT BLOCK
		// s.Storage.AddTXtoAddressRealTime(blockNumber, tx, tx.From)
		// s.Storage.AddTXtoAddressRealTime(blockNumber, tx, tx.From)
	}
}

func (s *synchronizer) getBlocks() (*models.BlockResponse, error) {
	var blockHeight interface{}
	lastParsedBlok := s.Storage.GetLastParsedBlock()
	if lastParsedBlok == 0 {
		blockHeight = "latest"
	} else {
		blockHeight = lastParsedBlok
	}

	rpcRequest, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{blockHeight, true},
		"id":      1,
	})
	if err != nil {
		return nil, fmt.Errorf("serialize error: %v", err)
	}

	res, err := http.Post(s.EthereumRpcURL, "application/json", bytes.NewBuffer(rpcRequest))
	if err != nil {
		return nil, fmt.Errorf("serialize error: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error read body of the response: %v", err)
	}

	var blockResponse *models.BlockResponse
	err = json.Unmarshal(body, &blockResponse)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	return blockResponse, nil
}
