package random

import (
	"testing"

	geth_client "github.com/ethereum/go-ethereum/ethclient"
)

func TestGetRandomWalletBalance(t *testing.T) {
	client, err := geth_client.Dial("wss://mainnet.gateway.tenderly.co")
	client2, err2 := geth_client.Dial("wss://bsc-rpc.publicnode.com")
	if err != nil {
		panic(err)
	}
	if err2 != nil {
		panic(err2)
	}
	balance, _, _ := getRandomWalletBalance([2]*geth_client.Client{client, client2})
	if balance != 0 {
		t.Errorf("got %q, wanted %q", balance, 0)
	}
}
