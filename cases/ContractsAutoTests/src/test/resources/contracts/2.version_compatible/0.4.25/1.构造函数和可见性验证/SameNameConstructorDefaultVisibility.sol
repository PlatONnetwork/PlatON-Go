pragma solidity ^0.4.25;

/**
 * 1.0.4.25版本验证同名函数构造函数定义，可见性未声明（默认public可见性）验证；
 * 2.函数可见性非强制声明验证。5种函数
 * (1) 函数默认可见性声明
 * (2) 函数public可见性声明
 * (3) 函数external可见性声明
 * (4) 函数internal可见性声明
 * (5) 函数private可见性声明
 * @author Albedo
 * @dev 2019/12/19
 **/
contract SameNameConstructorDefaultVisibility {
    uint256 conParam;
    //同名函数构造函数，默认可见性
    function SameNameConstructorVisibility(uint256 param) {
        conParam = param;
    }
    //默认可见性
    function defaultVisibility(uint256 param) view returns (uint256){
        conParam = param;
        return conParam;
    }

    //public可见性
    function publicVisibility(uint256 param) public view returns (uint256){
        conParam = param;
        return conParam;
    }
    //external可见性
    function externalVisibility(uint256 param) external view returns (uint256){
        conParam = param;
        return conParam;
    }
    //internal可见性
    function internalVisibility(uint256 param) internal view returns (uint256){
        conParam = param;
        return conParam;
    }
    //private可见性
    function privateVisibility(uint256 param) private returns (uint256){
        conParam = param;
        return conParam;
    }

    //private可见性验证
    function privateVisibilityCheck(uint256 param) public view returns (uint256){
        return privateVisibility(param);
    }

    //private可见性验证
    function internalVisibilityCheck(uint256 param) public view returns (uint256){
        return internalVisibility(param);
    }

}