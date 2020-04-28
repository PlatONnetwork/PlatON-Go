#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * 内置函数
 * platon_ecrecover
 * platon_ripemd160
 * platon_sha256
 */
CONTRACT CryptographicFunction:public platon::Contract{
	public:
		ACTION void init() {}
		
		// platon_ecrecover
		CONST Address call_platon_ecrecover(const h256 &hash, const bytes &signature){
		     Address result;
			 int32_t res = platon_ecrecover(hash,signature,result);
			 return result;
		}

		// platon_ripemd160
		CONST h160  call_platon_ripemd160(const bytes &data){
			return platon_ripemd160(data);
		}

		// platon_sha256
		CONST h256  call_platon_sha256(const bytes &data) {
			return platon_sha256(data);
		}
};


PLATON_DISPATCH(CryptographicFunction, (init)(call_platon_ecrecover)(call_platon_ripemd160)(call_platon_sha256))


