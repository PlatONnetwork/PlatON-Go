#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;


CONTRACT InitWithListParams : public platon::Contract{
    public:
    ACTION void init(const std::list<std::string>  &inList){
        slist.self() = inList;
    }

    ACTION void set_list(const std::list<std::string>  &inList){
        slist.self() = inList;
    }

    CONST std::list<std::string> get_list(){
        return slist.self();
    }

    ACTION void set_list_list(const std::list<std::list<std::string>>  &inlistlist){
        slistlist.self() = inlistlist;
    }

    CONST std::list<std::list<std::string>> get_list_list(){
        return slistlist.self();
    }

    //list add element
    ACTION void add_list_element(std::string &value){
        slist.self().push_back(value);
    }

    //返回最后一个元素
    CONST std::string get_list_last_element(){
        return slist.self().back();
    }

    //返回第一个元素
    CONST std::string get_list_first_element(){
        return slist.self().front();
    }

    //删除所有元素
    ACTION void list_clear(){
        slist.self().clear();
    }

    //删除一个元素
    ACTION void list_remove_element(std::string &element){
        slist.self().remove(element);
    }

    //删除最后一个元素
    ACTION void list_pop_back(){
        slist.self().pop_back();
    }

    //删除第一个元素
    ACTION void list_pop_front(){
        slist.self().pop_front();
    }

    //list中元素个数
    CONST uint8_t list_size(){
        return slist.self().size();
    }

    //删除list中重复的元素
    ACTION void list_unique(){
        slist.self().unique();
    }

    //两个list合并
    ACTION void list_merge(const std::list<std::string> &inList){
        auto& l = slist.self();
        l.insert(l.end(), inList.begin(), inList.end());
    }

    private:
    platon::StorageType<"list"_n, std::list<std::string>> slist;
    platon::StorageType<"listlist"_n, std::list<std::list<std::string>>> slistlist;
};

PLATON_DISPATCH(InitWithListParams, (init)(set_list)(get_list)(set_list_list)(get_list_list)
(add_list_element)(get_list_last_element)(get_list_first_element)(list_clear)(list_remove_element)(list_pop_back)(list_pop_front)(list_size)
(list_unique)(list_merge))
