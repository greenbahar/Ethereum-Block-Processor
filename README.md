Implementing `Parser` interface to hook this up to notifications service to notify about any incoming/outgoing transactions.

```go
type Parser interface {
	// last parsed block
	GetCurrentBlock() int

	// add address to observer
	Subscribe(address string) bool

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []Transaction
}
```

Using Ethereum JSONRPC to interact with Ethereum Blockchain via `https://cloudflare-eth.com`
