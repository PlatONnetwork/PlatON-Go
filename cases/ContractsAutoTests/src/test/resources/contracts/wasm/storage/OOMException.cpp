#include <platon/platon.hpp>
using namespace platon;


CONTRACT OOMException : public platon::Contract{
	public:

    ACTION void init() {}

    ACTION void memory_limit() {
        for (int i = 0; i < 20; i++) {
          int* p = (int*)malloc(1024 * 1024);
          *p = i;
          printf("%d\t\n", *p);
        }
    }
};

PLATON_DISPATCH(OOMException, (init)(memory_limit))