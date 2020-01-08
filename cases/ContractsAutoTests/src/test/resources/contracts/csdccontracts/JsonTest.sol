pragma solidity ^0.4.12;

import "./utillib/LibJson.sol";
import "./utillib/LibLog.sol";
import "./sysbase/OwnerNamed.sol";

contract JsonTest is OwnerNamed{
	using LibJson for *;

	string json = "";

	function JsonTest() {
		register("CsdcModule", "0.0.1.0", "JsonTest", "0.0.1.0");
	}

	function testPush(string _json) {
		LibJson.push(_json);
		// LibLog.log(_json.jsonRead(_key));
	}

	function testJsonRead(string _key) {
		LibLog.log(json.jsonRead(_key));
	}

	function testPop() {
		LibJson.pop();
	}
}