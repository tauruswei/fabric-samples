package nft

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

// NFT nft chaincode

type NFT struct {
	contractapi.Contract
}
type BasicNFT struct {
	TokenId    string  `json:"tokenId,omitempty"`
	Name       string  `json:"name,omitempty" binding:"required"`
	Label      string  `json:"label,omitempty"`
	UriPicture string  `json:"uriPicture,omitempty"` // 图片的网址
	UriVideo   string  `json:"uriVideo,omitempty"`   // 音频的网址
	Desc       string  `json:"desc,omitempty"`
	Status     int     `json:"status,omitempty"` // nft 状态   0 未上架； 1 正在卖； 2 已卖出
	Price      float64 `json:"price,omitempty"`
	Owner      string  `json:"owner,omitempty"`
	OwnerName  string  `json:"ownerName,omitempty"` // 所有者名称
	NickName   string  `json:"nickName ,omitempty"` // 提交交易的用户。前端传参数，会把 nickname 传递过来，该参数暂时没用
}

// 歌曲名称，演唱真，词作者，曲作者，确定歌曲的唯一性
type Song struct {
	BasicNFT
	MusicType               int    `json:"musicType,omitempty"`               //类型：1 歌曲，2 MV；3 微电影
	Number                  string `json:"number,omitempty"`                  // ISRC 号码
	CollectionName          string `json:"collectionName,omitempty"`          // 转接名称
	Singer                  string `json:"singer,omitempty"`                  // 表演者名称
	IssuingDate             string `json:"issuingDate,omitempty"`             // 发行时间
	SongWriter              string `json:"songWriter,omitempty"`              // 词作者
	Composer                string `json:"composer,omitempty"`                // 曲作者
	SongProportion          string `json:"songProportion,omitempty"`          // 词权比例
	ComposerProportion      string `json:"composerProportion,omitempty"`      // 曲权比例
	NeighboringRightsPropor string `json:"neighboringRightsPropor,omitempty"` // 邻接权权利比例
	Language                string `json:"language,omitempty"`                // 语种
	MV
}
type MV struct {
	Player            string `json:"player,omitempty"`            // 表演者
	Director          string `json:"director,omitempty"`          // 导演
	Producer          string `json:"producer,omitempty"`          // 制片人
	Time              string `json:"time,omitempty"`              // 授权方拥有的权利期限
	CopyrightProption string `json:"copyrightProption,omitempty"` // 著作权比例
}
type Contract struct {
	ConmmonContract
	DelegateStartDate     string `json:"delegateStartDate,omitempty"`     // 授权截止日期  2006-01-02 15:04:05
	DelegateEndtDate      string `json:"delegateEndDate,omitempty"`       // 授权开始日期
	DelegateRelation      string `json:"delegateRelation,omitempty"`      // 授权关系 代理
	DelegateType          string `json:"delegateType,omitempty"`          // 授权形式 独家/非独家
	DelegateRegion        string `json:"delegateRegion,omitempty"`        // 授权区域
	DelegateForm          string `json:"delegateForm,omitempty"`          // 授权形式
	CanTransferDelegation bool   `json:"canTransferDelegation,omitempty"` // 是否可以转授权
	Remark                string `json:"remark"`                          // 备注
}
type DelegateContract struct {
	ConmmonContract
	PublicationRight         string `json:"publicationRight,omitempty"`         // 发表权
	SignatureRight           string `json:"signatureRight,omitempty"`           // 署名权
	AmendmentRight           string `json:"amendmentRight,omitempty"`           // 修改权
	KeepIntegrityRight       string `json:"keepIntegrityRight,omitempty"`       // 保护作品完整权
	ReproductionRight        string `json:"reproductionRight,omitempty"`        // 复制权
	IssuingRight             string `json:"issuingRight,omitempty"`             // 发行权
	RentalRight              string `json:"rentalRight,omitempty"`              // 出租权
	ExhibitionRight          string `json:"exhibitionRight,omitempty"`          // 展览权
	PerformingRight          string `json:"performingRight,omitempty"`          // 表演权
	ShowRight                string `json:"showRight,omitempty"`                // 放映权
	BroadcastRight           string `json:"broadcastRight,omitempty"`           // 广播权
	NetworkTransmissionRight string `json:"networkTransmissionRight,omitempty"` // 信息网络传播权
	ShootRight               string `json:"shootRight,omitempty"`               // 摄制权
	AdaptRight               string `json:"adaptRight,omitempty"`               // 改编权
	TranslationRight         string `json:"translationRight,omitempty"`         // 翻译权
	CompilationRight         string `json:"compilationRight,omitempty"`         // 汇编权
	OtherRight               string `json:"otherRight,omitempty"`               // 应当由著作权人享有的其他权利
	Remark                   string `json:"remark,omitempty"`
}
type ConmmonContract struct {
	OwnerName      string `json:"ownerName,omitempty" binding:"required"`      // 客户名称
	OtherOwnerName string `json:"otherOwnerName,omitempty" binding:"required"` // 版权方名称
	ContractType   int    `json:"contractType,omitempty" binding:"required"`   // 合同类型 1 音乐合作  2  版权交易  3版权授权
	BasicNFT
}

type DataListResult struct {
	BookMark string        `json:"bookMark"`
	Count    int64         `json:"count"`
	DataList []interface{} `json:"dataList"`
}

/*
 * @Desc:
 * @Param:
 * @Return:
 */
func (nft *NFT) Init(ctx contractapi.TransactionContextInterface) {
	defaultFormat := "%{color}%{time:2006-01-02 15:04:05.000} %{shortfile:15s} [->] %{shortfunc:-10s} %{level:.4s} %{id:03x}%{color:reset} %{message}"
	defaultLevel := "DEBUG"
	InitLog(defaultFormat, defaultLevel)
	fmt.Println("init success")
}

/*
 * @Desc: 查询 歌曲的 721 token
 * @Param:
 * @Return:
 */
