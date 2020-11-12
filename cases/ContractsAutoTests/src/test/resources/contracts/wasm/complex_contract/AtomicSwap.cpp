#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * @author liweic
 * Atomex
 * */

class SafeMath
{
	public:
		ACTION u128 add(u128 a, u128 b)
        {
            u128 c = a + b;
            platon_assert(c >= a, "SafeMath add wrong value");
            return c;
        }

        ACTION u128 sub(u128 a, u128 b)
        {
            platon_assert(b <= a, "SafeMath sub wrong value");
            u128 c = a - b;
            return c;
        }
};

enum State { Empty, Initiateds, Redeemeds, Refundeds };

struct Swap{
    public:
        bytes hashedSecret;
        bytes secret;
        Address initiator;
        Address participant;
        uint64_t refundTimestamp;
        u128 value;
        u128 payoff;
        int state;
        // Swap(){}
        // Swap(bytes &hashedSecret, bytes &secret, Address &initiator, Address &participant, uint64_t &refundTimestamp, u128 &value, u128 &payoff):hashedSecret(hashedSecret),secret(secret),initiator(initiator),participant(participant),refundTimestamp(refundTimestamp),value(value),payoff(payoff), state(state){}
        PLATON_SERIALIZE(Swap,(hashedSecret)(secret)(initiator)(participant)(refundTimestamp)(value)(payoff)(state))
};

CONTRACT AtomicSwap : public platon::Contract, public SafeMath{

    private:
        platon::StorageType<"owner"_n,Address> owner;
        StorageType<"swap"_n, Swap> swap_;
        platon::StorageType<"swaps"_n,std::map<bytes,Swap>> swaps;

    public:
        PLATON_EVENT1(Initiated, bytes, Address, Address, uint64_t, u128, u128);
        PLATON_EVENT0(Added, bytes, Address, u128);
        PLATON_EVENT0(Redeemed, bytes, bytes);
        PLATON_EVENT0(Refunded, bytes);

    public:
        ACTION void init()
        {
            owner.self() = platon_caller();
        }

        ACTION void destruct()
        {
            platon_assert(platon_caller() == owner.self(), "only owner");
            platon_assert(platon_balance(owner.self()) == 0, "balance is not zero");
            platon_destroy(owner.self());
        }

        ACTION void initiate(bytes _hashedSecret, Address _participant, uint64_t _refundTimestamp, u128 _payoff)
        {
            platon_assert(swaps.self()[_hashedSecret].state == Empty, "swap for this hash is already initiated");
            platon_assert(_refundTimestamp > platon_timestamp(), "refundTimestamp has already passed");
            u128 initiateval = platon_call_value();

            swaps.self()[_hashedSecret].value = SafeMath::sub(initiateval, _payoff);
            swaps.self()[_hashedSecret].hashedSecret = _hashedSecret;
            swaps.self()[_hashedSecret].initiator = platon_caller();
            swaps.self()[_hashedSecret].participant = _participant;
            swaps.self()[_hashedSecret].refundTimestamp = _refundTimestamp;
            swaps.self()[_hashedSecret].payoff = _payoff;
            swaps.self()[_hashedSecret].state = Initiateds;

            PLATON_EMIT_EVENT1(Initiated, _hashedSecret, swaps.self()[_hashedSecret].participant, platon_caller(), swaps.self()[_hashedSecret].refundTimestamp, swaps.self()[_hashedSecret].value, swaps.self()[_hashedSecret].payoff);
        }

        ACTION void add(bytes _hashedSecret)
        {
            platon_assert(swaps.self()[_hashedSecret].state == Initiateds, "swap for this hash is empty or already spent");
            platon_assert(platon_timestamp() <= swaps.self()[_hashedSecret].refundTimestamp, "refundTime has already come");
            u128 addval = platon_call_value();
            swaps.self()[_hashedSecret].value = SafeMath::add(swaps.self()[_hashedSecret].value, addval);

            PLATON_EMIT_EVENT0(Added, _hashedSecret, platon_caller(), swaps.self()[_hashedSecret].value);
        }

        ACTION void redeem(bytes _hashedSecret, bytes _secret)
        {
            platon_assert(swaps.self()[_hashedSecret].state == Initiateds, "swap for this hash is empty or already spent");
            platon_assert(platon_timestamp() < swaps.self()[_hashedSecret].refundTimestamp, "refundTimestamp has already passed");
            //wasm没有ABI编解码函数
            //platon_assert(sha256(abi.encodePacked(sha256(abi.encodePacked(_secret)))) == _hashedSecret, "secret is not correct");
            swaps.self()[_hashedSecret].secret = _secret;
            swaps.self()[_hashedSecret].state = Redeemeds;
            PLATON_EMIT_EVENT0(Redeemed, _hashedSecret, _secret);
            platon_transfer(swaps.self()[_hashedSecret].participant, Energon(swaps.self()[_hashedSecret].value));
            if(swaps.self()[_hashedSecret].payoff > 0)
            {
                platon_transfer(platon_caller(), Energon(swaps.self()[_hashedSecret].payoff));
            }

            swaps.self().erase(_hashedSecret);
        }

        ACTION void refund(bytes _hashedSecret)
        {

            platon_assert(swaps.self()[_hashedSecret].state == Initiateds, "swap for this hash is empty or already spent");
            platon_assert(platon_timestamp() >= swaps.self()[_hashedSecret].refundTimestamp, "refundTimestamp has not passed");
            swaps.self()[_hashedSecret].state = Refundeds;
            PLATON_EMIT_EVENT0(Refunded, _hashedSecret);
            platon_transfer(swaps.self()[_hashedSecret].initiator, Energon(SafeMath::add(swaps.self()[_hashedSecret].value,swaps.self()[_hashedSecret].payoff)));
            swaps.self().erase(_hashedSecret);
        }

};

PLATON_DISPATCH(AtomicSwap, (init)(destruct)(initiate)(add)(redeem)(refund))