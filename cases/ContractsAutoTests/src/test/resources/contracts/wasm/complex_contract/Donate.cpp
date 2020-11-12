#define TESTNET
#include <platon/platon.hpp>
#include <string>
#include <vector>
using namespace platon;

/**
 * Donate is a base contract for managing a donation
*/
CONTRACT Donate : public platon::Contract{
	
	 private:
		platon::StorageType<"owner"_n,Address> owner;
		platon::StorageType<"charity"_n,Address> charity;
		platon::StorageType<"openingTime"_n, u128> openingTime;
		platon::StorageType<"closingTime"_n, u128> closingTime;
		platon::StorageType<"minVonAmount"_n, u128> minVonAmount;
		platon::StorageType<"maxVonAmount"_n, u128> maxVonAmount;
		platon::StorageType<"maxNumDonors"_n, u128> maxNumDonors;
		platon::StorageType<"donors"_n,std::vector<Address>> donors;
		platon::StorageType<"whitelist"_n, std::map<Address, bool>> whitelist;
		platon::StorageType<"paused"_n, bool> paused;
		
	 public:
		PLATON_EVENT1(Donated, Address, u128)
		PLATON_EVENT0(OwnershipTransferred, Address, Address)

		/**
	     * @param _charity Address where collected funds will be forwarded to
	     * @param _openingTime Donate opening time
	     * @param _closingTime Donate closing time
	     * @param _minVonAmount Minimun donation amount in von
	     * @param _maxVonAmount Maximum donation amount in von
	     * @param _maxNumDonors Maximum number of donors
	   */
		ACTION void init(const Address& _charity, u128 _openingTime, u128 _closingTime, u128 _minVonAmount, u128 _maxVonAmount, u128 _maxNumDonors) {
		    platon_assert(_charity != Address(0));
		    platon_assert(_closingTime > _openingTime);
		    platon_assert(_minVonAmount > 0);
            platon_assert(_maxVonAmount > _minVonAmount);
            platon_assert(_maxNumDonors > 0);
	
			owner.self() = platon_caller();
			charity.self() = _charity;
			openingTime.self() = _openingTime;
			closingTime.self() = _closingTime;
			minVonAmount.self() = _minVonAmount;
			maxVonAmount.self() = _maxVonAmount;
			maxNumDonors.self() = _maxNumDonors;
		}
		
        /**
		 * low level donation.
		 */
		ACTION void donate(const Address& _donor) {
			u128 val = platon_call_value();
			preValidateDonate(_donor, val);
		  
			donors.self().push_back(_donor);
			PLATON_EMIT_EVENT1(Donated, platon_caller(), val);
		}
        
        /**
         * Sets opening and closing times.
         */		 
		ACTION void setOpeningClosingTimes(u128 _openingTime, u128 _closingTime) {
			onlyOwner();
			platon_assert(_closingTime > _openingTime);
			
			openingTime.self() = _openingTime;
		}
	  
	    /**
		 * Adds single address to whitelist.
		 */
		ACTION void addToWhitelist(const Address& _donor) {
			onlyOwner();
			
			whitelist.self()[_donor] = true;
		}
	  
		/**
		 * Adds list of addresses to whitelist. Not overloaded due to limitations with truffle testing.
		 */
		ACTION void addManyToWhitelist(const std::vector<Address>& _donors) {
			onlyOwner();
			
			for (uint64_t i = 0; i < _donors.size(); i++) {
		       whitelist.self()[_donors[i]] = true;
			}
		}
	  
	    /**
		 * Removes single address from whitelist.
		 */
		ACTION void removeFromWhitelist(const Address& _donor) {
			onlyOwner();
			
			whitelist.self()[_donor] = false;
		}

		/**
		 * called by the owner to pause, triggers stopped state.
		 */
		ACTION void pause() {
			onlyOwner();
			platon_assert(paused.self() == false);
			
			paused.self() = true;
		}
	 
	    /**
		 * called by the owner to unpause, returns to normal state.
		 */
		ACTION void unpause() {
			onlyOwner();
			platon_assert(paused.self() == true);
			
			paused.self() = false;
		}
	  
	    /**
		 * Allows the current owner to transfer control of the contract to a newOwner.
		 */
		ACTION void transferOwnership(const Address& newOwner) {
			onlyOwner();
			platon_assert(newOwner != Address(0));
			owner.self() = newOwner;
		  
			PLATON_EMIT_EVENT0(OwnershipTransferred, owner.self(), newOwner);
		}
		
		CONST Address getOwner() {
			return owner.self();
		}
		
		CONST Address getCharity() {
			return charity.self();
		}
		
		CONST u128 getOpeningTime() {
			return openingTime.self();
		}
		
		CONST u128 getClosingTime() {
			return closingTime.self();
		}
		
		CONST u128 getMinVonAmount() {
			return minVonAmount.self();
		}
		
		CONST u128 getMaxVonAmount() {
			return maxVonAmount.self();
		}
		
		CONST u128 getMaxNumDonors() {
			return maxNumDonors.self();
		}
		
		CONST std::vector<Address> getDonors() {
			return donors.self();
		}
		
		CONST std::map<Address, bool> getWhitelist() {
			return whitelist.self();
		}
		
		CONST bool getPaused() {
			return paused.self();
		}
	  
	  private:
		ACTION void preValidateDonate(const Address& _donor ,u128 _vonAmount) {
			platon_assert(u128(platon_timestamp()) >= openingTime.self() && u128(platon_timestamp()) <= closingTime.self());
			platon_assert(!paused.self());
			platon_assert(donors.self().size() <= maxNumDonors.self());
			platon_assert(whitelist.self()[_donor]);
			platon_assert(_donor != Address(0));
			platon_assert(minVonAmount.self() <= _vonAmount && _vonAmount <= maxVonAmount.self());
		}
	
		ACTION void onlyOwner() {
			platon_assert(platon_caller() == owner.self());
		}
};

PLATON_DISPATCH(Donate, (init)(donate)(setOpeningClosingTimes)(addToWhitelist)(addManyToWhitelist)(removeFromWhitelist)(pause)(unpause)(transferOwnership)(getOwner)(getCharity)(getOpeningTime)(getClosingTime)(getMinVonAmount)(getMaxVonAmount)(getMaxNumDonors)(getDonors)(getWhitelist)(getPaused))