pragma solidity ^0.5.13;

contract RecursiveStorageMemoryComplex {
    struct Tree {
        uint256 data;
        Tree[] children;
    }
    Tree storageTree;

    constructor() public {
        storageTree.data = 0x42;
        storageTree.children.length = 2;
        storageTree.children[0].data = 0x4200;
        storageTree.children[1].data = 0x4201;
        storageTree.children[0].children.length = 3;
        for (uint i = 0; i < 3; i++)
            storageTree.children[0].children[i].data = 0x420000 + i;
        storageTree.children[1].children.length = 4;
        for (uint i = 0; i < 4; i++)
            storageTree.children[1].children[i].data = 0x420100 + i;
    }

    uint256[] result;

    function countData(Tree memory tree) internal returns (uint256 c) {
        c = 1;
        for (uint i = 0; i < tree.children.length; i++) {
            c += countData(tree.children[i]);
        }
    }

    function copyFromTree(Tree memory tree,uint256 offset) internal returns (uint256) {
        result[offset++] = tree.data;
        for (uint i = 0; i < tree.children.length; i++) {
            offset = copyFromTree(tree.children[i],offset);
        }
        return offset;
    }

    function run() public returns (uint256[] memory) {
        Tree memory memoryTree;
        memoryTree = storageTree;
        uint256 length = countData(memoryTree);
        result = new uint256[](length);
        copyFromTree(memoryTree, 0);
        return result;
    }

    function getRunResult() public view returns(uint256[] memory){
        return result;
    }
}