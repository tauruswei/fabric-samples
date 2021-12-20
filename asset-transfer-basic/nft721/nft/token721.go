package nft

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/common/util"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// NFT nft chaincode 参照 ERC1155
type NFT721 struct {
	contractapi.Contract
}

/*
 * @Desc:
 * @Param:
 * @Return:
 */
func (nft *NFT721) Init(ctx contractapi.TransactionContextInterface) {
	fmt.Println("init success")
}

// Mint 铸造 NFT
func (nft *NFT721) Mint(ctx contractapi.TransactionContextInterface, owner string, tokenID uint64, uri string) error {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(KeyPrefixNFTTokenIdToTokenOwner, []string{fmt.Sprintf("%d", tokenID)})
	defer iterator.Close()
	if iterator.HasNext() {
		return fmt.Errorf("tokenId = %d is already assigned", tokenID)
	}
	fmt.Printf("Mint token %d for %s ,uri: %s\n", tokenID, owner, uri)
	key, err := GetTokenURIKey(ctx, tokenID)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(key, []byte(uri))
	if err != nil {
		return err
	}
	// 将 nft token 和用户绑定
	err = nft.addToken(ctx, owner, tokenID)
	if err != nil {
		return err
	}
	// 修改 该用户的 token 数量
	err = nft.increaseToken(ctx, owner)
	if err != nil {
		return err
	}
	totalKey, _ := GetTokenCountKey(ctx)
	// 修改 token 总量
	return nft.calcCount(ctx, totalKey, true)
}

// Mint 铸造 DNFT
func (nft *NFT721) MintDNFT(ctx contractapi.TransactionContextInterface, tokenID,dtokenID,height uint64 ,expiration string) error {
	if !nft.canMintDnft(ctx, tokenID, expiration){
		return  fmt.Errorf("can not mint Dnft")
	}
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(KeyPrefixNFTDelegateToken, []string{fmt.Sprintf("%d-%d", tokenID,dtokenID)})
	if err != nil {
		return err
	}
	defer iterator.Close()
	if iterator.HasNext() {
		return fmt.Errorf("dtokenId = %d is already assigned", tokenID)
	}
	fmt.Printf("Mint dtoken %d for token %d \n", dtokenID, tokenID)
	heightKey, err := GetDtokenHeightKey(ctx, tokenID)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(heightKey, []byte(strconv.FormatUint(height,10)))
	if err != nil {
		return err
	}
	expirationKey, err := GetDtokenExpirationKey(ctx, tokenID)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(expirationKey, []byte(expiration))
	if err != nil {
		return err
	}
	// 将 delegated token 和 parent token 绑定
	return nft.addDToken(ctx, dtokenID, tokenID)
}

// Revoke 注销 DNFT
func (nft *NFT721) RevokeDNFT(ctx contractapi.TransactionContextInterface, dtokenID uint64) error {
	key, err := GetDtokenParentKey(ctx, dtokenID)
	if err != nil {
		return err
	}
	tokenIdBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return err
	}
	tokenId, err := strconv.ParseUint(string(tokenIdBytes), 10, 64)
	if err != nil {
		return err
	}
	owner, err := nft.OwnerOf(ctx,tokenId)
	if err != nil {
		return err
	}

	transfer := nft.canTransfer(ctx, owner, tokenId)

	if !transfer{
		return fmt.Errorf("can not revoke dtoken, dtokenId = %d",dtokenID)
	}

	err = ctx.GetStub().DelState(key)
	if err != nil {
		return err
	}
	heightKey, err := GetDtokenHeightKey(ctx, dtokenID)
	if err != nil {
		return err
	}
	err = ctx.GetStub().DelState(heightKey)
	if err != nil {
		return err
	}
	expirationKey, err := GetDtokenExpirationKey(ctx, dtokenID)
	if err != nil {
		return err
	}
	err = ctx.GetStub().DelState(expirationKey)
	return nil
}

// BalanceOf owner 的 NFT 数量
func (nft *NFT721) BalanceOf(ctx contractapi.TransactionContextInterface, owner string) (uint64, error) {
	key, err := GetTokenCountByOwnerKey(ctx, owner)
	if err != nil {
		return 0, err
	}
	res, err := ctx.GetStub().GetState(key)
	if err != nil {
		return 0, err
	}
	count, _ := strconv.Atoi(string(res))
	return uint64(count), nil
}

