#!/bin/bash

var=$RANDOM
name=\"marbles${var}\"
echo $name

args='{"chaincode": "marbles","args": ["initMarble",'$name',"blue","36","tom"]}'
echo $args

curl -X POST "http://29.225.33.177:8080//gateway/api/v1/channels/test/transactions" -H "accept: */*" -H "Content-Type: application/json" -d '{"chaincode": "marbles","args": ["initMarble",'$name',"blue","36","tom"]}'


#keyLabel=\"Baas98799899214734247020142072493612775431\"
#
#echo $keyLabel
#
#curl -X POST "http://47.95.204.66:34999/brilliance/netsign/sign" -H "accept: */*" -H "Content-Type: application/json" -d '{"keyLabel": '$keyLabel',"origBytes": "123","digestAlg": "123"}'




