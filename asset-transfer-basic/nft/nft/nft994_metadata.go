package nft

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// DTokenHeight 根据 dtokenID 获取 dtoken 的 height
func (nft *NFT) DTokenHeight(ctx contractapi.TransactionContextInterface, dtokenID uint64) (string, error) {
	key, err := GetDtokenHeightKey(ctx, dtokenID)
	if err != nil {
		return "", err
	}
	res, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// DTokenHeight 根据 dtokenID 获取 dtoken 的 expiration
func (nft *NFT) DTokenExpiration(ctx contractapi.TransactionContextInterface, dtokenID uint64) (string, error) {
	key, err := GetDtokenExpirationKey(ctx, dtokenID)
	if err != nil {
		return "", err
	}
	res, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
