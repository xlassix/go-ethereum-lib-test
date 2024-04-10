package random

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/rand"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	geth_hex_util "github.com/ethereum/go-ethereum/common/hexutil"
	geth_client "github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/crypto"
)

func createWallet(password string) (string, string) {
	key := keystore.NewKeyStore("./wallet", keystore.StandardScryptN*8, keystore.StandardScryptP)

	account, err := key.NewAccount(password)

	if err != nil {
		panic(err)
	}

	return account.Address.Hex(), account.URL.Path
}

func randomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()_"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func decryptWallet(path string, password string) (string, string) {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	account, err := keystore.DecryptKey(file, password)
	if err != nil {
		panic(err)
	}
	return account.Address.Hex(), string(crypto.FromECDSA(account.PrivateKey))
}
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

func run() {
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
			quit <- true
			break // Exit the loop in the main goroutine
		}
		fmt.Println("I received: ", r)
	}
}
