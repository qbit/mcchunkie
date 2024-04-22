#!/bin/sh

echo "line one is here" > /tmp/test
echo "line two is here" >> /tmp/test
echo "line <two> is here" >> /tmp/test

curl --basic --user got:butters  --data-urlencode file@/tmp/test -o - "http://localhost:8043/_got"
curl -X POST -H "Content-Type: application/json" --data @test_body.json --basic --user got:butters -o - "http://localhost:8043/_got/v2"
