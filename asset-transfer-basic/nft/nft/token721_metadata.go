package nft

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// NFT base metadata
const (
	NFTName   = "BENFT"
	NFTSymbol = "$B$"
)

// Name 获取 NFT 代币名称
func (nft *NFT) Name(ctx contractapi.TransactionContextInterface) string {
	return NFTName
}

// Symbol 获取 NFT 代币符号
func (nft *NFT) Symbol(ctx contractapi.TransactionContextInterface) string {
	return NFTSymbol
}

// TokenURI 根据 tokenID 获取 token 的元数据 URI
func (nft *NFT) TokenURI(ctx contractapi.TransactionContextInterface, tokenID uint64) (string, error) {
	key, err := GetTokenURIKey(ctx, tokenID)
	if err != nil {
		return "", err
	}
	res, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// TokenLabel 根据 tokenID 获取 token 的元数据 label
func (nft *NFT) TokenLabel(ctx contractapi.TransactionContextInterface, tokenID uint64) (string, error) {
	key, err := GetTokenLabelKey(ctx, tokenID)
	if err != nil {
		return "", err
	}
	state, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", err
	}
	return string(state), nil
}
