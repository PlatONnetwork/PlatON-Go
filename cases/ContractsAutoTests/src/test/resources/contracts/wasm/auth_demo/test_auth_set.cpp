#define TESTNET
#include <platon/platon.hpp>

#include <string>

using namespace platon;

CONTRACT TestAuthSet : public platon::Contract {
public:
        ACTION void init(const std::string &addr) {
                platon::set_owner(addr);
                DEBUG("Set owner", "owner", platon::owner().toString());
        }

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

ACTION void test_owner_p(const std::string& addr) {
        DEBUG("owner is ", "owner", platon::owner().toString(), "addr", addr);
        if (platon::is_owner(addr)) {
                DEBUG("Caller is owner", "addr",
                      platon::platon_origin_caller().toString(), "owner",
                      platon::owner().toString());
        } else {
                DEBUG("Caller is not owner", "addr",
                      platon::platon_origin_caller().toString(), "owner",
                      platon::owner().toString());
        }
}
};

PLATON_DISPATCH(TestAuthSet, (init)(test_owner)(test_owner_p));
