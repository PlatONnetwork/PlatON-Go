#include <platon/platon.hpp>
#include <string>
using namespace platon;

class person {
    public:
        person(){}
        person(const std::string &my_name):name(my_name){}
        std::string name;
        PLATON_SERIALIZE(person, (name))
};

extern char const person_vector[] = "person_vector";

CONTRACT InitWithParams : public platon::Contract{
   public:
      ACTION void init(const std::string  &init_name){
         info.self().push_back(person(init_name));
      }

      ACTION std::vector<person> add_person(const person &one_person){
          info.self().push_back(one_person);
          return info.self();
      }
      CONST std::vector<person> get_person(){
          return info.self();
      }

   private:
      platon::StorageType<person_vector, std::vector<person>> info;
};

PLATON_DISPATCH(InitWithParams, (init)(add_person)(get_person))
