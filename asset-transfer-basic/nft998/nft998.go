package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"log"
	"os"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	"benft/nft"
)

//func main() {
//	cc, err := contractapi.NewChaincode(new(nft.NFT))
//	if err != nil {
//		panic(err.Error())
//	}
//	if err := cc.Start(); err != nil {
//		fmt.Printf("Error starting NFT chaincode: %s", err)
//	}
//}

type serverConfig struct {
	CCID    string
	Address string
}

func main() {
	// See chaincode.env.example
	config := serverConfig{
		CCID:    os.Getenv("CHAINCODE_ID"),
		Address: os.Getenv("CHAINCODE_SERVER_ADDRESS"),
	}

	chaincode, err := contractapi.NewChaincode(new(nft.NFT))

	if err != nil {
		log.Panicf("error create nft721 chaincode: %s", err)
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
		log.Panicf("error starting nft721 chaincode: %s", err)
	}
}
