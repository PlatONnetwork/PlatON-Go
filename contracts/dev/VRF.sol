// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * 调用PlatON内置合约生成VRF随机数
 */
contract VRF {

  error InvalidRandomWords(uint32 numWords, uint256 returnValueLength);

  // VrfInnerContract 内置合约地址
  address vrfInnerContractAddr = 0x3000000000000000000000000000000000000001;

  // 32个字节表示uint256
  uint32 base = 32;

  /**
   * 调用 VrfInnerContract 内置合约生成VRF随机数
   * @param numWords 随机数的个数
   */
  function requestRandomWords(uint32 numWords) internal returns (uint256[] memory) {
    bytes memory data = abi.encode(numWords);
    bytes memory returnValue = assemblyCall(data, vrfInnerContractAddr);

    if (numWords * base != returnValue.length) {
        revert InvalidRandomWords(
            numWords,
            returnValue.length
        );
    }

    uint256[] memory randomWords = new uint256[](numWords);
    for(uint i = 0; i < numWords; i++) {
        uint start = i * base;
        randomWords[i] = sliceUint(returnValue, start);
    }

    return randomWords;
  }

  /**
   * delegatecall 合约
   * @param data 合约input data
   * @param addr 合约地址
   */
    function assemblyCall(bytes memory data, address addr) internal returns (bytes memory) {
        uint256 len = data.length;
        uint retsize;
        bytes memory resval;
        assembly {
            let result := delegatecall(gas(), addr, add(data, 0x20), len, 0, 0)
            retsize := returndatasize()
        }
        resval = new bytes(retsize);
        assembly {
            returndatacopy(add(resval, 0x20), 0, returndatasize())
        }
        return resval;
    }

    function sliceUint(bytes memory bs, uint start) internal pure returns (uint256) {
        require(bs.length >= start + 32, "slicing out of range");
        uint256 x;
        assembly {
            x := mload(add(bs, add(0x20, start)))
        }
        return x;
    }
}
