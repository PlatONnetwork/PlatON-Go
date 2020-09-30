#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
* 验证计算方法在合约里的实现
* 计算两个日期之间相差的月份
* @author liweic
*/

CONTRACT ComputeDate : public platon::Contract{
	public:
    ACTION void init(){}

    CONST int MonthsBetween2Date(const std::string& date1, const std::string& date2){
       //取出日期中的年月日
        int year1, month1;
        int year2, month2;
        year1 = atoi((date1.substr(0,4)).c_str());
    	month1 = atoi((date1.substr(5,2)).c_str());
    	year2 = atoi((date2.substr(0,4)).c_str());
    	month2 = atoi((date2.substr(5,2)).c_str());
        return 12*(year2-year1)+ month2 - month1;
    }
};

PLATON_DISPATCH(ComputeDate, (init)(MonthsBetween2Date))