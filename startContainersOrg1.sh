#! /bin/bash -e

# COUCHDB_USER=""
# COUCHDB_PASSWORD=""
AR_SOLUTION_HOME="/media/izolotov/f163581d-53e6-4529-80d8-b822a479c7ab/dev/hyperledger/ar-solution"


docker network rm overnet
docker network create --attachable --driver overlay overnet

docker run -d --rm -it --network="overnet" \
--name orderer \
-p 7050:7050 \
-e ORDERER_GENERAL_LOGLEVEL=debug \
-e ORDERER_GENERAL_LISTENADDRESS=0.0.0.0 \
-e ORDERER_GENERAL_GENESISMETHOD=file \
-e ORDERER_GENERAL_GENESISFILE=/etc/hyperledger/configtx/genesis.block \
-e ORDERER_GENERAL_LOCALMSPID=OrdererMSP \
-e ORDERER_GENERAL_LOCALMSPDIR=/etc/hyperledger/crypto/orderer/msp \
-e ORDERER_GENERAL_TLS_ENABLED=true \
-e ORDERER_GENERAL_TLS_PRIVATEKEY=/etc/hyperledger/crypto/orderer/tls/server.key \
-e ORDERER_GENERAL_TLS_CERTIFICATE=/etc/hyperledger/crypto/orderer/tls/server.crt \
-e ORDERER_GENERAL_TLS_ROOTCAS='[/etc/hyperledger/crypto/orderer/tls/ca.crt,/etc/hyperledger/crypto/peerOrg1/tls/ca.crt,/etc/hyperledger/crypto/peerOrg2/tls/ca.crt]' \
-v $AR_SOLUTION_HOME/artifacts/channel:/etc/hyperledger/configtx \
-v $AR_SOLUTION_HOME/artifacts/channel/crypto-config/ordererOrganizations/orderer/orderers/orderer/:/etc/hyperledger/crypto/orderer \
-v $AR_SOLUTION_HOME/artifacts/channel/crypto-config/peerOrganizations/Org1/peers/peer0-org1/:/etc/hyperledger/crypto/peerOrg1 \
-v $AR_SOLUTION_HOME/artifacts/channel/crypto-config/peerOrganizations/Org2/peers/peer0-org2/:/etc/hyperledger/crypto/peerOrg2 \
-w /opt/gopath/src/github.com/hyperledger/fabric/orderers \
hyperledger/fabric-orderer:x86_64-1.0.5 orderer


docker run -d --rm -it --network="overnet" \
--name ca-org1 \
-p 7054:7054 \
-e CORE_PEER_NETWORKID=ca-org1 \
-e FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server \
-e FABRIC_CA_SERVER_CA_CERTFILE=/etc/hyperledger/fabric-ca-server-config/ca-org1-cert.pem \
-e FABRIC_CA_SERVER_CA_KEYFILE=/etc/hyperledger/fabric-ca-server-config/3d5484e30da800d8e715f85fbc3e5ce0abce31036d8eba832d0658d5d340e160_sk \
-e FABRIC_CA_SERVER_TLS_ENABLED=true \
-e FABRIC_CA_SERVER_TLS_CERTFILE=/etc/hyperledger/fabric-ca-server-config/ca-org1-cert.pem \
-e FABRIC_CA_SERVER_TLS_KEYFILE=/etc/hyperledger/fabric-ca-server-config/3d5484e30da800d8e715f85fbc3e5ce0abce31036d8eba832d0658d5d340e160_sk \
-v /media/izolotov/f163581d-53e6-4529-80d8-b822a479c7ab/dev/hyperledger/ar-solution/artifacts/channel/crypto-config/peerOrganizations/Org1/ca/:/etc/hyperledger/fabric-ca-server-config \
hyperledger/fabric-ca:x86_64-1.0.5 fabric-ca-server start -b admin:adminpw -d


# docker service create --network="overnet" \
# --name couchdb-peer0-org1 \
# -p 5984:5984 \
# -e COUCHDB_USER= \
# -e COUCHDB_PASSWORD= \
# hyperledger/fabric-couchdb:x86_64-1.0.5


