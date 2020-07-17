#define TESTNET
#undef NDEBUG
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

		/// init value for address.
		ACTION void initAddress()
		{
			//tAddress.self() = make_address("0xf674172E619af9C09C126a568CF2838d243cE7F7");
			auto address_info = make_address("lax1w2kjkufl4g2v93xd94a0lewc75ufdr66rnzuw2");
			if(address_info.second) tAddress.self() = address_info.first;
		}

		/// set value for address.
		ACTION void setAddress(const std::string& input)
		{
		    auto address_info = make_address(input);
            if(address_info.second) tAddress.self() = address_info.first;
		}
		
		CONST std::string getAddress()
		{
			return tAddress.self().toString();
		}

		/// set value for u256.
		ACTION void setU256(uint64_t input)
		{
			tU256.self() = u128(input);
		}
		
		CONST std::string getU256()
		{
			return std::to_string(tU256.self());
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
		platon::StorageType<"su255"_n, u128> tU256;
		platon::StorageType<"sh255"_n, h256> tH256;

};

// (int8)(int64)(uint8t)(uint32t)(uint64t)(u128t)(u256t)
//(setInt8)(getInt8)(setInt32)(getInt32)(setInt64)(getInt64)
// (setUint8)(getUint8)(setUint32)(getUint32)(setUint64)(getUint64)
PLATON_DISPATCH(IntegerDataTypeContract_4,(init)(initAddress)
(setAddress)(getAddress)(setU256)(getU256)(setH256)(getH256))



