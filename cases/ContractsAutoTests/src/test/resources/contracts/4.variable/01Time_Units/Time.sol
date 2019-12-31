pragma solidity 0.5.13;
/**
* 对 Time 单位进行测试
* Time的单位关键字有seconds, minutes, hours, days, weeks, years，换算格式如下：
* 1 == 1 seconds
* 1 minutes == 60 seconds
* 1 hours == 60 minutes
* 1 days == 24 hours
* 1 weeks == 7 days
* 1 years == 365 days   0.5.0被弃用,故在0.5.13版本无法使用编译不通过
* 默认缺省单位是秒
* @author liweic
* @dev 2019/12/26 18:10
*/

contract Time {

    // 定义全局变量
    uint time;

    //返回当前时间的Unix时间戳和当前块的Unix时间戳差值,结果为0
    function testimeDiff() public returns (uint256){
        return block.timestamp - now;
    }
    
    //时间的默认缺省单位是秒
    function testTime() public{
      time = 100000000;
    }

    //时间加1秒结果100000001
    function tSeconds() public view returns(uint){
      return time + 1 seconds; 
    }
    
    //时间加1分钟结果100000060
    function tMinutes() public view returns(uint){
      return time + 1 minutes;
    }
    
    //时间加一小时结果100003600
    function tHours() public view returns(uint){
      return time + 1 hours;
    }
    
    //时间加一周结果为100604800
    function tWeeks() public view returns(uint){
      return time + 1 weeks;
    }

}