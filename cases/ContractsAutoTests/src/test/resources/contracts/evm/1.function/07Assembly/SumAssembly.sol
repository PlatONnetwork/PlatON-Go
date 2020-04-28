pragma solidity 0.5.13;
/**
 1.验证内联汇编在库中的使用
 2.验证汇编的操作码add，mul等
 * @author liweic
 * @dev 2020/01/08 14:30
 */
library Sum {
    function sumUsingInlineAssembly(uint[] memory _data) public pure returns (uint o_sum) {
        for (uint i = 0; i < _data.length; ++i) {
            assembly {
                o_sum := add(o_sum, mload(add(add(_data, 0x20), mul(i, 0x20))))
            }
        }
    }
}
contract SumAssembly {
    uint[] data;

    constructor() public {
        data.push(1);
        data.push(2);
        data.push(3);
        data.push(4);
        data.push(5);
    }
    function sum() external view returns(uint){
        return Sum.sumUsingInlineAssembly(data);
    }
}