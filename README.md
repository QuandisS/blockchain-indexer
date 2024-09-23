# Mini Ethereum indexer

## Usage
To run indexer use the following command:
`go run . run --rpc <your api url>`

Options:
1. `--rpc`: is only *__required__* option; your API URL, WebSockets supported only (for new block notification reasons; tested on Infura Ethereum mainnet)
2. `--start`: block number to start indexing with. *Default is __1__*
3. `--out`: output file name. *Default is __blocks.log__*
4. `--limit`: limits a number of blocks to be indexed. API credits are not endless, though :sweat_smile:. *Default is __5__*
5. `--workers`: sets a workers count. *Default is __5__*. This indexer uses a workerpool to index blocks