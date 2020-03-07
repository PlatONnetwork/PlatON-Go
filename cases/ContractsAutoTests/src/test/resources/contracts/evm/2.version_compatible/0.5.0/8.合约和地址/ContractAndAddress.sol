pragma solidity ^0.5.0;
/**
 * 08-合约和地址
 * 1-contract合约类型不再包括 address类型的成员函数，
 * 必须显式转换成 address地址类型才能使用 send(), transfer()， 
 * balance等与之相关的成员函数/变量成员
 * 2-address地址类型细分为address和address payable，
 * 只有address payable可以使用transfer()，send()函数，
 * address payable类型可以直接转换为 address类型，反之不能。
 * 但是 address x可以通过 address(uint160(x))，强制转换成 address payable类型。
 * 如果 contract A不具有 payable的fallback函数, 那么 address(A)是 address类型。
 * 如果 contract A具有 payable的fallback函数, 那么 address(A)是 address payable类型
 * 3-msg.sender属于address payable类型
 *
 * @author hudenian
 * @dev 2019/12/19 09:57
 */

import "./ContractAdress.sol";
import "./ContractPayable.sol";


contract ContractAndAddress{
    address payable to_pay;
    address to;
    address payableToAddress;
    address payable addressForcePayable;

    function payableOrNot() external {
        bool result;
        to_pay = address(new ContractPayable());      // 包含回退函数的合约contractPayable是一个 `address payable`
        to = address(new ContractAdress());          // `address payable`  可以显示转换成address
        //   to_pay.transfer(1);             // `transfer()` 是 `address payable` 的成员函数
        //   result = to_pay.send(1);        // `send()` 是 `address payable` 的成员函数
        to = address(new ContractAdress());          // 不包含声明为payable回退函数的合约可以显示转换`address`
        // to_pay = address(new B());    // 不包含声明为payable回退函数的合约不可以显示转换成                                              // `address payable`
        // to.transfer(1 ether);        // `transfer()` 不是`address`成员函数
        // result = to.send(1 ether);   // `send()` 不是`address`成员函数（写法不对）
        //address john;
        payableToAddress = to_pay;                  // Right,`address payable` can directly convert to `address`

        addressForcePayable = address(uint160(to));    // Right,`address` can forcibly convert to `address payable`

    }

    //获取payable类型的地址
    function getNonalPayableAddress() view public returns (address){
        return to_pay;
    }

    //获取非payable类型的地址
    function getNonalContractAddress() view public returns (address){
        return to;
    }

    //获取payable转换成address的地址
    function getPayableToAddress() view public returns (address){
        return payableToAddress;
    }

    //获取address强制转换成payable的地址
    function  getAddressToPayable() view public returns (address){
        return addressForcePayable;
    }
}