func (nft *NFT) QueryMusicNFTToken(ctx contractapi.TransactionContextInterface, musicType int, name, singer, songWriter, composer, player, director, producer string) (string, error) {
	var queryString string
	if musicType == 1 {
		if name != "" && singer != "" && songWriter != "" && composer != "" {
			queryString = fmt.Sprintf("{\"selector\":{\"name\":\"%s\",\"musicType\":%d,\"singer\":\"%s\",\"songWriter\":\"%s\",\"composer\":\"%s\"}}", name, musicType, singer, songWriter, composer)
		} else {
			return "", fmt.Errorf("we need 5 parameters to definitively query music tokenId, musicType=%d, name=%s, singer =%s, songWriter=%s, composer=%s", musicType, name, singer, songWriter, composer)
		}
	} else if musicType == 2 {
		if name != "" && player != "" && director != "" && producer != "" {
			queryString = fmt.Sprintf("{\"selector\":{\"name\":\"%s\",\"musicType\":%d,\"player\":\"%s\",\"director\":\"%s\",\"producer\":\"%s\"}}", name, musicType, player, director, producer)
		} else {
			return "", fmt.Errorf("we need 5 parameters to definitively query music tokenId, musicType=%d, name=%s, player =%s, director=%s, producer=%s", musicType, name, player, director, producer)
		}
	} else {
		return "", fmt.Errorf("musicType can not be null, musicType=%d", musicType)
	}

	logger.Debugf("queryString = %s", queryString)
	//queryString = "{\"selector\":{\"name\":\"string\",\"singer\":\"zjl\",\"songWriter\":\"wbt\",\"composer\":\"wbt\"}}"
	queryResults, err := getQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return "", err
	}
	if string(queryResults) == "[]" {
		logger.Error(GetErrorStackf(nil, "could not find any record in the ledger"))
		return "", fmt.Errorf("could not find any record in the ledger")
	}
	logger.Debugf("queryResults = %s", string(queryResults))
	song := []Song{}
	err = json.Unmarshal(queryResults, &song)
	if err != nil {
		logger.Error(GetErrorStackf(err, "json unmarshal error， queryResults = %s", string(queryResults)))
		return "", errors.WithMessagef(err, "json unmarshal error， queryResults = %s", string(queryResults))
	}
	logger.Debugf("song: %+v", song)
	if len(song) > 1 {
		logger.Error(GetErrorStackf(nil, "we find more than one music token in the blockchain, song=%+v, musicType=%d, name=%s, player =%s, director=%s, producer=%s", song, musicType, name, player, director, producer))
		return "", fmt.Errorf("we find more than one music token in the blockchain, song=%+v, musicType=%d, name=%s, player =%s, director=%s, producer=%s", song, musicType, name, player, director, producer)
	}

	return string(queryResults), nil
}

type TransferHistory struct {
	TxId      string `json:"txId"`
	Value     string `json:"value"`
	Timestamp int64  `json:"timestamp"`
}

/*
 * @Desc: 查询 721 token 的 history
 * @Param:
 * @Return:
 */
func (nft *NFT) QueryNFTTokenHistory(ctx contractapi.TransactionContextInterface, tokenId string) (string, error) {
	logger.Debugf("method = QueryNFTTokenHistory, tokenId = %s", tokenId)

	tokenKey, err := GetTokenIdKey(ctx, tokenId)
	if err != nil {
		return "", err
	}
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(tokenKey)
	//var buffer bytes.Buffer

	//buffer.WriteString("[")

	//bArrayMemberAlreadyWritten := false

	histories := []TransferHistory{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}
		history := TransferHistory{TxId: queryResponse.TxId, Value: string(queryResponse.Value), Timestamp: queryResponse.Timestamp.Seconds}
		histories = append(histories, history)
		// Add a comma before array members, suppress it for the first array member
		//if bArrayMemberAlreadyWritten == true {
		//	buffer.WriteString(",")
		//}
		logger.Debugf(string(queryResponse.Value))
		//bArrayMemberAlreadyWritten = true
	}
	//buffer.WriteString("]")
	//logger.Debugf("buffer = %+v", string(buffer.Bytes()))
	marshal, err := json.Marshal(histories)
	if err != nil {
		logger.Error(GetErrorStackf(err, "json marshal error, history = %+v", histories))
		return "", errors.WithMessagef(err, "json marshal error, history = %+v", histories)
	}
	logger.Debugf("history: %s", string(marshal))
	return string(marshal), nil
}

/*
 * @Desc: 查询 歌曲的 授权 994 token，要求 授权token的过期时间大于 当前时间，默认返回符合要求的第一条数据
 * @Param:
 * @Return:
 */
func (nft *NFT) QueryMusicDelegatedToken(ctx contractapi.TransactionContextInterface, rootTokenId, ownerName string) (string, error) {
	var queryString string
	if rootTokenId == "" {
		logger.Error(GetErrorStackf(nil, "rootTokenId can not be null"))
		return "", fmt.Errorf("rootTokenId can not be null")
	}
	if ownerName == "" {
		logger.Error(GetErrorStackf(nil, "ownerName can not be null"))
		return "", fmt.Errorf("ownerName can not be null")
	}

	queryString = fmt.Sprintf("{\"selector\":{\"rootTokenId\":\"%s\",\"ownerName\":\"%s\"}}", rootTokenId, ownerName)

	logger.Debugf("queryString = %s", queryString)
	//queryString = "{\"selector\":{\"name\":\"string\",\"singer\":\"zjl\",\"songWriter\":\"wbt\",\"composer\":\"wbt\"}}"
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
		if compare > 0 {
			marshal, err := json.Marshal(delegatedToken)
			if err != nil {
				logger.Error(GetErrorStackf(err, ""))
				return "", err
			}
			return string(marshal), nil
		}
	}
	logger.Debugf("queryResults = %s", string(queryResults))
	//song := []Song{}
	//err = json.Unmarshal(queryResults, &song)
	//if err != nil {
	//	logger.Error(GetErrorStackf(err, "json unmarshal error， queryResults = %s", string(queryResults)))
	//	return "", errors.WithMessagef(err, "json unmarshal error， queryResults = %s", string(queryResults))
	//}
	//logger.Debugf("song: %+v", song)
	//if len(song) > 1 {
	//	logger.Error(GetErrorStackf(nil, "we find more than one music token in the blockchain, song=%+v, musicType=%d, name=%s, player =%s, director=%s, producer=%s", song, musicType, name, player, director, producer))
	//	return "", fmt.Errorf("we find more than one music token in the blockchain, song=%+v, musicType=%d, name=%s, player =%s, director=%s, producer=%s", song, musicType, name, player, director, producer)
	//}

	return string(queryResults), nil
}

/*
 * @Desc: 铸造 音乐人合作/版权交易 合同NFT
 * @Param:
 * @Return:
 */
