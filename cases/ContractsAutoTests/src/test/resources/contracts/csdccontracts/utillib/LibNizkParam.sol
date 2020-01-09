pragma solidity ^0.4.12;

import "../utillib/LibString.sol";
import "../utillib/LibInt.sol";

library LibNizkParam {
    using LibInt for *;
    using LibString for *;
    using LibNizkParam for *;

    struct NizkParam{
        string      cipher1;
        string      cipher2;
        string      pais;
        string      balapubcipher;
        string      traapubcipher;
        string      trabpubcipher;
        string      apukkey;
        string      bpukkey;
        string      nizkpp;
    }

    function toJson(NizkParam storage _self) internal returns(string _strjson) {

        _strjson = "{";
        _strjson = _strjson.concat(_self.cipher1.toKeyValue("cipher1"), ",");
        _strjson = _strjson.concat(_self.cipher2.toKeyValue("cipher2"), ",");
        _strjson = _strjson.concat(_self.pais.toKeyValue("pais"), ",");
        _strjson = _strjson.concat(_self.balapubcipher.toKeyValue("balapubcipher"), ",");
        _strjson = _strjson.concat(_self.traapubcipher.toKeyValue("traapubcipher"), ",");
        _strjson = _strjson.concat(_self.trabpubcipher.toKeyValue("trabpubcipher"), ",");
        _strjson = _strjson.concat(_self.apukkey.toKeyValue("apukkey"), ",");
        _strjson = _strjson.concat(_self.bpukkey.toKeyValue("bpukkey"), ",");
        _strjson = _strjson.concat(_self.nizkpp.toKeyValue("nizkpp"), "}");
    }

    function jsonParse(NizkParam storage _self, string _strjson) internal returns(bool) {
        _self.cipher1 = _strjson.getStringValueByKey("cipher1");
        _self.cipher2 = _strjson.getStringValueByKey("cipher2");
        _self.pais = _strjson.getStringValueByKey("pais");
        _self.balapubcipher = _strjson.getStringValueByKey("balapubcipher");
        _self.traapubcipher = _strjson.getStringValueByKey("traapubcipher");
        _self.trabpubcipher = _strjson.getStringValueByKey("trabpubcipher");
        _self.apukkey = _strjson.getStringValueByKey("apukkey");
        _self.bpukkey = _strjson.getStringValueByKey("bpukkey");
        _self.nizkpp = _strjson.getStringValueByKey("nizkpp");

        return true;
    }

    function reset(NizkParam storage _self) internal{
        _self.cipher1 = "";
        _self.cipher2 = "";
        _self.pais = "";
        _self.balapubcipher = "";
        _self.traapubcipher = "";
        _self.trabpubcipher = "";
        _self.apukkey = "";
        _self.bpukkey = "";
        _self.nizkpp = "";
    }
}
