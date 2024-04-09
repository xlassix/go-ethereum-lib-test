package random

import (
	"testing"

	geth_client "github.com/ethereum/go-ethereum/ethclient"
)

func TestGetCurrentBlock(t *testing.T) {
	client, err := geth_client.Dial("wss://mainnet.gateway.tenderly.co")
	client2, err2 := geth_client.Dial("wss://bsc-rpc.publicnode.com")
	if err != nil {
		panic(err)
	}
	if err2 != nil {
		panic(err2)
	}
	blockNumber, _, _ := getCurrentBlock(client)
	blockNumber2, _, _ := getCurrentBlock(client2)

	if blockNumber > blockNumber2 {
		t.Errorf("%q, should be greater then  %q", blockNumber, blockNumber2)
	}
}
