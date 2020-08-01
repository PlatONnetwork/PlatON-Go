#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;


CONTRACT InitWithVector : public platon::Contract{
    public:
    ACTION void init(uint16_t &age){
        ageVector.self().push_back(age);
    }

    ACTION void add_vector(uint64_t &one_age){
        ageVector.self().push_back(one_age);
    }

    CONST uint64_t get_vector_size(){
        return ageVector.self().size();
    }

    CONST uint64_t get_vector(uint8_t index){
        return ageVector.self()[index];
    }

    //push_back 在数组的最后添加一个数据
    ACTION void vector_push_back_element(std::string &value){
        strvector.self().push_back(value);
    }
    //insert 在第index元素后面插入value(如果index超过vector最后一个元素，则插入最后面)
    ACTION void vector_insert_element(std::string &value,uint8_t index){
        if(index > strvector.self().size()){
            strvector.self().insert(strvector.self().end(),value);
        }else{
            strvector.self().insert(strvector.self().begin()+index,value);
        }

    }
    //pop_back 去掉数组的最后一个数据
    ACTION void vector_pop_back_element(){
        strvector.self().pop_back();
    }
    //vector size
    CONST uint8_t get_strvector_size(){
        return strvector.self().size();
    }
    //得到编号位置的数据
    CONST std::string get_vector_element_by_position(uint8_t index){
        return strvector.self().at(index);
    }


    //验证for 循环
    ACTION void vectorFor(const std::vector<std::string> &inVector) {
 	  for(auto iter = inVector.begin(); iter != inVector.end(); iter++) {
        DEBUG("InitWithVector", "setSet", *iter);
      }
      strvector.self() = inVector;
    }

    //验证switch语句
    ACTION void vectorCase(const std::vector<std::string> &inVector) {
        switch(inVector.size()){
            case 1: str.self()="1";break;
            case 2: str.self()="2";break;
            case 3: str.self()="3";break;
            case 4: str.self()="4";break;
            case 5: str.self()="5";break;
            default: str.self()="big than 5";break;
        }
    }

    CONST std::string get_vectorCase_result(){
        return str.self();
    }


    private:
    platon::StorageType<"agevector"_n, std::vector<uint64_t>> ageVector;
    platon::StorageType<"strvector"_n, std::vector<std::string>> strvector;
    platon::StorageType<"strvar"_n, std::string> str;
};

PLATON_DISPATCH(InitWithVector, (init)(add_vector)(get_vector_size)(get_vector)(vector_push_back_element)
(vector_insert_element)(vector_pop_back_element)(get_strvector_size)(get_vector_element_by_position)(vectorFor)(vectorCase)(get_vectorCase_result))
