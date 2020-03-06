#include <platon/platon.hpp>
#include <string>
using namespace platon;

class Ownable 
{
	public:
		platon::StorageType<"owner"_n, Address> owner;
	
	public:
		Ownable() {
			owner.self() = platon_caller();	
		}
		
		bool onlyOwner(){
			Address sender = platon_caller();
			if(sender != owner.self()){
				platon_revert();
				return false;
			}		
			return true;
		}
	
};
