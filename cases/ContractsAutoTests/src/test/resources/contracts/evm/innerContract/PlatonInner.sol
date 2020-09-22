pragma solidity ^0.5.13;

/*
 * network.platon.test.evm 跨合约调用 Platon系统合约
 * @author hudenian
 * @dev 2020/02/25
*/
contract PlatonInner {

    bytes returnValue;

    function assemblyCallppos(bytes memory data,address addr) public {
        uint256 len = data.length;
        uint retsize;
        bytes memory resval;
        assembly {
            if iszero(call(gas, addr, 0,  add(data, 0x20), len, 0, 0)) {
                invalid()
            }
            retsize := returndatasize()
        }
        resval = new bytes(retsize);
        assembly {
            returndatacopy(add(resval, 0x20), 0, returndatasize())
        }
        returnValue = resval;
    }

    function getReturnValue() public view returns(bytes memory ){
        return returnValue;
    }

}