package nft

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"strings"
	"time"
)

// todo name owner 需要建立索引
type NFT994 struct {
	TokenId       string `json:"tokenId,omitempty"`       // 授权token的id
	RootTokenId   string `json:"rootTokenId,omitempty"`   // 歌曲的 721token id
	ParentTokenId string `json:"parentTokenId,omitempty"` // parent dtoken id
	ContractId    string `json:"contractId,omitempty"`    // 合同的 token id
	Name          string `json:"name,omitempty"`          // 歌曲名称
	OwnerName     string `json:"ownerName,omitempty"`     // 版权方名称
	Owner         string `json:"owner,omitempty"`         // 版权方地址
	Expiration    string `json:"expiration,omitempty"`    // 授权截止日期  format：2006-01-02 15:04:05，该属性应对特例：一个合同可能存在授权不同期限的歌曲
}

/*
 * @Desc: 查询 歌曲的版权，要求歌曲版权过期时间大于当前时间
 * @Param:
 * @Return:
 */
func (nft *NFT) QueryMusicDelegatedTokens(ctx contractapi.TransactionContextInterface, rootTokenId string) (string, error) {
	//sender, err := nft.GetSender(ctx)
	//if err != nil {
	//	return "", err
	//}
	queryString := fmt.Sprintf("{\"selector\":{\"rootTokenId\":\"%s\"}}", rootTokenId)
	logger.Debugf("queryString = %s", queryString)
	queryResults, err := getQueryResultForQueryString(ctx, queryString)
	if err != nil {
		logger.Error(GetErrorStackf(err, ""))
		return "", err
	}

	if string(queryResults) == "[]" {
		logger.Error(GetErrorStackf(nil, "could find any record in the ledger"))
		return "", fmt.Errorf("could find any record in the ledger")
	}
	delegatedTokens := []NFT994{}
	results := []NFT994{}

	err = json.Unmarshal(queryResults, &delegatedTokens)
	if err != nil {
		logger.Error(GetErrorStackf(err, ""))
		return "", err
	}

	for _, delegatedToken := range delegatedTokens {
		timestamp := time.Now().Unix()                      //1504079553
		timeNow := time.Unix(timestamp, 0)                  //2017-08-30 16:19:19 +0800 CST
		timeString := timeNow.Format("2006-01-02 15:04:05") //2015-06-15 08:52:32
		compare := strings.Compare(timeString, delegatedToken.Expiration)
		//logger.Debugf("compare = %d",compare)
		if compare < 0 {
			results = append(results, delegatedToken)
		}
	}

	marshal, err := json.Marshal(results)
	if err != nil {
		logger.Error(GetErrorStackf(err, ""))
		return "", err
	}
	return string(marshal), nil
}

/**
 * @Author: fengxiaoxiao /13156050650@163.com
 * @Desc:
 * @Version: 1.0.0
 * @Date: 2021/12/22 1:57 下午
 */
