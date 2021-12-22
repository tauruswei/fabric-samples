package nft

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Mint 铸造 NFT
func (nft *NFT998) Mint(ctx contractapi.TransactionContextInterface, owner string, tokenID uint64, uri string) error {
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
	// todo 存放合同细节，每个合同属性也单独的作为 key-value 存放在区块链上

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
func (nft *NFT998) MintDNFT(ctx contractapi.TransactionContextInterface, tokenID,dtokenID,height uint64 ,expiration string) error {
	if !nft.canMintDnft(ctx, tokenID, expiration){
		return  fmt.Errorf("can not mint Dnft")
	}
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(KeyPrefixNFTDelegateTokenParent, []string{fmt.Sprintf("%d",dtokenID)})
	if err != nil {
		return err
	}
	defer iterator.Close()
	if iterator.HasNext() {
		return fmt.Errorf("dtokenId = %d is already assigned", dtokenID)
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

	owner, err := nft.OwnerOf(ctx, tokenID)
	if err != nil {
		return err
	}

	// 将 nft token 和用户绑定
	err = nft.addDToken(ctx, owner, tokenID,dtokenID)
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

// Revoke 注销 DNFT
func (nft *NFT998) RevokeDNFT(ctx contractapi.TransactionContextInterface, dtokenID uint64) error {
	// 先判断权限
	owner, err := nft.OwnerOf(ctx,dtokenID)
	if err != nil {
		return err
	}

	transfer := nft.canTransfer(ctx, owner, dtokenID)

	if !transfer{
		return fmt.Errorf("can not revoke dtoken, dtokenId = %d",dtokenID)
	}

	// 权限验证通过，删 key
	key, err := GetDtokenParentKey(ctx, dtokenID)
	if err != nil {
		return err
	}
	err = ctx.GetStub().DelState(key)
	if err != nil {
		return err
	}

	err = nft.delToken(ctx, owner, dtokenID)
	if err != nil {
		return err
	}
	err = nft.decreaseToken(ctx, owner)
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
	if err != nil {
		return err
	}
	totalKey, _ := GetTokenCountKey(ctx)

	return nft.calcCount(ctx, totalKey, false)
}

// BalanceOf owner 的 NFT 数量
func (nft *NFT998) BalanceOf(ctx contractapi.TransactionContextInterface, owner string) (uint64, error) {
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
func (nft *NFT998) OwnerOf(ctx contractapi.TransactionContextInterface, tokenID uint64) (string, error) {
	key, err := GetTokenOwnerKey(ctx, tokenID)
	if err != nil {
		return "", err
	}
	owner, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", err
	}
	return string(owner), nil
}

// SafeTransferFrom 根据 tokenID 将 NFT从 from 转移到 to
func (nft *NFT998) SafeTransferFrom(ctx contractapi.TransactionContextInterface, from string, to string, tokenID uint64, data []byte) error {
	if !nft.canTransfer(ctx, from, tokenID) {
		return errors.New("can not transfer")
	}
	return nft.TransferFrom(ctx, from, to, tokenID)
}

// TransferFrom 根据 tokenID 将 NFT从 from 转移到 to
func (nft *NFT998) TransferFrom(ctx contractapi.TransactionContextInterface, from string, to string, tokenID uint64) error {

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

// Approve 授予 approved 拥有 tokenID 的转移权力
func (nft *NFT998) Approve(ctx contractapi.TransactionContextInterface, approved string, tokenID uint64) error {
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
func (nft *NFT998) SetApprovalForAll(ctx contractapi.TransactionContextInterface, operator string, approved bool) error {
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
func (nft *NFT998) GetApproved(ctx contractapi.TransactionContextInterface, tokenID uint64) (string, error) {
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
func (nft *NFT998) IsApprovedForAll(ctx contractapi.TransactionContextInterface, owner string, operator string) (bool, error) {
	key, err := GetApprovedAllKey(ctx, owner, operator)
	if err != nil {
		return false, err
	}
	raw, err := ctx.GetStub().GetState(key)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(string(raw))
}

func (nft *NFT998) canTransfer(ctx contractapi.TransactionContextInterface, from string, tokenID uint64) bool {
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
	owner, err := nft.OwnerOf(ctx, tokenID)
	if err != nil {
		fmt.Printf("OwnerOf error: %s", err.Error())
		return false
	}
	// 所有资产授权操作人
	sender, _ := getSender(ctx)
	if err != nil {
		fmt.Printf("OwnerOf error: %s", err.Error())
		return false
	}
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
func (nft *NFT998) canMintDnft(ctx contractapi.TransactionContextInterface, tokenID uint64,expiration string) bool {
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
	heightKey, err := GetDtokenHeightKey(ctx, tokenID)
	if err != nil {
		fmt.Printf("get Dtoken Height Key error: %s", err.Error())
		return false
	}
	height, err := ctx.GetStub().GetState(heightKey)
	if err != nil {
		fmt.Printf("get Dtoken Height error: %s", err.Error())
		return false
	}
	if height==nil{
		return true
	}
	expirationKey, err := GetDtokenExpirationKey(ctx, tokenID)
	if err != nil {
		fmt.Printf("get Dtoken Expiration Key error: %s", err.Error())
		return false
	}
	expirationBytes, err := ctx.GetStub().GetState(expirationKey)
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

func (nft *NFT998) delToken(ctx contractapi.TransactionContextInterface, from string, tokenID uint64) error {
	key, err := GetTokenOwnerKey(ctx, tokenID)
	if err != nil {
		return err
	}
	return ctx.GetStub().DelState(key)

}

func (nft *NFT998) increaseToken(ctx contractapi.TransactionContextInterface, to string) error {
	key, err := GetTokenCountByOwnerKey(ctx, to)
	if err != nil {
		return err
	}
	return nft.calcCount(ctx, key, true)
}

func (nft *NFT998) decreaseToken(ctx contractapi.TransactionContextInterface, from string) error {
	key, err := GetTokenCountByOwnerKey(ctx, from)
	if err != nil {
		return err
	}
	return nft.calcCount(ctx, key, false)
}

// 将 nft token 与用户绑定
func (nft *NFT998) addToken(ctx contractapi.TransactionContextInterface, to string, tokenID uint64) error {
	key, err := GetTokenOwnerKey(ctx, tokenID)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(key, []byte(to))
}

// 将 delegated token 和 parent token 绑定
func (nft *NFT998) addDToken(ctx contractapi.TransactionContextInterface, to string,tokenID, dtokenID uint64) error {
	key, err := GetDtokenParentKey(ctx, dtokenID)
	if err != nil {
		return err
	}
	int64Str := strconv.FormatUint(tokenID, 10)
	err = ctx.GetStub().PutState(key, []byte(int64Str))
	if err != nil {
		return err
	}
	tokenOwnerKey, err := GetTokenOwnerKey(ctx, dtokenID)
	fmt.Printf("key = %s , owner = %s \n",tokenOwnerKey,to)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(tokenOwnerKey, []byte(to))
}

// 修改 nft token 数量
func (nft *NFT998) calcCount(ctx contractapi.TransactionContextInterface, key string, increase bool) error {
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

func (nft *NFT998) GetSender(ctx contractapi.TransactionContextInterface) (string, error) {
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
//func (nft *NFT998) SendNft721ToNft998(ctx contractapi.TransactionContextInterface,from string,tokenId uint64, childContractName string,childTokenId uint64) error {
//	transfer := nft.canTransfer(ctx, from, childTokenId)
//	if transfer{
//		_, err := TokenToOwner(ctx, tokenId)
//		if err != nil {
//			return err
//		}
//		index, err := ChildTokenIndex(ctx, tokenId, childContractName, childTokenId)
//		if index!=0{
//			return fmt.Errorf("Cannot send child token because it has already been received, tokenId = %d, childContractName = %s, childTokenId=%d",tokenId,childContractName,childTokenId)
//		}
//		childTokenIndexKey, err := GetChildTokenIndexKey(ctx, tokenId, childContractName, childTokenId)
//		if err != nil {
//			return err
//		}
//		err = ctx.GetStub().PutState(childTokenIndexKey, []byte(fmt.Sprintf("%s", childTokenId)))
//		if err != nil {
//			return err
//		}
//		childTokensKey, err := GetChildTokensKey(ctx, tokenId, childContractName, childTokenId)
//		if err != nil {
//			return err
//		}
//		err = ctx.GetStub().PutState(childTokensKey, []byte(fmt.Sprintf("%d",childTokenId)))
//		if err != nil {
//			return err
//		}
//		childTokenOwnerKey, err := GetChildTokenOwnerKey(ctx, childContractName, childTokenId)
//		if err != nil {
//			return err
//		}
//		err = ctx.GetStub().PutState(childTokenOwnerKey, []byte(fmt.Sprintf("%d", tokenId)))
//		if err != nil {
//			return err
//		}
//		childContractIndexKey, err := GetChildContractIndexKey(ctx, tokenId, childContractName)
//		if err != nil {
//			return err
//		}
//		childContractIndex, err := ctx.GetStub().GetState(childContractIndexKey)
//		if err != nil {
//			return err
//		}
//		if string(childContractIndex)== ""{
//			err := ctx.GetStub().PutState(childContractIndexKey, []byte(childContractName))
//			if err != nil {
//				return err
//			}
//			childContractsKey, err := GetChildContractsKey(ctx, tokenId, childContractName)
//			if err != nil {
//				return err
//			}
//			err = ctx.GetStub().PutState(childContractsKey, []byte(childContractName))
//			if err != nil {
//				return err
//			}
//		}
//	}
//	sender, err := getSender(ctx)
//	if err != nil {
//		return err
//	}
//	return fmt.Errorf("sender can not transfer, sender = %s, childTokenId = %d",sender,childTokenId)
//}
