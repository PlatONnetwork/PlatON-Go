#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;
 /**
   * 验证init函数能否重载
   * 1.当对init函数进行重载时，不管是否添加ACTION关键字编译都出错
   * 2.其它函数也不行进行重载
   * 3.二进制文件校验（删除前端两个字符看部署是否成功）
   *   执行结果：修改wasm文件可以编译成功，但是调用返回的结果为空
   * 编译：./platon-cpp vectortest.cpp -std=c++17
   * 打开ctool工具：./pWASM-ctool --config config.json 
   * 部署：deploy --wasm vectortest.wasm
   * 调用：invoke --addr 0x3a4B0C739F0F3fd9B11bee33997636c21e9b13Cd --func add_vector --params {"one_name":"1"}
   * 查询：call --addr 0x3a4B0C739F0F3fd9B11bee33997636c21e9b13Cd --func get_vector --params {"index":0}
   */
class person {
    public:
        person(){}
        person(const std::string &my_name):name(my_name){}
        std::string name;
        PLATON_SERIALIZE(person, (name))
};

//extern char const person_vector[] = "input_vector";

CONTRACT InitOverload : public platon::Contract{
    public:
    ACTION void init(){}

 /**
    //init函数重载编译失败
    void init(const std::string  &init_name){
      input_vector.self().push_back(person(init_name));
    }
*/    


    ACTION void add_vector(const std::string  &one_name){
        input_vector.self().push_back(person(one_name));
    }

    CONST uint64_t get_vector_size(){
        return input_vector.self().size();
    }

    CONST std::string get_vector(uint8_t index){
        return input_vector.self()[index].name;
    }
/**
    std::string get_vector(){
        return input_vector.self()[0].name;
    }
*/    

    private:
    platon::StorageType<"pvector"_n, std::vector<person>> input_vector;
};

PLATON_DISPATCH(InitOverload, (init)(add_vector)(get_vector_size)(get_vector))
