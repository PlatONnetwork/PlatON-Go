#!/bin/bash
VERSION=$1
CONTRACT=$2
TARGET=$3
SOLC=solc-${VERSION}
# echo "hello"
# echo "Choose solc version: ${VERSION}"
echo "Source solidity contract path: ${CONTRACT}"
echo "Compiled abi/bytecode file target path: ${CONTRACT}"
# echo "Enter solc binary dir ....."
cd ../solc
chmod a+x solc-${VERSION}
#if [ ! -f "$SOLC" ]; then
# echo "${SOLC} does not exist, pull it from server......"
# wget https://github.com/ethereum/solidity/releases/download/v${VERSION}/solc-static-linux
# mv solc-static-linux solc-${VERSION}
# chmod a+x solc-${VERSION}
#fi
# echo "Run solc command to compile contract ...."
version_num=0
array=(${VERSION//./ })
sum=0
len=${#array[@]}
for(( i=0;i<$len;i++))
do
    let sum+=$[10**($len-i)*${array[i]}]
done;
if [ "$sum" -ge "630" ]; then
  ./solc-${VERSION} -o ${TARGET} --evm-version istanbul --bin --abi --overwrite ${CONTRACT}
else
  ./solc-${VERSION} -o ${TARGET} --bin --abi --overwrite ${CONTRACT}
fi
