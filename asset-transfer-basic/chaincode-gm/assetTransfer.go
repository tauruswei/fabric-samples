/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"
	"os"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-gm/chaincode"
)

//func main() {
//	assetChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
//	if err != nil {
//		log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
//	}
//
//	if err := assetChaincode.Start(); err != nil {
//		log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
//	}
//}

type serverConfig struct {
	CCID    string
	Address string
}

func main() {
	// See chaincode.env.example

	os.Setenv("CHAINCODE_ID", "gm_1.0:5d46d431e3bee1af584bc9c788aec4267597fc6c39f48bee23400c8860cee793")
	os.Setenv("CHAINCODE_SERVER_ADDRESS", "192.168.2.150:9999")

	config := serverConfig{
		CCID:    os.Getenv("CHAINCODE_ID"),
		Address: os.Getenv("CHAINCODE_SERVER_ADDRESS"),
	}

	chaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})

	if err != nil {
		log.Panicf("error create asset-transfer-basic chaincode: %s", err)
	}

	server := &shim.ChaincodeServer{
		CCID:    config.CCID,
		Address: config.Address,
		CC:      chaincode,
		TLSProps: shim.TLSProperties{
			Disabled: true,
		},
	}

	log.Println("success")

	if err := server.Start(); err != nil {
		log.Panicf("error starting asset-transfer-basic chaincode: %s", err)
	}
}
