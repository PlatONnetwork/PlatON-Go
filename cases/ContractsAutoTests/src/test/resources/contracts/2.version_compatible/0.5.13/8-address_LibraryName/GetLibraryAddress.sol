pragma solidity 0.5.13;
/**
 * 1. address(LibraryName)：获取链接库的地址
 *
 * @author hudenian
 * @dev 2019/12/25 09:57
 *
 */

import "./UserLibrary.sol";

contract GetLibraryAddress{

    using UserLibrary for *;

    address userLibAddress;

    function setUserLibAddress() public{
        //根据library库的名称获取地址
        userLibAddress = address(UserLibrary);
    }

    function getUserLibAddress() view public returns(address adr){
        return userLibAddress;
    }
}