func (nft *NFT) CreateMusicContract(ctx contractapi.TransactionContextInterface, tokenId string, contractType int, data string) (string, error) {
	logger.Debugf("method = %s, contractType = %d, data = %s, tokenId = %s", "MintMusicContract", contractType, data, tokenId)
	tokenIdKey, err := GetContractKey(ctx, tokenId)
	if err != nil {
		return "", err
	}
	contractBytes, err := ctx.GetStub().GetState(tokenIdKey)
	if err != nil {
		return "", err
	}
	if contractBytes != nil {
		return "", fmt.Errorf("tokenId = %s is already assigned", tokenId)
	}
	sender, err := getSender(ctx)
	if err != nil {
		return "", err
	}
	contract := Contract{}
	err = json.Unmarshal([]byte(data), &contract)
	if err != nil {
		logger.Error(GetErrorStackf(err, "json unmarshal error， data = %s", data))
		return "", errors.WithMessagef(err, "json unmarshal error， data = %s", data)
	}
	logger.Debugf("contract = %+v", contract)

	//queryString := fmt.Sprintf("{\"selector\":{\"name\":\"%s\",\"ownerName\":\"%s\",\"otherOwnerName\":\"%s\"}}", contract.Name, contract.OwnerName, contract.OtherOwnerName)
	//
	//queryResults, err := getQueryResultForQueryString(ctx, queryString)
	//if err != nil {
	//	return "", err
	//}
	//if string(queryResults) != "[]" {
	//	logger.Debugf("queryResults = %s", string(queryResults))
	//	logger.Error(GetErrorStackf(err, "contract name already exist， data = %s,", data))
	//	return "", errors.WithMessagef(err, "contract name already exist， data = %s,", data)
	//}

	//song := &Song{}
	//mv := &MV{}
	//if contractType == 1 {
	//	err := json.Unmarshal([]byte(data), song)
	//	if err != nil {
	//		return err
	//	}
	//	song.Owner = sender
	//	//song.TokenId = tokenId
	//} else {
	//	err := json.Unmarshal([]byte(data), mv)
	//	if err != nil {
	//		return err
	//	}
	//	mv.Owner = sender
	//	//mv.TokenId = tokenId
	//}
	contract.Owner = sender
	marshal, err := json.Marshal(contract)
	if err != nil {
		return "", err
	}

	logger.Debugf("contract = %s", string(marshal))
	err = ctx.GetStub().PutState(tokenIdKey, marshal)
	if err != nil {
		return "", err
	}
	return tokenId, nil
}

/*
 * @Desc: 消息上链
 * @Param:
 * @Return:
 */
func (nft *NFT) CreateMusicMessage(ctx contractapi.TransactionContextInterface, tokenId, data string) error {
	logger.Debugf("method = %s, data = %s, tokenId = %s", "CreateMusicMessage", data, tokenId)
	tokenIdKey, err := GetMessageKey(ctx, tokenId)
	if err != nil {
		return err
	}
	msgBytes, err := ctx.GetStub().GetState(tokenIdKey)
	if err != nil {
		return err
	}
	if msgBytes != nil {
		return fmt.Errorf("duplicate key for tokenId = %s", tokenId)
	}

	err = ctx.GetStub().PutState(tokenIdKey, []byte(data))
	if err != nil {
		return err
	}
	return nil
}

/*
 * @Desc: 读消息
 * @Param:
 * @Return:
 */
func (nft *NFT) ReadMusicMessage(ctx contractapi.TransactionContextInterface, tokenId string) (string, error) {
	logger.Debugf("method = %s,  tokenId= %s", "ReadMusicMessage", tokenId)
	tokenIdKey, err := GetMessageKey(ctx, tokenId)
	if err != nil {
		return "", err
	}
	msgBytes, err := ctx.GetStub().GetState(tokenIdKey)
	if err != nil {
		return "", err
	}
	return string(msgBytes), nil
}

/*
 * @Desc: 铸造 版权授权 NFT
 * @Param:
 * @Return:
 */
func (nft *NFT) CreateMusicDelegatedContract(ctx contractapi.TransactionContextInterface, tokenId, data string) (string, error) {
	logger.Debugf("method = %s, data = %s", "CreateMusicDelegatedContract", data)
	tokenIdKey, err := GetContractKey(ctx, tokenId)
	if err != nil {
		return "", err
	}
	contractBytes, err := ctx.GetStub().GetState(tokenIdKey)
	if err != nil {
		return "", err
	}
	if contractBytes != nil {
		return "", fmt.Errorf("tokenId = %s is already assigned", tokenId)
	}
	sender, err := getSender(ctx)
	if err != nil {
		return "", err
	}

	//song := &Song{}
	//mv := &MV{}
	//if contractType == 1 {
	//	err := json.Unmarshal([]byte(data), song)
	//	if err != nil {
	//		return err
	//	}
	//	song.Owner = sender
	//	//song.TokenId = tokenId
	//} else {
	//	err := json.Unmarshal([]byte(data), mv)
	//	if err != nil {
	//		return err
	//	}
	//	mv.Owner = sender
	//	//mv.TokenId = tokenId
	//}
	contract := DelegateContract{}
	err = json.Unmarshal([]byte(data), &contract)
	if err != nil {
		return "", err
	}
	contract.Owner = sender
	marshal, err := json.Marshal(contract)
	if err != nil {
		return "", err
	}
	err = ctx.GetStub().PutState(tokenIdKey, marshal)
	if err != nil {
		return "", err
	}
	return tokenId, nil
}

/*
 * @Desc: 查询 合同详情
 * @Param:
 * @Return:
 */
func (nft *NFT) QueryContractDetail(ctx contractapi.TransactionContextInterface, tokenId string) (string, error) {
	logger.Debugf("method = QueryContractDetail, tokenId = %s", tokenId)

	// todo 权限的判断

	//err := nft.IsAuthorized(ctx, tokenId)
	//if err != nil {
	//	//logger.Error(GetErrorStackf(err, "contract name already exist， data = %s,", data))
	//	//return "",errors.WithMessagef(err, "contract name already exist， data = %s,", data)
	//	return "", err
	//}
	key, err := GetContractKey(ctx, tokenId)
	if err != nil {
		return "", err
	}
	state, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", err
	}
	logger.Debugf("contract = %s", string(state))
	return string(state), nil

}

/*
 * @Desc: 创建音乐 NFT
 * @Param:
 * @Return:
 */

