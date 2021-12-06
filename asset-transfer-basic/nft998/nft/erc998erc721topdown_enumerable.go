package nft

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

/**
 * @Author: fengxiaoxiao /13156050650@163.com
 * @Desc:
 * @Version: 1.0.0
 * @Date: 2021/12/6 10:30 下午
 */
//
//type ERC998ERC721TopDownEnumerable interface {
//	TotalChildContracts(tokenId uint64) (uint64,error);
//	ChildContractByIndex(tokenId,index uint64) (childContractName string,err error)
//	TotalChildTokens( tokenId uint64, childContractName string) (uint64,error)
//	ChildTokenByIndex( tokenId uint64,childContractName string,index uint64)  (childTokenId uint64,err error)
//}

/*
 * @Desc:
 * @Param:
 * @Return:
 */
func (nft *NFT998) TotalChildContracts(ctx contractapi.TransactionContextInterface, tokenId uint64) (uint64, error) {
	return 0, nil
}

/*
 * @Desc:
 * @Param:
 * @Return:
 */
func (nft *NFT998) TotalChildTokens(ctx contractapi.TransactionContextInterface, childContractNames string) (uint64, error) {
	return 0, nil
}

/*
 * @Desc:
 * @Param:
 * @Return:
 */
func (nft *NFT998) ChildContractByIndex(ctx contractapi.TransactionContextInterface, index uint64) (uint64, error) {
	return 0, nil
}

/*
 * @Desc:
 * @Param:
 * @Return:
 */
func (nft *NFT998) ChildTokenByIndex(ctx contractapi.TransactionContextInterface, childContractName string, index uint64) (uint64, error) {
	return 0, nil
}