// Mint 铸造 DNFT
func (nft *NFT) CreateMusicDelegatedNFT(ctx contractapi.TransactionContextInterface, tokenIds, dtokenIds []string, contractId, ownerName, publicKey string, expiration []string, messageTokenId string, userType int, musicIds []int) error {

	logger.Debugf("method = %s, contractId = %s", "CreateMusicDelegatedNFT", contractId)
	if len(tokenIds) != len(expiration) {
		detail, err := nft.QueryContractDetail(ctx, contractId)
		if err != nil {
			logger.Error(GetErrorStackf(err, "query contract detail error"))
			return errors.WithMessagef(err, "query contract detail error")
		}
		contract := &Contract{}
		err = json.Unmarshal([]byte(detail), contract)
		if err != nil {
			logger.Error(GetErrorStackf(err, "unmarshal error, contract = %s", detail))
			return errors.WithMessagef(err, "unmarshal error, contract = %s", detail)
		}
		for _, _ = range tokenIds {
			expiration = append(expiration, contract.DelegateEndtDate)
		}
	}

	if len(tokenIds) != len(dtokenIds) {
		return fmt.Errorf("tokenId and dtokenId are expected to have the same length, tokenId = %+v, dtokenId = %+v", tokenIds, dtokenIds)
	}

	for index, tokenId := range tokenIds {
		rootTokenId, err := nft.canMintDnft(ctx, tokenId, expiration[index])
		if err != nil {
			return err
		}
		dtokenId := dtokenIds[index]
		iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(KeyPrefixNFTDelegateToken, []string{fmt.Sprintf("%d", dtokenId)})
		if err != nil {
			return err
		}
		defer iterator.Close()
		if iterator.HasNext() {
			return fmt.Errorf("dtokenId = %s is already assigned", dtokenId)
		}
		logger.Debugf("Mint dtoken %s for token %s ", dtokenId, tokenId)

		//heightKey, err := GetDtokenHeightKey(ctx, tokenId)
		//if err != nil {
		//	return err
		//}
		//err = ctx.GetStub().PutState(heightKey, []byte(strconv.FormatUint(height, 10)))
		//if err != nil {
		//	return err
		//}
		//expirationKey, err := GetDtokenExpirationKey(ctx, tokenId)
		//if err != nil {
		//	return err
		//}
		//err = ctx.GetStub().PutState(expirationKey, []byte(expiration))
		//if err != nil {
		//	return err
		//}
		//根据公钥获取 owner 地址
		pubPemBytes, err := base64.StdEncoding.DecodeString(publicKey)
		if err != nil {
			return err
		}
		block, _ := pem.Decode(pubPemBytes)

		key, err := x509.ParsePKIXPublicKey(block.Bytes)

		pubKey, ok := key.(*ecdsa.PublicKey)
		if !ok {
			return fmt.Errorf("parse public key error，publicKey = %s", publicKey)
		}
		pubBytes := crypto.FromECDSAPub(pubKey)
		owner := common.BytesToAddress(crypto.Keccak256(pubBytes[1:])[12:]).String()
		//owner := ownerAddr

		logger.Debugf("address = %s", owner)

		rootNft := &NFT{}
		nftkey, err := GetTokenIdKey(ctx, rootTokenId)
		if err != nil {
			return err
		}
		rootNftBytes, err := ctx.GetStub().GetState(nftkey)
		if err != nil {
			return err
		}
		err = json.Unmarshal(rootNftBytes, rootNft)
		if err != nil {
			return err
		}
		dnft := NFT994{
			TokenId:       dtokenId,
			RootTokenId:   rootTokenId,
			ParentTokenId: tokenId,
			ContractId:    contractId,
			Name:          rootNft.Name,
			OwnerName:     ownerName,
			Owner:         owner,
			Expiration:    expiration[index],
		}
		logger.Debugf("dnft = %+v", dnft)
		dnftBytes, err := json.Marshal(dnft)
		if err != nil {
			return err
		}
		dnftKey, err := GetDtokenKey(ctx, dtokenId)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(dnftKey, dnftBytes)
		if err != nil {
			return err
		}
		//owner, err := nft.OwnerOf(ctx, tokenId)
		//if err != nil {
		//	return err
		//}
		//
		//// 将 nft token 和用户绑定
		//err = nft.addDToken(ctx, owner, tokenId, dtokenId)
		//if err != nil {
		//	return err
		//}
		//// 修改 该用户的 token 数量
		//err = nft.increaseToken(ctx, owner)
		//if err != nil {
		//	return err
		//}
		//totalKey, _ := GetTokenCountKey(ctx)
		//// 修改 token 总量
		//return nft.calcCount(ctx, totalKey, true)

	}
	return nil
}

// Revoke 注销 DNFT
func (nft *NFT) RevokeDNFT(ctx contractapi.TransactionContextInterface, dtokenID string) error {
	// 先判断权限
	//owner, err := nft.OwnerOf(ctx, dtokenID)
	//if err != nil {
	//	return err
	//}
	//
	//transfer := nft.canTransfer(ctx, owner, dtokenID)
	//
	//if !transfer {
	//	return fmt.Errorf("can not revoke dtoken, dtokenId = %d", dtokenID)
	//}

	// 权限验证通过，删 key
	key, err := GetDtokenKey(ctx, dtokenID)
	if err != nil {
		return err
	}
	return ctx.GetStub().DelState(key)

	//err = nft.delToken(ctx, owner, dtokenID)
	//if err != nil {
	//	return err
	//}
	//err = nft.decreaseToken(ctx, owner)
	//if err != nil {
	//	return err
	//}
	//
	//heightKey, err := GetDtokenHeightKey(ctx, dtokenID)
	//if err != nil {
	//	return err
	//}
	//err = ctx.GetStub().DelState(heightKey)
	//if err != nil {
	//	return err
	//}
	//expirationKey, err := GetDtokenExpirationKey(ctx, dtokenID)
	//if err != nil {
	//	return err
	//}
	//err = ctx.GetStub().DelState(expirationKey)
	//if err != nil {
	//	return err
	//}
	//totalKey, _ := GetTokenCountKey(ctx)
	//
	//return nft.calcCount(ctx, totalKey, false)
}
func (nft *NFT) canMintDnft(ctx contractapi.TransactionContextInterface, tokenID, expiration string) (string, error) {
	//todo

	// 首先进行权限的判断
	//owner, err := nft.OwnerOf(ctx, tokenID)
	//if err != nil {
	//	//fmt.Printf("get owner of tokenId=%d error: %s", tokenID, err.Error())
	//	return 0, fmt.Errorf("get owner of tokenId=%d error: %s", tokenID, err.Error())
	//}
	//
	//transfer := nft.canTransfer(ctx, owner, tokenID)
	//if !transfer {
	//	return 0, fmt.Errorf("token does not owned by from, tokenId = %d, from= %s ", tokenID, owner)
	//}

	// 其次判断 token 是不是 Dtoken,如果是 dtoken，要比较 expiration
	dtokenKey, err := GetDtokenKey(ctx, tokenID)
	if err != nil {
		fmt.Printf("get Dtoken Height Key error: %s", err.Error())
		return "", fmt.Errorf("get Dtoken Height Key error: %s", err.Error())
	}
	dtokenBytes, err := ctx.GetStub().GetState(dtokenKey)
	if err != nil {
		fmt.Printf("get dtoken error: %s", err.Error())
		return "", fmt.Errorf("get dtoken error: %s", err.Error())
	}
	if dtokenBytes == nil {
		tokenIdKey, err := GetTokenIdKey(ctx, tokenID)
		if err != nil {
			return "", err
		}
		state, err := ctx.GetStub().GetState(tokenIdKey)
		if err != nil {
			return "", err
		}
		if state != nil {
			return tokenID, nil
		}
	}

	//expirationKey, err := GetDtokenExpirationKey(ctx, tokenID)
	//if err != nil {
	//	fmt.Printf("get Dtoken Expiration Key error: %s", err.Error())
	//	return 0,false
	//}
	//expirationBytes, err := ctx.GetStub().GetState(expirationKey)
	//if err != nil {
	//	fmt.Printf("get Dtoken Expiration error: %s", err.Error())
	//	return 0,false
	//}
	logger.Debugf("dtoken = %s", string(dtokenBytes))

	dtoken := NFT994{}
	err = json.Unmarshal(dtokenBytes, &dtoken)
	if err != nil {
		return "", fmt.Errorf("unmarshal dtoken error: %s", err.Error())
	}

	//先把时间字符串格式化成相同的时间类型
	//t1, err1 := time.Parse("2006-01-02 15:04:05", dtoken.Expiration)
	//t2, err2 := time.Parse("2006-01-02 15:04:05", expiration)expiration
	//if err1 == nil && err2 == nil && t1.Before(t2) {
	//	fmt.Printf("dtoken's expiration = %s is after parent token's expiration = %s, error: %s", dtoken.Expiration, expiration, err.Error())
	//	return 0, fmt.Errorf("dtoken's expiration = %s is after parent token's expiration = %s, error: %s", dtoken.Expiration, expiration, err.Error())
	//}
	compare := strings.Compare(dtoken.Expiration, expiration)
	if compare < 0 {
		logger.Error(GetErrorStackf(nil, "dtoken's expiration = %s is after parent token's expiration = %s", expiration, dtoken.Expiration))
		return "", fmt.Errorf("dtoken's expiration = %s is after parent token's expiration = %s", expiration, dtoken.Expiration)

	}
	return dtoken.RootTokenId, nil
}

