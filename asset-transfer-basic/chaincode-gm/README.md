# 国密 chaincode as an external service

### 1、修改 core.yaml

  打开 `fabric-samples/config/core.yaml`修改 the field `externalBuilders` as the following:

   ```bash
   # 注意修改合约 chaincode 的地址
    externalBuilders:
        - path: /opt/gopath/src/github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-gm/sampleBuilder
          name: external-sample-builder
   ```

###2、修改 docker-compose-test-net.yaml

  修改 `test-network/docker/docker-compose-test-net.yaml`，容器 `peer0.org1.example.com`和 `peer0.org2.example.com`新增如下两行挂载:

   ```bash
    - ../..:/opt/gopath/src/github.com/hyperledger/fabric-samples
    - ../../config/core.yaml:/etc/hyperledger/fabric/core.yaml
   ```

###3、启动网络，并创建 channel

   ```bash
    bash network.sh up createChannel
   ```

###4、修改 connection.json(两种情况)
- ####（1）如果是编译器启动 合约 服务，合约地址应该设置为本机的ip

    ```bash
    # 192.168.2.150 是我本机的 ip
    {
      "address": "192.168.2.150:9999",
      "dial_timeout": "10s",
      "tls_required":false
    }
    ```

- ####（2）如果是 docker 启动合约服务

    ```bash
    # gm 是合约的名称
    {
      "address": "gm:9999",
      "dial_timeout": "10s",
      "tls_required":false
    }
    ```

###5、修改 metadata.json

  ```
  # gm 是 chaincode 的名字
  {
      "type": "external",
      "label": "gm_1.0"
  }
  ```

###6、重新打包

   ```bash
   export PKGNAME = gm.tgz
   cd fabric-samples/asset-transfer-basic/chaincode-gm
   tar cfz code.tar.gz connection.json
   # 注意压缩包的名称
   tar cfz $PKGNAME metadata.json code.tar.gz
   ```

###7、安装合约

   ```bash
    cd fabric-samples/test-network
    . ./scripts/envVar.sh
    export PATH=${PWD}/../bin:${PWD}:$PATH
    export FABRIC_CFG_PATH=$PWD/../config/
    setGlobals 1
    ../bin/peer lifecycle chaincode install ../asset-transfer-basic/chaincode-gm/$PKGNAME
    setGlobals 2
    ../bin/peer lifecycle chaincode install ../asset-transfer-basic/chaincode-gm/$PKGNAME
   ```

###8、设置环境变量

  安装合约时，控制台会输出 合约的 package id ，设置环境变零 PKGID

   ```bash
    export PKGID=gm_1.0:5d46d431e3bee1af584bc9c788aec4267597fc6c39f48bee23400c8860cee793
   ```

  设置环境变量 CONTRACT_NAME 作为合约的名称

   ```bash
    export CONTRACT_NAME=gm
   ```

###9、部署合约

   ```bash
    setGlobals 2
    peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --channelID mychannel --name $CONTRACT_NAME --version 1.0 --package-id $PKGID --sequence 1
    
    setGlobals 1
    peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --channelID mychannel --name $CONTRACT_NAME --version 1.0 --package-id $PKGID --sequence 1
    
    peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --channelID mychannel --name $CONTRACT_NAME --peerAddresses localhost:7051 --tlsRootCertFiles $PWD/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt --version 1.0 --sequence 1
   ```

###10、启动合约服务（两种情况）
####（1）编译器本地启动合约服务（便于调试）
- 修改合约 main 函数中的环境变量

    ```bash
    os.Setenv("CHAINCODE_ID","gm_1.0:5d46d431e3bee1af584bc9c788aec4267597fc6c39f48bee23400c8860cee793")
    os.Setenv("CHAINCODE_SERVER_ADDRESS","192.168.2.150:9999")
    ```

- 启动服务
####（2）docker 启动合约服务
- 修改chaincode.env

    ```
    # CHAINCODE_SERVER_ADDRESS must be set to the host and port where the peer can
    # connect to the chaincode server
    CHAINCODE_SERVER_ADDRESS=gm:9999
      
    # CHAINCODE_ID must be set to the Package ID that is assigned to the chaincode
    # on install. The `peer lifecycle chaincode queryinstalled` command can be
    # used to get the ID after install if required
    CHAINCODE_ID=gm_1.0:5d46d431e3bee1af584bc9c788aec4267597fc6c39f48bee23400c8860cee793
    ```

- 编译镜像

    ```bash
    cd fabric-samples/asset-transfer-basic/chaincode-gm
    docker buildx build --platform linux/amd64 -t gmchaincode .
    ```

- 启动镜像

    ```bash
    docker run -it --rm --name gm --hostname gm --env-file chaincode.env --network=docker_test gmchaincode
    ```


###11、合约调试

   ```bash
    peer chaincode invoke -n $CONTRACT_NAME -c '{"Args":["InitLedger"]}' -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel
   ```
###12、环境清理

   ```bash
    bash network.sh down
    echo y｜docker system prune
    docker volume rm $(docker volume ls -qf dangling=true)
   ```