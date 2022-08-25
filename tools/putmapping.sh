#!/bin/bash

curl localhost:9200/metadata -XPUT -d'{"mappings":{"objects":{"properties":{"name":{"type":"string","index":"not_analyzed"},"versions":{"type":"integer"},"size":{"type":"integer"},"hash":{"type":"string"}}}}}'
