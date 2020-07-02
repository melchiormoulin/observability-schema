#!/bin/sh
curl -XPUT "localhost:9200/_index_template/template_1?pretty" -H 'Content-Type: application/json' -d @template.json