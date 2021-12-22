package nft

import (
	"fmt"
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

func (nft *NFT) RootOwnerOf(ctx contractapi.TransactionContextInterface, tokenId uint64) (string, error) {
	return nft.RootOwnerOfChild(ctx, "", tokenId)
}

func (nft *NFT) RootOwnerOfChild(ctx contractapi.TransactionContextInterface, childContractName string, childTokenId uint64) (rootOwnerAddress string, err error) {
	if childContractName != "" {
		rootOwnerAddress, childTokenId, err = OwnerOfChild(ctx, childContractName, childTokenId)
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
		rootOwnerAddress, childTokenId, err = OwnerOfChild(ctx, rootOwnerAddress, childTokenId)
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

	// 获取 childtoken 在 contract name 中的 index
	//childTokenIndex, err := ChildTokenIndex(ctx, tokenId, childContractName, childTokenId)
	//if err != nil {
	//	return err
	//}
	//if childTokenIndex == 0 {
	//	return fmt.Errorf("child token not owned by token, childTokenid = %d, childContractName = %s, parentTokenId = %d", childTokenId, childContractName, tokenId)
	//}
	// remove child token
	//childTokenIndexKey, err := GetChildTokenIndexKey(ctx, tokenId, childContractName, childTokenId)
	//if err != nil {
	//	return err
	//}
	//err = ctx.GetStub().DelState(childTokenIndexKey)
	//if err != nil {
	//	return err
	//}
	childTokenOwnerKey, err := GetChildTokenOwnerKey(ctx, childContractName, childTokenId)
	if err != nil {
		return err
	}
	err = ctx.GetStub().DelState(childTokenOwnerKey)
	if err != nil {
		return err
	}
	//GetChildTokensKey(ctx,tokenId,childContractName)
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
		//childContractIndexKey, err := GetChildContractIndexKey(ctx, tokenId, childContractName)
		//if err != nil {
		//	return err
		//}
		//err = ctx.GetStub().DelState(childContractIndexKey)
		//if err != nil {
		//	return err
		//}
	}
	return nil
}

/*
 * @Desc: 将 child token 转移给某个地址
 * @Param:
 * @Return:
 */
func (nft *NFT) TransferChild(ctx contractapi.TransactionContextInterface, fromTokenId uint64, to, childContractName string, childTokenId uint64, data string) error {
	parentTokenId, err := ChildTokenOwner(ctx, childContractName, childTokenId)
	if err != nil {
		return err
	}
	if parentTokenId != fromTokenId {
		return fmt.Errorf("parent token does not own child token, parentTokenId = %d, childContractName = %s, childTokenId = %d", fromTokenId, childContractName, childTokenId)
	}
	//childTokenIndex, err := ChildTokenIndex(ctx, fromTokenId, childContractName, childTokenId)
	//if err != nil {
	//	return err
	//}
	//if childTokenIndex == 0 {
	//	return fmt.Errorf("child token index can not be 0, parentTokenId = %d, childContractName = %s, childTokenId = %d", fromTokenId, childContractName, childTokenId)
	//}
	//rootOwner, err := nft.RootOwnerOf(ctx, fromTokenId)
	//if err != nil {
	//	return err
	//}
	//sender, err := getSender(ctx)
	//if err != nil {
	//	return err
	//}
	//// todo
	//rootOwnerAndTokenIdToApprovedAddressKey, err := GetRootOwnerAndTokenIdToApprovedAddressKey(ctx, rootOwner, fromTokenId)
	//if err != nil {
	//	return err
	//}
	//rootOwnerAndTokenIdToApprovedAddress, err := ctx.GetStub().GetState(rootOwnerAndTokenIdToApprovedAddressKey)
	//if err != nil {
	//	return err
	//}
	//if rootOwner != sender && string(rootOwnerAndTokenIdToApprovedAddress) != sender {
	//	return fmt.Errorf("the sender does not have right to transfer child, sender = %s, parentTokenId = %d, childContractName = %s, childTokenId = %d", sender, fromTokenId, childContractName, childTokenId)
	//}

	//response := ctx.GetStub().InvokeChaincode(childContractName, util.ToChaincodeArgs("ReceiveFromNft998", to, strconv.FormatUint(childTokenId, 10), data), ctx.GetStub().GetChannelID())
	//if response.Status != 200 {
	//	return fmt.Errorf("nft998 发送到 nft721 失败，msg = %s, tokenId = %d, childContractName = %s, childTokenId = %d", response.Message, fromTokenId, childContractName, childTokenId)
	//}

	if childContractName == "" {
		key, err := GetTokenLabelKey(ctx, childTokenId)
		if err != nil {
			return err
		}
		state, err := ctx.GetStub().GetState(key)
		if err != nil {
			return err
		}
		childContractName = string(state)
	}
	err = nft.ReceiveFromNft998(ctx, to, childTokenId, data)
	if err != nil {
		return err
	}
	err = RemoveChild(ctx, fromTokenId, childContractName, childTokenId)
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
func (nft *NFT) SafeTransferChild(ctx contractapi.TransactionContextInterface, fromTokenId uint64, to, childContractName string, childTokenId uint64, data string) error {
	return nft.TransferChild(ctx, fromTokenId, to, childContractName, childTokenId, data)
}

/*
 * @Desc:
 * @Param:
 * @Return:
 */
//function transferChildToParent(uint256 _fromTokenId, address _toContract, uint256 _toTokenId, address _childContract, uint256 _childTokenId, bytes _data) external {
func (nft *NFT) TransferChildToParent(ctx contractapi.TransactionContextInterface, fromTokenId uint64, toContractName string, toTokenId uint64, childContractName string, childTokenId uint64, data string) {
	fmt.Println("not implement")
}

/*
 * @Desc:
 * @Param:
 * @Return:
 */
//function getChild(address _from, uint256 _tokenId, address _childContract, uint256 _childTokenId) external {

func (nft *NFT) GetChild(ctx contractapi.TransactionContextInterface, from string, tokenId uint64, childContractName string, childTokenId uint64) {
	fmt.Println("not implement")
}

/*
 * @Desc:
 * @Param:
 * @Return:
 */
//function onERC721Received(address _from, uint256 _childTokenId, bytes _data) external returns (bytes4) {

func (nft *NFT) OnERC721Received(ctx contractapi.TransactionContextInterface, childTokenId, tokenId uint64) {
	fmt.Println("not implement")
}

/*
 * @Desc:
 * @Param:
 * @Return:
 */
func (nft *NFT) ReceiveNft721(ctx contractapi.TransactionContextInterface, tokenId uint64, childContractName string, childTokenId uint64) error {

	_, err := TokenToOwner(ctx, tokenId)
	if err != nil {
		return err
	}
	//index, err := ChildTokenIndex(ctx, tokenId, childContractName, childTokenId)
	//if index!=0{
	//	return fmt.Errorf("Cannot send child token because it has already been received, tokenId = %d, childContractName = %s, childTokenId=%d",tokenId,childContractName,childTokenId)
	//}
	//childTokenIndexKey, err := GetChildTokenIndexKey(ctx, tokenId, childContractName, childTokenId)
	//if err != nil {
	//	return err
	//}
	//err = ctx.GetStub().PutState(childTokenIndexKey, []byte(fmt.Sprintf("%s", childTokenId)))
	if err != nil {
		return err
	}
	childTokensKey, err := GetChildTokensKey(ctx, tokenId, childContractName)
	if err != nil {
		return err
	}
	raw, err := ctx.GetStub().GetState(childTokensKey)
	if err != nil {
		return err
	}
	if raw != nil {
		return fmt.Errorf("childTokenId = %d already owned by tokenId = %d ", childTokenId, tokenId)
	}

	err = ctx.GetStub().PutState(childTokensKey, []byte(fmt.Sprintf("%d", childTokenId)))
	if err != nil {
		return err
	}
	childTokenOwnerKey, err := GetChildTokenOwnerKey(ctx, childContractName, childTokenId)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(childTokenOwnerKey, []byte(fmt.Sprintf("%d", tokenId)))
	if err != nil {
		return err
	}

	childContractsKey, err := GetChildContractsKey(ctx, tokenId, childContractName)
	if err != nil {
		return err
	}
	childContracts, err := ctx.GetStub().GetState(childContractsKey)
	if err != nil {
		return err
	}
	if childContracts != nil {
		//childContractsKey, err := GetChildContractsKey(ctx, tokenId, childContractName)
		//if err != nil {
		//	return err
		//}
		//err = ctx.GetStub().PutState(childContractsKey, []byte(childContractName))
		//if err != nil {
		//	return err
		//}
		return fmt.Errorf("childContract = %s already owned by tokenId = %d ", childContractName, tokenId)
	}

	err = ctx.GetStub().PutState(childContractsKey, []byte(childContractName))
	if err != nil {
		return err
	}
	return nil
}

func OwnerOfChild(ctx contractapi.TransactionContextInterface, childContractName string, childTokenId uint64) (parentTokenOwner string, parentTokenId uint64, err error) {
	//var childIndexKey string
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
		//childIndexKey, err = GetChildTokenIndexKey(ctx, parentTokenId, childContractName, childTokenId)
		//if err != nil {
		//	return "", 0, err
		//}
	}
	//childIndex, err := ctx.GetStub().GetState(childIndexKey)

	//if string(parentTokenIdBytes) == "" || string(childIndex) == "" {
	if string(parentTokenIdBytes) == "" {
		return "", 0, fmt.Errorf("child token does not have a parent token, childContractName = %s ,childTokenId = %d", childContractName, childTokenId)
	}
	parentTokenOwner, err = TokenToOwner(ctx, parentTokenId)
	if err != nil {
		return "", 0, err
	}

	return parentTokenOwner, parentTokenId, nil
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
