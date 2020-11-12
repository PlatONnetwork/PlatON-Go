#define TESTNET
#include <platon/platon.hpp>

using namespace platon;

CONTRACT TestAuth : public platon::Contract {
public:
  ACTION void init() { platon::set_owner(); }

  ACTION void test_owner() {
    DEBUG("owner is ", "owner", platon::owner().toString());
    if (platon::is_owner()) {
      DEBUG("Caller is owner", "caller",
            platon::platon_origin_caller().toString(), "owner",
            platon::owner().toString());
    } else {
      DEBUG("Caller is not owner", "caller",
            platon::platon_origin_caller().toString(), "owner",
            platon::owner().toString());
    }
  }

  ACTION void add_to_whitelist(const std::string& addr) {
    DEBUG("add to whitelist", "addr", addr);
    SysWhitelist w;
    w.Add(addr);
  }

  ACTION void is_in_whitelist(const std::string &addr) {
    DEBUG("check in whitelist", "addr", addr);
    SysWhitelist w;
    if (w.Exists(addr)) {
      DEBUG("address in whitelist", "addr", addr);
    }else {
      DEBUG("address not in whitelist", "addr", addr);
    }
  }

  ACTION void delete_from_whitelist(const std::string& addr) {
    DEBUG("Delete address from whitelist", "addr", addr);
    SysWhitelist w;
    w.Delete(addr);
  }
};

PLATON_DISPATCH(TestAuth,
                (init)(test_owner)(add_to_whitelist)(is_in_whitelist)(delete_from_whitelist));
