#include <platon/platon.hpp>
#include <vector>
#include <string>

using namespace platon;

CONTRACT call_precompile : public platon::Contract {
    public:
        ACTION void init(){}


};

PLATON_DISPATCH(cross_call, (init)(call_add_message)(get_vector_size))