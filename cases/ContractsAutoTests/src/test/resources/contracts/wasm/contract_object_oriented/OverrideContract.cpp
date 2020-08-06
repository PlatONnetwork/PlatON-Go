#define TESTNET
// Author: zjsunzone
// 验证合约的继承、重载功能；
#include <platon/platon.hpp>
using namespace platon;

// 基类
class Shape {
   protected:
      int32_t width, height;

   public:
		Shape( int32_t a=1, int32_t b=1)
		{
			width = a;
			height = b;
		}

		virtual uint32_t area()
		{
			return width * height;
		}

		uint32_t area2()
		{
			return 100 * 100;	
		}
};

// 合约类
CONTRACT OverrideContract: public platon::Contract, public Shape {
	public:
		ACTION void init() 
		{
			
		};
		
		// override area from super class.
		uint32_t area() 
		{
			return 100;	
		}

		CONST uint32_t getArea(uint64_t input)
		{	
			if(input == 1)
			{
				return area();		
			}

			if(input == 2)
			{
				return area2();		
			}
			return 0;
		}
};

PLATON_DISPATCH(OverrideContract, (init)(getArea))


