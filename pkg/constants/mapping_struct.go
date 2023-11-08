package constants

//es映射定义

// ADDR_MAPPING 风险名单信息映射
const ADDR_MAPPING = `
{
    "mappings":{
        "properties":{
			"addressId":{
				"type":"keyword"
			},
            "waAddr":{
                "type":"keyword"
            },
            "entityId":{
                "type":"keyword"
            },
            "balance":{
                "type":"double"
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
            "isContract":{
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
                        "type":"date",
                        "format":"yyyy-MM-dd HH:mm:ss"
                    },
                    "dsRules":{
                        "type":"nested",
                        "properties":{
                            "dsRuKey":{
                                "type":"keyword"
                            },
                            "dsRuType":{
                                "type":"keyword"
                            },
                            "dsRuCode":{
                                "type":"keyword"
                            },
                            "dsRuDesc":{
                                "type":"text"
                            },
                            "dsRuStatus":{
                                "type":"keyword"
                            },
                            "dsRuExpress":{
                                "type":"text"
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
                        "type":"nested",
                        "properties":{
                            "ruKey":{
                                "type":"keyword"
                            },
                            "ruType":{
                                "type":"keyword"
                            },
                            "ruCode":{
                                "type":"keyword"
                            },
                            "ruDesc":{
                                "type":"text"
                            },
                            "status":{
                                "type":"keyword"
                            },
                            "ruExpress":{
                                "type":"text"
                            }
                        }
                    }
,
            "riskChgHistory":{
                "type":"nested",
                "properties":{
                    "dateOfChange":{
                        "type":"date",
                        "format":"yyyy-MM-dd"
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
			"chain":{
				"type":"keyword"
			},
			"inputCount":{
				"type":"short"
			},
			"inputValue":{
				"type":"double"
			},
			"outputCount":{
				"type":"short"
			},
			"outputValue":{
				"type":"double"
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
              "type": "date",
              "format": "yyyy-MM-dd HH:mm:ss"
            },
            "blockHeight":{
                "type":"text"
            },
            "blockHash":{
                "type":"text"
            },
            "value":{
                "type":"double"
            },
			"valueText":{
				"type":"text"
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
                        "type":"double"
                    },
                    "valueText":{
				       "type":"text"
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
                        "type":"double"
                    },
					"valueText":{
						"type":"text"
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
                        "type":"double"
                    },
					"valueText":{
						"type":"text"
					},
                    "subTraces":{
                        "type":"short"
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
                "type":"nested",
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
                        "type":"double"
                    },
					"amountText":{
						"type":"text"
					}
                }
            },
            "rules":{
                        "type":"nested",
                        "properties":{
                            "ruKey":{
                                "type":"keyword"
                            },
                            "ruType":{
                                "type":"keyword"
                            },
                            "ruCode":{
                                "type":"keyword"
                            },
                            "ruDesc":{
                                "type":"text"
                            },
                            "status":{
                                "type":"keyword"
                            },
                            "ruExpress":{
                                "type":"text"
                            }
                        }
                    },
            "riskChgHistory":{
                "type":"nested",
                "properties":{
                    "dateOfChange":{
                         "type": "date",
                      "format": "yyyy-MM-dd"
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
             "type": "date",
            "format": "yyyy-MM-dd"
          },
          "mainEntry": {
            "type": "boolean"
          }
        }
      },
      "placeOfBirthList": {
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
             "type": "date",
             "format": "yyyy-MM-dd"
          },
          "expirationDate": {
             "type": "date",
            "format": "yyyy-MM-dd"
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
         "type": "date",
         "format": "yyyy-MM-dd"
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
                        "type":"nested",
                        "properties":{
                            "ruKey":{
                                "type":"keyword"
                            },
                            "ruType":{
                                "type":"keyword"
                            },
                            "ruCode":{
                                "type":"keyword"
                            },
                            "ruDesc":{
                                "type":"text"
                            },
                            "status":{
                                "type":"keyword"
                            },
                            "ruExpress":{
                                "type":"text"
                            }
                        }
                    },
            "riskChgHistory":{
                "type":"nested",
                "properties":{
                    "dateOfChange":{
                         "type": "date",
                         "format": "yyyy-MM-dd HH:mm:ss"
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
