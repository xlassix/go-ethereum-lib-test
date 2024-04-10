package random

import (
	"fmt"
	"testing"

	geth_client "github.com/ethereum/go-ethereum/ethclient"
)

func TestGetRandomWalletBalance(t *testing.T) {
	client, err := geth_client.Dial("wss://eth.drpc.org")
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

func TestCreateWallet(t *testing.T) {
	password := randomString(20)
	address, path := createWallet(password)

	fmt.Println(address, path, len(address))

	if len(address) != 42 {
		t.Errorf("Expected account to be of length %q, but got %q", 42, len(address))
	}
}

func TestDecryptWallet(t *testing.T) {
	password := randomString(20)
	address, path := createWallet(password)

	decryptAddress, _ := decryptWallet(path, password)

	fmt.Println(address, path, len(address))

	if len(address) != 42 {
		t.Errorf("Expected account to be of length %q, but got %q", 42, len(address))
	}

	if decryptAddress != address {
		t.Errorf("Expected account address %q, but got %q", 42, len(address))
	}
}
