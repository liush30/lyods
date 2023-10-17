package main

import (
	"log"
	"lyods-adsTool/es"
	"lyods-adsTool/pkg/constants"
	"lyods-adsTool/services/list"
)

func main() {
	esClient, err := es.CreateEsClient()
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v\n", err)
		return
	}
	//删除索引
	//deleteIndex(esClient, constants.ES_ADDRESS)
	//deleteIndex(esClient, constants.ES_TRANSACTION)
	//deleteIndex(esClient, constants.ES_ENTITY)
	// 创建不同的索引
	//createIndex(esClient, constants.ES_ADDRESS, constants.ADDR_MAPPING)
	//createIndex(esClient, constants.ES_TRANSACTION, constants.TRANS_MAPPING)
	//createIndex(esClient, constants.ES_ENTITY, constants.ENTITY_MAPPING)

	// 获取风险名单信息
	if err := list.GetAddrListByJSONOnBitcoin(constants.OPENSANCTIONS_URL, esClient); err != nil {
		log.Printf("Failed to get risk list from JSON source: %v\n", err)
	}
	if err := list.GetAddrListOnCsv(constants.UNISWAP_URL, esClient); err != nil {
		log.Printf("Failed to get risk list from CSV source: %v\n", err)
	}
	if err := list.GetAddrListOnXmlByElement(`D:\Code\GoProjec\lyods-adsTool\sdn.xml`, esClient); err != nil {
		log.Printf("Failed to get risk list from XML source: %v\n", err)
	}
}

func createIndex(client *es.ElasticClient, indexName, mapping string) {
	err := es.CreateIndex(client, indexName, mapping)
	if err != nil {
		log.Printf("Failed to create index %s: %v\n", indexName, err)
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
