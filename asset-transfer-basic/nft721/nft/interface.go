package nft

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

/*
interface ERC1155 {
	function safeTransferFrom(address _from, address _to, uint256 _id, uint256 _value, bytes calldata _data) external;
	function safeBatchTransferFrom(address _from, address _to, uint256[] calldata _ids, uint256[] calldata _values, bytes calldata _data) external;
	function balanceOf(address _owner, uint256 _id) external view returns (uint256);
	function balanceOfBatch(address[] calldata _owners, uint256[] calldata _ids) external view returns (uint256[] memory);
	function setApprovalForAll(address _operator, bool _approved) external;
	function isApprovedForAll(address _owner, address _operator) external view returns (bool);
}
*/

/*
interface ERC721 {
	function balanceOf(address _owner) external view returns (uint256);
	function ownerOf(uint256 _tokenId) external view returns (address);
	function safeTransferFrom(address _from, address _to, uint256 _tokenId, bytes data) external payable;
	function safeTransferFrom(address _from, address _to, uint256 _tokenId) external payable;
	function transferFrom(address _from, address _to, uint256 _tokenId) external payable;
	function approve(address _approved, uint256 _tokenId) external payable;
	function setApprovalForAll(address _operator, bool _approved) external;
	function getApproved(uint256 _tokenId) external view returns (address);
	function isApprovedForAll(address _owner, address _operator) external view returns (bool);
}

interface ERC721Metadata {
    function name() external view returns (string _name);
    function symbol() external view returns (string _symbol);
    function tokenURI(uint256 _tokenId) external view returns (string);
}

interface ERC721Enumerable  {
    function totalSupply() external view returns (uint256);
    function tokenByIndex(uint256 _index) external view returns (uint256);
    function tokenOfOwnerByIndex(address _owner, uint256 _index) external view returns (uint256);
}

*/

// ERC721 仿以太坊 erc721 协议 ERC721 接口  https://eips.ethereum.org/EIPS/eip-721
type ERC721 interface {
	BalanceOf(ctx contractapi.TransactionContextInterface, owner string) (uint64, error)
	OwnerOf(ctx contractapi.TransactionContextInterface, tokenID uint64) (string, error)
	SafeTransferFrom(ctx contractapi.TransactionContextInterface, from string, to string, tokenID uint64, data []byte) error
	TransferFrom(ctx contractapi.TransactionContextInterface, from string, to string, tokenID uint64) error
	Approve(ctx contractapi.TransactionContextInterface, approved string, tokenID uint64) error
	SetApprovalForAll(ctx contractapi.TransactionContextInterface, operator string, approved bool) error
	GetApproved(ctx contractapi.TransactionContextInterface, tokenID uint64) (string, error)
	IsApprovedForAll(ctx contractapi.TransactionContextInterface, owner string, operator string) (bool, error)
}

// ERC721Metadata 仿以太坊 erc721 协议 ERC721Metadata 接口 https://eips.ethereum.org/EIPS/eip-721
type ERC721Metadata interface {
	Name(ctx contractapi.TransactionContextInterface) string
	Symbol(ctx contractapi.TransactionContextInterface) string
	TokenURI(ctx contractapi.TransactionContextInterface, tokenID uint64) (string, error)
}

type ERC721Enumerable interface {
	TotalSupply() (uint64, error)
	TokenByIndex(index uint64) (uint64, error)
	TokenOfOwnerByIndex(owner string, index uint64) (uint64, error)
}