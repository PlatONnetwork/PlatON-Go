#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * 执行效率-循环调用
 * @author qcxiao
 **/

CONTRACT LoopCall : public platon::Contract {

    private:
        uint64_t sum;
    public:
        ACTION void init(){}
        ACTION void loopCallTest(uint64_t n) {
            for (int i = 0; i < n; i++) {
                sum += i;
            }
        }
	CONST uint64_t get_sum() {
	    return sum;
	}
};
PLATON_DISPATCH(LoopCall,(init)(loopCallTest)(get_sum))
