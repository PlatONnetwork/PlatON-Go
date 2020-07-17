#define TESTNET
// Author: zjsunzone
// Desc: 验证所有基础数据类型的入参、返回值等是否合规
#include <platon/platon.hpp>
#include <string>
using namespace platon;

CONTRACT IntegerDataTypeContract_2: public platon::Contract
{

	/// common
	public:
		ACTION void init()
		{
			// do something to init.
		}
	
	
	// ACTION
	public:
		/// to set value for int8.
		ACTION void setInt8(int8_t input)
		{
			tInt8.self() = input;
			DEBUG("Invoke setInt8", "input", input);
		}
		
		/// get the value from int8.
		CONST int8_t getInt8()
		{
			return tInt8.self();
		}
		
		/// to set value for int32.
		ACTION void setInt32(int32_t input)
		{
			tInt32.self() = input;
			DEBUG("Invoke setInt32", "input", input);
		}
		
		/// get the value from int32.
		CONST int32_t getInt32()
		{
			return tInt32.self();
		}

		/// to set value for int64.
		ACTION void setInt64(int64_t input)
		{
			tInt64.self() = input;
			DEBUG("Invoke setInt64", "input", input);
		}
		
		/// get the value from int64.
		CONST int64_t getInt64()
		{
			return tInt64.self();
		}

		/// to set value for uint8.
		ACTION void setUint8(uint8_t input)
		{
			tUint8.self() = input;
			DEBUG("Invoke setUint8", "input", input);
		}
		
		/// get the value from uint8.
		CONST uint8_t getUint8()
		{
			return tUint8.self();
		}
		
		/// to set value for uint32.
		ACTION void setUint32(uint32_t input)
		{
			tUint32.self() = input;
			DEBUG("Invoke setUint32", "input", input);
		}
		
		/// get the value from uint32.
		CONST uint32_t getUint32()
		{
			return tUint32.self();
		}
		
		/// to set value for uint64.
		ACTION void setUint64(uint64_t input)
		{
			tUint64.self() = input;
			DEBUG("Invoke setUint64", "input", input);
		}
		
		/// get the value from uint64.
		CONST uint64_t getUint64()
		{
			return tUint64.self();
		}
		
			
	private:
		platon::StorageType<"sint24"_n, int8_t> tInt8;
		platon::StorageType<"sint32"_n, int32_t> tInt32;
		platon::StorageType<"sint44"_n, int64_t> tInt64;
		platon::StorageType<"suint"_n, uint8_t> tUint8;
		platon::StorageType<"suint32"_n, uint32_t> tUint32;
		platon::StorageType<"suint44"_n, uint64_t> tUint64;

};

// (int8)(int64)(uint8t)(uint32t)(uint64t)(u128t)(u256t)
//(setInt8)(getInt8)(setInt32)(getInt32)(setInt64)(getInt64)
// (setUint8)(getUint8)(setUint32)(getUint32)(setUint64)(getUint64)
PLATON_DISPATCH(IntegerDataTypeContract_2,(init)
(setInt8)(getInt8)(setInt32)(getInt32)(setInt64)(getInt64)
(setUint8)(getUint8)(setUint32)(getUint32)(setUint64)(getUint64))



