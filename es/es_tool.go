package es

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"log"
	"lyods-adsTool/entities"
	"lyods-adsTool/pkg/constants"
	"lyods-adsTool/tool"
	"net"
	"net/http"
	"strings"
	"time"
)

var ElasticClient *elasticsearch.TypedClient

//var ctx context.Context

func init() {
	ElasticClient = createEsClient()
	//ctx = context.Background()
}

// createEsClient 创建Es客户端
func createEsClient() *elasticsearch.TypedClient {
	cfg := elasticsearch.Config{
		Addresses: []string{"https://localhost:9200"},
		Username:  "elastic",
		Password:  "123123",
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, //忽略验证
		},
	}
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		// Handle error
		log.Printf("Elastic 连接失败: %v\n", err.Error())
	} else {
		log.Println("Elastic 连接成功")
	}
	return client
}

// CreateIndex 创建索引
func CreateIndex(client *elasticsearch.Client, indexName, indexMapping string) error {
	var err error
	// 创建索引请求
	createIndexRequest := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  strings.NewReader(indexMapping),
	}
	// 执行创建索引请求
	res, err := createIndexRequest.Do(context.Background(), client)
	if err != nil {
		log.Fatalf("Error creating index: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Failed to create index: %s", res.Status())
	} else {
		fmt.Printf("Index '%s' created successfully.\n", indexName)
	}
	return nil
}

// DeleteIndexByName 根据索引值删除索引
//func DeleteIndexByName(indexName string) error {
//	//判断该索引是否存在
//	exits, err := ElasticClient.IndexExists(indexName).Do(ctx)
//	if err != nil {
//		log.Printf("Elastic IndexExists Error:%v\n", err.Error())
//		return err
//	}
//	//若索引不存在，返回错误信息
//	if !exits {
//		log.Printf("%s index not exits", indexName)
//		return errors.New("index not exits")
//	}
//	_, err = ElasticClient.DeleteIndex(indexName).Do(ctx)
//	if err != nil {
//		log.Printf("Delete %s index error:%s ", indexName, err.Error())
//		return err
//	}
//	return nil
//}

// Insert 新增1条数据并指定id
func Insert(indexName, id string, structBody interface{}) error {
	_, err := ElasticClient.Index(indexName).Id(id).Request(structBody).Do(context.Background())
	if err != nil {
		log.Printf("Fail elastic add,id is %s,error:%s\n", id, err.Error())
		return err
	}
	return nil
}

// IsExistById 按照ID查询数据信息是否存在
func IsExistById(indexName, id string) (bool, error) {
	return ElasticClient.Exists(indexName, id).Index(indexName).IsSuccess(nil)
}

// GetWalletAddrById 根据指定id,查询地址名单信息
func GetWalletAddrById(indexName, id string) (entities.WalletAddr, error) {
	var err error
	var addrInfo entities.WalletAddr
	//判断id是否存在
	isExists, err := IsExistById(indexName, id)
	if err != nil {
		log.Printf("GetWalletAddr:Elastic 查询%s是否存在失败：%v\n", id, err.Error())
		return entities.WalletAddr{}, err
	}
	//若不存在,直接返回为空
	if !isExists {
		return entities.WalletAddr{}, err
	}
	res, err := ElasticClient.Get(indexName, id).Do(context.Background())
	if err != nil {
		log.Printf("GetWalletAddr:Elastic 查询%s失败：%v\n", id, err.Error())
		return entities.WalletAddr{}, err
	}
	data, err := res.Source_.MarshalJSON()
	if err != nil {
		log.Printf("GetWalletAddr:MarshalJSON 失败:%v\n", err.Error())
		return entities.WalletAddr{}, err
	}
	err = json.Unmarshal(data, &addrInfo)
	if err != nil {
		log.Println("json unmarshal WalletAddr error:", err.Error())
		return entities.WalletAddr{}, err
	}
	return addrInfo, nil

}

// DeleteById 根据id删除指定数据
func DeleteById(indexName, id string) error {
	_, err := ElasticClient.Delete(indexName, id).Do(context.Background())
	if err != nil {
		log.Printf("DeleteById:Elastic 删除%s失败:%v\n", id, err.Error())
		return err
	}
	return nil
}

// GetAddrLevel 获得指定地址得监控层次
func GetAddrLevel(addr string) (int, error) {
	res, err := ElasticClient.Get(constants.ES_ADDRESS, addr).Source_("waRiskLevel").Do(context.Background())
	if err != nil {
		log.Println("get address level error:", err.Error())
		return -1, err
	}
	levelBytes, err := res.Source_.MarshalJSON()
	if err != nil {
		log.Println("response source marshal json error:", err.Error())
		return -1, err
	}
	val, err := jsonparser.GetInt(levelBytes, "waRiskLevel")
	if err != nil {
		tool.IsError(err, "jsonparser get waRiskLevel error")
		return -1, err
	}
	return int(val), err
}

