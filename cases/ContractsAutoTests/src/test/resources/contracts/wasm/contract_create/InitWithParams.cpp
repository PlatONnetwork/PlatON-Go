#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

class Person {
    public:
        Person(){}
        Person(const std::string &my_name):name(my_name){}
        std::string name;
        PLATON_SERIALIZE(Person, (name))
};

//extern char const person_vector[] = "person_vector";

CONTRACT InitWithParams : public platon::Contract{
   public:
      ACTION void init(const std::string &init_name){
         info.self().push_back(Person(init_name));
      }

      ACTION std::vector<Person> add_person(const Person &one_person){
          info.self().push_back(one_person);
          return info.self();
      }
      CONST std::vector<Person> get_person(){
          return info.self();
      }

   private:
      platon::StorageType<"pvector"_n, std::vector<Person>> info;
};

PLATON_DISPATCH(InitWithParams, (init)(add_person)(get_person))
