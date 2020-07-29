#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * @author hudenian
 * 众筹合约
 * */

CONTRACT CrowdFunding : public platon::Contract{

    private:
        platon::StorageType<"beneficiary"_n,Address> beneficiaryAddress; //受益人地址，设置为合约创建者
        platon::StorageType<"fundgoal"_n, uint64_t> fundingGoal; //众筹目标
        platon::StorageType<"raised"_n, u128> amountRaised; //已筹集金额数量， 单位是VON
        platon::StorageType<"deadline"_n, uint64_t> deadline; //截止时间
        platon::StorageType<"price"_n, u128> price; //代币价格
        platon::StorageType<"reachflg"_n, bool> fundingGoalReached;//达成众筹目标
        platon::StorageType<"closedflg"_n, bool> crowdsaleClosed;//众筹关闭

        platon::StorageType<"balance"_n, std::map<Address, u128>> balance;//保存众筹者对捐赈的金额
        platon::StorageType<"tokenmap"_n, std::map<Address, uint64_t>> tokenMap;//保存众筹者所拥有的代币数量
         
        //记录已接收的LAT通知
        PLATON_EVENT1(transfer1,std::string,Address,uint64_t)

        //转帐通知
        PLATON_EVENT1(transfer2,std::string,Address,uint64_t,bool)

    public:
        //第一个参数为众筹目标，第二个参数为持序时间，单位为秒
        ACTION void init(uint64_t _fundingGoalInlats,uint64_t _durationInMinutes){
            beneficiaryAddress.self() = platon_caller(); //收益人为合约发起者
            fundingGoal.self() = _fundingGoalInlats;
            deadline.self() = _durationInMinutes + platon_timestamp();
            price.self() = 500; //假设一个代币500Von
            crowdsaleClosed.self() = false;
        }

        //发起众筹
        ACTION void crowdFund(){
            Address caller = platon_caller();

            //众筹开关关闭后就不能再发起
            if(crowdsaleClosed.self() ){
                platon_revert();
            }
            u128 amount = platon_call_value();

            //捐款人的金额累加
            balance.self()[caller] += amount;

            //捐款总额累加
            amountRaised.self() += amount;

            //转帐操作，转多少代币给捐款人
            tokenMap.self()[caller]  += amount /price.self();

            PLATON_EMIT_EVENT1(transfer1, "transfer", caller, amount);
        }


        //检测众筹目标是否已经达到
        ACTION void checkGoalReached(){
            if(amountRaised.self() >= fundingGoal.self()){
                fundingGoalReached.self() = true;
                PLATON_EMIT_EVENT1(transfer1,"checkGoalReached",platon_caller(),amountRaised.self());
            }
            //关闭众筹
            crowdsaleClosed.self() = true;
        }

        /**
         * 收回资金
         *
         * 检查是否达到了目标或时间限制，如果有，并且达到了资金目标，
         * 将全部金额发送给受益人。如果没有达到目标，每个贡献者都可以退回
         * 他们贡献的金额
         */
         ACTION void safeWithdrawal(){
            //如果没有达成众筹目标
            if (!fundingGoalReached.self()) {
                //获取合约调用者已捐款余额
                uint64_t amount = balance.self()[beneficiaryAddress.self()];

                if (amount > 0) {
                    //返回合约发起者所有余额
                    platon_transfer(Address(platon_caller()), Energon(balance.self()[beneficiaryAddress.self()]));
                    PLATON_EMIT_EVENT1(transfer1,"checkGoalReached",platon_caller(),balance.self()[beneficiaryAddress.self()]);
                }
            }

            //如果达成众筹目标，并且合约调用者是受益人
            if (fundingGoalReached.self() && beneficiaryAddress.self() == platon_caller()) {

                //将所有捐款从合约中给受益人
                platon_transfer(Address(platon_caller()), Energon(amountRaised.self()));
                PLATON_EMIT_EVENT1(transfer1,"checkGoalReached",platon_caller(),amountRaised.self());
            }
         }
};

PLATON_DISPATCH(CrowdFunding, (init)(crowdFund)(checkGoalReached)(safeWithdrawal))

