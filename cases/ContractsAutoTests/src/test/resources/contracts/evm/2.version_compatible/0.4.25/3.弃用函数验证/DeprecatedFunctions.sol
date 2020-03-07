pragma solidity ^0.4.22;

/**
 * 验证0.5.0版本弃用但0.4.25版本仍生效函数
 * 1.callcode()（0.5.0版本已弃用，使用delegatecall()函数代替） 验证
 * 2.suicide() （0.5.0版本已弃用，使用selfdestruct()函数替用）验证
 * 3.sha3() （0.5.0版本已弃用，使用keccak256()函数代替）验证
 * 4.throw （0.5.0版本已弃用，使用异常函数验证）验证
 * @author Albedo
 * @dev 2019/12/19
 **/
contract DeprecatedFunctions {
    function functionCheck() public view returns (bool, bytes32){
        address _to = 0xca35b7d915458ef540ade6068dfe2f44e8fa733c;
        address _from = 0xd25ed029c093e56bc8911a07c46545000cbf37c6;
        //callcode() 验证
        bool callcodeResult = _from.callcode(_to);
        //sha3() 验证
        bytes32 strbytes = sha3("wangzhangxiong");
        //返回 false，0x49a40597e20d39bf568fe3296189c2d963951969c41761aabb19c402f8231695
        return (callcodeResult, strbytes);

    }

    function throwCheck(bool param) public view returns (bool) {
        //throw 验证
        if (!param) {throw;}
        return param;
    }

    function kill() public {
        address _to = 0xca35b7d915458ef540ade6068dfe2f44e8fa733c;
        //suicide() 验证
        suicide(_to);
    }
}