pragma solidity ^0.5.0;
/**
 * 09-其它
 * 10-验证0.5.0接口允许定义结构体与枚举(0.4.25版本会报错)
 *
 * @author hudenian
 * @dev 2019/12/24 09:57
 */

interface InterfaceEnableStructAndenum {

    //结构体
    struct Product{
        ProductCondition condition;
    }

    //枚举
    enum ProductCondition{
        New,
        Used
    }
}