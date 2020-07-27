#define TESTNET
#include <platon/platon.hpp>
using namespace platon;


CONTRACT OOMException : public platon::Contract{
    public:
    ACTION void init() {}

    ACTION void memory_limit() {
        int* p = NULL;
        for (int i = 0; i < 6 * 1024; i++) {
          p = (int*)malloc(1024 * 1024);
          *p = i;
          printf("i = %d\t\n", *p);
        }
        printf("i = %d\t\n", *p);
    }
};

PLATON_DISPATCH(OOMException, (init)(memory_limit))


