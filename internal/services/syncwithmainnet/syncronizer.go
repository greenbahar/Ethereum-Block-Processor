package syncwithmainnet

import (
	"bytes"
	"encoding/json"
	"ethereum-parser/internal/services/models"
	"ethereum-parser/pkg/utils"
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
	lastParsedBlok := s.Storage.GetLastParsedBlock()
	rpcRequest, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{lastParsedBlok, true},
		"id":      1,
	})
	if err != nil {
		log.Print("Serialize error")
	}

	res, err := http.Post(s.EthereumRpcURL, "application/json", bytes.NewBuffer(rpcRequest))
	if err != nil {
		log.Println("Serialize error")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("error read body of the response")
	}

	var blockResponse models.BlockResponse
	err = json.Unmarshal(body, &blockResponse)
	if err != nil {
		log.Println("2", err)
	}

	blockNumber, err := utils.ConvertHexToInt(blockResponse.Result.Number)
	if err != nil {
		log.Println("error in conversion hex to integer")
	}
	s.Storage.SetLastParsedBlock(blockNumber)

	// enhancement: worker pool
	for _, tx := range blockResponse.Result.Transactions {
		// AddTXtoAddress: FOR NON-REALTIME NOTIFICATION
		//s.Storage.AddTXtoAddress(tx, tx.From)
		//s.Storage.AddTXtoAddress(tx, tx.To)

		// AddTXtoAddressRealTime: FOR REALTIME NOTIFICATION for the inbound/outbound TXs according to CURRENT BLOCK
		s.Storage.AddTXtoAddressRealTime(blockNumber, tx, tx.From)
		s.Storage.AddTXtoAddressRealTime(blockNumber, tx, tx.From)
	}

	/*
		// LOG FOR POSTMAN TEST. Get the address from the stdout and use it for API calls
		storage := s.Storage.GetTXsPerAddressOfLatestBlockStorage()
		for key, val := range storage[s.Storage.GetLastParsedBlock()] {
			log.Println("GetTXsPerAddressOfLatestBlockStorage: ", key, val)
			break
		}
	*/
}
