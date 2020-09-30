#include <platon/platon.hpp>

class Gas {
 private:
  const char* name_ = nullptr;
  uint64_t gas_;

  PLATON_EVENT1(GasUsed, const std::string &, uint64_t)

 public:
  Gas(const char *name) : name_(name), gas_(platon_gas()) {}
  ~Gas()
  {
    emit();
  }

  void Reset(const char *name) {
    emit();
    name_ = name;
    gas_ = platon_gas();
  }

  void emit() {
    uint64_t cost = gas_ - platon_gas();
    PLATON_EMIT_EVENT1(GasUsed, name_, cost);
  }
};


// SIMULATOR
#define SIMULATOR_DEFINE(NAME)                                  \
  struct NAME##_simulator {                                     \
    using func_type = decltype(&NAME);                          \
    NAME##_simulator(const char* origin) : gasGuard_(origin) {} \
    func_type get_func() const { return &NAME; }                \
    Gas gasGuard_;                                              \
  };

#define NAMESPACE_SIMULATOR_DEFINE(NAMESPACE, NAME) \
  namespace NAMESPACE {                             \
  SIMULATOR_DEFINE(NAME)                            \
  }

#define SIMULATOR_CALL(NAME) NAME##_simulator(#NAME).get_func()


// test function

// platon_gas_price
SIMULATOR_DEFINE(platon_gas_price)
#define platon_gas_price SIMULATOR_CALL(platon_gas_price)

// platon_block_hash
SIMULATOR_DEFINE(platon_block_hash)
#define platon_block_hash SIMULATOR_CALL(platon_block_hash)

// platon_block_number
SIMULATOR_DEFINE(platon_block_number)
#define platon_block_number SIMULATOR_CALL(platon_block_number)

// platon_gas_limit
SIMULATOR_DEFINE(platon_gas_limit)
#define platon_gas_limit SIMULATOR_CALL(platon_gas_limit)

// platon_gas
SIMULATOR_DEFINE(platon_gas)
#define platon_gas SIMULATOR_CALL(platon_gas)

// platon_timestamp
SIMULATOR_DEFINE(platon_timestamp)
#define platon_timestamp SIMULATOR_CALL(platon_timestamp)

// platon_coinbase
SIMULATOR_DEFINE(platon_coinbase)
#define platon_coinbase SIMULATOR_CALL(platon_coinbase)

// platon_balance
SIMULATOR_DEFINE(platon_balance)
#define platon_balance SIMULATOR_CALL(platon_balance)

// platon_origin
SIMULATOR_DEFINE(platon_origin)
#define platon_origin SIMULATOR_CALL(platon_origin)

// platon_caller
SIMULATOR_DEFINE(platon_caller)
#define platon_caller SIMULATOR_CALL(platon_caller)

// platon_call_value
SIMULATOR_DEFINE(platon_call_value)
#define platon_call_value SIMULATOR_CALL(platon_call_value)

// platon_address
SIMULATOR_DEFINE(platon_address)
#define platon_address SIMULATOR_CALL(platon_address)

// platon_sha3
SIMULATOR_DEFINE(platon_sha3)
#define platon_sha3 SIMULATOR_CALL(platon_sha3)

// platon_caller_nonce
SIMULATOR_DEFINE(platon_caller_nonce)
#define platon_caller_nonce SIMULATOR_CALL(platon_caller_nonce)

// platon_transfer
SIMULATOR_DEFINE(platon_transfer)
#define platon_transfer SIMULATOR_CALL(platon_transfer)

// platon_set_state
SIMULATOR_DEFINE(platon_set_state)
#define platon_set_state SIMULATOR_CALL(platon_set_state)

// platon_get_state_length
SIMULATOR_DEFINE(platon_get_state_length)
#define platon_get_state_length SIMULATOR_CALL(platon_get_state_length)

// platon_get_state
SIMULATOR_DEFINE(platon_get_state)
#define platon_get_state SIMULATOR_CALL(platon_get_state)

// platon_get_input_length
SIMULATOR_DEFINE(platon_get_input_length)
#define platon_get_input_length SIMULATOR_CALL(platon_get_input_length)

