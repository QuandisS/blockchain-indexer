package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	client, err := ethclient.Dial("")
	if err != nil {
		log.Fatal(err)
	}

	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Последний блок: %d\n", blockNumber)

	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		log.Fatal(err)
	}

	for _, tx := range block.Transactions() {
		fmt.Printf("Хэш транзакции: %s\n", tx.Hash().Hex())
		fmt.Printf("Получатель: %s\n", tx.To().Hex())
		fmt.Printf("Сумма: %s\n", tx.Value().String())
		fmt.Printf("Лимит газа: %d\n", tx.Gas())
		fmt.Printf("Цена газа: %s\n", tx.GasPrice().String())
		fmt.Printf("Nonce: %d\n", tx.Nonce())
		fmt.Println()
	}
}
