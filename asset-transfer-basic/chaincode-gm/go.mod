module github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-gm

go 1.14

replace (
	github.com/hyperledger/fabric-chaincode-go => ./fabric-chaincode-go
	github.com/tjfoc/gmsm => ./tjfoc/gmsm
	github.com/tjfoc/gmtls => ./tjfoc/gmtls
)

require (
	github.com/golang/protobuf v1.5.2
	github.com/hyperledger/fabric-chaincode-go v0.0.0-20220720122508-9207360bbddd
	github.com/hyperledger/fabric-contract-api-go v1.2.0
	github.com/hyperledger/fabric-protos-go v0.0.0-20220613214546-bf864f01d75e
	github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go v0.0.0-20221121194459-22e1af493534
	github.com/stretchr/testify v1.8.0
	github.com/tjfoc/gmtls v0.0.0-00010101000000-000000000000 // indirect
)
