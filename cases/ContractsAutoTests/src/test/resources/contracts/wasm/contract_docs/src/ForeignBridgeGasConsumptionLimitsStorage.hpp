#include <platon/platon.hpp>
#include <string>
using namespace platon;

class ForeignBridgeGasConsumptionLimitsStorage 
{
	public:
		platon::StorageType<"gasLimitDepositRelay"_n, u128> gasLimitDepositRelay;
		platon::StorageType<"gasLimitWithdrawConfirm"_n, u128> gasLimitWithdrawConfirm;

	public:
		PLATON_EVENT1(GasConsumptionLimitsUpdated, u128, u128);
		
	public:
	    ACTION void setGasLimitDepositRelay(u128 gas) {
	        gasLimitDepositRelay.self() = gas;
	        PLATON_EMIT_EVENT1(GasConsumptionLimitsUpdated, gasLimitDepositRelay.self(), gasLimitWithdrawConfirm.self());
	    }

	    ACTION void setGasLimitWithdrawConfirm(u128 gas) {
	        gasLimitWithdrawConfirm.self() = gas;
	        PLATON_EMIT_EVENT1(GasConsumptionLimitsUpdated, gasLimitDepositRelay.self(), gasLimitWithdrawConfirm.self());
	    }	
};
