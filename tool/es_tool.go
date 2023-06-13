package tool

import (
	"context"
	"encoding/json"
	"github.com/buger/jsonparser"
	"github.com/olivere/elastic/v7"
	"log"
	"lyods-adsTool/entity"
	"lyods-adsTool/param"
)

var ElasticClient *elastic.Client
var ctx context.Context

func init() {
	ElasticClient = createEsClient("http://127.0.0.1:9200")
	ctx = context.Background()
	//创建addrList索引
	//createIndex(param.WA_ADDR_INDEX, param.ADDR_MAPPING)
}

// createEsClient 创建Es客户端
func createEsClient(urls ...string) *elastic.Client {
	// 创建ES client用于后续操作ES
	client, err := elastic.NewClient(
		// 设置ES服务地址，支持多个地址
		elastic.SetURL(urls...),
		// 设置基于http base auth验证的账号和密码
		elastic.SetBasicAuth("user", "secret"))
	if err != nil {
		// Handle error
		log.Printf("Elastic 连接失败: %v\n", err.Error())
	} else {
		log.Println("Elastic 连接成功")
	}
	return client
}

// CreateIndex 创建索引
func CreateIndex(indexName, indexMapping string) error {
	var err error
	// 首先检测下addrList索引是否存在
	exists, err := ElasticClient.IndexExists(indexName).Do(ctx)
	if err != nil {
		log.Printf("Elastic IndexExists Error:%v\n", err.Error())
		return err
	}
	if !exists {
		// addrList索引不存在，则创建一个
		_, err := ElasticClient.CreateIndex(indexName).BodyString(indexMapping).Do(ctx)
		if err != nil {
			log.Printf("Elastic 创建索引失败:%v\n", err.Error())
			return err
		}
		return nil
	}
	log.Println("该索引已存在")
	return nil
}

// Insert 新增1条数据
func Insert(indexName, id string, structBody interface{}) error {
	//uid := uuid.New().String()
	_, err := ElasticClient.Index().Index(indexName).Id(id).BodyJson(structBody).Do(ctx)
	if err != nil {
		log.Printf("Elastic 新增%s失败：%v\n", id, err.Error())
		return err
	}
	return nil
}

// IsExistById 按照ID查询数据信息是否存在
func IsExistById(indexName, id string) (bool, error) {
	return ElasticClient.Exists().Index(indexName).Id(id).Do(ctx)
}

// GetWalletAddrByte 根据指定id,查询地址名单信息
func GetWalletAddrByte(indexName, id string) (error, []byte) {
	var err error
	isExists, err := IsExistById(indexName, id)
	if err != nil {
		log.Printf("GetWalletAddr:Elastic 查询%s是否存在失败：%v\n", id, err.Error())
		return nil, nil
	}
	//若不存在,直接返回为空
	if !isExists {
		return nil, nil
	}
	get1, err := ElasticClient.Get().Index(indexName).Id(id).Do(ctx)
	if err != nil {
		log.Printf("GetWalletAddr:Elastic 查询%s失败：%v\n", id, err.Error())
		return nil, nil
	}
	data, err := get1.Source.MarshalJSON()
	if err != nil {
		log.Printf("GetWalletAddr:MarshalJSON 失败:%v\n", err.Error())
		return err, nil
	}
	return nil, data

}

// DeleteById 根据id删除指定数据
func DeleteById(indexName, id string) error {
	_, err := ElasticClient.Delete().Index(indexName).Id(id).Do(ctx)
	if err != nil {
		log.Printf("DeleteById:Elastic 删除%s失败:%v\n", id, err.Error())
		return err
	}
	return nil
}

// GetIndexDocNum 获取索引文档总数
func GetIndexDocNum(index string) (int64, error) {
	return ElasticClient.Count(index).Do(ctx)
}

// UpdateDsAddrSource 增加数据地址来源
func UpdateDsAddrSource(id string, addrSource entity.AdsDataSource) error {
	var err error
	err, addrByte := GetWalletAddrByte(param.WA_ADDR_INDEX, id)
	if err != nil {
		log.Printf("UpdateDsAddrSource:查询%s失败:%v", id, err.Error())
	}
	var addrSourceList []entity.AdsDataSource
	sourceDate, _, _, err := jsonparser.Get(addrByte, "adsDataSource")
	//err = json.Unmarshal(data, &walletAddr)
	if err != nil {
		log.Printf("UpdateDsAddrSource:jsonparser Get 失败:%v\n", err.Error())
		return err
	}
	err = json.Unmarshal(sourceDate, &addrSourceList)
	if err != nil {
		log.Printf("UpdateDsAddrSource:json 解析失败:%v\n", err.Error())
		return err
	}
	addrSourceList = append(addrSourceList, addrSource)
	_, err = ElasticClient.Update().Index(param.WA_ADDR_INDEX).Id(id).Doc(map[string]interface{}{"adsDataSource": addrSourceList}).Do(ctx)
	if err != nil {
		log.Printf("UpdateDsAddrSource:Elastic 更新%s失败：%v\n", id, err.Error())
		return err
	}
	return nil
}
