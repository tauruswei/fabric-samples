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

interface ERC998ERC721TopDown {

    /// @notice Get the root owner of tokenId.
    /// @param _tokenId The token to query for a root owner address
    /// @return rootOwner The root owner at the top of tree of tokens and ERC998 magic value.
    function rootOwnerOf(uint256 _tokenId) external view returns (bytes32 rootOwner);

    /// @notice Get the root owner of a child token.
    /// @param _childContract The contract address of the child token.
    /// @param _childTokenId The tokenId of the child.
    /// @return rootOwner The root owner at the top of tree of tokens and ERC998 magic value.
    function rootOwnerOfChild(address _childContract, uint256 _childTokenId) external view returns (bytes32 rootOwner);

    /// @notice Get the parent tokenId of a child token.
    /// @param _childContract The contract address of the child token.
    /// @param _childTokenId The tokenId of the child.
    /// @return parentTokenOwner The parent address of the parent token and ERC998 magic value
    /// @return parentTokenId The parent tokenId of _tokenId
    function ownerOfChild(address _childContract, uint256 _childTokenId) external view returns (bytes32 parentTokenOwner, uint256 parentTokenId);

    /// @notice A token receives a child token
    /// @param _operator The address that caused the transfer.
    /// @param _from The owner of the child token.
    /// @param _childTokenId The token that is being transferred to the parent.
    /// @param _data Up to the first 32 bytes contains an integer which is the receiving parent tokenId.
    function onERC721Received(address _operator, address _from, uint256 _childTokenId, bytes _data) external returns (bytes4);

    /// @notice Transfer child token from top-down composable to address.
    /// @param _fromTokenId The owning token to transfer from.
    /// @param _to The address that receives the child token
    /// @param _childContract The ERC721 contract of the child token.
    /// @param _childTokenId The tokenId of the token that is being transferred.
    function transferChild(uint256 _fromTokenId, address _to, address _childContract, uint256 _childTokenId) external;

    /// @notice Transfer child token from top-down composable to address.
    /// @param _fromTokenId The owning token to transfer from.
    /// @param _to The address that receives the child token
    /// @param _childContract The ERC721 contract of the child token.
    /// @param _childTokenId The tokenId of the token that is being transferred
    function safeTransferChild(uint256 _fromTokenId, address _to, address _childContract, uint256 _childTokenId) external;/// @notice Transfer child token from top-down composable to address.

    /// @param _fromTokenId The owning token to transfer from.
    /// @param _to The address that receives the child token
    /// @param _childContract The ERC721 contract of the child token.
    /// @param _childTokenId The tokenId of the token that is being transferred.
    /// @param _data Additional data with no specified format
    function safeTransferChild(uint256 _fromTokenId, address _to, address _childContract, uint256 _childTokenId, bytes _data) external;

    /// @notice Transfer bottom-up composable child token from top-down composable to other ERC721 token.
    /// @param _fromTokenId The owning token to transfer from.
    /// @param _toContract The ERC721 contract of the receiving token
    /// @param _toToken The receiving token
    /// @param _childContract The bottom-up composable contract of the child token.
    /// @param _childTokenId The token that is being transferred.
    /// @param _data Additional data with no specified format
    function transferChildToParent(uint256 _fromTokenId, address _toContract, uint256 _toTokenId, address _childContract, uint256 _childTokenId, bytes _data) external;

    /// @notice Get a child token from an ERC721 contract.
    /// @param _from The address that owns the child token.
    /// @param _tokenId The token that becomes the parent owner
    /// @param _childContract The ERC721 contract of the child token
    /// @param _childTokenId The tokenId of the child token

    // getChild function enables older contracts like cryptokitties to be transferred into a composable
    // The _childContract must approve this contract. Then getChild can be called.
    function getChild(address _from, uint256 _tokenId, address _childContract, uint256 _childTokenId) external;
}
interface ERC998ERC721TopDownEnumerable {
    function totalChildContracts(uint256 _tokenId) external view returns (uint256);
    function childContractByIndex(uint256 _tokenId, uint256 _index) external view returns (address childContract);
    function totalChildTokens(uint256 _tokenId, address _childContract) external view returns (uint256);
    function childTokenByIndex(uint256 _tokenId, address _childContract, uint256 _index) external view returns (uint256 childTokenId);
}

*/
// ERC998ERC721TopDown 仿以太坊 ERC998ERC721TopDown 协议接口  https://eips.ethereum.org/EIPS/eip-998
type ERC998ERC721TopDown interface {
	RootOwnerOf(ctx contractapi.TransactionContext, tokenId uint64) (string, error)
	RootOwnerOfChild(ctx contractapi.TransactionContext, childContractName string, childTokenId uint64) (string, error)
	OwnerOfChild(ctx contractapi.TransactionContext, childContractName string, childTokenId uint64) (string, uint64, error)
	OnERC721Received(ctx contractapi.TransactionContext, operator, from string, childTokenId uint64, data string) (string, error)
	TransferChild(ctx contractapi.TransactionContext, fromTokenId uint64, to, childContractName string, childTokenId uint64) error
	//SafeTransferChild(ctx contractapi.TransactionContext,fromTokenId uint64,to,childContractName string,childTokenId uint64)(error)
	SafeTransferChild(ctx contractapi.TransactionContext, fromTokenId uint64, to, childContractName string, childTokenId uint64, data string) error
	TransferChildToParent(ctx contractapi.TransactionContext, fromTokenId uint64, to, childContractName string, childTokenId uint64, data string) error
	GetChild(ctx contractapi.TransactionContext, from string, tokenId uint64, childContractName string, childTokenId uint64) error
}
type ERC998ERC721TopDownEnumerable interface {
	TotalChildContracts(tokenId uint64) (uint64, error)
	ChildContractByIndex(tokenId, index uint64) (childContractName string, err error)
	TotalChildTokens(tokenId uint64, childContractName string) (uint64, error)
	ChildTokenByIndex(tokenId uint64, childContractName string, index uint64) (childTokenId uint64, err error)
}

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
