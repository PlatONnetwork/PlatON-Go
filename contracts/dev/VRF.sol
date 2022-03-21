// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/** ****************************************************************************
 * 调用PlatON内置合约生成VRF随机数
 */
contract VRF {
  // 内置合约地址
  address vrfPrecompiled = 0x6168499c0cFfCaCD319c818142124B7A15E857ab;

  // Assumes the subscription is funded sufficiently.
  // function requestRandomWords(uint32 numWords) internal returns (uint256[] memory) {
  //   bytes memory data = abi.encode(numWords);
  //   bytes memory returnValue = assemblyCall(data, vrfPrecompiled);

  //   uint256[] memory randomWords = abi.decode(returnValue, (uint256[]));
  //   return randomWords;
  // }

  function requestRandomWords(uint32 numWords) internal returns (uint256[] memory) {
    uint256[] memory randomWords = new uint256[](1);
    return randomWords;
  }

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
}
