#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;


CONTRACT InitOverloadWithString : public platon::Contract{
    public:
    ACTION void init(const std::string  &initStr){
        str.self() = initStr;
    }

     //获取字符串长度
    CONST std::string get_string(){
        return str.self();
    }

    //获取字符串长度
   CONST uint8_t string_length(){
       return str.self().size();
   }

    //字符串连接
    CONST std::string string_splice(const std::string &spliceStr){
        return str.self() + spliceStr;
    }

    //字符串比较
    CONST int8_t string_compare(const std::string &strone,const std::string &strtwo){
        if(strone > strtwo){
            return 1;
        }else if (strone == strtwo){
            return 0;
        }else{
            return -1;
        }
    }

    //字符串倒置
    ACTION void string_reverse(const std::string &reverseStr){
         return reverse(str.self().begin(),str.self().end());
    }

    //字符串查找
    CONST int8_t string_find(const std::string &findStr){
         return str.self().find(findStr);
    }


    private:
    platon::StorageType<"stropt"_n, std::string> str;
};

PLATON_DISPATCH(InitOverloadWithString, (init)(get_string)(string_length)(string_splice)(string_compare)(string_reverse)(string_find))