func (nft *NFT) CreateMusicNFT(ctx contractapi.TransactionContextInterface, tokenId string, data string) (string, error) {
	logger.Debugf("method = CreateMusicNFT, tokenId = %s, data = %s", tokenId, data)
	tokenIdKey, err := GetTokenIdKey(ctx, tokenId)
	if err != nil {
		return "", err
	}
	nftTokenBytes, err := ctx.GetStub().GetState(tokenIdKey)
	if err != nil {
		return "", err
	}
	if nftTokenBytes != nil {
		return "", fmt.Errorf("tokenId = %s is already assigned", tokenId)
	}

	sender, err := nft.GetSender(ctx)
	if err != nil {
		return "", err
	}

	song := Song{}
	err = json.Unmarshal([]byte(data), &song)
	if err != nil {
		return "", err
	}

	var queryString string
	if song.MusicType == 1 {
		if song.Name != "" && song.Singer != "" && song.SongWriter != "" && song.Composer != "" {
			queryString = fmt.Sprintf("{\"selector\":{\"name\":\"%s\",\"musicType\":%d,\"singer\":\"%s\",\"songWriter\":\"%s\",\"composer\":\"%s\"}}", song.Name, song.MusicType, song.Singer, song.SongWriter, song.Composer)
		} else {
			return "", fmt.Errorf("we need 5 parameters to definitively query music tokenId, name=%s, musicType=%d, singer =%s, songWriter=%s, composer=%s", song.Name, song.MusicType, song.Singer, song.SongWriter, song.Composer)
		}
	} else if song.MusicType == 2 {
		if song.Name != "" && song.Player != "" && song.Director != "" && song.Producer != "" {
			queryString = fmt.Sprintf("{\"selector\":{\"name\":\"%s\",\"musicType\":%d,\"player\":\"%s\",\"director\":\"%s\",\"producer\":\"%s\"}}", song.Name, song.MusicType, song.Player, song.Director, song.Producer)
		} else {
			return "", fmt.Errorf("we need 5 parameters to definitively query music tokenId,  name=%s, musicType=%d, player =%s, director=%s, producer=%s", song.Name, song.MusicType, song.Player, song.Director, song.Producer)
		}
	} else {
		return "", fmt.Errorf("musicType can not be null, musicType=%d", song.MusicType)
	}
	logger.Debugf("queryString = %s", queryString)
	queryResults, err := getQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return "", err
	}
	if string(queryResults) != "[]" {
		logger.Error(GetErrorStackf(nil, "song already exists, tokenId = %s, data = %s", tokenId, data))
		return "", fmt.Errorf("song already exists, tokenId = %s, data = %s", tokenId, data)
	}
	song.TokenId = tokenId
	song.Owner = sender

	//iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(KeyPrefixNFT, []string{fmt.Sprintf("%d", tokenId)})
	//defer iterator.Close()
	//if iterator.HasNext() {
	//	return fmt.Errorf("tokenId = %d is already assigned", tokenId)
	//}
	//fmt.Printf("Mint token %d for %s ,uri: %s\n", tokenId, owner, uri)
	//nameKey, err := GetTokenNameKey(ctx, tokenId)
	//if err != nil {
	//	return err
	//}
	//err = ctx.GetStub().PutState(nameKey, []byte(name))
	//if err != nil {
	//	return err
	//}
	//labelKey, err := GetTokenLabelKey(ctx, tokenId)
	//if err != nil {
	//	return err
	//}
	//err = ctx.GetStub().PutState(labelKey, []byte(label))
	//if err != nil {
	//	return err
	//}
	//uriKey, err := GetTokenURIKey(ctx, tokenId)
	//if err != nil {
	//	return err
	//}
	//err = ctx.GetStub().PutState(uriKey, []byte(uri))
	//if err != nil {
	//	return err
	//}
	//descKey, err := GetTokenDescKey(ctx, tokenId)
	//if err != nil {
	//	return err
	//}
	//err = ctx.GetStub().PutState(descKey, []byte(desc))
	//if err != nil {
	//	return err
	//}

	//
	marshal, err := json.Marshal(song)
	if err != nil {
		return "", err
	}
	err = ctx.GetStub().PutState(tokenIdKey, marshal)
	if err != nil {
		return "", err
	}
	return tokenId, nil
	//if err != nil {
	//	return err
	//}
	//
	//// 将 nft token 和用户绑定
	//err = nft.addToken(ctx, owner, tokenId)
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

