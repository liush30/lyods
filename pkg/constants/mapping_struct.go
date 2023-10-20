package constants

//es映射定义

// ADDR_MAPPING 风险名单信息映射
const ADDR_MAPPING = `
{
    "mappings":{
        "properties":{
            "waAddr":{
                "type":"keyword"
            },
            "entityId":{
                "type":"keyword"
            },
            "waRiskLevel":{
                "type":"short"
            },
            "waChain":{
                "type":"keyword"
            },
            "isTrace":{ 
				"type":"boolean"
			},
           "isNeedTrace":{
				"type":"boolean"
			},
            "adsDataSource":{
                "type":"nested",
                "properties":{
                    "dsAddr":{
                        "type":"keyword"
                    },
                     "dsTransHash":{
                        "type":"text",
                        "fields":{
                            "hash":{
                                "type":"keyword"
                            }
                        }
                    },
                    "dsType":{
                        "type":"keyword"
                    },
                    "illustrate":{
                        "type":"text"
                    },
                    "time":{
                        "type":"date"
                    },
                    "dsRules":{
                        "type":"text",
                        "fields":{
                            "rule":{
                                "type":"keyword"
                            }
                        }
                    }
                }
            },
            "levelNumber":{
                "type":"nested",
                "properties":{
                    "level":{
                        "type":"short"
                    },
                    "number":{
                        "type":"short"
                    }
                }
            },
            "rules":{
                "type":"text",
                "fields":{
                    "rule":{
                        "type":"keyword"
                    }
                }
            },
            "riskChgHistory":{
                "type":"nested",
                "properties":{
                    "dateOfBirth":{
                        "type":"date"
                    },
                    "riskLevel":{
                        "type":"short"
                    },
                    "description":{
                        "type":"text"
                    }
                }
            }
        }
    }
}
`

// TRANS_MAPPING 交易映射
const TRANS_MAPPING = `
{
    "mappings":{
        "properties":{
            "hash":{
                "type":"keyword"
            },
            "riskLevel":{
                "type":"short"
            },
            "address":{
                "type":"keyword"
            },
            "size":{
                "type":"long"
            },
            "weight":{
                "type":"long"
            },
            "gasUsed":{
                "type":"text"
            },
            "isError":{
                "type":"text"
            },
            "errCode":{
                "type":"text"
            },
            "contractAddress":{
                "type":"keyword"
            },
            "functionName":{
                "type":"keyword"
            },
            "methodId":{
                "type":"keyword"
            },
            "traceId":{
                "type":"text"
            },
            "confirmations":{
                "type":"text"
            },
            "cumulativeGasUsed":{
                "type":"text"
            },
            "gasPrice":{
                "type":"text"
            },
            "lockTime":{
                "type":"long"
            },
            "txIndex":{
                "type":"keyword"
            },
            "doubleSpend":{
                "type":"boolean"
            },
            "time":{
                "type":"long"
            },
            "blockHeight":{
                "type":"text"
            },
            "blockHash":{
                "type":"text"
            },
            "value":{
                "type":"long"
            },
            "inputs":{
                "type":"nested",
                "properties":{
                    "sequence":{
                        "type":"long"
                    },
                    "witness":{
                        "type":"text"
                    },
                    "script":{
                        "type":"text"
                    },
                    "addr":{
                        "type":"keyword"
                    },
                    "spent":{
                        "type":"boolean"
                    },
                    "txIndex":{
                        "type":"text"
                    },
                    "value":{
                        "type":"long"
                    }
                }
            },
            "out":{
                "type":"nested",
                "properties":{
                    "spent":{
                        "type":"boolean"
                    },
                    "value":{
                        "type":"long"
                    },
                    "n":{
                        "type":"long"
                    },
                    "txIndex":{
                        "type":"text"
                    },
                    "script":{
                        "type":"text"
                    },
                    "addr":{
                        "type":"keyword"
                    }
                }
            },
            "internalTx":{
                "type":"nested",
                "properties":{
                    "id":{
                        "type":"keyword"
                    },
                    "traceAddress":{
                        "type":"text"
                    },
                    "fromAddr":{
                        "type":"keyword"
                    },
                    "toAddr":{
                        "type":"keyword"
                    },
                    "inputTx":{
                        "type":"text"
                    },
                    "outputTx":{
                        "type":"text"
                    },
                    "value":{
                        "type":"long"
                    },
                    "subTraces":{
                        "type":"long"
                    },
                    "callType":{
                        "type":"keyword"
                    }
                }
            },
            "logs":{
                "type":"nested",
                "properties":{
                    "address":{
                        "type":"keyword"
                    },
                    "eventInfo":{
                        "type":"text"
                    },
                    "topics":{
                        "type":"nested",
                        "properties":{
                            "key":{
                                "type":"keyword"
                            },
                            "value":{
                                "type":"keyword"
                            }
                        }
                    }
                }
            },
            "erc20Txn":{
                "properties":{
                    "fromAddr":{
                        "type":"keyword"
                    },
                    "toAddr":{
                        "type":"keyword"
                    },
                    "contractAddress":{
                        "type":"keyword"
                    },
                    "amount":{
                        "type":"long"
                    }
                }
            },
            "rules":{
                "type":"text",
                "fields":{
                    "rule":{
                        "type":"keyword"
                    }
                }
            },
            "riskChgHistory":{
                "type":"nested",
                "properties":{
                    "dateOfBirth":{
                        "type":"date"
                    },
                    "riskLevel":{
                        "type":"short"
                    },
                    "description":{
                        "type":"text"
                    }
                }
            }
        }
    }
}
`