/*
 * @Desc: 查询 nft 列表
 * @Param:
 * @Return:
 */
func (nft *NFT) DelegatedNFTListInContact(ctx contractapi.TransactionContextInterface, contractId string, pageSize int32, bookmark string) (string, error) {
	logger.Debugf("method = %s, contractId = %s", "DelegatedNFTListInContact", contractId)
	key, err := GetContractKey(ctx, contractId)
	if err != nil {
		return "", err
	}
	state, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", err
	}
	if state == nil {
		logger.Error(GetErrorStackf(nil, "could not find contract, contractId = %s", contractId))
		return "", fmt.Errorf("could not find contract, contractId = %s", contractId)
	}

	queryString := fmt.Sprintf("{\"selector\":{\"contractId\":\"%s\"}}", contractId)
	//queryString := "{\"selector\":{\"contractId\":{\"$exists\": true}}}"
	logger.Debugf("queryString = %s", queryString)
	// 多查询的一条数据 可以当作下次查询的 benchmark

	//queryResults, err := getQueryResultForQueryString(ctx, queryString)
	//if err != nil {
	//	logger.Error(GetErrorStackf(err, ""))
	//	return "", err
	//}
	//
	//if string(queryResults) == "[]" {
	//	logger.Error(GetErrorStackf(nil, "could find any record in the ledger"))
	//	return "", fmt.Errorf("could find any record in the ledger")
	//}
	//delegatedTokens := []NFT994{}
	//results := []NFT994{}
	//
	//err = json.Unmarshal(queryResults, &delegatedTokens)
	//if err != nil {
	//	logger.Error(GetErrorStackf(err, ""))
	//	return "", err
	//}
	//
	//for _, delegatedToken := range delegatedTokens {
	//	timestamp := time.Now().Unix()                      //1504079553
	//	timeNow := time.Unix(timestamp, 0)                  //2017-08-30 16:19:19 +0800 CST
	//	timeString := timeNow.Format("2006-01-02 15:04:05") //2015-06-15 08:52:32
	//	compare := strings.Compare(timeString, delegatedToken.Expiration)
	//	if compare > 0 {
	//		results = append(results, delegatedToken)
	//	}
	//}
	//
	//marshal, err := json.Marshal(results)
	//if err != nil {
	//	logger.Error(GetErrorStackf(err, ""))
	//	return "", err
	//}
	//return string(marshal), nil

	queryResults, err := getQueryResultForQueryStringWithPagination(ctx, queryString, pageSize, bookmark)
	if err != nil {
		return "", err
	}
	marshal, err := json.Marshal(queryResults)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}
