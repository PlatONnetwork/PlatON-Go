#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * 执行效率-递归调用
 * @author qcxiao
 **/
CONTRACT RecursionCall : public platon::Contract {

    private:
	    platon::StorageType<"test"_n, uint64_t> sum;
	public:
        ACTION void init(){}

        ACTION void call(uint64_t n) {
            if (sum < n) {
                ++sum;
                call(sum);
            }
        }

};
PLATON_DISPATCH(RecursionCall,(init)(call))