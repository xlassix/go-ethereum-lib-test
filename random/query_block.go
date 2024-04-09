package random

import (
	"context"
	"math/big"
	"strconv"

	geth_client "github.com/ethereum/go-ethereum/ethclient"
)

func getCurrentBlock(client *geth_client.Client) (int, string, string) {
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))

	if err != nil {
		panic(err)
	}
	return int(blockNumber), block.Hash().Hex(), strconv.FormatUint(block.Nonce(), 10)
}