// platon_get_input
SIMULATOR_DEFINE(platon_get_input)
#define platon_get_input SIMULATOR_CALL(platon_get_input)

// platon_get_call_output_length
SIMULATOR_DEFINE(platon_get_call_output_length)
#define platon_get_call_output_length SIMULATOR_CALL(platon_get_call_output_length)

// platon_get_call_output
SIMULATOR_DEFINE(platon_get_call_output)
#define platon_get_call_output SIMULATOR_CALL(platon_get_call_output)

// platon_return
SIMULATOR_DEFINE(platon_return)
#define platon_return SIMULATOR_CALL(platon_return)

// platon_revert
SIMULATOR_DEFINE(platon_revert)
#define platon_revert SIMULATOR_CALL(platon_revert)

// platon_panic
SIMULATOR_DEFINE(platon_panic)
#define platon_panic SIMULATOR_CALL(platon_panic)

// platon_debug
SIMULATOR_DEFINE(platon_debug)
#define platon_debug SIMULATOR_CALL(platon_debug)

// platon_call
SIMULATOR_DEFINE(platon_call)
#define platon_call SIMULATOR_CALL(platon_call)

// platon_delegate_call
SIMULATOR_DEFINE(platon_delegate_call)
#define platon_delegate_call SIMULATOR_CALL(platon_delegate_call)

// platon_destroy
SIMULATOR_DEFINE(platon_destroy)
#define platon_destroy SIMULATOR_CALL(platon_destroy)

// platon_migrate
SIMULATOR_DEFINE(platon_migrate)
#define platon_migrate SIMULATOR_CALL(platon_migrate)

// platon_clone_migrate
SIMULATOR_DEFINE(platon_clone_migrate)
#define platon_clone_migrate SIMULATOR_CALL(platon_clone_migrate)

// platon_event
SIMULATOR_DEFINE(platon_event)
#define platon_event SIMULATOR_CALL(platon_event)

// platon_ecrecover
SIMULATOR_DEFINE(platon_ecrecover)
#define platon_ecrecover SIMULATOR_CALL(platon_ecrecover)

// platon_ripemd160
SIMULATOR_DEFINE(platon_ripemd160)
#define platon_ripemd160 SIMULATOR_CALL(platon_ripemd160)

// platon_sha256
SIMULATOR_DEFINE(platon_sha256)
#define platon_sha256 SIMULATOR_CALL(platon_sha256)

// rlp_u128_size
SIMULATOR_DEFINE(rlp_u128_size)
#define rlp_u128_size SIMULATOR_CALL(rlp_u128_size)

// platon_rlp_u128
SIMULATOR_DEFINE(platon_rlp_u128)
#define platon_rlp_u128 SIMULATOR_CALL(platon_rlp_u128)

// rlp_bytes_size
SIMULATOR_DEFINE(rlp_bytes_size)
#define rlp_bytes_size SIMULATOR_CALL(rlp_bytes_size)

// platon_rlp_bytes
SIMULATOR_DEFINE(platon_rlp_bytes)
#define platon_rlp_bytes SIMULATOR_CALL(platon_rlp_bytes)

// rlp_list_size
SIMULATOR_DEFINE(rlp_list_size)
#define rlp_list_size SIMULATOR_CALL(rlp_list_size)

// platon_rlp_list
SIMULATOR_DEFINE(platon_rlp_list)
#define platon_rlp_list SIMULATOR_CALL(platon_rlp_list)

// platon_contract_code_length
SIMULATOR_DEFINE(platon_contract_code_length)
#define platon_contract_code_length SIMULATOR_CALL(platon_contract_code_length)

// platon_contract_code
SIMULATOR_DEFINE(platon_contract_code)
#define platon_contract_code SIMULATOR_CALL(platon_contract_code)

// platon_deploy
SIMULATOR_DEFINE(platon_deploy)
#define platon_deploy SIMULATOR_CALL(platon_deploy)

// platon_clone
SIMULATOR_DEFINE(platon_clone)
#define platon_clone SIMULATOR_CALL(platon_clone)