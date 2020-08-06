#define TESTNET
// Author: zjsunzone
// Desc: 验证所有基础数据类型的入参、返回值等是否合规
#include <platon/platon.hpp>
#include <string>
using namespace platon;


CONTRACT IntegerDataTypeContract_3: public platon::Contract
{

	/// common
	public:
		ACTION void init()
		{
			// do something to init.
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
		platon::StorageType<"sbyte"_n, char> tByte;
		platon::StorageType<"sbool"_n, bool> tBool;
		platon::StorageType<"sstring"_n, std::string> tString;
};

// (int8)(int64)(uint8t)(uint32t)(uint64t)(u128t)(u256t)
//(setInt8)(getInt8)(setInt32)(getInt32)(setInt64)(getInt64)
// (setUint8)(getUint8)(setUint32)(getUint32)(setUint64)(getUint64)
PLATON_DISPATCH(IntegerDataTypeContract_3,(init)
(setString)(getString)(setBool)(getBool)(setChar)(getChar))



