pragma solidity ^0.5.13;
/**
 * 添加对具有string或bytes键类型的mapping的获取器的支持
 * Add support for getters of mappings with string or bytes key types
 *
 * @author hudenian
 * @dev 2019/12/25 11:09
 */

contract StringmappingSupport {


    mapping(string =>string) amap;
    mapping(bytes9 => string) bmap;

    /**
     * 支持string为key的mapping数据存储
     */
    function setStringmapValue(string memory _key,string memory _value) public{
        amap[_key] = _value;
    }

    /**
     * 支持string为key的mapping数据获取
     */
    function getStringmapValue(string memory _key) public view returns(string memory){
        return amap[_key];
    }

    /**
     * 支持byte为key的mapping数据存储
     */
    function setByte32mapValue(bytes9 _key,string memory _value) public{
        bmap[_key] = _value;
    }

    /**
     * 支持byte为key的mapping数据获取
     */
    function getByte32mapValue(bytes9 _key) public view returns( string memory){
        return bmap[_key];
    }
}
