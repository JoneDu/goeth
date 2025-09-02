package biz

import (
	"context"
	"github.com/Bruce/goeth/task1/config"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func BlockInfo() {
	//查询区块
	//编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
	//实现查询指定区块号的区块信息，包括区块的哈希、时间戳、交易数量等。
	//输出查询结果到控制台。
	c := config.LoadConfig()
	// use ethclient dial to connect  Sepolia
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + c.INFURA_PK)
	if err != nil {
		log.Fatal(err)
	}

	blockNum := big.NewInt(5671744)
	block, err := client.BlockByNumber(context.Background(), blockNum)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Block Number:", block.Number().Uint64())
	log.Println("Block Time:", block.Time())
	log.Println("Block Difficulty:", block.Difficulty())
	log.Println("Transaction Count:", block.Transactions().Len())
	log.Println("Block Hash:", block.Hash().Hex())

}
