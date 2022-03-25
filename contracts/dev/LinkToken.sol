// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/LinkTokenInterface.sol";

contract LinkTokenInter is LinkTokenInterface
{
  /**
   * @inheritdoc LinkTokenInterface
   */
  function allowance(address owner, address spender) external override view returns (uint256 remaining){
    return uint256(0);
  }

  function approve(address spender, uint256 value) external override returns (bool success){
    return true;
  }

  function balanceOf(address owner) external override view returns (uint256 balance){
    return uint256(0);
  }

  function decimals() external override view returns (uint8 decimalPlaces){
    return uint8(0);
  }

  function decreaseApproval(address spender, uint256 addedValue) external override returns (bool success){
    return true;
  }

  function increaseApproval(address spender, uint256 subtractedValue) external override {

  }

  function name() external override view returns (string memory tokenName){
    return "";
  }

  function symbol() external override view returns (string memory tokenSymbol){
    return "";
  }

  function totalSupply() external override view returns (uint256 totalTokensIssued){
    return uint256(0);
  }

  function transfer(address to, uint256 value) external override returns (bool success){
    return true;
  }

  function transferAndCall(
    address to,
    uint256 value,
    bytes calldata data
  ) external override returns (bool success){
    return true;
  }

  function transferFrom(
    address from,
    address to,
    uint256 value
  ) external override returns (bool success){
    return true;
  }
}
