#define TESTNET
#include <platon/platon.hpp>
#include <string>
#include <platon/db/multi_index.hpp>
using namespace std;
using namespace platon;

/**
 * @author liweic
 * 多索引(MultiIndex类型)
 * */

struct Member {
    std::string name;
    uint8_t age;
    uint8_t sex;
    uint64_t $seq_;
    std::string Name() const { return name; }
    uint8_t Age() const { return age; }
    Member(){}
    Member(const std::string &my_name,uint8_t &my_age,uint8_t &my_sex):name(my_name),age(my_age),sex(my_sex){}
    PLATON_SERIALIZE(Member, (name)(age)(sex))
};

CONTRACT MultiIndexContract : public platon::Contract{
    private:
       platon::db::MultiIndex<"student"_n,
                  Member,
                  platon::db::IndexedBy<"indexname"_n,platon::db::IndexMemberFun<Member,std::string,&Member::Name,platon::db::IndexType::UniqueIndex>>,
                  platon::db::IndexedBy<"indexage"_n,platon::db::IndexMemberFun<Member,uint8_t,&Member::Age,platon::db::IndexType::NormalIndex>>
                  >  member_student;
    public:

        ACTION void init(){}

         //1、定义多索引类型,多索引插入数据emplace
        ACTION void addInitMultiIndex(const std::string &my_name,uint8_t &my_age,uint8_t &my_sex)
        {
            member_student.emplace(([&](auto &m) {
                                        m.age = my_age;
                                        m.name = my_name;
                                        m.sex = my_sex;
                                        }));
         }

        //2、验证：cbegin()多索引迭代器起始位置
        CONST bool getMultiIndexCbegin(const std::string &my_name)
        {
            auto iter = member_student.find<"indexname"_n>(my_name);
            if(iter == member_student.cbegin()){
                return true;
            }
            return false;
        }

         //3)、验证：cend()多索引迭代器结束位置
        CONST bool getMultiIndexCend(uint8_t &my_sex)
        {
            auto iter = member_student.find<"indexname"_n>(my_sex);
            if(iter == member_student.cend()){
                return true;
            }
            return false;
        }

        //4)、验证：count获取与索引值对应的数据的数量
        CONST uint8_t getMultiIndexCount(const uint8_t &my_age)
        {
            auto count = member_student.count<"indexage"_n>(my_age);
            return count;
        }

        //5)、验证：find多索引取值
        CONST std::string getMultiIndexFind(const std::string &my_name)
        {
              auto iter = member_student.find<"indexname"_n>(my_name);
              return iter->name;
        }

        //6)、验证：get_index获取非唯一索引的索引对象
        CONST bool getMultiIndexIndex(const uint8_t &my_age)
        {
            auto index = member_student.get_index<"indexage"_n>();
            for (auto it = index.cbegin(my_age); it != index.cend(my_age); ++it)
            {
                return true;
            }
            return false;
        }

        //7)、验证：modify基于迭代器修改数据
        ACTION void MultiIndexModify(const std::string &my_name)
        {
            member_student.modify(member_student.cbegin(), [&](auto &m) { m.name = my_name; });
        }

        //8)、验证：erase()多索引删除数据
        ACTION void MultiIndexErase(const std::string &my_name)
        {
          auto iter = member_student.find<"indexname"_n>(my_name);
          member_student.erase(iter);
        }

};

PLATON_DISPATCH(MultiIndexContract,(init)(addInitMultiIndex)(getMultiIndexCbegin)(getMultiIndexCend)(getMultiIndexCount)(getMultiIndexFind)(getMultiIndexIndex)(MultiIndexModify)(MultiIndexErase))