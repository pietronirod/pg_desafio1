#!/bin/bash

URL="http://localhost:8080"
API_KEY="123"

for i in {1..10}
do
	echo "Request #$i"
	curl -i -X GET $URL -H "API_KEY: $API_KEY"
	echo -e "\n"
done
