package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	geth_hex_util "github.com/ethereum/go-ethereum/common/hexutil"
	geth_client "github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/crypto"
)

func getRandomWalletBalance(clients [2]*geth_client.Client) (int, string, string) {

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	// fmt.Println("Private key", geth_hex_util.Encode(crypto.FromECDSA(privateKey))[2:])

	// Sign the message
	// signature, err := sk.Sign(message)
	// if err != nil {
	// 	panic(err)
	// }

	// // Verify the signature
	// valid, err := pk.Verify(signature, message)
	// if err != nil {
	// 	panic(err)
	// }

	publicKey := privateKey.Public()

	ecdsaPublicKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("Type assertion failed")
		return 0, geth_hex_util.Encode(crypto.FromECDSA(privateKey)), crypto.PubkeyToAddress(*ecdsaPublicKey).Hex()
	}
	// fmt.Println("Public key", geth_hex_util.Encode(crypto.FromECDSAPub(ecdsaPublicKey)))

	balance, err := clients[1].BalanceAt(context.Background(), crypto.PubkeyToAddress(*ecdsaPublicKey), nil)
	if err != nil {
		panic(err)
	}
	balance2, err := clients[1].BalanceAt(context.Background(), crypto.PubkeyToAddress(*ecdsaPublicKey), nil)
	if err != nil {
		panic(err)
	}
	return int(balance.Int64()) + int(balance2.Int64()), geth_hex_util.Encode(crypto.FromECDSA(privateKey)), crypto.PubkeyToAddress(*ecdsaPublicKey).Hex()
}

func getCurrentBlock(client geth_client.Client) (int, string, string) {
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))

	if err != nil {
		panic(err)
	}
	return int(blockNumber), block.Hash().Hex(), string(block.Nonce())
}

func main() {
	client, err := geth_client.Dial("wss://mainnet.gateway.tenderly.co")
	client2, err2 := geth_client.Dial("wss://bsc-rpc.publicnode.com")
	if err != nil {
		panic(err)
	}
	if err2 != nil {
		panic(err2)
	}

	quit := make(chan bool)
	res := make(chan int)

	go func() {
		idx := 0
		balance, privateKey, pub := getRandomWalletBalance([2]*geth_client.Client{client, client2})
		if balance > 0 { // Condition to stop the goroutine
			fmt.Println("Public key", privateKey, pub)
			panic("Found") // Exit the loop in the main goroutine
		}
		for {
			select {
			case <-quit:
				fmt.Println("Detected quit signal!")
				return
			default:
				fmt.Println("goroutine is doing stuff..")
				res <- idx
				idx++
			}
		}
	}()

	for r := range res {
		balance, privateKey, pub := getRandomWalletBalance([2]*geth_client.Client{client, client2})
		if balance > 0 { // Condition to stop the goroutine
			fmt.Println("Public key", privateKey, pub)
			getCurrentBlock
			quit <- true
			break // Exit the loop in the main goroutine
		}
		fmt.Println("I received: ", r)
	}
}
