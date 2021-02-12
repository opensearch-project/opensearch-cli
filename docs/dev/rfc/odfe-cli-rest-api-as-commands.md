# RFC: Support REST api as command

## Objective

This document describes how to support [REST api](https://www.elastic.co/guide/en/elasticsearch/reference/current/rest-apis.html) in a generic way.
As of now, Elasticsearch support api in different categories like index api, cat api, cluster api, document api, etc to
configure and access their features. It will take tremendous amount of time and effort to provide every api as native commands.
Hence, we will be categorizing those REST api based on GET/PUT/POST/DELETE  and allow users to perform their request 
in a generic way. This will support any future api without any additional support.
## Available commands

1. [get command](./odfe-cli-rest-api-as-commands.md#1-get-command)
2. [post command](./odfe-cli-rest-api-as-commands.md#2-post-command)
3. [put command](./odfe-cli-rest-api-as-commands.md#3-put-command)
4. [delete command](./odfe-cli-rest-api-as-commands.md#4-delete-command)

## Common parameters

The following options can be applied to all the odfe-cli commands. These parameters are supported by Elasticsearch as mentioned [here](https://www.elastic.co/guide/en/elasticsearch/reference/current/common-options.html).
In the future we will support more parameters or include multiple values that can be added as flag to any commands.



|Parameter name|Short hand	|Default|Usage|
|---------------|-----------|-------|-------|
|--pretty	    |-b	        | false	|if true, the response returned will be pretty formatted	|
|--output-format |-o	|json	|other value is yaml. The result to be returned in the (if supported) more readable yaml format.	|
|--human	| N/A	| false	|the output statistics are returned in a human friendly format. (e.g. `"exists_time": "1h"` or `"size": "1kb"`) instead of (e.g. `"exists_time_in_millis": 3600000` or `"size_in_bytes": 1024`). 	|
|--filter-path	|-f	| N/A	|filter the response fields returned by cluster. It accepts a comma separated list of filters expressed with the dot notation.  eg: `--filter-path=“took,hits.hits._id,hits.hits._score”`	|



## 1. get command

### Description

Use GET api to execute requests against Elasticsearch cluster. This command enables you to run any GET based REST api commands across all 
categories.

### Synopsis

```
odfe-cli curl get options [common parameters]
```

### Options

|Name|Required|Short hand|Usage|Separator for multiple values|
|---| :---: |:---:|---	|:---:|
| --path	|Y	|-P	|url-path	|N/A	|
| --query-params|N	|-q	|url query parameters	|&	|
| --headers	 |N	|-H	|pass additional information with request. It consists of case-insensitive name followed by a colon (`:`), then by its value.	|semi colon ( ; )	|
|--data	|N	|-d	|Sends the specified data in the command to Elasticsearch". If value starts with the letter @, the rest should be a file name to read the data from.	|N/A	|
|--data-binary	|N	|-b	|This is similar to "--data", except that newlines and carriage returns are preserved and conversions are never done. This is required if input file is compressed.	|N/A	|
|--help	|N	|-h	|Help for get command	|N/A	|

### Common Parameters

see [here](./odfe-cli-rest-api-as-commands.md#common-parameters)


**Note: users can escape separator like ‘\\;’ to use as part of value.**

### Example:

1. To get document count for an index

```
odfe-cli curl get --path "_cat/count/my-index-01" \
                 --query-params "v=true" \
                 --pretty
```

Response

```
epoch         timestamp        count
1612311592    12:24:24          100
```
2.  To  return the health status of a cluster

```
odfe-cli curl get --path "_cluster/health" --pretty
```

Response

```
{
    "cluster_name" : "odfe-cli",
    "status" : "green",
    "timed_out" : false,
    "number_of_nodes" : 2,
    "number_of_data_nodes" : 2,
    "active_primary_shards" : 1,
    "active_shards" : 1,
    "relocating_shards" : 0,
    "initializing_shards" : 0,
    "unassigned_shards" : 0,
    "delayed_unassigned_shards": 0,
    "number_of_pending_tasks" : 0,
    "number_of_in_flight_fetch": 0,
    "task_max_waiting_in_queue_millis": 0,
    "active_shards_percent_as_number": 100.0
}
```

 3.  To explain cluster allocation for index “my-index-01”, shard ‘0’


```
odfe-cli curl get --path "_cluster/allocation/explain" \
                 --data  '{
                    "index": "my-index-01",
                    "shard": 0,
                    "primary": false,
                    "current_node": "nodeA"                         
                  }'
```

## 2. post command

### Description

Use POST api to execute requests against Elasticsearch cluster. This command enables you to run any POST based REST api
commands across all categories.

### Synopsis

```
odfe-cli curl post [options] [common paramters]
              
```

### Options

see [here](./odfe-cli-rest-api-as-commands.md#Options)

### Common Parameters

see [here](./odfe-cli-rest-api-as-commands.md#common-parameters)

### Example:

1. To  change the allocation of shards in a cluster. Here, data is saved in a file “reroute.json”

```
$ cat reroute.json
{
  "commands": [
    {
        "move": {
           "index": "odfe-cli",
           "shard": 0,
           "from_node": "odfe-node1",
           "to_node": "odfe-node2"
        }
     },
     {
        "allocate_replica": {
            "index": "test",
            "shard": 1,
            "node": "odfe-node3"
         }
     }
    ]
}

$ odfe-cli curl post --path "_cluster/reroute" \
                    --data @reroute.json
```

2. To insert a document to an index 

```
$ odfe-cli curl post --path "my-index-01/_doc" \
                   --data '
                        {
                            "message": "insert document",
                            "ip": {
                                "address": "127.0.0.1"
                            }
                        }'
```



3. Search index with compressed data, also accept response in compressed format provided compression is enabled by setting “http.compression: true”
 ```
 $ odfe-cli curl post --path        "_search" \
                     --headers     "Content-Encoding : gzip;Accept-Encoding: gzip, deflate" \
                     --data-binary  @/tmp/req.txt.gz`
    
```    

## 3. put command

### Description

Use PUT api to execute requests against Elasticsearch. This command enables you to run any PUT based REST api
commands across all categories.

### Synopsis

```
odfe-cli curl put options [common paramters]]
              
```

### Options

see [here](./odfe-cli-rest-api-as-commands.md#Options)

### Common Parameters

see [here]((./odfe-cli-rest-api-as-commands.md#common-parameters))

### Example

1. Create a knn index

```
odfe-cli curl put --path "my-knn-index" ---pretty \
                --data '
                    {
                        "settings" : {
                            "number_of_shards" :   1,
                            "number_of_replicas" : 0,
                            "index.knn" : true
                        },
                        "mappings": {
                            "properties": {
                                "my_dense_vector": {
                                    "type": "knn_vector",
                                    "dimension": 2
                                },
                                "color" : {
                                    "type" : "keyword"
                                }
                            }
                        }
                    }'
```

2. Update cluster settings transiently

```
odfe-cli curl put --path             "_cluster/settings" \
                --query-parameters "flat_settings=true" --pretty \
                --data '
                {
                    "transient" : {
                        "indices.recovery.max_bytes_per_sec" : "20mb"
                    }
                }'
```

## 4. delete command

### Description

Use DELETE api to execute requests against Elasticsearch. This command enables you to run any DELETE based
REST api commands across all categories.

### Synopsis

```
odfe-cli curl delete [options] [common paramters]
              
```

### Options

|Name|Required|Short hand|Usage|Separator for multiple values|
|---| :---: |:---:|---	|:---:|
| --path	|Y	|-P	|url-path	|N/A	|
| --query-params|N	|-q	|url query parameters	|&	|
|--help	|N	|-h	|Help for get command	|N/A	|


### Common Parameters

see [here](./odfe-cli-rest-api-as-commands.md#common-parameters)


### Example

1. Delete a document from an index. 

```
odfe-cli curl delete --path         "my-index/_doc/1" \
                   --query-params  "routing=odfe-node1"
```

## Request for Comments:

We would like comments and feedback on the proposal for supporting REST api in our  ODFE CLI tool [here](https://github.com/opendistro-for-elasticsearch/odfe-cli/issues/37). Some specific questions we’re seeking feedback include

1. Can you share your use case (or api ) that you think will not be covered with above synopsis?
2. “curl” is chosen as a  category (or parent command) for get,put,post,delete commands, to  differentiate from plugin names. Do you prefer any other command name?

