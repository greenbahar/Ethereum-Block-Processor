package models

type BlockResponse struct {
	JsonRPC string      `json:"jsonrpc"`
	Result  BlockResult `json:"result"`
	Id      int         `json:"id"`
}

type BlockResult struct {
	Number           string        `json:"number"`
	Hash             string        `json:"hash"`
	MixHash          string        `json:"mixHash"`
	ParentHash       string        `json:"parentHash"`
	Nonce            string        `json:"nonce"`
	Sha3Uncles       string        `json:"sha3Uncles"`
	LogsBloom        string        `json:"logsBloom"`
	TransactionsRoot string        `json:"transactionsRoot"`
	StateRoot        string        `json:"stateRoot"`
	Miner            string        `json:"miner"`
	Difficulty       string        `json:"difficulty"`
	TotalDifficulty  string        `json:"totalDifficulty"`
	ExtraData        string        `json:"extraData"`
	Size             string        `json:"size"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Timestamp        string        `json:"timestamp"`
	Transactions     []Transaction `json:"transactions"`
	Uncles           []string      `json:"uncles"`
	ReceiptsRoot     string        `json:"receiptsRoot"`
}

type Transaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	TransactionIndex string `json:"transactionIndex"`
	Nonce            string `json:"nonce"`
	R                string `json:"r"`
	S                string `json:"s"`
	V                string `json:"v"`
	Type             string `json:"type"`
	YParity          string `json:"yParity"`
	//AccessList           string `json:"accessList"`
	ChainId              string `json:"chainId"`
	MaxFeePerGas         string `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
}
