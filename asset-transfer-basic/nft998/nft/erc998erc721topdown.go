package nft

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strconv"
)

/**
 * @Author: fengxiaoxiao /13156050650@163.com
 * @Desc: erc998erc721topdown 接口实现
 * @Version: 1.0.0
 * @Date: 2021/12/3 3:30 下午
 */
//type ERC998ERC721TopDown interface{
//	RootOwnerOf(ctx contractapi.TransactionContext,tokenId uint64)(string,error)
//	RootOwnerOfChild(ctx contractapi.TransactionContext,childContractName string,childTokenId uint64)(string,error)
//	OwnerOfChild(ctx contractapi.TransactionContext, childContractName string,childTokenId uint64)(string,uint64,error)
//	OnERC721Received(ctx contractapi.TransactionContext,operator,from string,childTokenId uint64,data string)(string,error)
//	TransferChild(ctx contractapi.TransactionContext,fromTokenId uint64,to,childContractName string,childTokenId uint64)(error)
//	//SafeTransferChild(ctx contractapi.TransactionContext,fromTokenId uint64,to,childContractName string,childTokenId uint64)(error)
//	SafeTransferChild(ctx contractapi.TransactionContext,fromTokenId uint64,to,childContractName string,childTokenId uint64,data string)(error)
//	TransferChildToParent(ctx contractapi.TransactionContext,fromTokenId uint64,to,childContractName string,childTokenId uint64,data string)(error)
//	GetChild( ctx contractapi.TransactionContext,from string,tokenId uint64, childContractName string,childTokenId uint64)(error)
//}

type NFT998 struct {
	contractapi.Contract
}

func (nft *NFT998) RootOwnerOf(ctx contractapi.TransactionContextInterface, tokenId uint64) (string, error) {
	return nft.RootOwnerOfChild(ctx, "", tokenId)
}
func (nft *NFT998) RootOwnerOfChild(ctx contractapi.TransactionContextInterface, childContractName string, childTokenId uint64) (rootOwnerAddress string, err error) {
	if childContractName != "" {
		rootOwnerAddress, childTokenId, err = nft.OwnerOfChild(ctx, childContractName, childTokenId)
		if err != nil {
			return "", err
		}
	} else {
		return TokenToOwner(ctx, childTokenId)
	}
	for {
		if rootOwnerAddress != "" {
			break
		}
		rootOwnerAddress, childTokenId, err = nft.OwnerOfChild(ctx, rootOwnerAddress, childTokenId)
		if err != nil {
			return "", err
		}
	}
	return
}

/*
 * @Desc: 删除 child token
 * @Param:
 * @Return:
 */
