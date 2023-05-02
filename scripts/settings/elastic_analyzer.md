## Settings for tag

````
{
  "settings": {
    "analysis": {
      "analyzer": {
        "tag_analyzer": {
          "type":      "custom",
          "tokenizer": "standard",
          "char_filter": [
            "tag_special_filter"
          ],
          "filter": [
            "lowercase",
            "asciifolding",
            "custom_stem",
            "unique"
          ]
        },
        "product_name_analyzer": {
          "type":      "custom",
          "tokenizer": "standard",
          "char_filter": [
            "single_char_filter",
            "useless_filter"
          ],
          "filter": [
            "lowercase",
            "asciifolding",
            "custom_stem",
            "stop",
            "unique",
            "remove_digit"
          ]
        }
      },
      "char_filter": {
        "tag_special_filter": {
          "type": "pattern_replace",
          "pattern": "[|&,']",
          "replacement": " "
        },
        "single_char_filter": {
          "type": "pattern_replace",
          "pattern": "(\\s[a-zA-Z]{1}\\s|\\s[0-9]{1}\\s|[,'\\\/!])",
          "replacement": " "
        },
        "useless_filter": {
          "type": "pattern_replace",
          "pattern": "([(]{0,1}[0-9]{0,}[.'\"-+,]{1}[0-9]{0,}[”'\"]{0,1}[)'\"]{0,1})",
          "replacement": " "
        }
      },
      "filter":{
        "custom_stem":{
            "type":"stemmer",
            "language":"english"
        },
        "remove_digit": {
            "type":"keep_types",
            "types": [ "<NUM>" ],
            "mode": "exclude"
        }
      }
    }
  },
  "mappings": {
    "tags": {
        "properties": {
            "tag": {
                "type":     "text",
                "analyzer": "tag_analyzer"
            },
            "product_id": {
                "type":     "text",
                "index":    "false"
            }
        }
    },
    "product": {
        "properties": {
            "name": {
                "type":     "text",
                "analyzer": "product_name_analyzer"
            },
            "id": {
                "type":     "text",
                "index":    "false"
            }
        }
    }
  }
}
````


## Index Setting for ecommerce_product
```
{
    "settings": {
        "analysis": {
            "analyzer": {
                "product_name_analyzer": {
                    "type":      "custom",
                    "tokenizer": "standard",
                    "char_filter": [
                        "single_char_filter",
                        "useless_filter"
                    ],
                    "filter": [
                        "lowercase",
                        "asciifolding",
                        "custom_stem",
                        "stop",
                        "unique",
                        "remove_digit"
                    ]
                }
            },
            "char_filter": {
                "single_char_filter": {
                    "type": "pattern_replace",
                    "pattern": "(\\s[a-zA-Z]{1}\\s|\\s[0-9]{1}\\s|[,'\\\/!])",
                    "replacement": " "
                },
                "useless_filter": {
                    "type": "pattern_replace",
                    "pattern": "([(]{0,1}[0-9]{0,}[.'\"-+,]{1}[0-9]{0,}[”'\"]{0,1}[)'\"]{0,1})",
                    "replacement": " "
                }
            },
            "filter":{
                "custom_stem":{
                    "type":"stemmer",
                    "language":"english"
                },
                "remove_digit": {
                    "type":"keep_types",
                    "types": [ "<NUM>" ],
                    "mode": "exclude"
                }
            }
        }
    }
}
```

## Index Setting for ecommerce_tag


```
{
    "settings": {
        "analysis": {
            "analyzer": {
                "tag_analyzer": {
                    "filter": [
                        "lowercase",
                        "asciifolding",
                        "custom_stem",
                        "unique"
                    ],
                    "char_filter": [
                        "tag_special_filter"
                    ],
                    "type": "custom",
                    "tokenizer": "standard"
                }
            },
            "char_filter": {
              "tag_special_filter": {
                  "pattern": "[|&,']",
                  "type": "pattern_replace",
                  "replacement": " "
              }
            }
            "filter":{
                "custom_stem":{
                    "type":"stemmer",
                    "language":"english"
                }
            }
        }
    }
}
```
