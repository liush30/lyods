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
	e := evm.CreateEvmClient(constants.HTTP_INFURA_ARB)
	arbClient := evm.EVMClient{
		Client:          e,
		Key:             []string{constants.ARB_KEY1, constants.ARB_KEY2},
		LastRequestTime: time.Now(),
		HTTPClient:      utils.CreateClient(),
		Chain:           constants.CHAIN_ARB,
	}
	cbClient := evm.ChainBaseClient{
		RequestCount:    0,
		LastRequestTime: time.Now(),
	}

	transList, _, err := arbClient.GetTxList(esClient, &cbClient, "0x7C2427FBf06370D4e0C6149794f756c91a1C6E11", "0")
	if err != nil {
		log.Println(err)
	}
	log.Println(len(transList))
}
