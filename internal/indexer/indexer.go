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
	defer client.Close()

	// TODO make a check for skipped blocks
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

	isBlockIndexedMap := &sync.Map{}
	for blockNum := start; blockNum < start+limit; blockNum++ {
		isBlockIndexedMap.Store(blockNum, false)
	}

	log.Println("Indexing from block: ", start, " to block: ", start+limit-1)
	blockToIndexNumCh := make(chan uint64)

	wg := &sync.WaitGroup{}
	for i := 0; i < workerNum; i++ {
		log.Println("Starting worker: ", i)
		wg.Add(1)
		go startWorker(client, blockToIndexNumCh, wg, isBlockIndexedMap)
	}

	for blockNum := start; (blockNum <= lastBlockNum) && (blockNum < start+limit); blockNum++ {
		blockToIndexNumCh <- blockNum
	}

	newBlocksToIndexCount := int64(start + limit - lastBlockNum - 1)

newBlocksLoop:
	for {
		if newBlocksToIndexCount <= 0 {
			break
		}
		log.Println("Waiting for new block...")
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

	// retrying skipped blocks
	for {
		skippedBlocksNums := make([]uint64, 0)

		isBlockIndexedMap.Range(func(key, value interface{}) bool {
			logStatus := "indexed"
			if !value.(bool) {
				logStatus = "skipped (not indexed)"
			}
			log.Println("Checking block: ", key.(uint64), " status: ", logStatus)
			blockNum := key.(uint64)
			if val, ok := isBlockIndexedMap.Load(blockNum); ok && !val.(bool) {
				skippedBlocksNums = append(skippedBlocksNums, blockNum)
			}
			return true
		})

		if len(skippedBlocksNums) == 0 {
			break
		}

		log.Println("Retrying skipped blocks...")
		blockToIndexNumCh := make(chan uint64)
		wg = &sync.WaitGroup{}

		for i := 0; (i < workerNum) && (i < len(skippedBlocksNums)); i++ {
			log.Println("Starting worker: ", i)
			wg.Add(1)
			go startWorker(client, blockToIndexNumCh, wg, isBlockIndexedMap)
		}

		for _, blockNum := range skippedBlocksNums {
			log.Println("Retrying block: ", blockNum)
			blockToIndexNumCh <- blockNum
		}
		close(blockToIndexNumCh)
		wg.Wait()
	}

	log.Println("Done!")
}

func startWorker(client *ethclient.Client, blockToIndexNumCh <-chan uint64, wg *sync.WaitGroup, isBlockIndexedMap *sync.Map) {
	for blockNum := range blockToIndexNumCh {
		// skip block if it is not planned for indexing
		if _, ok := isBlockIndexedMap.Load(blockNum); !ok {
			continue
		}
		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNum)))
		if err != nil {
			log.Println("Failed to get block: ", blockNum, " error: ", err)
			continue
		}

		log.Println("Indexing block: ", blockNum, " hash: ", block.Hash().Hex())
		isBlockIndexedMap.Store(blockNum, true)
	}
	wg.Done()
}
