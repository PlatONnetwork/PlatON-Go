// Author: zjsunzone
// 斐波拉契合约验证
#include <platon/platon.hpp>
#include <string>
using namespace platon;

CONTRACT Fibonacci: public platon::Contract
{

	public:
		PLATON_EVENT1(Notify, uint64_t, uint64_t)

	public:
		ACTION void init()
		{
			// do something to init.
		}
		
		ACTION void fibonacci_notify(uint64_t number)
		{
			uint64_t result = fibonacci(number);
			PLATON_EMIT_EVENT1(Notify, number, result); 
		}
		
		CONST uint64_t fibonacci_call(uint64_t number)
		{
			if(number == 0){
				return 0;			
			} else if(number == 1){
				return 1;			
			} else {
				return fibonacci(number-1) + fibonacci(number-2);		
			}
		}

};


PLATON_DISPATCH(Fibonacci,(init)(fibonacci_notify)(fibonacci_call))