// OwnerOf 根据 tokenID 返回其所有人地址
func (nft *NFT721) OwnerOf(ctx contractapi.TransactionContextInterface, tokenID uint64) (string, error) {
	key, err := GetTokenOwnerKey(ctx, tokenID)
	fmt.Println(fmt.Sprintf("key = %s", key))
	if err != nil {
		return "", err
	}
	owner, err := ctx.GetStub().GetState(key)
	fmt.Println(fmt.Sprintf("owner = %s", owner))
	if err != nil {
		return "", err
	}
	return string(owner), nil
}

// SafeTransferFrom 根据 tokenID 将 NFT从 from 转移到 to
func (nft *NFT721) SafeTransferFrom(ctx contractapi.TransactionContextInterface, from string, to string, tokenID uint64, data []byte) error {
	if !nft.canTransfer(ctx, from, tokenID) {
		return errors.New("can not transfer")
	}
	return nft.TransferFrom(ctx, from, to, tokenID)
}

// TransferFrom 根据 tokenID 将 NFT从 from 转移到 to
func (nft *NFT721) TransferFrom(ctx contractapi.TransactionContextInterface, from string, to string, tokenID uint64) error {

	err := nft.delToken(ctx, from, tokenID)
	if err != nil {
		return err
	}
	err = nft.decreaseToken(ctx, from)
	if err != nil {
		return err
	}
	err = nft.addToken(ctx, to, tokenID)
	if err != nil {
		return err
	}
	err = nft.increaseToken(ctx, to)
	if err != nil {
		return err
	}
	return nil
}

// TransferFrom 根据 tokenID 将 NFT 从 998 转移到 721
func (nft *NFT721) ReceiveFromNft998(ctx contractapi.TransactionContextInterface, to string, tokenID uint64, data string) error {
	err := nft.addToken(ctx, to, tokenID)
	if err != nil {
		return err
	}
	err = nft.increaseToken(ctx, to)
	if err != nil {
		return err
	}
	return nil
}

// Approve 授予 approved 拥有 tokenID 的转移权力
func (nft *NFT721) Approve(ctx contractapi.TransactionContextInterface, approved string, tokenID uint64) error {
	key, err := GetTokenApprovedKey(ctx, tokenID)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(key, []byte(approved))
}

/*
 * @Desc: owner 的所有 NFT 都可以由 operator  来控制
 * @Param:
 * @Return:
 */
func (nft *NFT721) SetApprovalForAll(ctx contractapi.TransactionContextInterface, operator string, approved bool) error {
	sender, err := getSender(ctx)
	if err != nil {
		return err
	}
	key, err := GetApprovedAllKey(ctx, sender, operator)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(key, []byte(strconv.FormatBool(approved)))
}

