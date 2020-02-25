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

    private:
    platon::StorageType<"agevector"_n, std::vector<uint64_t>> ageVector;
};

PLATON_DISPATCH(InitWithVector, (init)(add_vector)(get_vector_size)(get_vector))
