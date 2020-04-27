#include <platon/platon.hpp>
using namespace platon;


CONTRACT OOMException : public platon::Contract{
    public:
    ACTION void init() {}

    ACTION void memory_limit() {
        int* p = NULL;
        uint64_t number = platon_block_number();
        for (int i = 0; i < 6 * 1024; i++) {
          p = (int*)malloc(1024 * 1024);
          *p = i;
          printf("i = %d,block_number = %llu\t\n", *p, number);
        }
        printf("i = %d,block_number = %llu\t\n", *p, number);
    }
};

PLATON_DISPATCH(OOMException, (init)(memory_limit))

