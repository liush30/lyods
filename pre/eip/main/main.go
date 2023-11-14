package main

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"lyods-adsTool/es"
	"lyods-adsTool/pkg/constants"
	"lyods-adsTool/pkg/utils"
	"lyods-adsTool/services/evm"
	"time"
)

//import (
//	"fmt"
//	"github.com/ethereum/go-ethereum/crypto"

//	"log"
//	"lyods-adsTool/pkg/constants"
//	"lyods-adsTool/pkg/utils"
//	"lyods-adsTool/services/evm"
//	"time"
//)

func main() {
	esClient, err := es.CreateEsClient()
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v\n", err)
		return
	}
	cbClient := evm.ChainBaseClient{
		RequestCount:    0,
		LastRequestTime: time.Now(),
	}
	e := evm.CreateEvmClient(constants.HTTP_ALCHEMY_ETH)
	ethClient := evm.EVMClient{
		Client:          e,
		Key:             []string{constants.ETH_KEY1, constants.ETH_KEY2},
		LastRequestTime: time.Now(),
		HTTPClient:      utils.CreateClient(),
		Chain:           constants.CHAIN_ETH,
	}
	list, _, err := ethClient.GetTxList(esClient, &cbClient, "0xa0e1c89Ef1a489c9C7dE96311eD5Ce5D32c20E4B", "0")
	if err != nil {
		log.Println(err)
	}
	log.Println(len(list))
	//e, err := bscClient.TransactionReceipt(context.Background(), common.HexToHash("0xbb270aef05636363ceaab06b2a2e4ee85ca611eb363080edc590448bc20bbed5"))
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Println(e)
}
