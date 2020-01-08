/**
* @file SystemTest.sol
* @author liaoyan
* @time 2017-04-17
* @desc the defination of SystemTest contract
*/

pragma solidity ^0.4.2 ;

import "../sysbase/OwnerNamed.sol";

contract SystemTest is OwnerNamed {
    
    using LibString for *;
    using LibInt for *;

    mapping(address=>uint) counter;

    function SystemTest() {
        register("SystemModuleManager","0.0.1.0","SystemTest","0.0.1.0");
    }

    function count() public {
        counter[msg.sender]++;
    }

    function resetCounter() public {
        counter[msg.sender] = 0;
    }

    function getCounter() public constant returns(string _json) {
        _json = _json.concat("{");
        _json = _json.concat( msg.sender.toKeyValue("addr"), "," );
        _json = _json.concat( counter[msg.sender].toKeyValue("counter"));
        _json = _json.concat("}");
    }
}
