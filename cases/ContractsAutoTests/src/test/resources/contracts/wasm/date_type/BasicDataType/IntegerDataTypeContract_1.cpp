#define TESTNET
// Author: zjsunzone
// Desc: 验证所有基础数据类型的入参、返回值等是否合规
#include <platon/platon.hpp>
#include <string>
using namespace platon;


CONTRACT IntegerDataTypeContract_1: public platon::Contract
{

	/// common
	public:
		ACTION void init()
		{
			// do something to init.
		}
	
	/// integer data type.
	public: 
		/// int8 返回验证
		/// range: -32768 到 32767
		CONST short int int8()
		{
			return 3;
		}

		/// int32
		/// range: -2147483648 到 2147483647
		CONST int int32()
		{
			return 2;
		}
	
		/// int64
		/// range: -9,223,372,036,854,775,808 到 9,223,372,036,854,775,807
		CONST long long int64()
		{
			return 200;
		}
		
		/// uint8_t
		/// range: 
		CONST uint8_t uint8t(uint8_t input)
		{
			return input * 2;
		} 

		/// uint32_t
		CONST uint32_t uint32t(uint32_t input)
		{
			return input * 2;
		}
		
		/// uint64_t
		CONST uint64_t uint64t(uint64_t input)
		{
			return input * 2;
		}
		

		/// u128
		CONST std::string u128t(uint64_t input)
		{
			u128 u = u128(input);
			return std::to_string(u);
		}		

		/// u256
		CONST std::string u256t(uint64_t input)
		{
			u128 u = u128(input);
			return std::to_string(u);
		}



};

// (int8)(int64)(uint8t)(uint32t)(uint64t)(u128t)(u256t)
//(setInt8)(getInt8)(setInt32)(getInt32)(setInt64)(getInt64)
// (setUint8)(getUint8)(setUint32)(getUint32)(setUint64)(getUint64)
PLATON_DISPATCH(IntegerDataTypeContract_1,(init)
(int8)(int64)(uint8t)(uint32t)(uint64t)(u128t)(u256t))



