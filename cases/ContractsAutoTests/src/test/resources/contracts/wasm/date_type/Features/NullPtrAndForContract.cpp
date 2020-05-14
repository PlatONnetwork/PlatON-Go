#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;
using namespace std;

/**
 * @author qudong
 * 测试合约Nullptr 和 序列for循环
 * 1、Nullptr关键字
 *  其出现是为了解决NULL表示空指针在C++中具有二义性的问题
 *  NULL在C++中就是0；Nullptr表示为空指针
 * 2、序列for循环
 *
 * */

CONTRACT NullPtrAndForContract : public platon::Contract{

   private:
      platon::StorageType<"stringmap"_n, std::map<std::string,std::string>> storage_string_map;
      platon::StorageType<"arr"_n,std::array<std::string,10>> storage_array_string;
      platon::StorageType<"vector1"_n, std::vector<string>> vector_clothes;

    public:
       ACTION void init(){
       }

       /**
        * 1、Nullptr关键字
        */
       //1)、nullprt赋值
       CONST bool get_nullptr(){
          uint32_t *p = nullptr;//赋值空指针
          uint32_t *q = NULL;//定义null
          bool equal = (p == q);
          return equal;
       }
       //2)、nullprt赋值不同类型,转换成任何指针类型和bool布尔类型
       //但不能转换成整数
       CONST bool get_nullptr_one(){
             uint32_t *p1 = nullptr;   //编译正常,赋值空指针
             char *p2 = nullptr;       //编译正常,赋值空指针
             int64_t *p3 = nullptr;    //编译正常,赋值空指针
             std::string *p4 = nullptr;//编译正常,赋值空指针
             bool p5 = (bool)nullptr;  //p5为false
             //int8_t p6 = nullptr;    //编译异常，nullptr不能转换整型
             return p5;
        }
        //3)、验证NULL在C++中就是0；Nullptr表示为空指针
        CONST std::string set_nullptr_overload(){
              auto testNullptr = [](void* i) -> std::string{
                  return "is nullptr";
              };
             auto testNull = [](int8_t i) -> std::string{
                  return "is null";
              };
              return testNullptr(nullptr);
        }
        /**
         * 2、序列for循环
         * 在C++中for循环能够使用相似java的简化的for循环,能够用于遍历数组,容器,string
         */
       //1)、遍历map容器
       CONST std::string get_foreach_map(){
            std::string msg;
            storage_string_map.self().insert(pair<std::string,std::string>("one","1"));
            storage_string_map.self().insert(pair<std::string,std::string>("two","2"));
            storage_string_map.self().insert(pair<std::string,std::string>("three","3"));
            storage_string_map.self().insert(pair<std::string,std::string>("four","4"));
            for (auto itr: storage_string_map.self()){
                msg += itr.first + ",";
            }
            return msg;
       }
     //2)、遍历map容器
      CONST uint32_t get_foreach_array(){
           uint32_t sum;
           uint32_t array[10] = {0, 1, 2, 3, 4, 5, 6, 7, 8, 9};
           for (auto itr: array){
               sum += itr;
           }
           return sum;
      }

};

PLATON_DISPATCH(NullPtrAndForContract,(init)
(get_nullptr)(get_nullptr_one)
(set_nullptr_overload)(get_foreach_map)
(get_foreach_array)
)
