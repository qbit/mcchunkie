#!/bin/sh

PW="$(op item get got-test --field password)"

echo "test one two three" > /tmp/test

curl --basic --user got:${PW}  --data-urlencode file@/tmp/test -o - "https://suah.dev/_got"
curl -X POST -H "Content-Type: application/json" --data @test_body.json --basic --user got:${PW} -o - "https://suah.dev/_got/v2"

