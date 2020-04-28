#include <platon/platon.hpp>
#include <string>
#include <platon/db/multi_index.hpp>
using namespace std;
using namespace platon;

/**
 * @author qudong
 * 合约引用类型多索引(MultiIndex类型)：多索引支持唯一索引和普通索引。
 * 唯一索引应该放在参数的第一位。该结构需要提供与索引字段对应的get函数。
 *
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

CONTRACT ReferenceDataTypeMultiIndexContract : public platon::Contract{
    private:
       platon::db::MultiIndex<"student"_n,
                  Member,
                  platon::db::IndexedBy<"indexname"_n,platon::db::IndexMemberFun<Member,std::string,&Member::Name,platon::db::IndexType::UniqueIndex>>,
                  platon::db::IndexedBy<"indexage"_n,platon::db::IndexMemberFun<Member,uint8_t,&Member::Age,platon::db::IndexType::NormalIndex>>
                  >  member_student;
    public:
        ACTION void init(){}
         /**
         * 1、定义多索引类型
         **/
         //1)、多索引插入数据
         ACTION void addInitMultiIndex(const std::string &my_name,uint8_t &my_age,uint8_t &my_sex){
                member_student.emplace(([&](auto &m) {
                                         m.age = my_age;
                                         m.name = my_name;
                                         m.sex = my_sex;
                                        }));
         }
         //2)、验证：find（）多索引取值(查询年龄为10)
         CONST uint8_t getMultiIndexFind(const uint8_t &my_age){
              auto iter = member_student.find<"indexage"_n>(my_age);
              return iter->age;
          }

        //3)、验证：cbegin()多索引迭代器起始位置
         CONST bool getMultiIndexCbegin(){
            auto iter = member_student.find<"indexage"_n>(uint8_t(10));
            if(iter == member_student.cbegin()){
                return true;
            }
            return false;
         }

      //4)、验证：erase()多索引删除数据
      ACTION void deleteMultiIndexErase(){
          auto iter = member_student.find<"indexage"_n>(uint8_t(10));
          member_student.erase(iter);
      }

};

PLATON_DISPATCH(ReferenceDataTypeMultiIndexContract,(init)(addInitMultiIndex)(getMultiIndexFind)
               (getMultiIndexCbegin)(deleteMultiIndexErase))


