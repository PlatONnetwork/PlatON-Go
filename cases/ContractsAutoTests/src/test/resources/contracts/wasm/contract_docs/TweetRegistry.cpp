#define TESTNET
// Author: zjsunzone
#include <platon/platon.hpp>
#include <string>
using namespace platon;


CONTRACT TweetRegistry: public platon::Contract
{
	private:
		// mappings to look up account names, account ids and addresses
		platon::StorageType<"addrmapname"_n, std::map<Address, std::string>> _addressToAccountName;
		platon::StorageType<"idmapadr"_n, std::map<u128, Address>> _accountIdToAccountAddress;
		platon::StorageType<"namemapaddr"_n, std::map<std::string, Address>> _accountNameToAddress;

		// might be insteesting to see how many people use the system.
		platon::StorageType<"numberofaccount"_n, u128> _numberOfAccounts;
		
		// owner
		platon::StorageType<"_registryAdmin"_n, Address> _registryAdmin;

		// allowed to administrate accounts only, not everything.
		platon::StorageType<"_accountAdmin"_n, Address> _accountAdmin;

		// if a newer version of this registry is available, force users to use it.
		platon::StorageType<"_dsabled"_n, bool> _registrationDisabled;
	
	public:
		ACTION void init()
		{
			Address sender = platon_caller();
			_registryAdmin.self() = sender;
			_accountAdmin.self() = sender;	// can be changed later.
			_numberOfAccounts.self() = u128(0);
			_registrationDisabled.self() = false;
		}
		
		ACTION int registry(const std::string& name, const Address& accountAddress) {
			int result = 0;
			if(_accountNameToAddress.self()[name] != Address(0)){
				// name already token
				result = -1;
			} else if (_addressToAccountName.self()[accountAddress].length() != 0) {
				// account address is already registered
				result = -2;
			} else if (name.length() >= 64){
				// name too long
				result = -3;			
			} else if (_registrationDisabled.self()) {
				// registry is disabled because a newer version is available
				result = -4;			
			} else {
				_addressToAccountName.self()[accountAddress] = name;
				_accountNameToAddress.self()[name] = accountAddress;
				_accountIdToAccountAddress.self()[_numberOfAccounts] = accountAddress;
				_numberOfAccounts.self()++;
				result = 0; // success.			
			}
			return result;
		}

		CONST u128 getNumberOfAccounts() {
			return _numberOfAccounts.self();		
		}

		CONST Address getAddressOfName(const std::string& name){
			return _accountNameToAddress.self()[name];		
		}
		
		CONST std::string getNameOfAddress(const Address& addr) {
			return _addressToAccountName.self()[addr];		
		}

		CONST Address getAddressOfId(const u128& id) {
			return _accountIdToAccountAddress.self()[id];		
		}

		ACTION std::string unregister() {
			Address sender = platon_caller();
			std::string unregisteredAccountName = _addressToAccountName.self()[sender];
			_addressToAccountName.self()[sender] = "";
			_accountNameToAddress.self()[unregisteredAccountName] = Address(0);
			return unregisteredAccountName;		
		}
	
		ACTION void adminUnregister(const std::string& name) {
			Address sender = platon_caller();
			if(sender == _registryAdmin.self() || 
				sender == _accountAdmin.self() ){
				Address addr = _accountNameToAddress.self()[name];
				_addressToAccountName.self()[addr] = "";
				_accountNameToAddress.self()[name] = Address(0);
			}		
		} 
		
		ACTION bool adminSetRegistrationDisable(bool registrationDisabled){
			// currently, the code of the registry can not be updatd once it is
			// deployed. if a newer version of the registry is available, account
			// registration can be disabled
			Address sender = platon_caller();
			if(sender == _registryAdmin.self()){
				_registrationDisabled.self() = registrationDisabled;			
				return true;			
			}	
			return false;	
		}

		CONST bool getRegistrationDisabled() {
			return _registrationDisabled.self();		
		}
		
		ACTION void adminSetAccountAdministrator(const Address& accountAdmin) {
			Address sender = platon_caller();
			if(sender == _registryAdmin.self()){
				_accountAdmin.self() = accountAdmin;
			}		
		}

		ACTION void adminRetrieveDonations(){
			Address sender = platon_caller();
			if(sender == _registryAdmin.self()){
				Address caddr = platon_address();
				Energon e = platon_balance(caddr);
				platon_transfer(_registryAdmin.self(), e);
			}	
		}

		ACTION void adminDeleteRegistry() {
			Address sender = platon_caller();
			if(sender == _registryAdmin.self()){
				// this is a predefined function, it deletes the contract 
				// and returns all funds to the admin's address.
				platon_destroy(_registryAdmin.self());	
			}
		}

};

PLATON_DISPATCH(TweetRegistry,(init)(registry)(getNumberOfAccounts)(getAddressOfName)
(getNameOfAddress)(getAddressOfId)(unregister)(adminUnregister)(adminSetRegistrationDisable)
(adminSetAccountAdministrator)(adminRetrieveDonations)(adminDeleteRegistry)(getRegistrationDisabled))



