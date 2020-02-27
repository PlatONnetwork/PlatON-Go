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

    private:
    platon::StorageType<"list"_n, std::list<std::string>> slist;
    platon::StorageType<"listlist"_n, std::list<std::list<std::string>>> slistlist;
};

PLATON_DISPATCH(InitWithListParams, (init)(set_list)(get_list)(set_list_list)(get_list_list))
