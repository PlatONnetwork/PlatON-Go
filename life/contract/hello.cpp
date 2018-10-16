#include <stdlib.h>
#include <string.h>
#include <print.hpp>

#define PLATON_ABI(NAME, MEMBER)
PLATON_ABI(platon::Token, transfer)

namespace platon {
    typedef char* address;

    class Contract {
        public:
            Contract(){}
            virtual void init() = 0;
    };

    class Token : public Contract {
        public:
            Token(){}
            virtual void init() {
            }
            void test(){
                char c[10000];
                for (unsigned long i = 0; i < sizeof(c)/sizeof(char); i++){
                    c[i] = i;
                }
            }
        public:
            int transfer(address from, address to, int asset) {
                // char a[1000];
                // a[999]= 11;
                // char b[1000];
                // b[999]= 22;
                // char c[1000];
                // c[999]= 33;
                print_f("from:% to:% asset: % \n", from, to, asset);
                char b[10000];
                for (unsigned long i = 0; i < sizeof(b)/sizeof(char); i++){
                    b[i] = i;
                }
                test();
                strlen(from);
                atoi("sdf");
		return 88;
                // print("hello");
            }
    };
}