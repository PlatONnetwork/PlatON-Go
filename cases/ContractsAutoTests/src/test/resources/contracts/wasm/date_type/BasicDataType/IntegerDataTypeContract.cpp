#include <platon/platon.hpp>
#include <string>
using namespace platon;


// 定义基础数据类型的存储
extern char const sint8[] = "sint8";
extern char const sint32[] = "sint32";
extern char const sint64[] = "sint64";
extern char const sbyte[] = "sbyte";
extern char const sstring[] = "string";
extern char const sbool[] = "sbool";
extern char const suint8[] = "suint8";
extern char const suint32[] = "suint32";
extern char const suint64[] = "suint64";



/**
 * 验证所有基础数据类型的入参、返回值等是否合规
 */
CONTRACT IntegerDataTypeContract: public platon::Contract
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
		CONST long int int64()
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
			return to_string(u);
		}		

		/// u256
		CONST std::string u256t(uint64_t input)
		{
			u256 u = u256(input);
			return to_string(u);
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
		
		/// To set value for string.
		ACTION void setString(const std::string& input)
		{
			tString.self() = input;		
		}
		
		/// get the value from string.
		CONST std::string getString()
		{
			return tString.self();
		}
		
		/// To set value for bool.
		ACTION void setBool(bool input)
		{
			tBool.self() = input;		
		}
		
		/// get the value from bool.
		CONST bool getBool()
		{
			return tBool.self();
		}

		/// To set value for char.
		ACTION void setChar(char input)
		{
			tByte.self() = input;		
		}
		
		/// get the value from char.
		CONST char getChar()
		{
			return tByte.self();
		}
			
	private:
		platon::StorageType<sint8, int8_t> tInt8;
		platon::StorageType<sint32, int32_t> tInt32;
		platon::StorageType<sint64, int64_t> tInt64;
		platon::StorageType<sstring, std::string> tString;
		platon::StorageType<suint8, uint8_t> tUint8;
		platon::StorageType<suint32, uint32_t> tUint32;
		platon::StorageType<suint64, uint64_t> tUint64;
		platon::StorageType<sbyte, char> tByte;
		platon::StorageType<sbool, bool> tBool;

		

};

PLATON_DISPATCH(IntegerDataTypeContract,(init)(int8)(int64)(uint8t)(uint32t)(uint64t)(u128t)(u256t)
(setInt8)(getInt8)(setInt32)(getInt32)(setInt64)(getInt64)
(setUint8)(getUint8)(setUint32)(getUint32)(setUint64)(getUint64)
(setString)(getString)(setBool)(getBool)(setChar)(getChar))