// GetApproved 返回 tokenID 的授权地址
func (nft *NFT721) GetApproved(ctx contractapi.TransactionContextInterface, tokenID uint64) (string, error) {
	key, err := GetTokenApprovedKey(ctx, tokenID)
	if err != nil {
		return "", err
	}
	raw, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// IsApprovedForAll 查询 owner 的 NFT 转移权力是否授予 operator
func (nft *NFT721) IsApprovedForAll(ctx contractapi.TransactionContextInterface, owner string, operator string) (bool, error) {
	key, err := GetApprovedAllKey(ctx, owner, operator)
	if err != nil {
		return false, err
	}
	raw, err := ctx.GetStub().GetState(key)
	if err != nil {
		return false, err
	}
	if raw == nil {
		return false, nil
	}
	return strconv.ParseBool(string(raw))
}
// 判断是不是有操作 token 的权限
func (nft *NFT721) canTransfer(ctx contractapi.TransactionContextInterface, from string, tokenID uint64) bool {
	owner, err := nft.OwnerOf(ctx, tokenID)
	if err != nil {
		fmt.Printf("OwnerOf error: %s", err.Error())
		return false
	}
	if owner != from {
		fmt.Errorf("token does not owned by from, tokenId = %d, from= %s ", tokenID, from)
		return false
	}
	isOwner, err := checkSender(ctx, from)
	if err != nil {
		fmt.Printf("checkSender error: %s", err.Error())
		return false
	}
	// NFT721 拥有者
	if isOwner {
		return true
	}
	// NFT721 授权操作人
	approved, err := nft.GetApproved(ctx, tokenID)
	if err != nil {
		fmt.Printf("GetApproved error: %s", err.Error())
		return false
	}
	if approved == from {
		return true
	}

	sender, _ := getSender(ctx)
	// 所有资产授权操作人
	approvedAll, err := nft.IsApprovedForAll(ctx, owner, sender)
	if err != nil {
		fmt.Printf("IsApprovedForAll error: %s", err.Error())
		return false
	}
	if approvedAll {
		return true
	}
	return false
}
func (nft *NFT721) canMintDnft(ctx contractapi.TransactionContextInterface, tokenID uint64,expiration string) bool {
	// 首先进行权限的判断
	owner, err := nft.OwnerOf(ctx, tokenID)
	if err != nil {
		fmt.Printf("get owner of tokenId=%d error: %s", tokenID,err.Error())
		return false
	}

	transfer := nft.canTransfer(ctx, owner, tokenID)
	if !transfer{
		return transfer
	}

	// 其次判断 token 是不是 Dtoken,如果是 dtoken，要比较 expiration
	key, err := GetDtokenHeightKey(ctx, tokenID)
	if err != nil {
		fmt.Printf("get Dtoken Height Key error: %s", err.Error())
		return false
	}
	height, err := ctx.GetStub().GetState(key)
	if err != nil {
		fmt.Printf("get Dtoken Height error: %s", err.Error())
		return false
	}
	if height==nil{
		return true
	}
	key, err = GetDtokenExpirationKey(ctx, tokenID)
	if err != nil {
		fmt.Printf("get Dtoken Expiration Key error: %s", err.Error())
		return false
	}
	expirationBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		fmt.Printf("get Dtoken Expiration error: %s", err.Error())
		return false
	}
	//先把时间字符串格式化成相同的时间类型
	t1, err1 := time.Parse("2006-01-02 15:04:05", string(expirationBytes))
	t2, err2 := time.Parse("2006-01-02 15:04:05", expiration)
	if err1 == nil && err2 == nil && t1.Before(t2) {
		fmt.Printf("dtoken's expiration = %s is after parent token's expiration = %s, error: %s", string(expirationBytes),expiration,err.Error())
		return false
	}
	return true
}

func (nft *NFT721) delToken(ctx contractapi.TransactionContextInterface, from string, tokenID uint64) error {
	key, err := GetTokenOwnerKey(ctx, tokenID)
	if err != nil {
		return err
	}
	return ctx.GetStub().DelState(key)

}


func (nft *NFT721) delDtoken(ctx contractapi.TransactionContextInterface, from string, tokenID uint64) error {
	key, err := GetTokenOwnerKey(ctx, tokenID)
	if err != nil {
		return err
	}
	return ctx.GetStub().DelState(key)

}

func (nft *NFT721) increaseToken(ctx contractapi.TransactionContextInterface, to string) error {
	key, err := GetTokenCountByOwnerKey(ctx, to)
	if err != nil {
		return err
	}
	return nft.calcCount(ctx, key, true)
}

func (nft *NFT721) decreaseToken(ctx contractapi.TransactionContextInterface, from string) error {
	key, err := GetTokenCountByOwnerKey(ctx, from)
	if err != nil {
		return err
	}
	return nft.calcCount(ctx, key, false)
}

// 将 nft token 与用户绑定
func (nft *NFT721) addToken(ctx contractapi.TransactionContextInterface, to string, tokenID uint64) error {
	key, err := GetTokenOwnerKey(ctx, tokenID)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(key, []byte(to))
}
// 将 delegated token 和 parent token 绑定
func (nft *NFT721) addDToken(ctx contractapi.TransactionContextInterface, dtokenID, tokenID uint64) error {
	key, err := GetDtokenParentKey(ctx, dtokenID)
	if err != nil {
		return err
	}
	int64Str := strconv.FormatUint(tokenID, 10)
	return ctx.GetStub().PutState(key, []byte(int64Str))
}

// 修改 nft token 数量
func (nft *NFT721) calcCount(ctx contractapi.TransactionContextInterface, key string, increase bool) error {
	var calcF func(int) int
	calcFIncrease := func(count int) int {
		count++
		return count
	}
	calFDecrease := func(count int) int {
		count--
		return count
	}
	if increase {
		calcF = calcFIncrease
	} else {
		calcF = calFDecrease
	}
	res, err := ctx.GetStub().GetState(key)
	if err != nil {
		return err
	}
	if len(res) == 0 {
		res = []byte("0")
	}
	count, err := strconv.Atoi(string(res))
	if err != nil {
		return err
	}
	newCount := strconv.Itoa(calcF(count))
	return ctx.GetStub().PutState(key, []byte(newCount))
}

