module github.com/hyperledger/fabric-chaincode-go

go 1.14

replace (
	github.com/tjfoc/gmsm => ../tjfoc/gmsm
	github.com/tjfoc/gmtls => ../tjfoc/gmtls
)

require (
	github.com/golang/protobuf v1.5.2
	github.com/hyperledger/fabric-protos-go v0.0.0-20220315113721-7dc293e117f7
	github.com/tjfoc/gmsm v1.2.0
	github.com/tjfoc/gmtls v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.26.0
)
