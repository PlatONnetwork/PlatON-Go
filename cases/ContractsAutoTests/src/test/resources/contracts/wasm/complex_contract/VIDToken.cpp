#define TESTNET
#include <platon/platon.hpp>

using namespace platon;

CONTRACT VIDToken : public Contract {
 public:
  PLATON_EVENT0(TransferEv, Address, Address, u128);
  PLATON_EVENT0(ApprovalEv, Address, Address, u128);
  PLATON_EVENT0(BurnEv, Address, u128);
  PLATON_EVENT0(FreezeEv, Address, bool);
  PLATON_EVENT0(ValidateFileEv, u128, std::string);
  PLATON_EVENT0(ValidatePublisherEv, Address, bool, std::string);
  PLATON_EVENT0(ValidateWalletEv, Address, bool, std::string);
  PLATON_EVENT0(LogEventEv, u128, std::string);

 public:
  PLATON_EVENT0(PauseEv);
  PLATON_EVENT0(UnpauseEv);

  ACTION void pause() {
    platon_assert(is_owner(), "O1- Owner only function");
    platon_assert(!paused_.get());

    paused_.self() = true;
    PLATON_EMIT_EVENT0(PauseEv);
  }

  ACTION void unpause() {
    platon_assert(is_owner(), "O1- Owner only function");
    platon_assert(paused_.get());

    paused_.self() = false;
    PLATON_EMIT_EVENT0(UnpauseEv);
  }

 public:
  ACTION void init() {
    set_owner();
    paused_.self() = false;
    validation_price_.self() = 7000000000000000000;
    validation_wallet_.self() = platon_caller();
    verify_wallet_.self()[platon_caller()] = true;
    total_supply_.self() = kInitialSupply * kUnit;
    balances_.self()[platon_caller()] = total_supply_.get();
    PLATON_EMIT_EVENT0(TransferEv, Address(), owner(), kInitialSupply);

    DEBUG("init contract", "owner", owner().toString(), "balance",
          balances_.get().at(platon_caller()));
  }

  ACTION bool Transfer(const std::string& to_addr, u128 value) {
    auto sender = platon_caller();
    Address to;
    auto address_info = make_address(to_addr);
    if(address_info.second){
      to = address_info.first;
    }
    DEBUG("transfer", "to", to.toString());
    platon_assert(!paused_.get());
    platon_assert(to != sender, "T1- Recipient can not be the same as sender");
    bool empty_addr = false;
    if (to) {
      empty_addr = true;
    }
    platon_assert(empty_addr, "T2- Please check the recipient address");
    platon_assert(Balance(sender) >= value,
                  "T3- The balance of sender is too low");
    platon_assert(!Frozen(sender), "T4- The wallet of sender is frozen");
    platon_assert(!Frozen(to), "T5- The wallet of recipient is frozen");

    balances_.self()[sender] = Balance(sender) - value;
    if (balances_.self().find(to) != balances_.self().end()) {
      balances_.self()[to] = Balance(to) + value;
    } else {
      balances_.self()[to] = value;
    }

    PLATON_EMIT_EVENT0(TransferEv, sender, to, value);
    return true;
  }

  ACTION bool TransferFrom(const std::string& from_addr,
                           const std::string& to_addr, u128 value) {
    Address from;
    Address to;
    auto address_info = make_address(from_addr);
    if(address_info.second){
      from = address_info.first;
    }

    auto address_info2 = make_address(to_addr);
    if(address_info2.second){
      to = address_info2.first;
    }

    auto sender = platon_caller();
    bool empty_addr = false;
    if (to) {
      empty_addr = true;
    }
    platon_assert(!paused_.get());
    platon_assert(empty_addr, "TF1- Please check the recipient address");
    platon_assert(Balance(from) >= value,
                  "TF2- The balance of sender is too low");
    platon_assert(Allowed(from, sender) >= value,
                  "TF3- The allowed of sender if too low");
    platon_assert(!Frozen(from), "TF4- The wallet of sender is frozen");
    platon_assert(!Frozen(to), "TF5- The wallet of recipient is frozen");

    balances_.self()[from] = Balance(from) - value;
    balances_.self()[to] = Balance(to) + value;

    allowed_.self()[from][sender] = allowed_.self()[from][sender] - value;

    PLATON_EMIT_EVENT0(TransferEv, from, to, value);
    return true;
  }

  CONST u128 BalanceOf(const std::string& owner_addr) {
    DEBUG("balance of ", "owner_addr", owner_addr);

    Address to;
    auto address_info = make_address(owner_addr);
    if(address_info.second){
      to = address_info.first;
    }

    return Balance(to);
  }

  ACTION bool Approve(const std::string& spender_addr, u128 value) {
    Address spender;
    auto address_info = make_address(spender_addr);
    if(address_info.second){
      spender = address_info.first;
    }

    platon_assert(!paused_.get());
    auto sender = platon_caller();
    platon_assert(!paused_.get());
    platon_assert(value == 0 || (Allowed(sender, spender) == 0),
                  "A1- Reset allowance to first");

    auto& allowed = allowed_.self();
    if (allowed.find(sender) == allowed.end()) {
      std::map<Address, u128> allow;
      allow[spender] = value;
      allowed[sender] = allow;
    } else {
      allowed[sender][spender] = value;
    }

    PLATON_EMIT_EVENT0(ApprovalEv, sender, spender, value);
    return 0;
  }

  ACTION bool IncreaseApproval(const std::string& spender_addr,
                               u128 added_value) {
    Address spender;
    auto address_info = make_address(spender_addr);
    if(address_info.second){
      spender = address_info.first;
    }
    platon_assert(!paused_.get());
    auto sender = platon_caller();

    auto& allowed = allowed_.self();
    if (allowed.find(sender) == allowed.end()) {
      std::map<Address, u128> allow;
      allow[spender] = added_value;
      allowed[sender] = allow;
    } else {
      auto& allow = allowed.at(sender);
      if (allow.find(spender) != allow.end()) {
        allow[spender] = allow.at(spender) + added_value;
      } else {
        allow[spender] = added_value;
      }
    }

    PLATON_EMIT_EVENT0(ApprovalEv, sender, spender,
                       allowed_.self().at(sender).at(spender));
    return true;
  }

  ACTION bool DecreaseApproval(const std::string& spender_addr,
                               u128 subtracted_value) {
    Address spender;
    auto address_info = make_address(spender_addr);
    if(address_info.second){
      spender = address_info.first;
    }
    platon_assert(!paused_.get());
    auto sender = platon_caller();

    auto& allowed = allowed_.self();
    if (allowed.find(sender) == allowed.end()) {
      std::map<Address, u128> allow;
      allow[spender] = 0;
      allowed[sender] = allow;
    } else {
      auto& allow = allowed.at(sender);
      if (allow.find(spender) != allow.end()) {
        auto c = allow.at(spender) - subtracted_value;
        DEBUG("decrease approval", "c", c, "allow", allow.at(spender), "sub",
              subtracted_value);
        platon_assert(allow.at(spender) >= subtracted_value &&
                      c <= allow.at(spender));
        allow[spender] = c;
      } else {
        allow[spender] = 0;
      }
    }

    PLATON_EMIT_EVENT0(ApprovalEv, sender, spender,
                       allowed_.self()[sender][spender]);
    return true;
  }

  CONST u128 Allowance(const std::string& owner, const std::string& spender) {
    Address addr1;
    auto address_info1 = make_address(owner);
    if(address_info1.second){
      addr1 = address_info1.first;
    }

    Address addr2;
    auto address_info2 = make_address(spender);
    if(address_info2.second){
      addr2 = address_info2.first;
    }


    return Allowed(addr1, addr2);
  }

  struct TKN {
    Address sender;
    u128 value;
    bytes data;
    byte sig[4];
  };

  ACTION bool TokenFallback(const std::string& from_addr, u128 value,
                            const std::string& data) {
    Address from;
    auto address_info = make_address(from_addr);
    if(address_info.second){
      from = address_info.first;
    }
    TKN tkn;
    tkn.sender = from;
    tkn.value = value;
    tkn.data.insert(tkn.data.begin(), data.begin(), data.end());
    uint32_t u = uint32_t(data[3]) + (uint32_t(data[2]) << 8) +
                 (uint32_t(data[1]) << 16) + (uint32_t(data[0]) << 24);
    memcpy(tkn.sig, &u, sizeof(u));
    return true;
  }

  ACTION void TransferToken(const std::string& token_addr_s, u128 tokens) {
    Address token_addr;
    auto address_info = make_address(token_addr_s);
    if(address_info.second){
      token_addr = address_info.first;
    }
    platon_assert(is_owner(), "O1- Owner only function");
    platon_assert(owner() != token_addr,
                  "T1- Recipient can not be the same as sender");
    platon_assert(Balance(token_addr) >= tokens,
                  "T3- The balance of sender is too low");
    platon_assert(!Frozen(token_addr), "T4- The wallet of sender is frozen");

    balances_.self()[token_addr] = Balance(token_addr) - tokens;
    balances_.self()[owner()] = Balance(owner()) + tokens;

    PLATON_EMIT_EVENT0(TransferEv, token_addr, owner(), tokens);
  }

  ACTION bool Burn(u128 value) {
    platon_assert(is_owner(), "O1- Owner only function");

    auto sender = platon_caller();
    platon_assert(value <= Balance(sender),
                  "B1- The balance of burner is too low");

    balances_.self()[sender] = Balance(sender) - value;
    total_supply_.self() = total_supply_.get() - value;

    PLATON_EMIT_EVENT0(BurnEv, sender, value);
    return true;
  }

  ACTION bool Freeze(const std::string& addr_s, bool state) {
    Address addr;
    auto address_info = make_address(addr_s);
    if(address_info.second){
      addr = address_info.first;
    }
    platon_assert(is_owner(), "O1- Owner only function");

    frozen_account_.self()[addr] = state;

    PLATON_EMIT_EVENT0(FreezeEv, addr, state);
    return true;
  }

  ACTION bool ValidatePublisher(const std::string& addr_s, bool state,
                                const std::string& publisher) {
    Address addr;
    auto address_info = make_address(addr_s);
    if(address_info.second){
      addr = address_info.first;
    }
    platon_assert(is_owner(), "O1- Owner only function");

    verify_publisher_.self()[addr] = state;

    PLATON_EMIT_EVENT0(ValidatePublisherEv, addr, state, publisher);

    return true;
  }

  ACTION bool ValidateWallet(const std::string& addr_s, bool state,
                             const std::string& wallet) {
    Address addr;
    auto address_info = make_address(addr_s);
    if(address_info.second){
      addr = address_info.first;
    }
    platon_assert(is_owner(), "O1- Owner only function");

    verify_wallet_.self()[addr] = state;

    PLATON_EMIT_EVENT0(ValidateWalletEv, addr, state, wallet);

    return true;
  }

  ACTION bool ValidateFile(const std::string& to_addr, u128 payment,
                           const std::string& data, bool store, bool log) {
    DEBUG("validate file", "to_addr", to_addr, "payment", payment, "data", data,
          "store", store, "log", log, "price", validation_price_.get());
    Address to;
    auto address_info = make_address(to_addr);
    if(address_info.second){
      to = address_info.first;
    }
    auto sender = platon_caller();
    platon_assert(!paused_.get());
    platon_assert(payment >= validation_price_.get(),
                  "V1- Insufficient payment provided");

    bool verify_publisher = false;
    if (verify_publisher_.self().find(sender) != verify_publisher_.self().end()) {
      verify_publisher = verify_publisher_.self().at(sender);
    }
    platon_assert(verify_publisher, "V2- Unverified publisher address");

    bool frozen = false;
    if (frozen_account_.self().find(sender) != frozen_account_.self().end()) {
      frozen = frozen_account_.self().at(sender);
    }
    platon_assert(!frozen, "V3- The wallet of publisher is frozen");
    platon_assert(data.size() == 64, "V4- Invalid hash provided");

    auto verify_wallet = false;
    if (verify_wallet_.self().find(to) != verify_wallet_.self().end()) {
      verify_wallet = verify_wallet_.self().at(to);
    }
    frozen = false;
    if (frozen_account_.self().find(to) != frozen_account_.self().end()) {
      frozen = true;
    }
    if (!verify_wallet || frozen) {
      to = validation_wallet_.get();
    }

    u128 index = 0;
    std::string file_hash(data);

    if (store) {
      if (file_index_.self().size() > 0) {
        platon_assert(
            file_hashes_.self().find(file_hash) == file_hashes_.self().end(),
            "V5- This hash was previously validated");
      }

      file_index_.self().push_back(file_hash);
      file_hashes_.self()[file_hash] = FStruct{file_index_.self().size() - 1};
      index = file_hashes_.self().at(file_hash).index;
    }

    if (Allowed(to, sender) >= payment) {
      allowed_.self()[to][sender] = allowed_.self()[to][sender] - payment;
    } else {
      balances_.self()[sender] = Balance(sender) - payment;
      balances_.self()[to] = Balance(to) + payment;
    }

    if (log) {
      PLATON_EMIT_EVENT0(ValidateFileEv, index, file_hash);
    }
    return true;
  }

  CONST bool VerifyFile(const std::string& file_hash) {
    if (file_index_.self().size() == 0) {
      return false;
    }

    u128 index = 0;
    if (file_hashes_.self().find(file_hash) == file_hashes_.self().end()) {
      return false;
    }
    index = file_hashes_.self().at(file_hash).index;

    auto fh = file_index_.self()[index];

    if (fh.size() != file_hash.size()) {
      return false;
    }

    for (size_t i = 0; i < fh.size(); i++) {
      if (fh[i] != file_hash[i]) {
        return false;
      }
    }
    return true;
  }

  ACTION void SetPrice(u128 new_price) {
    platon_assert(is_owner(), "O1- Owner only function");

    validation_price_.self() = new_price;
  }

  ACTION void SetWallet(const std::string& new_wallet_s) {
    Address new_wallet;
    auto address_info = make_address(new_wallet_s);
    if(address_info.second){
      new_wallet = address_info.first;
    }
    platon_assert(is_owner(), "O1- Owner only function");
    validation_wallet_.self() = new_wallet;
  }

  ACTION bool ListFiles(u128 start_at, u128 stop_at) {
    platon_assert(is_owner(), "O1- Owner only function");

    if (file_index_.self().size() == 0) {
      return false;
    }

    u128 max_index = file_index_.self().size() - 1;

    platon_assert(start_at <= max_index, "L1- Please select a valid start");
    if (stop_at > 0) {
      platon_assert(stop_at > start_at && stop_at <= max_index,
                    "L2- Please selct a valid stop");
    } else {
      stop_at = max_index;
    }

    for (u128 i = start_at; i < stop_at; i++) {
      PLATON_EMIT_EVENT0(LogEventEv, i, file_index_.get()[i]);
    }
    return true;
  }

 private:
  u128 Balance(const Address& addr) {
    const auto& balances = balances_.self();
    if (balances.find(addr) != balances.end()) {
      DEBUG("balance", "addr", addr.toString(), "balance", balances.at(addr));
      return balances.at(addr);
    }
    DEBUG("balance", "addr", addr.toString(), "balance", 0);
    return 0;
  }

  bool Frozen(const Address& addr) {
    const auto& frozen_account = frozen_account_.self();
    if (frozen_account.find(addr) != frozen_account.end()) {
      return frozen_account.at(addr);
    }
    return false;
  }

  u128 Allowed(const Address& owner, const Address& allow) {
    const auto& allowed = allowed_.self();
    if (allowed.find(owner) == allowed.end()) {
      return 0;
    }

    const auto& allower = allowed.at(owner);
    if (allower.find(allow) != allower.end()) {
      return allower.at(allow);
    }
    return 0;
  }

 private:
  StorageType<"total_supply"_n, u128> total_supply_;
  StorageType<"paused"_n, bool> paused_;

  Map<"balances"_n, Address, u128> balances_;
  Map<"allowed"_n, Address, std::map<Address, u128>> allowed_;
  Map<"frozen_account"_n, Address, bool> frozen_account_;
  Map<"verify_publisher"_n, Address, bool> verify_publisher_;
  Map<"verify_wallet"_n, Address, bool> verify_wallet_;

  struct FStruct {
    u128 index;
    PLATON_SERIALIZE(FStruct, (index));
  };
  Map<"file_hashes"_n, std::string, FStruct> file_hashes_;
  Vector<"file_index"_n, std::string> file_index_;

  // const static std::string kName = "V-ID Token";
  const static uint8_t decimals = 18;
  // const static std::string kSymbol = "VIDT";
  const static u128 kInitialSupply = 100000000;
  const static u128 kUnit = 1000000000000000000;

  StorageType<"validation_price"_n, u128> validation_price_;
  StorageType<"validation_wallet"_n, Address> validation_wallet_;
};

PLATON_DISPATCH(
    VIDToken,
    (init)(pause)(unpause)(Transfer)(TransferFrom)(BalanceOf)(Approve)(
        IncreaseApproval)(DecreaseApproval)(Allowance)(TokenFallback)(
        TransferToken)(Burn)(Freeze)(ValidatePublisher)(ValidateWallet)(
        ValidateFile)(VerifyFile)(SetPrice)(SetWallet)(ListFiles));