// ENTITY_MAPPING 实体映射
const ENTITY_MAPPING = `
{
  "mappings": {
    "properties": {
	  "entityId": {
		"type": "keyword"
	  },
      "isIndividual": {
        "type": "boolean"
      },
      "riskLevel": {
		"type": "short"
      },
      "name": {
        "type": "text"
      },
      "akaList": {
         "type": "text",
         "fields": {
           "akaName": {
             "type": "keyword"
           }
         }
     },
      "addressList": {
         "type": "nested",
        "properties": {
          "country": {
            "type": "text"
          },
          "stateOrProvince": {
            "type": "text"
          },
          "city": {
            "type": "text"
          },
          "other":{
             "type": "text",
             "fields": {
               "address": {
               "type": "keyword"
              }
            }
          }
        }
  },
      "dateOfBirthList": {
        "type": "nested",
        "properties": {
          "dateOfBirth": {
            "type": "text"
          },
          "mainEntry": {
            "type": "boolean"
          }
        }
      },
      "placeOfBirth": {
        "type": "nested",
        "properties": {
          "placeOfBirth": {
            "type": "text"
          },
          "mainEntry": {
            "type": "boolean"
          }
        }
      },
      "gender": {
        "type": "keyword"
      },
      "emailList": {
         "type": "text",
         "fields": {
           "email": {
             "type": "keyword"
           }
         }
     },
      "websiteList": {
         "type": "text",
         "fields": {
           "website": {
             "type": "keyword"
           }
         }
     },
      "phoneNumberList": {
         "type": "text",
         "fields": {
           "phoneNumber": {
             "type": "keyword"
           }
         }
     },
      "idList": {
        "type": "nested",
        "properties": {
          "idType": {
            "type": "keyword"
          },
          "idNumber": {
            "type": "keyword"
          },
          "idCountry": {
            "type": "keyword"
          },
          "issueDate": {
            "type": "text"
          },
          "expirationDate": {
            "type": "text"
          }
        }
      },
      "nationalityList": {
        "type": "nested",
        "properties": {
          "country": {
            "type": "keyword"
          },
          "mainEntry": {
            "type": "boolean"
          }
        }
      },
      "organizationType": {
        "type": "keyword"
      },
      "citizenshipList": {
        "type": "nested",
        "properties": {
          "country": {
            "type": "keyword"
          },
          "mainEntry": {
            "type": "boolean"
          }
        }
      },
      "orgEstDate": {
        "type": "text"
      },
      "otherInfo": {
        "type": "nested",
        "properties": {
          "type": {
            "type": "keyword"
          },
          "info": {
            "type": "text"
          }
        }
      },
            "rules":{
                "type":"text",
                "fields":{
                    "rule":{
                        "type":"keyword"
                    }
                }
            },
            "riskChgHistory":{
                "type":"nested",
                "properties":{
                    "dateOfBirth":{
                        "type":"date"
                    },
                    "riskLevel":{
                        "type":"short"
                    },
                    "description":{
                        "type":"text"
                    }
                }
            }
    }
  }
}
`
