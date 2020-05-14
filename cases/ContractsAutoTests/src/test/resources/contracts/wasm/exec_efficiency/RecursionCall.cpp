#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * 执行效率-递归调用
 * @author qcxiao
 **/
CONTRACT RecursionCall : public platon::Contract {

    private:
	    platon::StorageType<"test"_n, uint64_t> sum = 0;
	public:
        ACTION void init(){}

        ACTION void call(uint64_t n) {
            if (sum < n) {
                ++sum;
                call(n);
            }
        }

        CONST uint64_t get_sum() {
            return sum.self();
        }
};
PLATON_DISPATCH(RecursionCall,(init)(call)(get_sum))