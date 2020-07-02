#!/bin/sh
curl -X POST "localhost:9200/filebeat-001/_doc/?pretty" -H 'Content-Type: application/json' -d'
{
  "user_id" : "kimchy",
  "port": 8080,
  "@timestamp" : "2020-06-24T17:47:12",
  "test":"mytest",
  "random_nb":42,
  "message" : "trying out Elasticsearch"
}
'
sleep 3
curl -XGET "localhost:9200/filebeat-001/_search" | jq .
