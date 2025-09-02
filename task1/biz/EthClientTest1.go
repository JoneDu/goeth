package biz

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/Bruce/goeth/task1/config"
	"github.com/Bruce/goeth/task1/contract"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

// 查询区块
// 编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
// 实现查询指定区块号的区块信息，包括区块的哈希、时间戳、交易数量等。
// 输出查询结果到控制台。
func BlockInfo() {

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

//发送交易
//准备一个 Sepolia 测试网络的以太坊账户，并获取其私钥。
//编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
//构造一笔简单的以太币转账交易，指定发送方、接收方和转账金额。
//对交易进行签名，并将签名后的交易发送到网络。
//输出交易的哈希值。

func TransferEth() {
	c := config.LoadConfig()
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + c.INFURA_PK)
	if err != nil {
		log.Fatal(err)
	}
	//defer client.Close()

	// 解析私钥
	privateKey, err := crypto.HexToECDSA(c.Ak1)
	if err != nil {
		log.Fatal(err)
	}

	// get sender public key address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// get receiver address
	toAddress := common.HexToAddress(c.APk2)

	// get Sender Nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	// transfer eth
	amount := big.NewInt(10000000000000000) // 0.01 eth
	// gasLimit
	gasLimit := uint64(21000)
	// gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// build transaction
	tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)

	// get chain ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// sign transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	// send transaction
	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		log.Fatal(err)
	}
	log.Println("tx sent: ", signedTx.Hash().Hex())

}

func GoCounterGet() {
	c := config.LoadConfig()
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + c.INFURA_PK)
	if err != nil {
		log.Fatal(err)
	}
	//defer client.Close()
	// New Contract
	counterContract, err := contract.NewContract(common.HexToAddress("0x9Dd442BD234c3085525ba0EAd55162528FCa06eD"), client)
	if err != nil {
		log.Fatal(err)
	}
	// call get method
	count, err := counterContract.Get(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("count: ", count)
}

func GoCounterInc() {
	c := config.LoadConfig()
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + c.INFURA_PK)
	if err != nil {
		log.Fatal(err)
	}
	//defer client.Close()

	//privateKey
	privateKey, err := crypto.HexToECDSA(c.Ak1)
	if err != nil {
		log.Fatal(err)
	}
	// get from address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// get nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	gasLimit := uint64(30000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	//get chanid
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// create transaction
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice

	// contract address
	counterContract, err := contract.NewContract(common.HexToAddress("0x9Dd442BD234c3085525ba0EAd55162528FCa06eD"), client)
	if err != nil {
		log.Fatal(err)
	}
	tx, err := counterContract.Inc(auth)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tx sent: ", tx.Hash().Hex())

	// wait for transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tx mined: ", receipt.Status)
	fmt.Println("tx gas used: ", receipt.GasUsed)
	fmt.Println("tx confirmed in block: ", receipt.BlockNumber.String())
}
func SubIncrementEvents() {
	c := config.LoadConfig()
	client, err := ethclient.Dial("wss://sepolia.infura.io/ws/v3/" + c.INFURA_PK)
	if err != nil {
		log.Fatal(err)
	}
	counterContract, err := contract.NewContract(common.HexToAddress("0x9Dd442BD234c3085525ba0EAd55162528FCa06eD"), client)
	if err != nil {
		log.Fatal(err)
	}
	// create watch options
	watchOpts := &bind.WatchOpts{Context: context.Background()}
	// create event channel
	eventChan := make(chan *contract.ContractIncrement)
	sub, err := counterContract.WatchIncrement(watchOpts, eventChan)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	log.Println("Waiting for events...")
	for {
		select {
		case event := <-eventChan:
			log.Println("Event Count:", event.Count)
			log.Println("Event block num:", event.Raw.BlockNumber)
			log.Println("Event tx hash:", event.Raw.TxHash.Hex())
			log.Printf("  Block Hash: %s", event.Raw.BlockHash.Hex())
		case err := <-sub.Err():
			log.Println("Error:", err)
			return
		}
	}
}