func checkSender(ctx contractapi.TransactionContextInterface, address string) (bool, error) {
	sender, err := getSender(ctx)
	if err != nil {
		return false, err
	}
	if sender == address {
		return true, nil
	}
	return false, nil
}

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

func (nft *NFT721) GetSender(ctx contractapi.TransactionContextInterface) (string, error) {
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

/*
 * @Desc:
 * @Param:
 * @Return:
 */
func (nft *NFT721) SendNft721ToNft998(ctx contractapi.TransactionContextInterface, from string, tokenId uint64, parentContractName, childContractName string, childTokenId uint64) error {
	transfer := nft.canTransfer(ctx, from, childTokenId)
	if transfer {
		//_, err := TokenToOwner(ctx, tokenId)
		//if err != nil {
		//	return err
		//}

		//index, err := ChildTokenIndex(ctx, tokenId, childContractName, childTokenId)
		//if index!=0{
		//	return fmt.Errorf("Cannot send child token because it has already been received, tokenId = %d, childContractName = %s, childTokenId=%d",tokenId,childContractName,childTokenId)
		//}
		//childTokenIndexKey, err := GetChildTokenIndexKey(ctx, tokenId, childContractName, childTokenId)
		//if err != nil {
		//	return err
		//}
		//err = ctx.GetStub().PutState(childTokenIndexKey, []byte(fmt.Sprintf("%s", childTokenId)))
		//if err != nil {
		//	return err
		//}
		//childTokensKey, err := GetChildTokensKey(ctx, tokenId, childContractName, childTokenId)
		//if err != nil {
		//	return err
		//}
		//raw, err := ctx.GetStub().GetState(childTokensKey)
		//if err != nil {
		//	return err
		//}
		//if raw != nil {
		//	return fmt.Errorf("childTokenId = %d already owned by tokenId = %d ", childTokenId, tokenId)
		//}
		//
		//err = ctx.GetStub().PutState(childTokensKey, []byte(fmt.Sprintf("%d", childTokenId)))
		//if err != nil {
		//	return err
		//}
		//childTokenOwnerKey, err := GetChildTokenOwnerKey(ctx, childContractName, childTokenId)
		//if err != nil {
		//	return err
		//}
		//err = ctx.GetStub().PutState(childTokenOwnerKey, []byte(fmt.Sprintf("%d", tokenId)))
		//if err != nil {
		//	return err
		//}
		//
		//childContractsKey, err := GetChildContractsKey(ctx, tokenId, childContractName)
		//if err != nil {
		//	return err
		//}
		//childContracts, err := ctx.GetStub().GetState(childContractsKey)
		//if err != nil {
		//	return err
		//}
		//if childContracts != nil {
		//	//childContractsKey, err := GetChildContractsKey(ctx, tokenId, childContractName)
		//	//if err != nil {
		//	//	return err
		//	//}
		//	//err = ctx.GetStub().PutState(childContractsKey, []byte(childContractName))
		//	//if err != nil {
		//	//	return err
		//	//}
		//	return fmt.Errorf("childContract = %s already owned by tokenId = %d ",childContractName,tokenId)
		//}
		//
		//err = ctx.GetStub().PutState(childContractsKey, []byte(childContractName))
		//if err != nil {
		//	return err
		//}
		response := ctx.GetStub().InvokeChaincode(parentContractName, util.ToChaincodeArgs("ReceiveNft721", strconv.FormatUint(tokenId, 10), childContractName, strconv.FormatUint(childTokenId, 10)), ctx.GetStub().GetChannelID())
		if response.Status != 200 {
			return fmt.Errorf("nft7 发送到 nft998 失败，msg = %s, parentContractName = %s, tokenId = %d, childContractName = %s, childTokenId = %d", response.Message, parentContractName, tokenId, childContractName, childTokenId)
		}
		err := nft.delToken(ctx, from, childTokenId)
		if err != nil {
			return err
		}
		err = nft.decreaseToken(ctx, from)
		if err != nil {
			return err
		}
		return nil
	}
	sender, err := getSender(ctx)
	if err != nil {
		return err
	}
	return fmt.Errorf("sender can not transfer, sender = %s, childTokenId = %d", sender, childTokenId)
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
