#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

CONTRACT GuessingWasm : public platon::Contract{

   private:
        platon::StorageType<"endblock"_n,uint64_t> endBlock; //竞猜截止块高
        platon::StorageType<"closedflg"_n, bool> guessingClosed;//竞猜是否已开奖
        platon::StorageType<"baseunit"_n, u128> baseUnit; //最小转金额
        platon::StorageType<"balance"_n, u128> balance; //竞猜总金额
        platon::StorageType<"avgmount"_n, u128> averageAmount; //每个人获奖金额

        platon::StorageType<"gussinger"_n, std::map<Address, u128>> gussingerLat;//每个竞猜者对应的金额
        platon::StorageType<"guesseridx"_n, std::map<uint64_t,Address>> indexOfgussinger;//每个竞猜者对应的下标（大于5个lat就给他分配一个随机数）
        platon::StorageType<"winnermap"_n, std::map<Address, uint64_t>> winnerMap;//中奖者对应中奖号码个数
        platon::StorageType<"indexkey"_n, uint64_t> indexKey;//自增序列
        platon::StorageType<"provector"_n,std::vector<Address>> winnerAddresses;//中奖者地址
        platon::StorageType<"createaddr"_n,Address> createAddress;//合约创建者

   public:
        //竞猜成功通知
        PLATON_EVENT1(transfer1,std::string,Address,u128)

        //记录已接收的LAT通知
        PLATON_EVENT1(transfer2,std::string,Address,u128)

        /**
         * 初始化构造函数
         *
         * @param _endBlock 竞猜截止块高（达到此块高后不能再往合约转账）
         */ 
        ACTION void init(uint64_t _endBlock){
           createAddress.self() = platon_caller();
           endBlock.self() = _endBlock; 
           indexKey.self() = 0;//默认下标为0
           baseUnit.self() = (5_LAT).Get();
           guessingClosed.self() = false;
        }

        /**
         * 竞猜(带上金额)
         */ 
        ACTION void guessingWithLat(){
           DEBUG("Guessing", "guessingWithLat", platon_call_value());
           //判断转账金额要大于5lat 
           u128 amount = platon_call_value();
           platon_assert( amount >= baseUnit.self(), "bad platon_call_value");

           uint64_t currentBlock_number = platon_block_number();//获取当前块高
           platon_assert( endBlock.self() >= currentBlock_number, "not reach end block number");

           indexOfgussinger.self()[indexKey.self()] = platon_caller();
           indexKey.self()++;

           //竞猜人的金额累加(可以投多次)
           gussingerLat.self()[platon_caller()] += amount;

           //竞猜总额累加 
           balance.self() +=amount;

           //竞猜成功通知
           PLATON_EMIT_EVENT1(transfer1,"guessing success",platon_caller(),amount); 
         }

         ACTION void draw(){
            DEBUG("Guessing", "draw", platon_caller());
            //只有合约创建者可以开奖
            if(!guessingClosed.self() && createAddress.self() == platon_caller() && indexKey.self() > 0){
                  uint64_t random = 99999;
                  uint64_t drawIndex = random%indexKey.self();

                  if(indexKey.self()<100){
                     getwinners(drawIndex,10);
                  }else if(indexKey.self()<10000){
                     getwinners(drawIndex,100);
                  }else{
                     getwinners(drawIndex,1000);
                  }

                  //每个中奖者可以分到的金额
                  averageAmount.self() = balance.self()/winnerAddresses.self().size();

                  //向中奖者转账
                  for(uint64_t j=0;j<winnerAddresses.self().size();j++){
                     //中奖者中奖票号统计
                     winnerMap.self()[winnerAddresses.self()[j]] =winnerMap.self()[winnerAddresses.self()[j]]+1;
              
                     platon_transfer(winnerAddresses.self()[j], Energon(averageAmount.self()));
                  }

                  guessingClosed.self() = true;
            }
         }


          ACTION void getwinners(uint64_t drawIndex,uint64_t times){
             uint64_t postfix = drawIndex%times;
               if(postfix ==0){
                     for(uint64_t i=0;i<indexKey.self();i++){
                        if((i-postfix)%times == 0){
                           winnerAddresses.self().push_back(indexOfgussinger.self()[i]);
                        }
                     }
               }else{
                     for(uint64_t i=0;i<indexKey.self();i++){
                        if(i%times != 0 && (i-postfix)%times == 0){
                           winnerAddresses.self().push_back(indexOfgussinger.self()[i]);
                        }
                     }
               }
          }

         //查看当前合约中的余额
         CONST u128 getBalance(){
            return balance.self();
         }

         //查看一共有几个中奖号码
         CONST uint64_t getWinnerCount(){
            return winnerAddresses.self().size();
         }

         // 获取所有中奖人地址（可能有重复，调用方可以对此进行合并）
         CONST std::vector<Address> getWinnerAddresses(){
            return winnerAddresses.self();
         }

         // 获取当前参与者中奖次数
         CONST uint64_t getMyGuessCodes(Address &address){
            return winnerMap.self()[address];
         }

         //查看当前下标
         CONST uint64_t getIndexKey(){
            return indexKey.self();
         }
 
};

PLATON_DISPATCH(GuessingWasm, (init)(guessingWithLat)(draw)(getwinners)(getBalance)(getWinnerCount)(getWinnerAddresses)(getMyGuessCodes)(getIndexKey))
