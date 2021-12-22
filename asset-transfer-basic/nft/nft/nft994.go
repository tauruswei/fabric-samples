package nft

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strconv"
)

/**
 * @Author: fengxiaoxiao /13156050650@163.com
 * @Desc:
 * @Version: 1.0.0
 * @Date: 2021/12/22 1:57 下午
 */
// Mint 铸造 DNFT
func (nft *NFT) MintDNFT(ctx contractapi.TransactionContextInterface, tokenID, dtokenID, height uint64, expiration string) error {
	if !nft.canMintDnft(ctx, tokenID, expiration) {
		return fmt.Errorf("can not mint Dnft")
	}
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(KeyPrefixNFTDelegateTokenParent, []string{fmt.Sprintf("%d", dtokenID)})
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
	err = ctx.GetStub().PutState(heightKey, []byte(strconv.FormatUint(height, 10)))
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
	err = nft.addDToken(ctx, owner, tokenID, dtokenID)
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
func (nft *NFT) RevokeDNFT(ctx contractapi.TransactionContextInterface, dtokenID uint64) error {
	// 先判断权限
	owner, err := nft.OwnerOf(ctx, dtokenID)
	if err != nil {
		return err
	}

	transfer := nft.canTransfer(ctx, owner, dtokenID)

	if !transfer {
		return fmt.Errorf("can not revoke dtoken, dtokenId = %d", dtokenID)
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
