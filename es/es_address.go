package es

//风险来源地址相关es操作
import (
	"context"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"log"
	"lyods-adsTool/domain"
	"lyods-adsTool/pkg/constants"
)

// AddDsAddrSource 增加指定风险地址的数据地址来源
func AddDsAddrSource(id string, addrSource domain.AdsDataSource) error {
	var err error
	updateData := map[string]any{
		"source": `ctx._source.adsDataSource.add(params.adsDataSource)`,
		"params": map[string]any{
			"adsDataSource": addrSource,
		},
	}
	req := update.NewRequest()
	req.Script = updateData
	_, err = ElasticClient.Update(constants.ES_ADDRESS, id).Request(req).Do(context.Background())
	if err != nil {
		log.Printf("AddDsAddrSource:Elastic 更新%s失败：%v\n", id, err.Error())
		return err
	}
	return nil
}
