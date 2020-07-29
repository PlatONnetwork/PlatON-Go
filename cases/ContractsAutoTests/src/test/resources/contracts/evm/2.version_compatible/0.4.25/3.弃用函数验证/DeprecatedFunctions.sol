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
        address _to = "lax1eg6m0kg4gk802s9ducrgml30gn505ueuswqu73";
        address _from = "lax16f0dq2wqj0jkhjy3rgruge29qqxt7d7x4zvmg4"; // lat16f0dq2wqj0jkhjy3rgruge29qqxt7d7x6875x6
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
        address _to = "lax1eg6m0kg4gk802s9ducrgml30gn505ueuswqu73"; // lat1eg6m0kg4gk802s9ducrgml30gn505ueultjns7
        //suicide() 验证
        suicide(_to);
    }
}