// Mint 铸造 NFT
func (nft *NFT) Mint(ctx contractapi.TransactionContextInterface, owner string, tokenId, name, label, uriPicture, uriVideo, desc string) error {

	tokenIdKey, err := GetTokenIdKey(ctx, tokenId)
	if err != nil {
		return err
	}
	nftTokenBytes, err := ctx.GetStub().GetState(tokenIdKey)
	if err != nil {
		return err
	}
	if nftTokenBytes != nil {
		return fmt.Errorf("tokenId = %d is already assigned", tokenId)
	}

	//iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(KeyPrefixNFT, []string{fmt.Sprintf("%d", tokenId)})
	//defer iterator.Close()
	//if iterator.HasNext() {
	//	return fmt.Errorf("tokenId = %d is already assigned", tokenId)
	//}
	//fmt.Printf("Mint token %d for %s ,uri: %s\n", tokenId, owner, uri)
	//nameKey, err := GetTokenNameKey(ctx, tokenId)
	//if err != nil {
	//	return err
	//}
	//err = ctx.GetStub().PutState(nameKey, []byte(name))
	//if err != nil {
	//	return err
	//}
	//labelKey, err := GetTokenLabelKey(ctx, tokenId)
	//if err != nil {
	//	return err
	//}
	//err = ctx.GetStub().PutState(labelKey, []byte(label))
	//if err != nil {
	//	return err
	//}
	//uriKey, err := GetTokenURIKey(ctx, tokenId)
	//if err != nil {
	//	return err
	//}
	//err = ctx.GetStub().PutState(uriKey, []byte(uri))
	//if err != nil {
	//	return err
	//}
	//descKey, err := GetTokenDescKey(ctx, tokenId)
	//if err != nil {
	//	return err
	//}
	//err = ctx.GetStub().PutState(descKey, []byte(desc))
	//if err != nil {
	//	return err
	//}

	nft1 := BasicNFT{
		TokenId:    tokenId,
		Name:       name,
		Label:      label,
		UriPicture: uriPicture,
		UriVideo:   uriVideo,
		Desc:       desc,
		Owner:      owner,
	}
	marshal, err := json.Marshal(nft1)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(tokenIdKey, marshal)
	//if err != nil {
	//	return err
	//}
	//
	//// 将 nft token 和用户绑定
	//err = nft.addToken(ctx, owner, tokenId)
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

/*
 * @Desc: 查询 nft 列表
 * @Param:
 * @Return:
 */
func (nft *NFT) NFTList(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string) (string, error) {
	if bookmark == "nil" {
		bookmark = ""
	}
	logger.Debugf("method = %s, pageSize = %d, bookmark = %s", "NFTList", pageSize, bookmark)

	//queryString := fmt.Sprintf("{\"selector\":{\"status\":%t}}", false)
	queryString := "{\"selector\":{\"label\":{\"$exists\": true}}}"
	logger.Debugf("queryString = %s", queryString)
	// 多查询的一条数据 可以当作下次查询的 benchmark
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

/*
 * @Desc: 查询 用户的 NFT 列表
 * @Param:
 * @Return:
 */
func (nft *NFT) NFTListForUser(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string) (string, error) {
	if bookmark == "nil" {
		bookmark = ""
	}
	logger.Debugf("method = %s, pageSize = %d, bookmark = %s", "NFTListForUser", pageSize, bookmark)

	sender, err := getSender(ctx)
	if err != nil {
		return "", err
	}
	//queryString := fmt.Sprintf("{\"selector\":{\"owner\":\"%s\",\"label\":{\"$eq\": \"\"},\"contractType\":{\"$exists\":false}}}", sender)
	//queryString := fmt.Sprintf("{\"selector\":{\"$and\":[{\"owner\":{\"$eq\":\"%s\"}},{\"label\":{\"$exists\": true}},{\"contractType\":{\"$exists\":false}}]}}", sender)
	queryString := fmt.Sprintf("{\"selector\":{\"$and\":[{\"owner\":{\"$eq\":\"%s\"}},{\"label\":{\"$exists\": true}},{\"contractType\":{\"$exists\":false}}]}}", sender)

	logger.Debugf("queryString = %s", queryString)

	// 多查询的一条数据 可以当作下次查询的 benchmark

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

/*
 * @Desc: 查询 用户的 合同 列表
 * @Param:
 * @Return:
 */
func (nft *NFT) ContractListForUser(ctx contractapi.TransactionContextInterface, pageSize int32, userName, bookmark string) (string, error) {
	if bookmark == "nil" {
		bookmark = ""
	}
	logger.Debugf("method = %s, pageSize = %d, bookmark = %s", "ContractListForUser", pageSize, bookmark)

	//queryString := fmt.Sprintf("{\"selector\":[{\"$or\": [{\"ownerName\": {\"$eq\": \"%s\"}}, {\"otherOwnerName\": {\"$eq\": \"%s\"}}]},{\"contractType\":{\"$gt\":0}}]}", userName, userName)
	//queryString := fmt.Sprintf("{\"selector\":[{\"$or\": [{\"ownerName\": {\"$eq\": \"%s\"}}, {\"otherOwnerName\": {\"$eq\": \"%s\"}}]},{\"contractType\":{\"$gt\":0}}]}", userName, userName)
	//queryString := fmt.Sprintf("{\"selector\":{\"contractType\":{\"$gt\":0}}}")
	queryString := fmt.Sprintf("{\"selector\":{\"$and\":[{\"$or\": [{\"ownerName\": {\"$eq\": \"%s\"}}, {\"otherOwnerName\": {\"$eq\": \"%s\"}}]},{\"contractType\":{\"$gt\":0}}]}}", userName, userName)

	logger.Debugf("queryString = %s", queryString)
	// 多查询的一条数据 可以当作下次查询的 benchmark
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

/*
 * @Desc: 查询 合同 列表
 * @Param:
 * @Return:
 */
func (nft *NFT) ContractList(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string) (string, error) {
	if bookmark == "nil" {
		bookmark = ""
	}
	logger.Debugf("method = %s, pageSize = %d, bookmark = %s", "ContractList", pageSize, bookmark)

	resultsIterator, responseMetadata, err := ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(KeyPrefixContract, nil, pageSize, bookmark)

	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	result, err := constructQueryResponseFromIteratorPage(resultsIterator)
	if err != nil {
		return "", err
	}

	addPaginationMetadataToQueryResultsPage(result, responseMetadata)

	logger.Debugf("queryResults = %+v", result)

	marshal, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

// =========================================================================================
// getQueryResultForQueryStringWithPagination executes the passed in query string with
// pagination info. Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryStringWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int32, bookmark string) (*DataListResult, error) {

	//fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	result, err := constructQueryResponseFromIteratorPage(resultsIterator)
	if err != nil {
		return nil, err
	}

	addPaginationMetadataToQueryResultsPage(result, responseMetadata)

	logger.Debugf("queryResults = %+v", result)

	//fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", bufferWithPaginationInfo.String())
	//logger.Debugf("- query result = %+v",result)
	return result, nil
}

// ===========================================================================================
// addPaginationMetadataToQueryResults adds QueryResponseMetadata, which contains pagination
// info, to the constructed query results
// ===========================================================================================
func addPaginationMetadataToQueryResults(buffer *bytes.Buffer, responseMetadata *pb.QueryResponseMetadata) *bytes.Buffer {

	buffer.WriteString("[{\"ResponseMetadata\":{\"RecordsCount\":")
	buffer.WriteString("\"")
	buffer.WriteString(fmt.Sprintf("%v", responseMetadata.FetchedRecordsCount))
	buffer.WriteString("\"")
	buffer.WriteString(", \"Bookmark\":")
	buffer.WriteString("\"")
	buffer.WriteString(responseMetadata.Bookmark)
	buffer.WriteString("\"}}]")

	return buffer
}

// ===========================================================================================
// addPaginationMetadataToQueryResults adds QueryResponseMetadata, which contains pagination
// info, to the constructed query results
// ===========================================================================================
func addPaginationMetadataToQueryResultsPage(result *DataListResult, responseMetadata *pb.QueryResponseMetadata) {
	result.Count = int64(responseMetadata.FetchedRecordsCount)
	result.BookMark = responseMetadata.Bookmark
}

/*
 * @Desc: 修改 nft 的状态
 * @Param:
 * @Return:
 */

func (nft *NFT) SellNft(ctx contractapi.TransactionContextInterface, tokenId string, price float64) error {
	//todo 规则细化
	err := nft.IsAuthorized(ctx, tokenId)
	if err != nil {
		return err
	}

	//todo 调用 approve 接口，授权给平台管理员用户

	tokenIdKey, err := GetTokenIdKey(ctx, tokenId)
	if err != nil {
		return err
	}

	nftTokenBytes, err := ctx.GetStub().GetState(tokenIdKey)
	if err != nil {
		return err
	}
	nftToken := &NFT{}
	err = json.Unmarshal(nftTokenBytes, nftToken)
	if err != nil {
		return err
	}
	//
	//nftToken.Status = true
	//nftToken.Price = price

	marshal, err := json.Marshal(nftToken)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(tokenIdKey, marshal)
}

// ===== Example: Ad hoc rich query ========================================================
// queryMarbles uses a query string to perform a query for marbles.
// Query string matching state database syntax is passed in and executed as is.
// Supports ad hoc queries that can be defined at runtime by the client.
// If this is not desired, follow the queryMarblesForOwner example for parameterized queries.
// Only available on state databases that support rich query (e.g. CouchDB)
// =========================================================================================
func (nft *NFT) QueryNfts(ctx contractapi.TransactionContextInterface, queryString string) (string, error) {

	//   0
	// "queryString"
	//if len(args) < 1 {
	//	return nil,errors.New("Incorrect number of arguments. Expecting 1")
	//}
	//
	//queryString := args[0]

	queryResults, err := getQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return "", err
	}

	if string(queryResults) == "[]" {
		logger.Error(GetErrorStackf(nil, "could find any record in the ledger"))
		return "", fmt.Errorf("could find any record in the ledger")
	}
	return string(queryResults), nil
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]byte, error) {

	//fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

// constructQueryResponseFromIterator constructs a JSON array containing query results from
// a given result iterator
// ===========================================================================================
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		//buffer.WriteString("{\"Key\":")
		//buffer.WriteString("\"")
		//buffer.WriteString(queryResponse.Key)
		//buffer.WriteString("\"")
		//
		//buffer.WriteString(", \"Record\":")
		//// Record is a JSON object, so we write as-is
		//buffer.WriteString(string(queryResponse.Value))
		//buffer.WriteString("}")
		//buffer.WriteString("{")
		buffer.WriteString(string(queryResponse.Value))
		//buffer.WriteString("}")

		bArrayMemberAlreadyWritten = true

		//logger.Debugf("queryResponse.Value = %s",string(queryResponse.Value))
	}
	buffer.WriteString("]")
	logger.Debugf("buffer = %+v", string(buffer.Bytes()))

	return &buffer, nil
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForPage(ctx contractapi.TransactionContextInterface, queryString string) (*DataListResult, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	dataListResult, err := constructQueryResponseFromIteratorPage(resultsIterator)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	logger.Debugf("dataListResult = %+v", dataListResult)

	return dataListResult, nil
}

// constructQueryResponseFromIterator constructs a JSON array containing query results from
// a given result iterator
// ===========================================================================================
func constructQueryResponseFromIteratorPage(resultsIterator shim.StateQueryIteratorInterface) (*DataListResult, error) {
	// buffer is a JSON array containing QueryResults
	//var buffer bytes.Buffer
	dataListResult := DataListResult{}
	//buffer.WriteString("[")
	//
	//bArrayMemberAlreadyWritten := false
	//for resultsIterator.HasNext() {
	//	queryResponse, err := resultsIterator.Next()
	//	if err != nil {
	//		return nil, err
	//	}
	//	// Add a comma before array members, suppress it for the first array member
	//	if bArrayMemberAlreadyWritten == true {
	//		buffer.WriteString(",")
	//	}
	//	//buffer.WriteString("{\"Key\":")
	//	//buffer.WriteString("\"")
	//	//buffer.WriteString(queryResponse.Key)
	//	//buffer.WriteString("\"")
	//	//
	//	//buffer.WriteString(", \"Record\":")
	//	//// Record is a JSON object, so we write as-is
	//	//buffer.WriteString(string(queryResponse.Value))
	//	//buffer.WriteString("}")
	//	//buffer.WriteString("{")
	//	buffer.WriteString(string(queryResponse.Value))
	//	//buffer.WriteString("}")
	//
	//	bArrayMemberAlreadyWritten = true
	//}
	//buffer.WriteString("]")
	nft := BasicNFT{}
	nft994 := NFT994{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		logger.Debugf("queryResponse.Value = %s", string(queryResponse.Value))
		// todo
		if strings.Contains(string(queryResponse.Value), "rootTokenId") {
			err = json.Unmarshal(queryResponse.Value, &nft994)
			dataListResult.DataList = append(dataListResult.DataList, nft994)
		} else {
			err = json.Unmarshal(queryResponse.Value, &nft)
			dataListResult.DataList = append(dataListResult.DataList, nft)
		}
		if err != nil {
			return nil, err
		}
	}
	return &dataListResult, nil
}

// BalanceOf owner 的 NFT 数量
func (nft *NFT) BalanceOf(ctx contractapi.TransactionContextInterface, owner string) (uint64, error) {
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

// OwnerOf 根据 tokenId 返回其所有人地址
func (nft *NFT) OwnerOf(ctx contractapi.TransactionContextInterface, tokenId uint64) (string, error) {
	key, err := GetTokenOwnerKey(ctx, tokenId)
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

// SafeTransferFrom 根据 tokenId 将 NFT从 from 转移到 to
func (nft *NFT) SafeTransferFrom(ctx contractapi.TransactionContextInterface, from string, to string, tokenId uint64, data string) error {
	if !nft.canTransfer(ctx, from, tokenId) {
		return errors.New("can not transfer")
	}
	return nft.TransferFrom(ctx, from, to, tokenId)
}

// SafeTransferFrom 根据 tokenId 将 NFT从 from 转移到 to
func (nft *NFT) SafeTransferFromCouchdb(ctx contractapi.TransactionContextInterface, from string, to string, tokenId string, data string) error {
	logger.Debugf("method = SafeTransferFromCouchdb, from = %s, to = %s, tokenId = %s, data = %s", from, to, tokenId, data)

	if !nft.canTransferCouchdb(ctx, from, tokenId) {
		return errors.New("can not transfer")
	}
	return nft.TransferFromCouchdb(ctx, from, to, tokenId, data)
}

// TransferFrom 根据 tokenId 将 NFT从 from 转移到 to
func (nft *NFT) TransferFrom(ctx contractapi.TransactionContextInterface, from string, to string, tokenId uint64) error {

	err := nft.delToken(ctx, from, tokenId)
	if err != nil {
		return err
	}
	err = nft.decreaseToken(ctx, from)
	if err != nil {
		return err
	}
	err = nft.addToken(ctx, to, tokenId)
	if err != nil {
		return err
	}
	err = nft.increaseToken(ctx, to)
	if err != nil {
		return err
	}

	ownerTokensHistorykey, err := GetOwnerTokensHistoryKey(ctx, from)
	if err != nil {
		return err
	}
	tokenIdBytes, err := ctx.GetStub().GetState(ownerTokensHistorykey)

	return ctx.GetStub().PutState(ownerTokensHistorykey, []byte(fmt.Sprintf("%s-%s", string(tokenIdBytes), strconv.FormatUint(tokenId, 10))))
}

// TransferFromCouchdb 根据 tokenId 将 NFT从 from 转移到 to
func (nft *NFT) TransferFromCouchdb(ctx contractapi.TransactionContextInterface, from, to, tokenId, nickName string) error {
	logger.Debugf("method = TransferFromCouchdb, from = %s, to = %s, tokenId = %s", from, to, tokenId)

	key, err := GetTokenIdKey(ctx, tokenId)
	if err != nil {
		logger.Error(GetErrorStackf(err, fmt.Sprintf("get token key error, tokenId = %s", tokenId)))
		return errors.WithMessagef(err, fmt.Sprintf("get token key error, tokenId = %s", tokenId))
	}

	data, err := ctx.GetStub().GetState(key)
	if err != nil || data == nil {
		logger.Error(GetErrorStackf(err, fmt.Sprintf("could find any record in the ledger, tokenId = %s", tokenId)))
		return errors.WithMessagef(err, fmt.Sprintf("could find any record in the ledger, tokenId = %s", tokenId))
	}

	nft721 := BasicNFT{}
	err = json.Unmarshal(data, &nft721)
	if err != nil {
		logger.Error(GetErrorStackf(err, "json unmarshal error， data = %s", data))
		return errors.WithMessagef(err, "json unmarshal error， data = %s", data)
	}

	nft721.Owner = to
	nft721.OwnerName = nickName

	marshal, err := json.Marshal(nft721)
	if err != nil {
		logger.Error(GetErrorStackf(err, "json marshal error， nft721 = %+v", nft721))
		return errors.WithMessagef(err, "json marshal error， nft721 = %+v", nft721)
	}

	return ctx.GetStub().PutState(key, marshal)

}

// 获取 owner 曾经拥有的 nft token
func (nft *NFT) GetOwnerTokensHistory(ctx contractapi.TransactionContextInterface, from string, to string, tokenId uint64) error {

	err := nft.delToken(ctx, from, tokenId)
	if err != nil {
		return err
	}
	err = nft.decreaseToken(ctx, from)
	if err != nil {
		return err
	}
	err = nft.addToken(ctx, to, tokenId)
	if err != nil {
		return err
	}
	err = nft.increaseToken(ctx, to)
	if err != nil {
		return err
	}

	ownerTokensHistorykey, err := GetOwnerTokensHistoryKey(ctx, from)
	if err != nil {
		return err
	}
	tokenIdBytes, err := ctx.GetStub().GetState(ownerTokensHistorykey)

	return ctx.GetStub().PutState(ownerTokensHistorykey, []byte(fmt.Sprintf("%s-%s", string(tokenIdBytes), strconv.FormatUint(tokenId, 10))))

}

// TransferFrom 根据 tokenId 将 NFT 从 998 转移到 721
func (nft *NFT) ReceiveFromNft998(ctx contractapi.TransactionContextInterface, to string, tokenId uint64, data string) error {
	err := nft.addToken(ctx, to, tokenId)
	if err != nil {
		return err
	}
	err = nft.increaseToken(ctx, to)
	if err != nil {
		return err
	}
	return nil
}

// Approve 授予 approved 拥有 tokenId 的转移权力
func (nft *NFT) Approve(ctx contractapi.TransactionContextInterface, approved string, tokenId uint64) error {
	key, err := GetTokenApprovedKey(ctx, tokenId)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(key, []byte(approved))
}

// ApproveCouchdb 授予 approved 拥有 tokenId 的转移权力
func (nft *NFT) ApproveCouchdb(ctx contractapi.TransactionContextInterface, approved string, tokenId string) error {
	logger.Debugf("method = %s, approved = %s, tokenId = %s", "ApproveCouchdb", approved, tokenId)

	// 保证 owner 和 sender 一致
	sender, err := getSender(ctx)
	if err != nil {
		logger.Error(GetErrorStackf(err, "get sender error"))
		return errors.WithMessagef(err, "get sender error")
	}
	key, err := GetTokenIdKey(ctx, tokenId)
	if err != nil {
		logger.Error(GetErrorStackf(err, "create token key error, tokenId = %s", tokenId))
		return errors.WithMessagef(err, "create token key error, tokenId = %s", tokenId)
	}
	state, err := ctx.GetStub().GetState(key)
	if err != nil {
		logger.Error(GetErrorStackf(err, "could find any record in the ledger"))
		return errors.WithMessagef(err, "could find any record in the ledger")
	}
	nft721 := BasicNFT{}
	err = json.Unmarshal(state, &nft721)
	if err != nil {
		logger.Error(GetErrorStackf(err, "json unmarshal error, data = %s", string(state)))
		return errors.WithMessagef(err, "json unmarshal error, data = %s", string(state))
	}

	if nft721.Owner != sender {
		logger.Error(GetErrorStackf(nil, "token does not owned by sender, tokenId = %s, sender= %s ", tokenId, sender))
		return fmt.Errorf("token does not owned by sender, tokenId = %s, sender= %s ", tokenId, sender)
	}

	key1, err := GetTokenApprovedKeyCouchdb(ctx, tokenId)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(key1, []byte(approved))
}

/*
 * @Desc: owner 的所有 NFT 都可以由 operator  来控制
 * @Param:
 * @Return:
 */
func (nft *NFT) SetApprovalForAll(ctx contractapi.TransactionContextInterface, operator string, approved bool) error {
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

// GetApproved 返回 tokenId 的授权地址
func (nft *NFT) GetApproved(ctx contractapi.TransactionContextInterface, tokenId uint64) (string, error) {
	key, err := GetTokenApprovedKey(ctx, tokenId)
	if err != nil {
		return "", err
	}
	raw, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// GetApproved 返回 tokenId 的授权地址
func (nft *NFT) GetApprovedCouchdb(ctx contractapi.TransactionContextInterface, tokenId string) (string, error) {
	key, err := GetTokenApprovedKeyCouchdb(ctx, tokenId)
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
func (nft *NFT) IsApprovedForAll(ctx contractapi.TransactionContextInterface, owner string, operator string) (bool, error) {
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
func (nft *NFT) canTransfer(ctx contractapi.TransactionContextInterface, from string, tokenId uint64) bool {
	owner, err := nft.OwnerOf(ctx, tokenId)
	if err != nil {
		fmt.Printf("OwnerOf error: %s", err.Error())
		return false
	}
	if owner != from {
		fmt.Errorf("token does not owned by from, tokenId = %d, from= %s ", tokenId, from)
		return false
	}
	isOwner, err := checkSender(ctx, from)
	if err != nil {
		fmt.Printf("checkSender error: %s", err.Error())
		return false
	}
	// NFT 拥有者
	if isOwner {
		return true
	}
	// NFT 授权操作人
	approved, err := nft.GetApproved(ctx, tokenId)
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

// 判断是不是有操作 token 的权限, 数据库 couchdb
func (nft *NFT) canTransferCouchdb(ctx contractapi.TransactionContextInterface, from string, tokenId string) bool {
	key, err := GetTokenIdKey(ctx, tokenId)
	if err != nil {
		logger.Error(GetErrorStackf(err, "create token key error, tokenId = %s", tokenId))
		return false
	}
	state, err := ctx.GetStub().GetState(key)
	if err != nil {
		logger.Error(GetErrorStackf(err, "could find any record in the ledger"))
		return false
	}
	nft721 := BasicNFT{}
	err = json.Unmarshal(state, &nft721)
	if err != nil {
		logger.Error(GetErrorStackf(err, "json unmarshal error, data = %s", string(state)))
		return false
	}

	if nft721.Owner != from {
		logger.Error(GetErrorStackf(nil, "token does not owned by from, tokenId = %s, from= %s ", tokenId, from))
		return false
	}

	isOwner, err := checkSender(ctx, from)
	if err != nil {
		logger.Error(GetErrorStackf(err, "checkSender error"))
		return false
	}
	// NFT 拥有者
	if isOwner {
		return true
	}

	// NFT 授权操作人
	approved, err := nft.GetApprovedCouchdb(ctx, tokenId)
	if err != nil {
		logger.Error(GetErrorStackf(err, "GetApproved error"))
		return false
	}
	sender, _ := getSender(ctx)
	if approved == sender {
		return true
	}
	// 所有资产授权操作人
	approvedAll, err := nft.IsApprovedForAll(ctx, nft721.Owner, sender)
	if err != nil {
		logger.Error(GetErrorStackf(err, "IsApprovedForAll error"))
		return false
	}
	if approvedAll {
		return true
	}
	return false
}

// 判断是不是有操作 token 的权限
func (nft *NFT) canContractTransfer(ctx contractapi.TransactionContextInterface, from string, tokenId uint64, contractTokenId uint64) bool {
	owner, err := nft.OwnerOf(ctx, tokenId)
	if err != nil {
		fmt.Printf("OwnerOf error: %s", err.Error())
		return false
	}
	if owner != from {
		fmt.Errorf("token does not owned by from, tokenId = %d, from= %s ", tokenId, from)
		return false
	}
	isOwner, err := checkSender(ctx, from)
	if err != nil {
		fmt.Printf("checkSender error: %s", err.Error())
		return false
	}
	// NFT 拥有者
	if isOwner {
		return true
	}
	// NFT 授权操作人
	approved, err := nft.GetApproved(ctx, tokenId)
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

func (nft *NFT) delToken(ctx contractapi.TransactionContextInterface, from string, tokenId uint64) error {
	key, err := GetTokenOwnerKey(ctx, tokenId)
	if err != nil {
		return err
	}
	return ctx.GetStub().DelState(key)

}

func (nft *NFT) delDtoken(ctx contractapi.TransactionContextInterface, from string, tokenId uint64) error {
	key, err := GetTokenOwnerKey(ctx, tokenId)
	if err != nil {
		return err
	}
	return ctx.GetStub().DelState(key)

}

func (nft *NFT) increaseToken(ctx contractapi.TransactionContextInterface, to string) error {
	key, err := GetTokenCountByOwnerKey(ctx, to)
	if err != nil {
		return err
	}
	return nft.calcCount(ctx, key, true)
}

func (nft *NFT) decreaseToken(ctx contractapi.TransactionContextInterface, from string) error {
	key, err := GetTokenCountByOwnerKey(ctx, from)
	if err != nil {
		return err
	}
	return nft.calcCount(ctx, key, false)
}

// 将 nft token 与用户绑定
func (nft *NFT) addToken(ctx contractapi.TransactionContextInterface, to string, tokenId uint64) error {
	key, err := GetTokenOwnerKey(ctx, tokenId)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(key, []byte(to))
}

// 将 delegated token 和 parent token 绑定
func (nft *NFT) addDToken(ctx contractapi.TransactionContextInterface, to string, tokenId, dtokenId uint64) error {
	key, err := GetDtokenParentKey(ctx, dtokenId)
	if err != nil {
		return err
	}
	int64Str := strconv.FormatUint(tokenId, 10)
	err = ctx.GetStub().PutState(key, []byte(int64Str))
	if err != nil {
		return err
	}
	tokenOwnerKey, err := GetTokenOwnerKey(ctx, dtokenId)
	fmt.Printf("key = %s , owner = %s \n", tokenOwnerKey, to)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(tokenOwnerKey, []byte(to))
}

// 修改 nft token 数量
func (nft *NFT) calcCount(ctx contractapi.TransactionContextInterface, key string, increase bool) error {
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

func (nft *NFT) GetSender(ctx contractapi.TransactionContextInterface) (string, error) {
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
func (nft *NFT) SendNft721ToNft998(ctx contractapi.TransactionContextInterface, from string, tokenId uint64, parentContractName, childContractName string, childtokenId uint64) error {
	transfer := nft.canTransfer(ctx, from, childtokenId)
	if transfer {
		//_, err := TokenToOwner(ctx, tokenId)
		//if err != nil {
		//	return err
		//}

		//index, err := ChildTokenIndex(ctx, tokenId, childContractName, childtokenId)
		//if index!=0{
		//	return fmt.Errorf("Cannot send child token because it has already been received, tokenId = %d, childContractName = %s, childtokenId=%d",tokenId,childContractName,childtokenId)
		//}
		//childTokenIndexKey, err := GetChildTokenIndexKey(ctx, tokenId, childContractName, childtokenId)
		//if err != nil {
		//	return err
		//}
		//err = ctx.GetStub().PutState(childTokenIndexKey, []byte(fmt.Sprintf("%s", childtokenId)))
		//if err != nil {
		//	return err
		//}
		//childTokensKey, err := GetChildTokensKey(ctx, tokenId, childContractName, childtokenId)
		//if err != nil {
		//	return err
		//}
		//raw, err := ctx.GetStub().GetState(childTokensKey)
		//if err != nil {
		//	return err
		//}
		//if raw != nil {
		//	return fmt.Errorf("childtokenId = %d already owned by tokenId = %d ", childtokenId, tokenId)
		//}
		//
		//err = ctx.GetStub().PutState(childTokensKey, []byte(fmt.Sprintf("%d", childtokenId)))
		//if err != nil {
		//	return err
		//}
		//childTokenOwnerKey, err := GetChildTokenOwnerKey(ctx, childContractName, childtokenId)
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
		//response := ctx.GetStub().InvokeChaincode(parentContractName, util.ToChaincodeArgs("ReceiveNft721", strconv.FormatUint(tokenId, 10), childContractName, strconv.FormatUint(childtokenId, 10)), ctx.GetStub().GetChannelID())
		//if response.Status != 200 {
		//	return fmt.Errorf("nft7 发送到 nft998 失败，msg = %s, parentContractName = %s, tokenId = %d, childContractName = %s, childtokenId = %d", response.Message, parentContractName, tokenId, childContractName, childtokenId)
		//}
		if childContractName == "" {
			label, err := nft.TokenLabel(ctx, childtokenId)
			if err != nil {
				return err
			}
			childContractName = label
		}
		err := nft.ReceiveNft721(ctx, tokenId, childContractName, childtokenId)
		if err != nil {
			return err
		}

		err = nft.delToken(ctx, from, childtokenId)
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
	return fmt.Errorf("sender can not transfer, sender = %s, childtokenId = %d", sender, childtokenId)
}

/*
 * @Desc: 校验权限
 * @Param:
 * @Return:
 */
func (nft *NFT) IsAuthorized(ctx contractapi.TransactionContextInterface, tokenId string) error {
	// 首先进行权限的判断
	//owner, err := nft.OwnerOf(ctx, tokenId)
	//if err != nil {
	//	//fmt.Printf("get owner of tokenId=%d error: %s", tokenId, err.Error())
	//	return fmt.Errorf("get owner of tokenId=%d error: %s", tokenId, err.Error())
	//}
	//
	//transfer := nft.canTransfer(ctx, owner, tokenId)
	//if !transfer {
	//	return fmt.Errorf("token does not owned by from, tokenId = %d, from= %s ", tokenId, owner)
	//}
	return nil
}