// UpdateDsAddrNumber 更新标记次数
func UpdateDsAddrNumber(id, parAddr string, number int) error {
	updateData := map[string]any{
		"source": `def targets = ctx._source.adsDataSource.findAll
        (addr -> addr.dsAddr == params.dsAddr); for(addr in targets) 
        { addr.number = params.number }`, //修改adsDataSource中指定的父级来源地址的标记次数
		"params": map[string]any{
			"dsAddr": parAddr,
			"number": number,
		},
	}
	req := update.NewRequest()
	req.Script = updateData
	_, err := ElasticClient.Update(constants.ES_ADDRESS, id).Request(req).Do(context.Background())
	if err != nil {
		log.Printf("AddDsAddrSource:Elastic 更新%s失败：%v\n", id, err.Error())
		return err
	}
	return nil
}

// UpdateAddrLevel 更新地址的风险层次
func UpdateAddrLevel(id string, newLevel int) error {
	updateData := map[string]any{
		"source": `ctx._source.waRiskLevel=params.waRiskLevel`,
		"params": map[string]any{
			"waRiskLevel": newLevel,
		},
	}
	req := update.NewRequest()
	req.Script = updateData
	_, err := ElasticClient.Update(constants.ES_ADDRESS, id).Request(req).Do(context.Background())
	if err != nil {
		log.Printf("UpdateAddrLevel:Elastic 更新%s\twaRiskLevel失败：%v\n", id, err.Error())
		return err
	}
	return nil
}

//// UpdateOrAddSourceOnTrans 在交易为转出交易的前提下，将其转入地址名单存入风险名单信息中
//func UpdateOrAddSourceOnTrans(inputAddr, parentAddr string) error {
//	//判断该地址是否存在
//	isExist, err := IsExistById(constants.ES_ADDRESS, inputAddr)
//	if err != nil {
//		tool.IsError(err, "IsExistById Error :")
//		return err
//	}
//	dsAddrInfo := entities.AdsDataSource{
//		DsAddr: parentAddr, //将父级地址作为来源地址
//		DsType: constants.DS_TYPE_ADDRESS,
//		Number: constants.INIT_NUMBER,
//	}
//	parentLevel, err := GetAddrLevel(parentAddr) //获得父级监控层次
//	if err != nil {
//		log.Println("get parentLevel error!")
//		return err
//	}
//	//该地址已经存在，更新数据来源信息
//	if isExist {
//		//获取该地址信息
//		//var addrInfo entities.WalletAddr
//		addrInfo, err := GetWalletAddrById(param.WA_ADDR_INDEX, inputAddr)
//		if err != nil {
//			tool.IsError(err, "GetWalletAddr get address info")
//			return err
//		}
//		//判断该地址记来源（即父地址）是否已经存在，存在则更新标次数，不存在新增来源地址
//		isExitsParentAddress, parentIndex := utils.IsExitsAddrSource(addrInfo.DsAddr, parentAddr)
//		//若该来源地址已经存在，更新标记次数
//		if isExitsParentAddress {
//			//更新标记次数
//			err = UpdateDsAddrNumber(inputAddr, parentAddr, int(addrInfo.DsAddr[parentIndex].Number+1))
//			if err != nil {
//				return err
//			}
//			//来源地址不存在，新增来源地址
//		} else {
//			addrInfo.DsAddr = append(addrInfo.DsAddr, dsAddrInfo)
//			//添加地址来源
//			err = AddDsAddrSource(inputAddr, dsAddrInfo)
//			if err != nil {
//				return err
//			}
//		}
//		//判断自身层级级别高于当前层级，所高于层级不变，否则更新最新层级（数字越小，层级越高）
//		if addrInfo.WaRiskLevel > uint(parentLevel+1) {
//			err = UpdateAddrLevel(inputAddr, parentLevel+1)
//			if err != nil {
//				return err
//			}
//		}
//		log.Printf("%s信息已存在,添加该数据来源\n", inputAddr)
//	} else {
//		//该地址不存在，新增风险名单信息
//		addrInfo := entities.WalletAddr{
//			WaAddr:      inputAddr,
//			WaRiskLevel: uint(parentLevel + 1),
//			WaChain:     utils.GetChainByTicker(constants.TICKER_BTC),
//			WaTicker:    constants.TICKER_BTC,
//			DsAddr: []entities.AdsDataSource{
//				dsAddrInfo,
//			},
//		}
//		err = Insert(param.WA_ADDR_INDEX, inputAddr, addrInfo)
//		if err != nil {
//			log.Println("insert addr_list error:", err.Error())
//			return err
//		}
//	}
//	return nil
//}
