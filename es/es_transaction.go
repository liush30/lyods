package es

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"lyods-adsTool/domain"
	"lyods-adsTool/pkg/constants"
	"strings"
)

// AddAddressData 向指定文档ID的es中添加地址数据
func (es *ElasticClient) AddAddressData(documentID, addressID, address string) (domain.EsTrans, error) {
	//log.Println(documentID)
	//根据documentID查询数据信息
	getRes, err := es.Get(constants.ES_TRANSACTION, documentID).Do(context.Background())
	//getRes, err := getReq.Do(context.Background(), es)
	if err != nil {
		return domain.EsTrans{}, fmt.Errorf("error getting document: %s", err)
	}
	data, err := getRes.Source_.MarshalJSON()
	if err != nil {
		log.Printf("GetWalletAddr: MarshalJSON failed: %v\n", err.Error())
		return domain.EsTrans{}, err
	}
	var existingData domain.EsTrans
	err = json.Unmarshal(data, &existingData)
	if err != nil {
		log.Println("json unmarshal WalletAddr error:", err.Error())
		return domain.EsTrans{}, err
	}
	//如果信息已经存在，直接返回交易信息
	for _, item := range existingData.AddressListId {
		if item == addressID {
			return existingData, nil
		}
	}
	// 在现有数据上添加新数据
	existingData.AddressList = append(existingData.AddressList, address)
	existingData.AddressListId = append(existingData.AddressListId, addressID)

	// 将更新后的数据转换为JSON
	jsonData, err := json.Marshal(existingData)
	if err != nil {
		return domain.EsTrans{}, fmt.Errorf("error marshaling JSON: %s", err)
	}
	// 创建更新请求
	updateReq := esapi.IndexRequest{
		Index:      constants.ES_TRANSACTION,
		DocumentID: documentID,
		Body:       strings.NewReader(string(jsonData)),
		Refresh:    "true",
	}
	// 执行更新请求
	updateRes, err := updateReq.Do(context.Background(), es)
	if err != nil {
		return domain.EsTrans{}, fmt.Errorf("error updating document: %s", err)
	}
	defer updateRes.Body.Close()

	if updateRes.IsError() {
		return domain.EsTrans{}, fmt.Errorf("error updating document: %s", updateRes.String())
	}

	return existingData, nil
}
