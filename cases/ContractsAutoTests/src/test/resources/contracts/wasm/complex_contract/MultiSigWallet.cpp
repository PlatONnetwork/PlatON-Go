#define TESTNET
#include <platon/platon.hpp>
#include <vector>
#include <string>

using namespace platon;

/**
 * MultiSigWallet.sol合约
 * */

CONTRACT MultiSigWallet : public platon::Contract{

    private:
        platon::StorageType<"maxowner"_n, uint64_t> MAX_OWNER_COUNT;
        platon::StorageType<"nonce"_n, uint64_t> nonce;
        platon::StorageType<"threshold"_n, uint64_t> threshold;
        platon::StorageType<"ownercount"_n, uint64_t> ownersCount;
        platon::StorageType<"isowner"_n,std::map<Address,bool>> isOwner;
        platon::StorageType<"tuple"_n,std::tuple<uint64_t,bytes,bytes>> vsr_tuple;
        platon::StorageType<"eventname"_n, std::string> eventName;
        platon::StorageType<"target"_n, std::string> target_addr;

    public:
        PLATON_EVENT1(OwnerAdded,std::string,Address)
        PLATON_EVENT1(OwnerRemoved,std::string,uint64_t)
        PLATON_EVENT1(ThresholdChanged,std::string,uint64_t)
        PLATON_EVENT1(Executed,std::string,Address,uint64_t,std::string)
        PLATON_EVENT1(Received,std::string,uint64_t,Address)

        ACTION void init(uint64_t _threshold,std::set<Address> _owners){
            target_addr.self() = "0x01"; //ecrecover(hash, v, r, s) 函数对应内置合约地址
            platon_assert(_owners.size() > 0, "MSW: Not enough or too many owners");
//            platon_assert(1>2,"_owners.size() is ", _owners.size());
            platon_assert(_threshold > 0 && _threshold <= _owners.size(), "MSW: Invalid threshold");
            ownersCount.self() = _owners.size();
            threshold.self() = _threshold;

            MAX_OWNER_COUNT.self() = 10;

            std::string name = "OwnerAdded";
            for(auto iter = _owners.begin(); iter != _owners.end(); iter++) {
                isOwner.self()[*iter] = true;
                PLATON_EMIT_EVENT1(OwnerAdded,name,*iter);
            }
             PLATON_EMIT_EVENT1(ThresholdChanged,name,_threshold);
        }

        ACTION void execute(Address _to, uint64_t _value, bytes _data, bytes _signatures,uint64_t value,uint64_t gas){
//            uint64_t v;
//            bytes r;
//            bytes s;

            uint64_t count = sizeof(_signatures) / 65;
            DEBUG("MultiSigWallet execute ", "count", count);
            DEBUG("MultiSigWallet execute ", "threshold.self()", threshold.self());
            nonce.self() += 1;
            uint64_t valid;
            Address lastSigner = Address(0);
            for(uint64_t i = 0; i < count; i++) {
                Address recovered = Address(0);

                auto address_info = make_address(target_addr.self());
                if(address_info.second){
                     if (platon_call(address_info.first, _signatures, value, gas)) {
                        platon::bytes ret;
                        size_t len = platon_get_call_output_length();
                        ret.resize(len);
                        platon_get_call_output(ret.data());
                        std::string str = toHex(ret);

                        auto address_info = make_address(str);
                        if(address_info.second){
                          recovered = address_info.first;
                        }
                    }
                }

                platon_assert(recovered > lastSigner, "MSW: Badly ordered signatures"); // make sure signers are different
                lastSigner = recovered;
                if(isOwner.self()[recovered]) {
                    uint64_t gas = 4712388;
                    valid += 1;
                    if(valid >= threshold.self()) {
                        if (platon_call(_to, _data, _value, gas)) {
                            DEBUG("cross call contract cross_call_ppos_send success", "address", _to);
                        }else{
                            DEBUG("cross call contract cross_call_ppos_send fail", "address", _to);
                        }
                        std::string name = "Executed";
                        PLATON_EMIT_EVENT1(Executed,name,_to, _value, toHex(_signatures));
                        return;
                    }
                }
            }
            platon_panic();
        }

//        void splitSignature(bytes _signatures,uint64_t _index){
//            uint64_t r;
//            bytes s;
//            bytes v;
//
//            //r calc
//            uint64_t rmul = 65 * _index;
//            uint64_t radd = 32 + rmul;
//
//
//
//            // assembly {
//            //     r := mload(add(_signatures, add(0x20,mul(0x41,_index))))
//            //     s := mload(add(_signatures, add(0x40,mul(0x41,_index))))
//            //     v := and(mload(add(_signatures, add(0x41,mul(0x41,_index)))), 0xff)
//            // }
//            // require(v == 27 || v == 28, "MSW: Invalid v");
//            // std::tuple<uint64_t,bytes,bytes> restuple = make_tuple(r,s,v);
//            vsr_tuple.self() = make_tuple(r,s,v);
//            // return restuple;
//        }

        ACTION void addOwner(Address _owner){
            platon_assert(ownersCount.self() < MAX_OWNER_COUNT.self(), "ownersCount is:",ownersCount.self(),"MAX_OWNER_COUNT is:",MAX_OWNER_COUNT.self(),"MSW: MAX_OWNER_COUNT reached");
            platon_assert(isOwner.self()[_owner] == false, "MSW: Already owner");
            ownersCount.self() += 1;
            isOwner.self()[_owner] = true;
            std::string name = "OwnerAdded";
            PLATON_EMIT_EVENT1(OwnerAdded,name,_owner);
        }

        ACTION void removeOwner(Address _owner){
            platon_assert(ownersCount.self() > threshold.self(), "MSW: Too few owners left");
            platon_assert(isOwner.self()[_owner] == true, "MSW: Not an owner");
            ownersCount.self() -= 1;
            isOwner.self().erase(_owner);
            std::string name = "OwnerRemoved";
            PLATON_EMIT_EVENT1(OwnerRemoved,name,ownersCount.self());
        }

        ACTION void changeThreshold(uint64_t _newThreshold){
            platon_assert(_newThreshold > 0 && _newThreshold <= ownersCount.self(), "MSW: Invalid new threshold");
            threshold.self() = _newThreshold;
            std::string name = "OwnerRemoved";
            PLATON_EMIT_EVENT1(OwnerRemoved,name,_newThreshold);
        }

         CONST std::map<Address,bool> getIsOwner(){
            return isOwner.self();
         }



};

PLATON_DISPATCH(MultiSigWallet, (init)(execute)(addOwner)(removeOwner)(changeThreshold)(getIsOwner))

