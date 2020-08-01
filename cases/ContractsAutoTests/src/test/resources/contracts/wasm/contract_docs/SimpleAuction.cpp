#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * 拍卖合约
 * */

CONTRACT SimpleAuction : public platon::Contract{

    public:
        PLATON_EVENT1(HighestBidIncreased,Address,u128)
        PLATON_EVENT1(AuctionEnded,Address,uint64_t)

    private:
        platon::StorageType<"beneficiary"_n,Address> beneficiary;
        platon::StorageType<"auctionEnd"_n,uint64_t> auctionEnd;
        platon::StorageType<"highestBidder"_n,Address> highestBidder;
        platon::StorageType<"amount"_n,uint64_t> amount;
        platon::StorageType<"highestBid"_n,u128> highestBid;
        platon::StorageType<"value"_n,u128> callvalue;
        platon::StorageType<"mappending"_n,std::map<Address,uint64_t>> pendingReturns;
        platon::StorageType<"ended"_n,bool> ended;

    public:
        ACTION void init(uint64_t _biddingTime, Address _beneficiary)
        {
            beneficiary.self() = _beneficiary;
            auctionEnd.self() = _biddingTime + platon_timestamp();
        }

        /// 对拍卖进行出价，具体的出价随交易一起发送
        /// 如果没有在拍卖中胜出，则返还出价
        ACTION void bid()
        {
            callvalue.self() = platon_call_value();
            if(platon_timestamp() < auctionEnd.self())
            {
                platon_panic();
            }

            if(callvalue.self() > highestBid.self())
            {
                platon_panic();
            }

            if(highestBid.self() != 0)
            {
                pendingReturns.self()[highestBidder.self()] += highestBid.self();
            }

            highestBidder.self() = platon_caller();
            highestBid.self() = callvalue.self();
            PLATON_EMIT_EVENT1(HighestBidIncreased, platon_caller(), callvalue.self());
        }

        /// 取回出价（当该出价已被超越）
        ACTION bool withdraw()
        {
            amount.self() = pendingReturns.self()[platon_caller()];
            if(amount.self() > 0)
            {
                pendingReturns.self()[platon_caller()] = 0;

                if(platon_transfer(Address(platon_caller()), Energon(amount.self())) != 0)
                {
                    pendingReturns.self()[platon_caller()] = amount.self();
                    return false;
                }
            }
            return true;
        }

        /// 结束拍卖，并把最高的出价发送给受益人
        ACTION void auctionEnded()
        {
            // 1. 条件
            if(platon_timestamp() >= auctionEnd.self())
            {
                platon_revert();
            }

            // 2. 生效
            if(!ended) {
                platon_revert();
            }

            //3.交互
            ended.self() = true;
            PLATON_EMIT_EVENT1(AuctionEnded, highestBidder.self(), highestBid.self());
            platon_transfer(Address(beneficiary.self()), Energon(amount.self()));
        }

};

PLATON_DISPATCH(SimpleAuction, (init)(bid)(withdraw)(auctionEnded))
