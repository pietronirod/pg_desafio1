#!/bin/bash

URL="http://localhost:8080"
API_KEY="123"

for i in {1..7}
do
	echo "Request #$i"
	curl -i -X GET $URL 
	echo -e "\n"
done
