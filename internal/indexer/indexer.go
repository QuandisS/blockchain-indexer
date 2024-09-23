package indexer

import (
	"context"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
)

const workerNum = 5

func Start(rpc string, start int, out string, limit int) {
	client, err := ethclient.Dial(rpc)

	if err != nil {
		log.Fatal("Failed to connect with Ethereum client: ", err)
	}

	lastBlockNum, err := client.BlockNumber(context.Background())

	if err != nil {
		log.Fatal("Failed to get last block number: ", err)
	}

	if uint64(start+limit) <= lastBlockNum {
		log.Println("Indexing from block: ", start, " to block: ", start+limit)
		blockToIndexNumCh := make(chan int)

		wg := &sync.WaitGroup{}
		for i := 0; i < workerNum; i++ {
			log.Println("Starting worker: ", i)
			wg.Add(1)
			go startWorker(client, blockToIndexNumCh, wg)
		}

		for blockNum := start; blockNum < start+limit; blockNum++ {
			blockToIndexNumCh <- blockNum
		}

		close(blockToIndexNumCh)

		wg.Wait()
		log.Println("Indexing finished")
	}
}

func startWorker(client *ethclient.Client, blockToIndexNumCh <-chan int, wg *sync.WaitGroup) {
	for blockNum := range blockToIndexNumCh {
		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNum)))
		if err != nil {
			log.Println("Failed to get block: ", blockNum, " error: ", err)
			continue
		}

		log.Println("Indexing block: ", blockNum, " hash: ", block.Hash().Hex())
	}
	wg.Done()
}
