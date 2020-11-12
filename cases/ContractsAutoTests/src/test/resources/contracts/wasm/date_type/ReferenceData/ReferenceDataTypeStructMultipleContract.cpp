#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;

/**
 * @author qudong
 * 合约引用类型结构体(struct类型)：允许存储不同类型的数据项

 * 测试验证功能点：
 * 2、定义struct类型包含引用类型及基本类型
 * */

struct group{
    public:
       uint64_t groupID;
       std::string groupName;
       std::array<std::string,10> subGroupNameArray;
       std::map<std::string,std::string> subGroupScoreMap;
       group(){}
       PLATON_SERIALIZE(group,(groupID)(groupName)(subGroupNameArray)(subGroupScoreMap))
};

CONTRACT ReferenceDataTypeStructMultipleContract : public platon::Contract{

    private:
      platon::StorageType<"group"_n,group> storage_struct_group;

    public:
        ACTION void init(){}

        /**
         * 2、定义struct类型包含引用类型及基本类型
         */
          //1)、赋值基本类型
        ACTION void setGroupValue(const std::string &myGroupName,const uint64_t &myGroupId){
             storage_struct_group.self().groupID = myGroupId;
             storage_struct_group.self().groupName = myGroupName;
        }
       CONST std::string getGroupName(){
              return storage_struct_group.self().groupName;
        }
        //2)、赋值引用类型
        ACTION void setGroupArrayValue(const std::string &oneValue,const std::string &twoValue){
             storage_struct_group.self().subGroupNameArray[0] = oneValue;
             storage_struct_group.self().subGroupNameArray[1] = twoValue;
        }
        CONST std::string getGroupArrayIndexValue(const uint32_t &index){
             return storage_struct_group.self().subGroupNameArray[index];
        }
};

PLATON_DISPATCH(ReferenceDataTypeStructMultipleContract, (init)(setGroupValue)(getGroupName)
               (setGroupArrayValue)(getGroupArrayIndexValue))

