pragma solidity ^0.4.12;


contract IRegisterApplyManager {

    //用户注册申请
    function insert(string _json) public;

    //更新注册申请
    function update(string _json) public;

    //用户注册审核
    function audit(string _json) public;

    //获取自动审核开关
    function getAutoAuditSwitch() constant public returns (uint);

    //更新自动审核开关
    function updateAutoAuditSwitch(uint code) public;

    //根据ID查询申请信息
    function findById(string _applyId) constant public returns (string);

    //根据uuid查询申请信息
     function findByUuid(string _uuid) constant public returns (string);

    //分页查询申请列表
    function listByCondition(string _name, string _mobile, uint _certType, uint _pageSize, uint _pageNo, string _auditStatus, uint _accountStatus) constant public returns (string);
}
