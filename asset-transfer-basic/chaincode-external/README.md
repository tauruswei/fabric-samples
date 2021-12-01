# Asset-Transfer-Basic as an external service

This sample provides an introduction to how to use external builder and launcher scripts to run chaincode as an external service to your peer. For more information, see the [Chaincode as an external service](https://hyperledger-fabric.readthedocs.io/en/latest/cc_launcher.html) topic in the Fabric documentation.

**Note:** each organization in a real network would need to setup and host their own instance of the external service. For simplification purpose, in this sample we use the same instance for both organizations.

## Setting up the external builder and launcher

打开 `fabric-samples/config/core.yaml` 
修改 the field `externalBuilders` as the following:
```
externalBuilders:
    - path: /opt/gopath/src/github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-external/sampleBuilder
      name: external-sample-builder
```
This configuration sets the name of the external builder as `external-sample-builder`, and the path of the builder to the scripts provided in this sample. Note that this is the path within the peer container, not your local machine.

打开 `test-network/docker/docker-compose-test-net.yaml`
容器 `peer0.org1.example.com` 和 `peer0.org2.example.com` 新增如下两行挂载:

```
        - ../..:/opt/gopath/src/github.com/hyperledger/fabric-samples
        - ../../config/core.yaml:/etc/hyperledger/fabric/core.yaml
```

This will mount the fabric-sample builder into the peer container so that it can be found at the location specified in the config file,
and override the peer's core.yaml config file within the fabric-peer image so that the config file modified above is used.

## Packaging and installing Chaincode

The Asset-Transfer-Basic external chaincode requires two environment variables to run, `CHAINCODE_SERVER_ADDRESS` and `CHAINCODE_ID`, which are described and set in the `chaincode.env` file.

The peer needs a corresponding `connection.json` configuration file so that it can connect to the external Asset-Transfer-Basic service.

The address specified in the `connection.json` must correspond to the `CHAINCODE_SERVER_ADDRESS` value in `chaincode.env`, which is `asset-transfer-basic.org1.example.com:9999` in our example.

First, create a `code.tar.gz` archive containing the `connection.json` file:

```
cd fabric-samples/asset-transfer-basic/chaincode-external
tar cfz code.tar.gz connection.json
```

Then, create the chaincode package, including the `code.tar.gz` file and the supplied `metadata.json` file:

```
tar cfz asset-transfer-basic-external.tgz metadata.json code.tar.gz
```

You are now ready to use the external chaincode. We will use the `test-network` sample to get a network setup and make use of it.

## 修改 test-network/configtx/configtx.yaml

```
Application: &ApplicationDefaults

    # Organizations is the list of orgs which are defined as participants on
    # the application side of the network
    Organizations:

    # Policies defines the set of policies at this level of the config tree
    # For Application policies, their canonical path is
    #   /Channel/Application/<PolicyName>
    Policies:
        Readers:
            Type: ImplicitMeta
            Rule: "ANY Readers"
        Writers:
            Type: ImplicitMeta
            Rule: "ANY Writers"
        Admins:
            Type: ImplicitMeta
            Rule: "ANY Admins"
        LifecycleEndorsement:
            Type: ImplicitMeta
            Rule: "ANY Endorsement"
        Endorsement:
            Type: ImplicitMeta
            Rule: "ANY Endorsement"
```

## Starting the test network

In a different terminal, from the `test-network` sample directory starts the network using the following command:

```
cd fabric-samples/test-network
./network.sh up createChannel -c mychannel -ca
```

This starts the test network and creates the channel. We will now proceed to installing our external chaincode package.

## Installing the external chaincode

We can't use the `test-network/network.sh` script to install our external chaincode so we will have to do a bit more work by hand but we can still leverage part of the test-network scripts to make this easier.

First, get the functions to setup your environment as needed by running the following command (this assumes you are still in the `test-network` directory):

```
. ./scripts/envVar.sh
```

安装 `asset-transfer-basic-external.tar.gz` chaincode on org1:

```
export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
setGlobals 1
../bin/peer lifecycle chaincode install ../asset-transfer-basic/chaincode-external/asset-transfer-basic-external.tgz
```

setGlobals simply defines a bunch of environment variables suitable to act as one organization or another, org1 or org2.

安装 `asset-transfer-basic-external.tar.gz` chaincode on org2:

```
setGlobals 2
../bin/peer lifecycle chaincode install ../asset-transfer-basic/chaincode-external/asset-transfer-basic-external.tgz
```

This will output the chaincode pakage identifier such as `basic_1.0:0262396ccaffaa2174bc09f750f742319c4f14d60b16334d2c8921b6842c090` that you will need to use in the following commands.

新建环境变量 `PKGID=basic_1.0:0262396ccaffaa2174bc09f750f742319c4f14d60b16334d2c8921b6842c090`

```
export PKGID=basic_1.0:0262396ccaffaa2174bc09f750f742319c4f14d60b16334d2c8921b6842c090
```

你也可以通过下面的命令查询 `package-id`:

```
setGlobals 1
../bin/peer lifecycle chaincode queryinstalled --peerAddresses localhost:7051 --tlsRootCertFiles organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
```

切换到目录： `fabric-samples/asset-transfer-basic/chaincode-external` ,修改 `chaincode.env` ,set `CHAINCODE_ID` euqal to the
 chaincode `package-id` obtained above.


## Running the Asset-Transfer-Basic external service

切换到 `fabric-samples/asset-transfer-basic/chaincode-external` 目录,编译合约镜像:

对于 amd 系统 执行下面的命令编译镜像

```
docker build -t hyperledger/asset-transfer-basic .
```
对于 Apple Silicon 系统 执行下面的命令编译镜像
```
docker buildx build --platform linux/amd64 -t hyperledger/asset-transfer-basic .
```

启动合约容器:

```
docker run -it --rm --name asset-transfer-basic.org1.example.com --hostname asset-transfer-basic.org1.example.com --env-file chaincode.env --network=docker_test hyperledger/asset-transfer-basic
```

This will start the container and start the external chaincode service within it.

## 部署合约


```
setGlobals 2
../bin/peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --channelID mychannel --name basic --version 1.0 --package-id $PKGID --sequence 1

setGlobals 1
../bin/peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --channelID mychannel --name basic --version 1.0 --package-id $PKGID --sequence 1

../bin/peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --channelID mychannel --name basic --peerAddresses localhost:7051 --tlsRootCertFiles $PWD/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt --version 1.0 --sequence 1
```

This approves the chaincode definition for both orgs and commits it using org1. This should result in an output similar to:

```
2020-08-05 15:41:44.982 PDT [chaincodeCmd] ClientWait -> INFO 001 txid [6bdbe040b99a45cc90a23ec21f02ea5da7be8b70590eb04ff3323ef77fdedfc7] committed with status (VALID) at localhost:7051
2020-08-05 15:41:44.983 PDT [chaincodeCmd] ClientWait -> INFO 002 txid [6bdbe040b99a45cc90a23ec21f02ea5da7be8b70590eb04ff3323ef77fdedfc7] committed with status (VALID) at localhost:9051
```

Now that the chaincode is deployed to the channel, and started as an external service, it can be used as normal.

## Using the Asset-Transfer-Basic external chaincode

切换到 `fabric-samples/test-network` 目录，执行下面的命令

```
peer chaincode invoke -n basic -c '{"Args":["InitLedger"]}' -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel
peer chaincode invoke -n basic -c '{"Args":["ReadAsset","asset1"]}' -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel
```

## 清理网络,防止有残留数据
```
bash network.sh down
docker system prune
docker volume rm $(docker volume ls -qf dangling=true)
```