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


contract ContractPayable{
    function () external payable {} //回退函数（没有名字、没有参数、没有返回值）
}