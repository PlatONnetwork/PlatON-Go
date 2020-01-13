pragma solidity ^0.5.13;

library BaseStaticLibrary {
    function compare(uint self, uint value) internal pure returns (bool)
    {
        if (self<value)
            return false;
        return true;
    }
}
