#include <stdlib.h>
#include <string.h>
#include <platon.hpp>


namespace platon {
    typedef char* address;


    class Token : public Contract {
        public:
            Token(){}
            virtual void init() {
            }

        public:
            void transfer(address from, char to[20], int asset) {
                print_f("from:% to:% asset: % \n", from, to, asset);
            }
    };

}

PLATON_ABI(platon::Token, transfer)
//platon autogen begin
extern "C" { 
void transfer(char * from,char * to,int asset) {
platon::Token Token_platon;
Token_platon.transfer(from,to,asset);
}
void init() {
platon::Token Token_platon;
Token_platon.init();
}

}
//platon autogen end