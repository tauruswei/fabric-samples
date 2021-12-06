package main

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

/**
 * @Author: WeiBingtao/13156050650@163.com
 * @Version: 1.0
 * @Description:
 * @Date: 2021/11/1 下午9:10
 */
func main() {
	simpleContract := new(SimpleContract)

	cc, err := contractapi.NewChaincode(simpleContract)

	if err != nil {
		panic(err.Error())
	}

	if err := cc.Start(); err != nil {
		panic(err.Error())
	}
}
