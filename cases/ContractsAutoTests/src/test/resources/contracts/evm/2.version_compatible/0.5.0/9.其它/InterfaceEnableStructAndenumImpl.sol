pragma solidity ^0.5.0;
/**
 * 09-其它
 * 10-验证0.5.0接口允许定义结构体与枚举(0.4.25版本会报错)
 *
 * @author hudenian
 * @dev 2019/12/24 09:57
 */

import "./InterfaceEnableStructAndenum.sol";

contract  InterfaceEnableStructAndenumImpl is InterfaceEnableStructAndenum{

    //构建Product结构体传入ProductCondition，Used则返回1，ProductCondition.New返回0
    function getProductCondition() public pure returns (ProductCondition){
        Product memory product = Product(ProductCondition.Used);
        return product.condition;
    }

}