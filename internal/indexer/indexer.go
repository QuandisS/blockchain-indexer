package indexer

import (
	"context"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const workerNum = 5

func Start(rpc string, start uint64, out string, limit uint64) {
	client, err := ethclient.Dial(rpc)

	if err != nil {
		log.Fatal("Failed to connect with Ethereum client: ", err)
	}

	lastBlockNum, err := client.BlockNumber(context.Background())

	if err != nil {
		log.Fatal("Failed to get last block number: ", err)
	}

	newBlockChan := make(chan *types.Header)
	var sub ethereum.Subscription
	if uint64(start+limit) > lastBlockNum {
		sub, err = client.SubscribeNewHead(context.Background(), newBlockChan)
		if err != nil {
			log.Fatal("Failed to subscribe new blocks (use WebSockets connection instead?): ", err)
		}
	}

	log.Println("Indexing from block: ", start, " to block: ", start+limit)
	blockToIndexNumCh := make(chan uint64)

	wg := &sync.WaitGroup{}
	for i := 0; i < workerNum; i++ {
		log.Println("Starting worker: ", i)
		wg.Add(1)
		go startWorker(client, blockToIndexNumCh, wg)
	}

	for blockNum := start; (blockNum <= lastBlockNum) && (blockNum <= start+limit); blockNum++ {
		blockToIndexNumCh <- blockNum
	}

	newBlocksToIndexCount := int64(start + limit - lastBlockNum)

newBlocksLoop:
	for {
		if newBlocksToIndexCount <= 0 {
			break
		}
		log.Println("Waiting for new blocks...")
		select {
		case newBlock := <-newBlockChan:
			log.Println("New block received: ", newBlock.Number.Uint64())
			blockToIndexNumCh <- newBlock.Number.Uint64()
			newBlocksToIndexCount--
			if newBlocksToIndexCount == 0 {
				break newBlocksLoop
			}

		case err := <-sub.Err():
			log.Println("Subscription error: ", err)
			break newBlocksLoop
		}
	}

	close(blockToIndexNumCh)
	wg.Wait()
	log.Println("Indexing finished")

}

func startWorker(client *ethclient.Client, blockToIndexNumCh <-chan uint64, wg *sync.WaitGroup) {
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
