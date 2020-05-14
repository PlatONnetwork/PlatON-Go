#define TESTNET
#undef NDEBUG
#include "platon/contract.hpp"
//#include "platon/common.h"
#include "platon/assert.hpp"
#include "platon/dispatcher.hpp"
#include "platon/panic.hpp"
#include "platon/print.hpp"
#include "platon/storage.hpp"
#include "platon/storagetype.hpp"

namespace platon {

CONTRACT SolSimulation : public platon::Contract {
 public:
  /*
 *  solidity 合约存储有固定32字节key-value 组成，写入的数据有32条，每条数据的相应字节为非零，其他为零
 *  step 是将32个数据分成几组，如分成8组则step为4，fragment 为从第几组开始写入
 */
  ACTION void init(size_t step, size_t fragment) {
    platon_assert(32 % step == 0, "step is illegal:", step);
    platon_assert(step * fragment < 32, "fragment is illegal", "step", step, "fragment", fragment);
    *counter_ = 1;
    for (size_t i = 0; i < 32; i++) {
      write(i, *counter_);
    }
    *step_ = step;
    *fragment_ = fragment;
    DEBUG("fragment", fragment, "step", step);
  }

  ACTION void action() {
    DEBUG("counter:", *counter_, "step:", *step_, "fragment:", *fragment_, "scope:", *step_ * (*fragment_), *step_ * (*fragment_ + 1));
    for (size_t i = *step_ * (*fragment_); i < *step_ * (*fragment_ + 1); i++) {
      check(i, *counter_ % 32);
    }

    *fragment_ = (*step_*(*fragment_ + 1)  ) >= 32 ? 0 : (*fragment_ + 1);
    DEBUG("fragment", *fragment_);

    *counter_ = (*counter_+1) % 255;
    for (size_t i = *step_ * (*fragment_); i < *step_ * (*fragment_ + 1); i++) {
      write(i, *counter_ % 32);
    }
    DEBUG("counter:", *counter_, "step:", *step_, "fragment:", *fragment_);
  }

  ACTION void debug() {
    println("counter:", *counter_, "step:", *step_, "fragment:", *fragment_);
  }
  
  CONST uint64_t getCounter() {
    return *counter_;
  }
  
  CONST size_t getStep() {
    return *step_;
  }
  
  CONST size_t getFragment() {
    return *fragment_;
  }

 private:
  void write(size_t pos, uint64_t num) {
    uint8_t key[32] = {0};
    key[pos] = pos;
    uint8_t value[32] = {0};
    value[pos] = num;
    platon_set_state(key, sizeof(key), value, sizeof(value));
  }

  void check(size_t pos, uint64_t except) {
    uint8_t key[32] = {0};
    key[pos] = pos;
    uint8_t value[32] = {0};
    size_t length = platon_get_state_length(key, sizeof(key));
    platon_assert(length == 32, "get state length:", length, "except:32");
    platon_get_state(key, sizeof(key), value, sizeof(value));


    for (size_t i = 0; i < 32; i++) {
      if (i == pos) {
        platon_assert(value[i] == except, "get value pos:", pos,
                      "value:", value[i], "except:", except);
      } else {
        platon_assert(value[i] == 0, "get value pos:", pos, "value:", value[i],
                      "except:0");
      }
    }
  }

 private:
  StorageType<"counter"_n, uint64_t> counter_;
  StorageType<"section"_n, size_t> step_;
  StorageType<"fragment"_n, size_t> fragment_;
};


PLATON_DISPATCH(SolSimulation, (init)(action)(debug))
}