func RemoveChild(ctx contractapi.TransactionContextInterface, tokenId uint64, childContractName string, childTokenId uint64) error {

	childTokenIndex, err := ChildTokenIndex(ctx, tokenId, childContractName, childTokenId)
	if err != nil {
		return err
	}
	if childTokenIndex == 0 {
		return fmt.Errorf("child token not owned by token, childTokenid = %d, childContractName = %s, parentTokenId = %d", childTokenId, childContractName, tokenId)
	}
	// remove child token
	childTokenIndexKey, err := GetChildTokenIndexKey(ctx, tokenId, childContractName, childTokenId)
	if err != nil {
		return err
	}
	err = ctx.GetStub().DelState(childTokenIndexKey)
	if err != nil {
		return err
	}
	childTokenOwnerKey, err := GetChildTokenOwnerKey(ctx, childContractName, childTokenId)
	if err != nil {
		return err
	}
	err = ctx.GetStub().DelState(childTokenOwnerKey)
	if err != nil {
		return err
	}
	childTokensKey, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTChildTokens, []string{fmt.Sprintf("%d", tokenId), childContractName, fmt.Sprintf("%d", childTokenId)})
	if err != nil {
		return err
	}
	err = ctx.GetStub().DelState(childTokensKey)
	if err != nil {
		return err
	}
	// remove child contract
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(KeyPrefixNFTChildTokens, []string{fmt.Sprintf("%d", tokenId), childContractName})
	defer iterator.Close()
	next := iterator.HasNext()
	if !next {
		childContractsKey, err := ctx.GetStub().CreateCompositeKey(KeyPrefixNFTChildContracts, []string{fmt.Sprintf("%d", tokenId), childContractName})
		if err != nil {
			return err
		}
		err = ctx.GetStub().DelState(childContractsKey)
		if err != nil {
			return err
		}
		childContractIndexKey, err := GetChildContractIndexKey(ctx, tokenId, childContractName)
		if err != nil {
			return err
		}
		err = ctx.GetStub().DelState(childContractIndexKey)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
 * @Desc:
 * @Param:
 * @Return:
 */
func (nft *NFT998) TransferChild(ctx contractapi.TransactionContextInterface, fromTokenId uint64, to, childContractName string, childTokenId uint64, data string) error {
	parentTokenId, err := ChildTokenOwner(ctx, childContractName, childTokenId)
	if err != nil {
		return err
	}
	if parentTokenId != fromTokenId {
		return fmt.Errorf("parent token does not own child token, parentTokenId = %d, childContractName = %s, childTokenId = %d", fromTokenId, childContractName, childTokenId)
	}
	childTokenIndex, err := ChildTokenIndex(ctx, fromTokenId, childContractName, childTokenId)
	if err != nil {
		return err
	}
	if childTokenIndex == 0 {
		return fmt.Errorf("child token index can not be 0, parentTokenId = %d, childContractName = %s, childTokenId = %d", fromTokenId, childContractName, childTokenId)
	}
	rootOwner, err := nft.RootOwnerOf(ctx, fromTokenId)
	if err != nil {
		return err
	}
	sender, err := getSender(ctx)
	if err != nil {
		return err
	}
	rootOwnerAndTokenIdToApprovedAddressKey, err := GetRootOwnerAndTokenIdToApprovedAddressKey(ctx, rootOwner, fromTokenId)
	if err != nil {
		return err
	}
	rootOwnerAndTokenIdToApprovedAddress, err := ctx.GetStub().GetState(rootOwnerAndTokenIdToApprovedAddressKey)
	if err != nil {
		return err
	}
	if rootOwner != sender && string(rootOwnerAndTokenIdToApprovedAddress) != sender {
		return fmt.Errorf("the sender does not have right to transfer child, sender = %s, parentTokenId = %d, childContractName = %s, childTokenId = %d", sender, fromTokenId, childContractName, childTokenId)
	}

	err = RemoveChild(ctx, fromTokenId, childContractName, childTokenId)
	if err != nil {
		return err
	}

	// todo 合约之间互相调用, 修改 nft 721   transfer  接口，接收  data string 参数

	err = nft.addToken(ctx, to, childTokenId)
	if err != nil {
		return err
	}
	err = nft.increaseToken(ctx, to)
	if err != nil {
		return err
	}

	return nil
}

/*
 * @Desc:
 * @Param:
 * @Return:
 */
func (nft *NFT998) SafeTransferChild(ctx contractapi.TransactionContextInterface, fromTokenId uint64, to, childContractName string, childTokenId uint64, data string) error {
	return nft.TransferChild(ctx, fromTokenId, to, childContractName, childTokenId, data)
}

/*
 * @Desc:
 * @Param:
 * @Return:
 */
//function transferChildToParent(uint256 _fromTokenId, address _toContract, uint256 _toTokenId, address _childContract, uint256 _childTokenId, bytes _data) external {
func (nft *NFT998) TransferChildToParent(ctx contractapi.TransactionContextInterface, fromTokenId uint64, toContractName string, toTokenId uint64, childContractName string, childTokenId uint64, data string) {
	fmt.Println("not implement")
}

/*
 * @Desc:
 * @Param:
 * @Return:
 */
//function getChild(address _from, uint256 _tokenId, address _childContract, uint256 _childTokenId) external {

func (nft *NFT998) GetChild(ctx contractapi.TransactionContextInterface, from string, tokenId uint64, childContractName string, childTokenId uint64) {
	fmt.Println("not implement")
}

/*
 * @Desc:
 * @Param:
 * @Return:
 */
//function onERC721Received(address _from, uint256 _childTokenId, bytes _data) external returns (bytes4) {

func (nft *NFT998) OnERC721Received(ctx contractapi.TransactionContextInterface, childTokenId, tokenId uint64) {
	fmt.Println("not implement")
}

/*
 * @Desc: get the parent token id of the specified child token
 * @Param:
 * @Return:
 */
func ChildTokenOwner(ctx contractapi.TransactionContextInterface, childcontractName string, childTokenId uint64) (uint64, error) {
	childTokenOwnerKey, err := GetChildTokenOwnerKey(ctx, childcontractName, childTokenId)
	if err != nil {
		return 0, err
	}
	tokenId, err := ctx.GetStub().GetState(childTokenOwnerKey)
	if err != nil {
		return 0, err
	}
	parentTokenId, err := strconv.ParseUint(string(tokenId), 10, 64)
	if err != nil {
		return 0, err
	}
	return parentTokenId, nil
}
func (nft *NFT998) OwnerOfChild(ctx contractapi.TransactionContextInterface, childContractName string, childTokenId uint64) (parentTokenOwner string, parentTokenId uint64, err error) {
	var childIndexKey string
	key, err := GetChildTokenOwnerKey(ctx, childContractName, childTokenId)
	if err != nil {
		return "", 0, err
	}
	parentTokenIdBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", 0, err
	}
	if string(parentTokenIdBytes) != "" {
		parentTokenId, err = strconv.ParseUint(string(parentTokenIdBytes), 10, 64)
		if err != nil {
			return "", 0, err
		}
		childIndexKey, err = GetChildTokenIndexKey(ctx, parentTokenId, childContractName, childTokenId)
		if err != nil {
			return "", 0, err
		}
	}
	childIndex, err := ctx.GetStub().GetState(childIndexKey)

	if string(parentTokenIdBytes) == "" || string(childIndex) == "" {
		return "", 0, fmt.Errorf("child token does not have a parent token, childContractName = %s ,childTokenId = %d", childContractName, childTokenId)
	}
	parentTokenOwner, err = TokenToOwner(ctx, parentTokenId)
	if err != nil {
		return "", 0, err
	}

	return parentTokenOwner, parentTokenId, nil
}
func TokenToOwner(ctx contractapi.TransactionContextInterface, tokenId uint64) (string, error) {
	parentTokenToOwnerKey, err := GetTokenIdToTokenOwnerKey(ctx, tokenId)
	if err != nil {
		return "", err
	}
	tokenOwnerAddress, err := ctx.GetStub().GetState(parentTokenToOwnerKey)
	if tokenOwnerAddress == nil {
		return "", fmt.Errorf("token does not hava a owner, tokenId = %d", tokenId)
	}
	return string(tokenOwnerAddress), nil
}

/*
 * @Desc:  获取 chilid token 在 child contract 中的 index
 * @Param:
 * @Return:
 */
func ChildTokenIndex(ctx contractapi.TransactionContextInterface, tokenId uint64, childContractName string, childTokenId uint64) (uint64, error) {
	childTokenIndexKey, err := GetChildTokenIndexKey(ctx, tokenId, childContractName, childTokenId)
	if err != nil {
		return 0, err
	}
	childTokenIndexBytes, err := ctx.GetStub().GetState(childTokenIndexKey)
	if err != nil {
		return 0, err
	}
	childTokenIndex, err := strconv.ParseUint(string(childTokenIndexBytes), 10, 64)
	return childTokenIndex, nil
}

/*
 * @Desc:  获取 chilid token 在 child contract 中的 index
 * @Param:
 * @Return:
 */
func getSender(ctx contractapi.TransactionContextInterface) (string, error) {
	cert, err := ctx.GetClientIdentity().GetX509Certificate()
	if err != nil {
		return "", err
	}
	if pubkey, ok := cert.PublicKey.(*ecdsa.PublicKey); ok {
		addr := crypto.PubkeyToAddress(*pubkey)
		return addr.String(), nil
	}
	return "", errors.New("not found")
}
