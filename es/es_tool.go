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
	"lyods-adsTool/config"
	"lyods-adsTool/domain"
	"lyods-adsTool/pkg/constants"
	"net"
	"net/http"
	"strings"
	"time"
)

type ElasticClient struct {
	*elasticsearch.TypedClient
}

//var ctx context.Context

// CreateEsClient 创建Elasticsearch客户端
func CreateEsClient() (*ElasticClient, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{config.ES_URL},
		Username:  config.ES_USERNAME,
		Password:  config.ES_PWD,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
	}
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		log.Printf("Elastic 连接失败: %v\n", err.Error())
		return nil, err
	}
	log.Println("Elastic 连接成功")
	return &ElasticClient{client}, nil
}

// CreateIndex 创建索引
func CreateIndex(es *ElasticClient, indexName, indexMapping string) error {
	var err error
	isExits, err := IndexExists(es, indexName)
	if err != nil {
		return fmt.Errorf("error checking index existence: %s", err)
	}
	//判断索引是否已经存在
	if isExits {
		log.Printf("Index %s already exists\n", indexName)
		return nil
	}
	// 创建索引请求
	createIndexRequest := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  strings.NewReader(indexMapping),
	}
	// 执行创建索引请求
	res, err := createIndexRequest.Do(context.Background(), es)
	if err != nil {
		return fmt.Errorf("error creating index: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error creating index: %s", res)
	} else if err != nil {
		return err
	}
	return nil
}

// IndexExists 判断索引是否存在
func IndexExists(es *ElasticClient, indexName string) (bool, error) {
	req := esapi.IndicesExistsRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	if res.StatusCode == 200 {
		return true, nil
	}

	return false, nil
}

// DeleteIndexByName 根据索引值删除索引
func DeleteIndexByName(es *ElasticClient, indexName string) error {
	req := esapi.IndicesDeleteRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting index: %s", res.String())
	}
	return nil
}

// Insert 新增1条数据并指定id
func (c *ElasticClient) Insert(indexName, id string, structBody interface{}) error {
	_, err := c.Index(indexName).Id(id).Request(structBody).Do(context.Background())
	if err != nil {
		log.Printf("Fail elastic add, id is %s, error: %s\n", id, err.Error())
		return err
	}
	return nil
}

// IsExistById 按照ID查询数据信息是否存在
func (c *ElasticClient) IsExistById(indexName, id string) (bool, error) {
	return c.Exists(indexName, id).Index(indexName).IsSuccess(nil)
}

// GetWalletAddrById 根据指定id,查询地址名单信息
func (c *ElasticClient) GetWalletAddrById(indexName, id string) (domain.WalletAddr, error) {
	var addrInfo domain.WalletAddr
	isExists, err := c.IsExistById(indexName, id)
	if err != nil {
		log.Printf("GetWalletAddr: Elasticsearch query failed for %s: %v\n", id, err.Error())
		return domain.WalletAddr{}, err
	}
	if !isExists {
		return domain.WalletAddr{}, err
	}
	res, err := c.Get(indexName, id).Do(context.Background())
	if err != nil {
		log.Printf("GetWalletAddr: Elasticsearch query for %s failed: %v\n", id, err.Error())
		return domain.WalletAddr{}, err
	}
	data, err := res.Source_.MarshalJSON()
	if err != nil {
		log.Printf("GetWalletAddr: MarshalJSON failed: %v\n", err.Error())
		return domain.WalletAddr{}, err
	}
	err = json.Unmarshal(data, &addrInfo)
	if err != nil {
		log.Println("json unmarshal WalletAddr error:", err.Error())
		return domain.WalletAddr{}, err
	}
	return addrInfo, nil
}

// DeleteById 根据id删除指定数据
func (c *ElasticClient) DeleteById(indexName, id string) error {
	_, err := c.Delete(indexName, id).Do(context.Background())
	if err != nil {
		log.Printf("DeleteById:Elastic 删除%s失败:%v\n", id, err.Error())
		return err
	}
	return nil
}

// GetAddrLevel 获得指定地址得监控层次
func (c *ElasticClient) GetAddrLevel(addr string) (int, error) {
	res, err := c.Get(constants.ES_ADDRESS, addr).Source_("waRiskLevel").Do(context.Background())
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
		//bitcoin.IsError(err, "jsonparser get waRiskLevel error")
		return -1, err
	}
	return int(val), err
}

// UpdateDsAddrNumber 更新标记次数
func (c *ElasticClient) UpdateDsAddrNumber(id, parAddr string, number int) error {
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
	_, err := c.Update(constants.ES_ADDRESS, id).Request(req).Do(context.Background())
	if err != nil {
		log.Printf("AddDsAddrSource:Elastic 更新%s失败：%v\n", id, err.Error())
		return err
	}
	return nil
}

// UpdateAddrLevel 更新地址的风险层次
func (c *ElasticClient) UpdateAddrLevel(id string, newLevel int) error {
	updateData := map[string]any{
		"source": `ctx._source.waRiskLevel=params.waRiskLevel`,
		"params": map[string]any{
			"waRiskLevel": newLevel,
		},
	}
	req := update.NewRequest()
	req.Script = updateData
	_, err := c.Update(constants.ES_ADDRESS, id).Request(req).Do(context.Background())
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
//	dsAddrInfo := domain.AdsDataSource{
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
//		//var addrInfo domain.WalletAddr
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
//		addrInfo := domain.WalletAddr{
//			WaAddr:      inputAddr,
//			WaRiskLevel: uint(parentLevel + 1),
//			WaChain:     utils.GetChainByTicker(constants.TICKER_BTC),
//			WaTicker:    constants.TICKER_BTC,
//			DsAddr: []domain.AdsDataSource{
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