docker run -d --rm -it --network="overnet" \
--name couchdb-peer0-org1 \
-p 5984:5984 \
-e COUCHDB_USER= \
-e COUCHDB_PASSWORD= \
-e CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=overnet \
hyperledger/fabric-couchdb:x86_64-1.0.5


docker run -d --rm -it --network="overnet" \
--link orderer:orderer \
--link couchdb-peer0-org1:couchdb-peer0-org1 \
--name peer0-org1 \
-p 7051:7051 \
-p 7053:7053 \
-p 7052:7052 \
-e CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock \
-e CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=overnet \
-e CORE_LOGGING_LEVEL=DEBUG \
-e CORE_PEER_GOSSIP_USELEADERELECTION=true \
-e CORE_PEER_GOSSIP_ORGLEADER=false \
-e CORE_PEER_GOSSIP_SKIPHANDSHAKE=true \
-e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/crypto/peer/msp \
-e CORE_PEER_TLS_ENABLED=true \
-e CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/crypto/peer/tls/server.key \
-e CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/crypto/peer/tls/server.crt \
-e CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/crypto/peer/tls/ca.crt \
-e CORE_PEER_ID=peer0-org1 \
-e CORE_PEER_CHAINCODELISTENADDRESS=peer0-org1:7052 \
-e CORE_PEER_LOCALMSPID=Org1MSP \
-e CORE_PEER_ADDRESSAUTODETECT=true \
-e CORE_LEDGER_STATE_STATEDATABASE=CouchDB \
-e CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=couchdb-peer0-org1:5984 \
-e CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME= \
-e CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD= \
-w /opt/gopath/src/github.com/hyperledger/fabric/peer \
-v /var/run/:/host/var/run/ \
-v $AR_SOLUTION_HOME/artifacts/channel/crypto-config/peerOrganizations/Org1/peers/peer0-org1/:/etc/hyperledger/crypto/peer \
hyperledger/fabric-peer:x86_64-1.0.5 peer node start


docker run -d --rm -it --network="overnet" \
--name couchdb-peer1-org1 \
-p 5985:5984 \
-e COUCHDB_USER= \
-e COUCHDB_PASSWORD= \
-e CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=overnet \
hyperledger/fabric-couchdb:x86_64-1.0.5


docker run -d --rm -it --network="overnet" \
--link orderer:orderer \
--link couchdb-peer1-org1:couchdb-peer1-org1 \
--name peer1-org1 \
-p 7056:7051 \
-p 7058:7053 \
-p 7057:7052 \
-e CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock \
-e CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=overnet \
-e CORE_LOGGING_LEVEL=DEBUG \
-e CORE_PEER_GOSSIP_USELEADERELECTION=true \
-e CORE_PEER_GOSSIP_ORGLEADER=false \
-e CORE_PEER_GOSSIP_SKIPHANDSHAKE=true \
-e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/crypto/peer/msp \
-e CORE_PEER_TLS_ENABLED=true \
-e CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/crypto/peer/tls/server.key \
-e CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/crypto/peer/tls/server.crt \
-e CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/crypto/peer/tls/ca.crt \
-e CORE_PEER_ID=peer1-org1 \
-e CORE_PEER_CHAINCODELISTENADDRESS=peer1-org1:7052 \
-e CORE_PEER_LOCALMSPID=Org1MSP \
-e CORE_PEER_ADDRESSAUTODETECT=true \
-e CORE_LEDGER_STATE_STATEDATABASE=CouchDB \
-e CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=couchdb-peer1-org1:5984 \
-e CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME= \
-e CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD= \
-w /opt/gopath/src/github.com/hyperledger/fabric/peer \
-v /var/run/:/host/var/run/ \
-v $AR_SOLUTION_HOME/artifacts/channel/crypto-config/peerOrganizations/Org1/peers/peer1-org1/:/etc/hyperledger/crypto/peer \
hyperledger/fabric-peer:x86_64-1.0.5 peer node start
