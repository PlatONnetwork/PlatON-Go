pragma solidity ^0.5.0;
/**
 * 10-链上函数
 * 1- 0.5.0版本函数 keccak256() 代替 0.4.25版本函数 sha3()
 *
 * @author hudenian
 * @dev 2019/12/19 09:57
 *
 */

contract Sha3AndKeccake256 {

    uint256 afterSha256value;

    /**
     * keccak256a
     */
    function keccak(string memory sha256value) public{
        afterSha256value = uint256(keccak256(abi.encodePacked(sha256value, sha256value)));
    }

    function getKeccak256Value() view public returns(uint256){
        return afterSha256value;
    }

    /**
     * sha3 只能在0.4.x版本上使用
     */
    // function sha(string memory sha256value)  public returns(uint256){
    //      return uint256(sha3(sha256value, "c"));
    // }
}
