#!/bin/sh
curl -XPUT "localhost:9200/_index_template/template_1?pretty" -H 'Content-Type: application/json' -d @../../examples/elasticsearch/output/template.json