#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

jq --version > /dev/null 2>&1
if [ $? -ne 0 ]; then
	echo "Please Install 'jq' https://stedolan.github.io/jq/ to execute this script"
	echo
	exit 1
fi
starttime=$(date +%s)

echo "POST request Enroll on Org1  ..."
echo
ORG1_TOKEN=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=Jim&orgName=org1')
echo $ORG1_TOKEN
ORG1_TOKEN=$(echo $ORG1_TOKEN | jq ".token" | sed "s/\"//g")
echo
echo "ORG1 token is $ORG1_TOKEN"
echo
echo "POST request Enroll on Org2 ..."
echo
ORG2_TOKEN=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=Barry&orgName=org2')
echo $ORG2_TOKEN
ORG2_TOKEN=$(echo $ORG2_TOKEN | jq ".token" | sed "s/\"//g")
echo
echo "ORG2 token is $ORG2_TOKEN"
echo
echo
echo "POST request Create channel  ..."
echo
curl -s -X POST \
  http://localhost:4000/channels \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"channelName":"mychannel",
	"channelConfigPath":"../artifacts/channel/mychannel.tx"
}'
echo
echo
sleep 5
echo "POST request Join channel on Org1"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/peers \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer1","peer2"]
}'
echo
echo

echo "POST request Join channel on Org2"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/peers \
  -H "authorization: Bearer $ORG2_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer1","peer2"]
}'
echo
echo


echo "POST Install test chaincode on Org1"
echo
curl -s -X POST \
  http://localhost:4000/chaincodes \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
        "peers": ["peer1", "peer2"],
        "chaincodeName":"cctest",
        "chaincodePath":"github.com/ar_solution",
        "chaincodeVersion":"v0"
}'
echo
echo

echo "POST Install test chaincode on Org2"
echo
curl -s -X POST \
  http://localhost:4000/chaincodes \
  -H "authorization: Bearer $ORG2_TOKEN" \
  -H "content-type: application/json" \
  -d '{
        "peers": ["peer1", "peer2"],
        "chaincodeName":"cctest",
        "chaincodePath":"github.com/ar_solution",
        "chaincodeVersion":"v0"
}'
echo
echo


echo "POST instantiate test chaincode on peer1 of Org1"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
        "chaincodeName":"cctest",
        "chaincodeVersion":"v0",
        "args":[]
}'
echo
echo

echo "POST invoke newSaleDoc"
echo
TRX_ID=$(curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/cctest \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
        "fcn":"newSaleDoc",
        "args":["fd5c0e62-53bf-11e8-a305-8e9b35eee675","0000-000001","Org1","2018-04-02T00:00:00","Org1","Org2","35000","0","Продажа (0000-000001 от 01.05.2018)"]
}')
echo "Transacton ID is $TRX_ID"
echo
echo

echo "POST invoke newPurchaseDoc"
echo
TRX_ID=$(curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/cctest \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
        "fcn":"newPurchaseDoc",
        "args":["d2755510-529e-11e8-a305-8e9b35eee675","0000-000002","Org1","2018-04-02T00:00:00","Org1","Org2","0","63900","Приход (0000-000537 от 02.04.2018)"]
}')
echo "Transacton ID is $TRX_ID"
echo
echo

echo "POST invoke newExpenseDoc"
echo
TRX_ID=$(curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/cctest \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
        "fcn":"newExpenseDoc",
        "args":["b55a0e3e-541c-11e8-a305-8e9b35eee675","0000-000003","Org1","2018-04-02T00:00:00","Org1","Org2","63900","0","Оплата (0000-000001 от 02.04.2018)"]
}')
echo "Transacton ID is $TRX_ID"
echo
echo

echo "POST invoke newAdmissionDoc"
echo
TRX_ID=$(curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/cctest \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
        "fcn":"newAdmissionDoc",
        "args":["fd5c0e69-53bf-11e8-a305-8e9b35eee675","0000-000004","Org1","2018-04-02T00:00:00","Org1","Org2","0","35000","Оплата (0000-000001 от 09.05.2018)"]
}')
echo "Transacton ID is $TRX_ID"
echo
echo

echo "GET query 'all purchase docs'"
echo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/cctest?peer=peer1&fcn=getAllPurchaseDocs&args=%5B%5D" \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json"
echo
echo

echo "GET query 'all expense docs'"
echo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/cctest?peer=peer1&fcn=getAllExpenseDocs&args=%5B%5D" \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json"
echo
echo

echo "GET query 'all sale docs'"
echo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/cctest?peer=peer1&fcn=getAllSaleDocs&args=%5B%5D" \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json"
echo
echo

echo "GET query 'all admission docs'"
echo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/cctest?peer=peer1&fcn=getAllAdmissionDocs&args=%5B%5D" \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json"
echo
echo

echo "Total execution time : $(($(date +%s)-starttime)) secs ..."
