module github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go

go 1.14

replace (
	github.com/hyperledger/fabric-chaincode-go => ./fabric-chaincode-go
	github.com/tjfoc/gmsm => ./tjfoc/gmsm
	github.com/tjfoc/gmtls => ./tjfoc/gmtls
)

require (
	github.com/golang/protobuf v1.3.2
	github.com/hyperledger/fabric-chaincode-go v0.0.0-20210718160520-38d29fabecb9
	github.com/hyperledger/fabric-contract-api-go v1.1.1
	github.com/hyperledger/fabric-protos-go v0.0.0-20201028172056-a3136dde2354
	github.com/stretchr/testify v1.5.1
	github.com/tjfoc/gmtls v0.0.0-00010101000000-000000000000 // indirect
)
