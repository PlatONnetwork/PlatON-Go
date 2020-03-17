#include <platon/platon.hpp>
#include <string>
using namespace platon;

class BridgeDeploymentAddressStorage 
{
	public:
		platon::StorageType<"deployedAtBlock"_n, u128> deployedAtBlock;

	public:
	    ACTION BridgeDeploymentAddressStorage() {
	        deployedAtBlock.self() = u128(platon_block_number());
	    }	
};
