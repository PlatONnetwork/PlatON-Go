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
struct Key {
  std::string prefix;
  uint64_t seq;
  PLATON_SERIALIZE(Key, (prefix)(seq))
};
CONTRACT WasmStorage : public platon::Contract {
 public:
  /*
   * 测试不同数据长度，16b, 32b, 64b, 64Kb, 500Kb
   * 测试数据的修改
   * 测试数据删除
   * 测试单合约大数据量
   */
  ACTION void init() {
    auto block_num = platon_block_number();
    auto timestamp = platon_timestamp();
    for (auto k : kInternal) {
      states_->emplace(std::make_pair(k.first, block_num));
      timestamp_->emplace(std::make_pair(k.first, timestamp));
      counter_->emplace(std::make_pair(k.first, 0));
      DEBUG("kInternal:", k.first);
      println("kInternal:", k.first);
    }
    *random_ = 1234567890;
  }

  ACTION void random_data() {
     auto current = platon_block_number();
     auto timestamp = platon_timestamp();
     random(current, timestamp);
  }

  ACTION void action() {
    auto current = platon_block_number();
    auto timestamp = platon_timestamp();
    bool hit = false;
    uint64_t index = 0;
    platon_assert(counter_->size() != 0 , "counter removed");
    DEBUG("current:", current, "timestamp:", timestamp);
    for (auto k : kInternal) {
      auto number = (*states_)[k.first];
      if (current > number + k.second) {
        if ((*counter_)[k.first] != 0) {
          check(k.first, (*timestamp_)[k.first], (*counter_)[k.first], index);
        }
        uint64_t counter = ++(*counter_)[k.first];
        write(k.first, timestamp, counter, index);
        (*timestamp_)[k.first] = timestamp;
        (*states_)[k.first] = current;
        hit = true;
        break;
      }
      index++;
    }

    if (!hit) {
      random(current, timestamp);
    }
  }

  ACTION void debug() {
    for (auto k : kInternal) {
      println("internal:", k.first, "states:", (*states_)[k.first],
              "timestamp:", (*timestamp_)[k.first],
              "counter:", (*counter_)[k.first], "random:", *random_);
    }
  }

 private:
  void write(const std::string &prefix, uint64_t timestamp, uint64_t counter,
             uint64_t index) {
    DEBUG("write prefix:", prefix, "timestamp:", timestamp, "counter:", counter,
          "index:", index);

    std::vector<uint64_t> vec(32, counter);
    Key key{.prefix = prefix, .seq = 1};
    set_state(key, vec);
    uint64_t length = counter % kMaxStringLength;

    key.seq = 2;
    if (length == 0) {
      del_state(key);
    }
    std::string str(length, kNumber[index]);
    set_state(key, str);
  }

  void check(const std::string &prefix, uint64_t timestamp, uint64_t counter,
             uint64_t index) {
    DEBUG("check prefix:", prefix, "timestamp:", timestamp, "counter:", counter,
          "index:", index);
    std::vector<uint64_t> vec;
    Key key{.prefix = prefix, .seq = 1};
    get_state(key, vec);
    platon_assert(vec.size() == 32, "vector size error size:", vec.size(),
                  "except:32");
    for (auto k : vec) {
      platon_assert(k == counter, "vector is not equal value:", k,
                    "except:", counter);
    }

    uint64_t length = counter % kMaxStringLength;
    key.seq = 2;
    if (length == 0) {
      platon_assert(!has_state(key), "string had exists");
      return;
    }
    std::string str;
    get_state(key, str);
    platon_assert(str.length() == length,
                  "string length error length:", str.length(),
                  "except:", counter);
    for (auto s : str) {
      platon_assert(s == kNumber[index], "string  error s:", s,
                    "except:", kNumber[index]);
    }
  }

  void random(uint64_t number, uint64_t timestamp) {
    DEBUG("write random data:", *random_, "number:", number,
          "timestamp:", timestamp);
    uint64_t seq = *random_;
    for (auto k : kRandomSize) {
      DEBUG("write random size:", k);
      uint8_t *value = (uint8_t *)malloc(k);
      value[0] = timestamp;
      platon_set_state((uint8_t *)&seq, sizeof(seq), value, k);
      ++seq;
    }
    *random_ = seq;
    DEBUG("write random data:", *random_);
  }

 private:
  StorageType<"states"_n, std::map<std::string, uint64_t>> states_;
  StorageType<"counter"_n, std::map<std::string, uint64_t>> counter_;
  StorageType<"timestamp"_n, std::map<std::string, uint64_t>> timestamp_;
  StorageType<"random"_n, uint64_t> random_;

  const std::vector<std::pair<std::string, uint64_t>> kInternal{
      {"internal200", 200},
      {"internal100", 100},
      {"internal50", 50},
  };

  const std::vector<uint64_t> kRandomSize{16, 32, 64, 128, 256};
  const std::vector<uint64_t> kSize{16, 32, 64};
  const char kNumber[7] = {'0', '1', '2', '3', '4', '5', '6'};
  const uint64_t kMaxStringLength = 512;
};

PLATON_DISPATCH(WasmStorage, (init)(action)(random_data)(debug))
}  // namespace platon
