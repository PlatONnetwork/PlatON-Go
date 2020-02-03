pragma solidity ^0.5.13;

contract StructDataType {
    struct A {
        mapping(uint=>uint) m;
    }
    struct B {
        mapping(uint=>uint) m;
        uint x;
    }
    struct C {
        mapping(uint=>uint)[] ma;
    }
    struct D {
        A[] a;
    }
    A storageA;
    B storageB;
    C storageC;
    D storageD;
    constructor() public {
        storageA.m[1] = 2;
        storageB.m[3] = 4;
        storageB.x = 5;
        storageC.ma.length = 6;
        storageD.a.length = 7;
    }

    uint storageAmValue;
    uint storageBxValue;
    uint memoryBxValue;
    uint storageCmaLen;
    uint memoryD1aLen;
    uint memoryD2aLen;


    function run() public returns (uint, uint, uint, uint, uint, uint) {
        A memory memoryA = A();
        B memory memoryB = B(2);
        C memory memoryC = C();
        D memory memoryD1 = D(new A[](9));
        D memory memoryD2 = storageD;
        storageA = memoryA;
        storageB = memoryB;
        storageC = memoryC;

        storageAmValue = storageA.m[1];
        storageBxValue = storageB.x;
        memoryBxValue = memoryB.x;
        storageCmaLen = storageC.ma.length;
        memoryD1aLen =  memoryD1.a.length;
        memoryD2aLen = memoryD2.a.length;

        return (
        storageA.m[1],
        storageB.x,
        memoryB.x,
        storageC.ma.length,
        memoryD1.a.length,
        memoryD2.a.length
        );
    }

    //执行run后的结果
    function getRunValue() public view returns(uint, uint, uint, uint, uint, uint) {
        return (storageAmValue,storageBxValue,memoryBxValue,storageCmaLen,memoryD1aLen,memoryD2aLen);
    }
}
