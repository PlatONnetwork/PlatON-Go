#include <platon/platon.hpp>
#include <string>
using namespace platon;

/// Due to nature of bridge operations it makes sense to have the same value
/// of gas consumption limits which will distributed among all validators serving
/// particular bridge. This approach introduces few advantages:
/// --- new bridge instances will pickup limits from the contract instead of
///     looking at the configuration file (this configuration parameters could be
///     depricated)
/// --- as soon as upgradable bridge contract is implemented these limits needs
///     to be updated every time the contract is upgraded. Validators could get
///     an event that limits updated and use new values to send transactions.
class HomeBridgeGasConsumptionLimitsStorage 
{
	public:
		platon::StorageType<"gasLimitWithdrawRelay"_n, u128> gasLimitWithdrawRelay;

	public:
		PLATON_EVENT0(GasConsumptionLimitsUpdated, u128);
		
	public:
	    ACTION void setGasLimitWithdrawRelay(u128 gas) {
	        gasLimitWithdrawRelay.self() = gas;
	        PLATON_EMIT_EVENT0(GasConsumptionLimitsUpdated, gasLimitWithdrawRelay.self());
	    }
};
