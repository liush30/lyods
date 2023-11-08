package es

//风险来源地址相关es操作
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"log"
	"lyods-adsTool/domain"
	"lyods-adsTool/pkg/constants"
	"strings"
)

// AddDsAddrSource 增加指定风险地址的数据地址来源
func (c *ElasticClient) AddDsAddrSource(id string, addrSource domain.AdsDataSource) error {
	var err error
	updateData := map[string]any{
		"source": `ctx._source.adsDataSource.add(params.adsDataSource)`,
		"params": map[string]any{
			"adsDataSource": addrSource,
		},
	}
	req := update.NewRequest()
	req.Script = updateData
	_, err = c.Update(constants.ES_ADDRESS, id).Request(req).Do(context.Background())
	if err != nil {
		log.Printf("AddDsAddrSource:Elastic 更新%s失败：%v\n", id, err.Error())
		return err
	}
	return nil
}

// UpdateBalanceByID 修改指定实体的余额
func (c *ElasticClient) UpdateBalanceByID(docID string, newBalance float64) error {
	// 构建更新请求
	body := map[string]interface{}{
		"doc": map[string]interface{}{
			"balance": newBalance,
		},
	}
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req := esapi.UpdateRequest{
		Index:      constants.ES_ADDRESS,
		DocumentID: docID,
		Body:       strings.NewReader(string(data)),
	}

	// 执行更新请求
	res, err := req.Do(context.Background(), c)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error updating document: %s", res.Status())
	}

	return nil
}
