#!/bin/sh

set -x

PW="testing1234"

echo "test one two three" > /tmp/test

curl --basic --user "got:${PW}"  --data-urlencode file@/tmp/test -o - "http://localhost:8043/_got"
curl -X POST -H "Content-Type: application/json" --data @test_body.json --basic --user "got:${PW}" -o - "http://localhost:8043/_got/v2"

