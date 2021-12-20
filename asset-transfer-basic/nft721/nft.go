package main

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"nft721/nft"
)

func main() {
	cc, err := contractapi.NewChaincode(new(nft.NFT721))
	if err != nil {
		panic(err.Error())
	}
	if err := cc.Start(); err != nil {
		fmt.Printf("Error starting NFT chaincode: %s", err)
	}
	fmt.Println()
}

//type serverConfig struct {
//	CCID    string
//	Address string
//}
//
//func main() {
//	// See chaincode.env.example
//	config := serverConfig{
//		CCID:    os.Getenv("CHAINCODE_ID"),
//		Address: os.Getenv("CHAINCODE_SERVER_ADDRESS"),
//	}
//
//	chaincode, err := contractapi.NewChaincode(new(nft.NFT721))
//
//	if err != nil {
//		log.Panicf("error create nft7 chaincode: %s", err)
//	}
//
//	server := &shim.ChaincodeServer{
//		CCID:    config.CCID,
//		Address: config.Address,
//		CC:      chaincode,
//		TLSProps: shim.TLSProperties{
//			Disabled: true,
//		},
//	}
//
//	log.Println("success")
//
//	if err := server.Start(); err != nil {
//		log.Panicf("error starting nft7 chaincode: %s", err)
//	}
//}
