pragma solidity ^0.5.13;
/**
 *  事件验证
 * （1）event关键字声明事件验证
 * （2）indexed关键字定义事件索引验证
 * （3）emit关键字触发事件验证
 * （4）anonymous关键字定义匿名事件验证
 * （5）JavaScript API配合回调（使用 Web3 监听事件）事件验证
 * @author Albedo
 * @dev 2019/12/19
 **/
contract EventCallContract {
    //事件声明(事件名称以大写字母开头，以区别于函数)
    event Increment(address who);
    //添加索引(如果indexed声明索引超过3个，编译异常：超过3个indexed标记的元素)
    event Deposit(
        address indexed _from,
    //   bytes32 indexed _idl,
    //  bytes32 indexed _idm,
        uint _value
    );

    //匿名事件
    event Anonymous(uint256 _id) anonymous;
    //一般事件触发
    function emitEvent() public returns (uint256 count){
        //emit 事件触发
        emit Increment(msg.sender);
        count += 1;
    }

    function indexedEvent() public returns (uint256 count){
        //emit 事件触发
        emit Deposit(msg.sender, 12);
        count += 1;
    }

    function anonymousEvent() public returns (uint256 count){
        count += 1;
        //emit 事件触发
        emit Anonymous(count);
    }
    //外部：JavaScript API配合来回调--定义事件及触发事件——使用 Web3 监听事件

    event BoolEvent(bool result);
    //函数多事件监听验证
    function testBool() public{
        emit BoolEvent(false);
        emit BoolEvent(true);
        emit Increment(msg.sender);
        emit Anonymous(12);
        emit Deposit(msg.sender, 12);
    }
}