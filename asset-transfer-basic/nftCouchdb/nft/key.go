package nft

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// state db key prefix
const (
	KeyPrefixMessage         = "KEYMessage"    // message key name
	KeyPrefixNFT             = "KEYToken"      // token name
	KeyPrefixContract        = "KEYContract"   // contract
	KeyPrefixNFTName         = "KEYTokenName"  // token name
	KeyPrefixNFTLabel        = "KEYTokenLabel" // token label
	KeyPrefixNFTURI          = "KEYTokenURI"   // token uri
	KeyPrefixNFTDesc         = "KEYTokenDesc"  // token desc
	KeyPrefixNFTCount        = "KEYTokenCount"
	KeyPrefixNFTCountByOwner = "KEYTokenCountOfOwner"
	KeyPrefixNFTOwner        = "KEYTokenOwner"
	//KeyPrefixNFTOnSale       = "KEYTokenOnSale"
	//KeyPrefixNFTOnSaleCount  = "KEYTokenOnSaleCount"
	KeyPrefixNFTOwnerHistory = "KEYTokenOwnerHistory"
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

	//
	KeyPrefixNFTDelegateToken = "KeyDelegateToken"

	// dtokenID => tokenId
	KeyPrefixNFTDelegateTokenParent = "KeyDelegateTokenParent"

	// dtokenID => height
	KeyPrefixNFTDelegateTokenHeight = "KeyDelegateTokenHeight"

	// dtokenID => expiration
	KeyPrefixNFTDelegateTokenExpiration = "KeyDelegateTokenExpiration"
)

/*
 * @Desc: 获取 message key
 * @Param:
 * @Return:
 */
func GetMessageKey(ctx contractapi.TransactionContextInterface, messageId string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixMessage, []string{messageId})
	if err != nil {
		return "", err
	}
	return key, nil
}

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
 * @Desc: get the key of TokenId
 * @Param: tokenId
 * @Return:
 */
func GetTokenIdKey(ctx contractapi.TransactionContextInterface, tokenId string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFT, []string{fmt.Sprintf("%s", tokenId)})
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
func GetChildTokensKey(ctx contractapi.TransactionContextInterface, tokenId uint64, childContractName string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTChildTokens, []string{fmt.Sprintf("%d", tokenId), childContractName})
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

// GetTokenNameKey .
func GetTokenNameKey(ctx contractapi.TransactionContextInterface, tokenID uint64) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTName, []string{fmt.Sprintf("%d", tokenID)})
	if err != nil {
		return "", err
	}
	return key, nil
}

// GetTokenLabelKey .
func GetTokenLabelKey(ctx contractapi.TransactionContextInterface, tokenID uint64) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTLabel, []string{fmt.Sprintf("%d", tokenID)})
	if err != nil {
		return "", err
	}
	return key, nil
}

// GetTokenURIKey .
func GetTokenURIKey(ctx contractapi.TransactionContextInterface, tokenID uint64) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTURI, []string{fmt.Sprintf("%d", tokenID)})
	if err != nil {
		return "", err
	}
	return key, nil
}

// GetTokenDescKey .
func GetTokenDescKey(ctx contractapi.TransactionContextInterface, tokenID uint64) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTDesc, []string{fmt.Sprintf("%d", tokenID)})
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

// GetTokenApprovedKey .
func GetTokenApprovedKeyCouchdb(ctx contractapi.TransactionContextInterface, tokenID string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTApprove, []string{fmt.Sprintf("%s", tokenID)})
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

// GetDtokenParentKey .
func GetDtokenParentKey(ctx contractapi.TransactionContextInterface, dtokenID uint64) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTDelegateTokenParent, []string{fmt.Sprintf("%d", dtokenID)})
	if err != nil {
		return "", err
	}
	return key, nil
}

// GetDtokenHeightKey .
func GetDtokenHeightKey(ctx contractapi.TransactionContextInterface, dtokenID uint64) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTDelegateTokenHeight, []string{fmt.Sprintf("%d", dtokenID)})
	if err != nil {
		return "", err
	}
	return key, nil
}

// GetDtokenExpirationKey .
func GetDtokenExpirationKey(ctx contractapi.TransactionContextInterface, dtokenID uint64) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTDelegateTokenExpiration, []string{fmt.Sprintf("%d", dtokenID)})
	if err != nil {
		return "", err
	}
	return key, nil
}

// GetTokenOwnerHistoryKey .
func GetOwnerTokensHistoryKey(ctx contractapi.TransactionContextInterface, owner string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTOwnerHistory, []string{fmt.Sprintf("%s", owner)})
	if err != nil {
		return "", err
	}
	return key, nil
}

// GetDtokenKey .
func GetDtokenKey(ctx contractapi.TransactionContextInterface, dtokenID string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTDelegateToken, []string{fmt.Sprintf("%s", dtokenID)})
	if err != nil {
		return "", err
	}
	return key, nil
}

// GetContractKey .
func GetContractKey(ctx contractapi.TransactionContextInterface, tokenId string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(KeyPrefixContract, []string{fmt.Sprintf("%s", tokenId)})
	if err != nil {
		return "", err
	}
	return key, nil
}
