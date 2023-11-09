package main

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"lyods-adsTool/es"
	"lyods-adsTool/pkg/constants"
	"lyods-adsTool/pkg/utils"
	"lyods-adsTool/services/bitcoin"
	"lyods-adsTool/services/evm"
	"lyods-adsTool/services/list"
	"time"
)

func main() {
	esClient, err := es.CreateEsClient()
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v\n", err)
		return
	}
	//str, err := utils.DateChange("08 Apr 2011")
	//if err != nil {
	//	log.Fatalf("Failed to create Elasticsearch client: %v\n", err)
	//	return
	//}
	//fmt.Println(str)
	//删除索引
	deleteIndex(esClient, constants.ES_ADDRESS)
	deleteIndex(esClient, constants.ES_TRANSACTION)
	deleteIndex(esClient, constants.ES_ENTITY)
	//创建不同的索引
	createIndex(esClient, constants.ES_ADDRESS, constants.ADDR_MAPPING)
	createIndex(esClient, constants.ES_TRANSACTION, constants.TRANS_MAPPING)
	createIndex(esClient, constants.ES_ENTITY, constants.ENTITY_MAPPING)
	bitClient := bitcoin.BitClient{
		RequestCount:    0,
		LastRequestTime: time.Now(),
		Httpclient:      utils.CreateClient(),
	}
	//client := utils.CreateClient()
	e := evm.CreateEthClient()
	ethClient := evm.EthClient{
		Client:          e,
		Key:             []string{constants.ETH_KEY1, constants.ETH_KEY2},
		LastRequestTime: time.Now(),
		HTTPClient:      utils.CreateClient(),
	}
	cbClient := evm.ChainBaseClient{
		RequestCount:    0,
		LastRequestTime: time.Now(),
	}
	rClient := list.RClient{
		EsClient: esClient,
		EtClient: &ethClient,
		CbClient: &cbClient,
		BtClient: &bitClient,
	}
	// 获取风险名单信息
	//if err := list.GetAddrListByJSONOnBitcoin(constants.OPENSANCTIONS_URL, &bitClient, esClient); err != nil {
	//	log.Printf("Failed to get risk list from JSON source: %v\n", err)
	//}
	//if err := list.GetAddrListOnCsv(constants.UNISWAP_URL, esClient); err != nil {
	//	log.Printf("Failed to get risk list from CSV source: %v\n", err)
	//}
	if err := rClient.GetAddrListOnXmlByElement(`D:\Code\GoProjec\lyods-adsTool\sdn.xml`); err != nil {
		log.Printf("Failed to get risk list from XML source: %v\n", err)
	}
	//_, _, err = rClient.EtClient.GetTxListOnEth(rClient.EsClient, rClient.CbClient, "0xF67721A2D8F736E75a49FdD7FAd2e31D8676542a", "0")
	//if err != nil {
	//	log.Printf("Failed to get risk list from XML source: %v\n", err)
	//}
	//interTx, err := evm.GetTraceTransaction(rClient.CbClient, "0xfbe35150a2bd11d2fa1b7ebd25a960433fc6801341a4612151f433cedbaa260d")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//log.Println(interTx)
	//fmt.Println(len(list))
	//dbClient := db.GetDb()
	//defer dbClient.Close()
	////添加更新记录
	//err := db.AddUpdateRecord(dbClient, domain.UpdateLog{
	//	LogKey:     utils.GenerateUuid(),

	//	UpdateDate: time.Now().Format(time.DateTime),
	//	UpdateName: "opensanctions",
	//})
	//if err != nil {
	//	log.Printf("Failed to add opensanctions record: %v\n", err)
	//}
	//err = db.AddUpdateRecord(dbClient, domain.UpdateLog{
	//	LogKey:     utils.GenerateUuid(),
	//	UpdateDate: time.Now().Format(time.DateTime),
	//	UpdateName: "uniswap",
	//})
	//if err != nil {
	//	log.Printf("Failed to add uniswap record: %v\n", err)
	//}
	//err = db.AddUpdateRecord(dbClient, domain.UpdateLog{
	//	LogKey:     utils.GenerateUuid(),
	//	UpdateDate: time.Now().Format(time.DateTime),
	//	UpdateName: "sdn",
	//})
	//if err != nil {
	//	log.Printf("Failed to add sdn record: %v\n", err)
	//}
	//ethClient := evm.EthClient{
	//	evm.CreateEthClient(),
	//}
	//blockNumber, err := ethClient.GetLatestBlockNumber()
	//if err != nil {
	//	log.Printf("Failed to get latest block number: %v\n", err)
	//}
	//err = db.AddUpdateRecord(dbClient, domain.UpdateLog{
	//	LogKey:       utils.GenerateUuid(),
	//	UpdateDate:   time.Now().Format(time.DateTime),
	//	UpdateName:   "evm-block",
	//	UpdateRecord: strconv.Itoa(int(blockNumber)),
	//})
	//if err != nil {
	//	log.Printf("Failed to add evm-block record: %v\n", err)
	//}
	//btcBlock, err := bitcoin.GetLatestBlockNumber(*utils.CreateClient())
	//if err != nil {
	//	log.Printf("Failed to get latest block number: %v\n", err)
	//}
	//err = db.AddUpdateRecord(dbClient, domain.UpdateLog{
	//	LogKey:       utils.GenerateUuid(),
	//	UpdateDate:   time.Now().Format(time.DateTime),
	//	UpdateName:   "btc-block",
	//	UpdateRecord: strconv.Itoa(int(btcBlock)),
	//})
	//if err != nil {
	//	log.Printf("Failed to add btc-block record: %v\n", err)
	//}
}

func createIndex(client *es.ElasticClient, indexName, mapping string) {
	err := es.CreateIndex(client, indexName, mapping)
	if err != nil {
		log.Fatalf("Failed to create index %s: %v\n", indexName, err)
	} else {
		log.Printf("Index '%s' created successfully.\n", indexName)
	}
}
func deleteIndex(client *es.ElasticClient, indexName string) {
	err := es.DeleteIndexByName(client, indexName)
	if err != nil {
		log.Printf("Failed to delete index %s: %v\n", indexName, err.Error())
	} else {
		log.Printf("Index '%s' deleted successfully.\n", indexName)
	}
}
