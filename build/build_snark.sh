#!/usr/bin/env bash

DISTRO=""
Get_Dist_Name()
{
    if grep -Eqii "CentOS" /etc/issue || grep -Eq "CentOS" /etc/*-release; then
        DISTRO='CentOS'
        PM='yum'
    elif grep -Eqi "Red Hat Enterprise Linux Server" /etc/issue || grep -Eq "Red Hat Enterprise Linux Server" /etc/*-release; then
        DISTRO='RHEL'
        PM='yum'
    elif grep -Eqi "Aliyun" /etc/issue || grep -Eq "Aliyun" /etc/*-release; then
        DISTRO='Aliyun'
        PM='yum'
    elif grep -Eqi "Fedora" /etc/issue || grep -Eq "Fedora" /etc/*-release; then
        DISTRO='Fedora'
        PM='yum'
    elif grep -Eqi "Debian" /etc/issue || grep -Eq "Debian" /etc/*-release; then
        DISTRO='Debian'
        PM='apt'
    elif grep -Eqi "Ubuntu 18.04" /etc/issue || grep -Eq "Ubuntu 18.04" /etc/*-release; then
        DISTRO='Ubuntu18'
        PM='apt'
    elif grep -Eqi "Ubuntu 16.04" /etc/issue || grep -Eq "Ubuntu 16.04" /etc/*-release; then
        DISTRO='Ubuntu16'
        PM='apt'
    elif grep -Eqi "Raspbian" /etc/issue || grep -Eq "Raspbian" /etc/*-release; then
        DISTRO='Raspbian'
        PM='apt'
    else
        DISTRO='unknow'
    fi
#    echo $DISTRO;
}
Get_Dist_Name
echo $DISTRO

if [ ! -f "build/build_snark.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

root=`pwd`
root=$root/life/resolver

if [ "`ls $root/vc`" = "" ]; then
    git submodule add https://github.com/PlatONnetwork/libcsnark.git life/resolver/vc
fi

if [ "`ls $root/vc/build`" = "" ]; then
    cd $root/vc
    mkdir -p build
    git submodule update --init --recursive
fi

if [ "`ls $root/libcsnark`" = "" ]; then
    cd $root
    mkdir -p libcsnark
fi

# Build vc
SF_BUILD=$root/vc/build
MAKE="make"
if [ `expr substr $(uname -s) 1 5` = "Linux" ]; then
    if [ "$DISTRO" = "Ubuntu16" ]; then
        echo "Ubuntu 16.04 install lib"
        sudo apt-get install llvm-6.0-dev llvm-6.0 libclang-6.0-dev
        sudo apt-get install libgmpxx4ldbl libgmp-dev libprocps4-dev
        sudo apt-get install libboost-all-dev libssl-dev
    elif [ "$DISTRO" = "Ubuntu18" ]; then
        echo "Ubuntu 18.04 install lib"
        sudo apt-get install llvm-6.0-dev llvm-6.0 libclang-6.0-dev
        sudo apt-get install libgmpxx4ldbl libgmp-dev libprocps-dev
        sudo apt-get install libboost-all-dev libssl-dev
    elif [ "$DISTRO" = "CentOS" ]; then
            sudo yum install -y llvm clang gmp procps
    else
        echo "not support system $DISTRO"
    fi
else
    echo "not support system $(uname -s)"
    exit 0
fi

cd $SF_BUILD
#$MAKE clean
cmake ../ -DMONTGOMERY_OUTPUT=OFF -DBINARY_OUTPUT=OFF
$MAKE
cp ./src/libcsnark.a ../../libcsnark/libcsnark.a
cp ./depends/libsnark/depends/libff/libff/libff.a ../../libcsnark/libff.a
cp ./depends/libsnark/libsnark/libsnark.a ../../libcsnark/libsnark.a

