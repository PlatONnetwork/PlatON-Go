pragma solidity ^0.5.13;
/**
 * 验证0.5.13版本允许public函数覆盖external函数
 *
 * @author hudenian
 * @dev 2019/12/25 11:09
 */
import "./AnimalExternal.sol";

contract PersonPublic is AnimalExternal {
    /**
     * 0.5.13 允许public函数覆盖external函数
     *
     */
    function birthDay() view public returns (string memory) {

      return "2020-12-15";
    }

}
