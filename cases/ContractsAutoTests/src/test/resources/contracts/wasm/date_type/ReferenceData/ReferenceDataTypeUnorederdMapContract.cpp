#define TESTNET
#include <platon/platon.hpp>
#include <string>
#include <unordered_map>
using namespace std;
using namespace platon;



/**
 * @author qudong
 * 合约引用类型(std::unorederd_map类型)：无序map容器,
 * 关联性：是一个关联容器，其中的元素根据键来引用，而不是根据索引引用
 * 无序性：元素不会根据其键值或映射值按任何特定顺序排序，而是根据哈希值组织
 * 唯一性：元素的键是唯一的
 *
 * */

CONTRACT ReferenceDataTypeUnorderMapContract : public platon::Contract{
    private:
      platon::StorageType<"unordermap"_n, std::unordered_map<std::string,std::string>> storage_string_map;
    public:
      ACTION void init(){}
      // 1)、unorderMap类型初始化赋值，编译异常，暂不支持此类型
      ACTION void init_unorder_map(){
            storage_string_map.self() = {{"apple","red"},{"lemon","yellow"}};
      }
};
PLATON_DISPATCH(ReferenceDataTypeUnorderMapContract,(init)(init_unorder_map))
