pragma solidity ^0.4.12;
/**
* file CsdcBaseInterface.sol
* author yiyating
* time 2017-05-16
* desc the defination of CsdcBaseInterface
*/


contract CsdcBaseInterface {
    function findNameById(address _id) constant returns (uint);
    function hasTodo(address _userId) constant returns (bool);
}
