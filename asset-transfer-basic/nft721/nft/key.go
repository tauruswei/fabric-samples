package nft

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// state db key prefix
const (
	ContractName              = "nft7"
	KeyPrefixNFTURI          = "KEYTokenURI" // token uri
	KeyPrefixNFTCount        = "KEYTokenCount"
	KeyPrefixNFTCountByOwner = "KEYTokenCountOfOwner"
	KeyPrefixNFTOwner        = "KEYTokenOwner"
	KeyPrefixNFTApprove      = "KEYTokenApprove"
	KeyPrefixNFTApproveAll   = "KEYTokenApproveAll"

	// tokenId => token owner
	KeyPrefixNFTTokenIdToTokenOwner = "KEYTokenIdToTokenOwner"

	// tokenId => []child contract
	KeyPrefixNFTChildContracts = "KEYChildContracts"

	// tokenId => (child address => contract index)
	KeyPrefixNFTChildContractIndex = "KEYChildContractIndex"

	// tokenId => (child address => [] childtoken)
	KeyPrefixNFTChildTokens = "KEYChildTokens"

	// child address => childId => tokenId
	KeyPrefixNFTChildTokenOwner = "KEYChildTokenOwner"

	// tokenId => (child address => (child token => child index)
	KeyPrefixNFTChildTokenIndex = "KEYChildTokenIndex"

	// root token owner address => (tokenId => approved address)
	//mapping(address => mapping(uint256 => address))
	KeyPrefixNFTRootOwnerAndTokenIdToApprovedAddress = "KEYRootOwnerAndTokenIdToApprovedAddress"
)

/*
 * @Desc: get the key of RootOwnerAndTokenIdToApprovedAddress
 * @Param: rootOwner  tokenId
 * @Return:
 */
func GetRootOwnerAndTokenIdToApprovedAddressKey(ctx contractapi.TransactionContextInterface, rootTokenOwner string, tokenId uint64) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTRootOwnerAndTokenIdToApprovedAddress, []string{rootTokenOwner, fmt.Sprintf("%d", tokenId)})
	if err != nil {
		return "", err
	}
	return key, nil
}

/*
 * @Desc: get the key of TokenIdToTokenOwner
 * @Param: tokenId
 * @Return:
 */
func GetTokenIdToTokenOwnerKey(ctx contractapi.TransactionContextInterface, tokenId uint64) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTTokenIdToTokenOwner, []string{fmt.Sprintf("%d", tokenId)})
	if err != nil {
		return "", err
	}
	return key, nil
}

/*
 * @Desc:
 * @Param:
 * @Return:
 */
func GetChildContractIndexKey(ctx contractapi.TransactionContextInterface, tokenid uint64, childcontractName string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTChildContractIndex, []string{fmt.Sprintf("%d", tokenid), childcontractName})
	if err != nil {
		return "", err
	}
	return key, nil
}

/*
 * @Desc:
 * @Param:
 * @Return:
 */
func GetChildContractsKey(ctx contractapi.TransactionContextInterface, tokenId uint64, childContractName string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTChildContracts, []string{fmt.Sprintf("%d", tokenId), childContractName})
	if err != nil {
		return "", err
	}
	return key, nil
}

/*
 * @Desc: get the key of childTokens
 * @Param:
 * @Return:
 */
func GetChildTokensKey(ctx contractapi.TransactionContextInterface, tokenId uint64, childContractName string, childTokenId uint64) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTChildTokens, []string{fmt.Sprintf("%d", tokenId), childContractName, fmt.Sprintf("%d", childTokenId)})
	if err != nil {
		return "", err
	}
	return key, nil
}

/*
 * @Desc: get the key of ChildTokenOwner
 * @Param:
 * @Return:
 */
func GetChildTokenOwnerKey(ctx contractapi.TransactionContextInterface, childContractName string, childTokenId uint64) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTChildTokenOwner, []string{childContractName, fmt.Sprintf("%d", childTokenId)})
	if err != nil {
		return "", err
	}
	return key, err
}

/*
 * @Desc: get the key of ChildTokenIndex
 * @Param:
 * @Return:
 */
func GetChildTokenIndexKey(ctx contractapi.TransactionContextInterface, tokenId uint64, childContractName string, childTokenId uint64) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTChildTokenIndex, []string{fmt.Sprintf("%d", tokenId), childContractName, fmt.Sprintf("%d", childTokenId)})
	if err != nil {
		return "", err
	}
	return key, err
}

// GetTokenURIKey .
func GetTokenURIKey(ctx contractapi.TransactionContextInterface, tokenID uint64) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTURI, []string{fmt.Sprintf("%d", tokenID)})
	if err != nil {
		return "", err
	}
	return key, nil
}

// GetTokenCountKey .
func GetTokenCountKey(ctx contractapi.TransactionContextInterface) (string, error) {
	return KeyPrefixNFTCount, nil
}

// GetTokenOwnerKey .
func GetTokenOwnerKey(ctx contractapi.TransactionContextInterface, tokenID uint64) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTTokenIdToTokenOwner, []string{fmt.Sprintf("%d", tokenID)})
	if err != nil {
		return "", err
	}
	return key, nil
}

// GetTokenApprovedKey .
func GetTokenApprovedKey(ctx contractapi.TransactionContextInterface, tokenID uint64) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTApprove, []string{fmt.Sprintf("%d", tokenID)})
	if err != nil {
		return "", err
	}
	return key, nil
}

// GetApprovedAllKey ...
func GetApprovedAllKey(ctx contractapi.TransactionContextInterface, owner, operator string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTApproveAll, []string{owner, operator})
	if err != nil {
		return "", err
	}
	return key, nil
}

// GetTokenCountByOwnerKey .
func GetTokenCountByOwnerKey(ctx contractapi.TransactionContextInterface, owner string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTCountByOwner, []string{owner})
	if err != nil {
		return "", err
	}
	return key, nil
}
