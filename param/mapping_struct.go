package param

// ADDR_MAPPING 风险名单信息index
const ADDR_MAPPING = `
{
  "mappings": {
    "properties": {
      "waAddr": {
        "type": "keyword"
      },
      "waRiskLevel": {
        "type": "short"
      },
      "waChain": {
        "type": "keyword"
      },
      "waTicker": {
        "type": "keyword"
      },
      "adsDataSource": {
         "type": "nested",
         "properties": {
            "dsAddr":    { "type": "keyword"  },
            "dsType": { "type": "keyword"  },
            "number":     { "type": "short"   }
          }
      }
    }
  }
}`
const TRANS_MAPPING = `
{
  "mappings": {
    "properties": {
      "hash": {
        "type": "keyword"
      },
      "txType": {
        "type": "keyword"
      },
      "size": {
        "type": "long"
      },
      "weight": {
        "type": "long"
      },
      "isError": {
        "type": "keyword"
      },
      "errCode": {
        "type": "keyword"
      },
      "internalTxType": {
        "type": "keyword"
      },
      "contractAddress": {
        "type": "keyword"
      },
      "functionName": {
        "type": "keyword"
      },
      "methodId": {
        "type": "keyword"
      },
       "traceId": {
        "type": "keyword"
      },
      "confirmations": {
        "type": "keyword"
      },
      "tokenName": {
        "type": "keyword"
      },
      "tokenSymbol": {
        "type": "keyword"
      },
      "tokenDecimal": {
        "type": "keyword"
      },
      "gasPrice": {
        "type": "text"
      },
      "lock_time": {
        "type": "long"
      },
      "tx_index": {
        "type": "keyword"
      },
      "double_spend": {
        "type": "boolean"
      },
      "time": {
        "type": "long"
      },
     "block_height": {
        "type": "long"
      },
      "blockHash": {
        "type": "text"
      },
      "inputs": {
         "type": "nested",
         "properties": {
            "sequence": {"type": "long"},
            "witness":    { "type": "text"  },
            "script": { "type": "text"  },
            "addr":     { "type": "keyword"   },
            "spent": {"type": "boolean" },
            "tx_index": {"type": "long"},
            "value": {"type": "long"}
          }
      },
      "out": {
         "type": "nested",
         "properties": {
            "spent": {"type": "boolean" },
            "value": {"type": "long"},
            "n": {"type": "long"},
            "tx_index": {"type": "long"},
            "script": { "type": "text"  },
            "addr":     { "type": "keyword"   }
          }
      }
    }
  }
}
`
