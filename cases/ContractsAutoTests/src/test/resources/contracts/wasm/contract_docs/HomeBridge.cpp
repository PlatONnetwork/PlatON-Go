#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

#include "src/BridgeDeploymentAddressStorage.hpp"
#include "src/HomeBridgeGasConsumptionLimitsStorage.hpp"
#include "src/Message.hpp"
#include "src/Helpers.hpp"

class HomeBridge: public platon::Contract, public BridgeDeploymentAddressStorage, public HomeBridgeGasConsumptionLimitsStorage
{
	public:
		/// Number of authorities signatures required to withdraw the money.
	    ///
	    /// Must be lesser than number of authorities.
		platon::StorageType<"requiredSignatures"_n, u128> requiredSignatures;
		/// The gas cost of calling `HomeBridge.withdraw`.
	    ///
	    /// Is subtracted from `value` on withdraw.
	    /// recipient pays the relaying authority for withdraw.
	    /// this shuts down attacks that exhaust authorities funds on home chain.
	    platon::StorageType<"estimatedGasCostOfWithdraw"_n, u128> estimatedGasCostOfWithdraw;

	    /// Contract authorities.
	    platon::StorageType<"authorities"_n, std::vector<Address>> authorities;

	    /// Used foreign transaction hashes.
	    platon::StorageType<"withdraws"_n, std::map<h256, bool>> withdraws;

	public:
		/// Event created on money deposit.
	    PLATON_EVENT1(Deposit, Address, u128);

	    /// Event created on money withdraw.
	    PLATON_EVENT1(Withdraw, Address, u128);

	public:
		ACTION void init(u128 requiredSignaturesParam, std::vector<Address> authoritiesParam, u128 estimatedGasCostOfWithdrawParam) {
			if(requiredSignaturesParam == u128(0)){
				platon_revert();
			}
			if(requiredSignaturesParam > u128(authoritiesParam.size())){
				platon_revert();
			}
	        requiredSignatures.self() = requiredSignaturesParam;
	        authorities.self() = authoritiesParam;
	        estimatedGasCostOfWithdraw.self() = estimatedGasCostOfWithdrawParam;
		}
		
		/// final step of a withdraw.
	    /// checks that `requiredSignatures` `authorities` have signed of on the `message`.
	    /// then transfers `value` to `recipient` (both extracted from `message`).
	    /// see message library above for a breakdown of the `message` contents.
	    /// `vs`, `rs`, `ss` are the components of the signatures.

	    /// anyone can call this, provided they have the message and required signatures!
	    /// only the `authorities` can create these signatures.
	    /// `requiredSignatures` authorities can sign arbitrary `message`s
	    /// transfering any ether `value` out of this contract to `recipient`.
	    /// bridge users must trust a majority of `requiredSignatures` of the `authorities`.
	    ACTION void withdraw(std::vector<uint8_t> vs, std::vector<h256> rs, std::vector<h256> ss, bytes message) {
	        if(message.size() != 116){
	        	platon_revert();
	        }

	        // check that at least `requiredSignatures` `authorities` have signed `message`
	        if(!Helpers::hasEnoughValidSignatures(message, vs, rs, ss, authorities.self(), requiredSignatures.self())){
				platon_revert();
	        }

	        Address recipient = Message::getRecipient(message);
	        u128 value = Message::getValue(message);
	        h256 hash = Message::getTransactionHash(message);
	        u128 homeGasPrice = Message::getHomeGasPrice(message);

	        // if the recipient calls `withdraw` they can choose the gas price freely.
	        // if anyone else calls `withdraw` they have to use the gas price
	        // `homeGasPrice` specified by the user initiating the withdraw.
	        // this is a security mechanism designed to shut down
	        // malicious senders setting extremely high gas prices
	        // and effectively burning recipients withdrawn value.
	        // see https://github.com/paritytech/parity-bridge/issues/112
	        // for further explanation.
	        Address sender = platon_caller();
	        if(recipient != sender && platon_gas_price() != homeGasPrice){
	        	platon_revert();
	        }

	        // The following two statements guard against reentry into this function.
	        // Duplicated withdraw or reentry.
	        if(withdraws.self()[hash]){
	        	platon_revert();
	        }
	        // Order of operations below is critical to avoid TheDAO-like re-entry bug
	        withdraws.self()[hash] = true;

	        u128 estimatedWeiCostOfWithdraw = estimatedGasCostOfWithdraw.self() * homeGasPrice;

	        // charge recipient for relay cost
	        u128 valueRemainingAfterSubtractingCost = value - estimatedWeiCostOfWithdraw;

	        // pay out recipient
	        platon_transfer(recipient, Energon(valueRemainingAfterSubtractingCost));

	        // refund relay cost to relaying authority
	        platon_transfer(sender, Energon(estimatedWeiCostOfWithdraw));

	        PLATON_EMIT_EVENT1(Withdraw, recipient, valueRemainingAfterSubtractingCost);
	    }


};

PLATON_DISPATCH(HomeBridge,(init)(withdraw)(setGasLimitWithdrawRelay));
