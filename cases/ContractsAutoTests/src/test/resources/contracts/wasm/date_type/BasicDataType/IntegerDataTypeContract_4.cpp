// Author: zjsunzone
// Desc: 验证所有基础数据类型的入参、返回值等是否合规
#include <platon/platon.hpp>
#include <string>
using namespace platon;


CONTRACT IntegerDataTypeContract_4: public platon::Contract
{

	/// common
	public:
		ACTION void init()
		{
			// do something to init.
		}

		/// set value for address.
		ACTION void setAddress(const std::string& input)
		{
			tAddress.self() = Address(input);
		}
		
		CONST std::string getAddress()
		{
			return tAddress.self().toString();
		}

		/// set value for u256.
		ACTION void setU256(uint64_t input)
		{
			tU256.self() = u256(input);
		}
		
		CONST std::string getU256()
		{
			return to_string(tU256.self());
		}

		/// set value for h256.
		ACTION void setH256(const std::string& input)
		{
			tH256.self() = h256(input);
		}
		
		CONST std::string getH256()
		{
			return tH256.self().toString();
		}
			
	private:
		platon::StorageType<"saddress"_n, Address> tAddress;
		platon::StorageType<"su255"_n, u256> tU256;
		platon::StorageType<"sh255"_n, h256> tH256;

};

// (int8)(int64)(uint8t)(uint32t)(uint64t)(u128t)(u256t)
//(setInt8)(getInt8)(setInt32)(getInt32)(setInt64)(getInt64)
// (setUint8)(getUint8)(setUint32)(getUint32)(setUint64)(getUint64)
PLATON_DISPATCH(IntegerDataTypeContract_4,(init)
(setAddress)(getAddress)(setU256)(getU256)(setH256)(getH256))



