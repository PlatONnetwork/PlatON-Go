#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * 针对系统链上函数的调用
 */
CONTRACT InnerFunction_2:public platon::Contract{
	public:
		ACTION void init() {}

		/// 主币转账
		/// define: int32_t platon_transfer(const Address& addr, const Energon& amount);
		ACTION void transfer(const std::string& addr, uint64_t amount) {
			if(amount == 0){
				DEBUG("Transfer failed", "address", addr, "amount", amount);
				return;
			}

			auto address_info = make_address(addr);
			if(address_info.second){
			    platon_transfer(address_info.first, Energon(amount));
			}
		}
		
		/// 获取消息携带的value(fix) 
		/// define: u128 platon_call_value();
		CONST std::string value() { 
			u128 val = platon_call_value();
			return std::to_string(val);		
		}

		/// sha3操作
		/// define: h256 platon_sha3(const bytes& data);
		CONST std::string sha3(const std::string& str) {
			bytes data;
			data.insert(data.begin(), str.begin(), str.end());
			h256 hash = platon::platon_sha3(data);
			return hash.toString();
		} 

		/// 设置函数返回值
		/// define: template <typename T> void platon_return(const T& t);
		CONST void rreturn() {
			std::string str = "hello";
			platon_return(str);
		}

		/// 终止交易 panic, 交易完成，合约执行异常
		/// define: void platon_panic();
		ACTION void panic() {
			platon_panic();		
		}

		/// 终止交易 revert
		/// define: void platon_revert();
		ACTION void revert(int64_t flag) {
			if(flag == 1){
				platon_revert();
			}		
		} 

		/// 合约销毁 destroy, 销毁后检测余额
		/// define: bool platon_destroy(const Address& addr);
		ACTION void destroy(const std::string& addr) {
		    auto address_info = make_address(addr);
		    if(address_info.second){
		        platon_destroy(address_info.first);
		    }
		}
		
		/// 消息的原始发送者origin
		/// define: Address platon_origin();
		CONST Address origin() {
			Address ori = platon::platon_origin();
			return ori;		
		}

		/// compile test
		/// summary: compile success.
		std::string compile(){
			return "compile";		
		}

		CONST Address addr(){

            Address address;
		    auto address_info = make_address("lax1fyeszufxwxk62p46djncj86rd553skpptsj8v6");
            if(address_info.second){
                address = address_info.first;
            }
			return address;
		}
	
};

PLATON_DISPATCH(InnerFunction_2, (init)(addr)(transfer)(value)(sha3)(rreturn)(panic)(revert)(destroy)(origin))


