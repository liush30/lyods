package tool

import (
	"github.com/elastic/go-elasticsearch/v7"
)

var client *elasticsearch.Client

// AddBulk 批量添加进ES
// see more: https://www.elastic.co/cn/blog/the-go-client-for-elasticsearch-working-with-data
//func AddBulk(index string, numWorkers, flushBytes int, jsonData []string) error {
//	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
//		Index:         index,            // The default index name
//		Client:        client,           // The Elasticsearch client
//		NumWorkers:    numWorkers,       // The number of worker goroutines
//		FlushBytes:    flushBytes,       // The flush threshold in bytes
//		FlushInterval: 30 * time.Second, // The periodic flush interval
//	})
//	if err != nil {
//		return err
//	}
//
//	var countSuccessful uint64 = 0
//
//	// 组合，具体参考ES Bulk API：https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-bulk.html#bulk-api-request-body
//	for _, str := range jsonData {
//		err := bulkIndexer.Add(context.Background(), esutil.BulkIndexerItem{
//			// Action field configures the operation to perform (index, create, delete, update)
//			Action: "index",
//			// DocumentID is the (optional) document ID
//			// DocumentID: strconv.Itoa(a.ID),
//			// Body is an `io.Reader` with the payload
//			Body: strings.NewReader(str),
//			// OnSuccess is called for each successful operation
//			OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
//				atomic.AddUint64(&countSuccessful, 1)
//			},
//			// OnFailure is called for each failed operation
//			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
//				if err != nil {
//					log.Printf("ERROR: %s", err)
//				} else {
//					log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
//				}
//			},
//		})
//		if err != nil {
//			return err
//		}
//	}
//
//	start := time.Now()
//	// Close waits until all added items are flushed and closes the indexer.
//	if err := bulkIndexer.Close(context.Background()); err != nil {
//		return err
//	} else {
//		// 统计时间
//		biStats := bulkIndexer.Stats()
//		dur := time.Since(start)
//
//		if biStats.NumFailed > 0 {
//			//log.Println()
//			//logger.Warnf(
//			//	"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
//			//	Comma(int64(biStats.NumFlushed)),
//			//	Comma(int64(biStats.NumFailed)),
//			//	dur.Truncate(time.Millisecond),
//			//	Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
//			//)
//		} else {
//			//logger.Infof(
//			//	"Successfully indexed [%s] documents in %s (%s docs/sec)",
//			//	Comma(int64(biStats.NumFlushed)),
//			//	dur.Truncate(time.Millisecond),
//			//	Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
//			//)
//		}
//	}
//
//	return nil
//}
