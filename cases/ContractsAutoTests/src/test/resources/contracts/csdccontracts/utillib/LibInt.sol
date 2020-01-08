pragma solidity ^0.4.12;
/**
* file LibInt.sol
* author liaoyan
* time 2016-11-29
* desc the defination of LibInt libary
*/


library LibInt {
    using LibInt for *;

    function toString(uint _self, uint width) internal returns (string _ret) {
        uint len = 0;
        uint tmp = _self;
        while (tmp > 0) {
            tmp /= 10;
            len++;
        }
        
        if (len < width) {
            len = width;
        }
        _ret = new string(len);
        
        int8 i = int8(len-1);
        while (_self > 0) {
            bytes(_ret)[uint(i)] = byte(_self%10+0x30);
            --i;
            _self /= 10;
        }
        
        while (i >= 0) {
            bytes(_ret)[uint(i)] = byte(0x30);
            --i;
        }
    }

    function toString(uint _self) internal returns (string _ret) {
        if (_self == 0) {
            return "0";
        }

        uint8 len = 0;
        uint tmp = _self;
        while (tmp > 0) {
            tmp /= 10;
            len++;
        }
        
        _ret = new string(len);

        uint8 i = len-1;
        while (_self > 0) {
            bytes(_ret)[i--] = byte(_self%10+0x30);
            _self /= 10;
        }
    }

    function toHexString(uint _self) internal returns (string _ret) {
        if (_self == 0) {
            return "0x0";
        }

        uint8 len = 2;
        uint tmp = _self;
        while (tmp > 0) {
            tmp /= 16;
            len++;
        }
        
        _ret = new string(len);

        uint8 i = len-1;
        while (_self > 0) {
            if (_self%16 > 9) {
                    bytes(_ret)[i--] = byte(_self%16+0x61-0xa);
                } else {
                    bytes(_ret)[i--] = byte(_self%16+0x30);
                }
            
            _self /= 16;
        }
        
        bytes(_ret)[0] = byte(0x30);
        bytes(_ret)[1] = byte(0x78);
    }

    function toHexString64(uint _self) internal returns (string _ret) {      
        _ret = new string(66);
        bytes(_ret)[0] = '0';
        bytes(_ret)[1] = 'x';

        for (uint8 i=65; i>=2; --i) {
            uint8 digit = uint8(_self&0x0F);
            _self /= 16;

            if (digit < 10)
                bytes(_ret)[i] = byte(digit+0x30);
            else
                bytes(_ret)[i] = byte(digit-10+0x61);
        }
    }

    function toString(int _self) internal returns (string _ret) {
        if (_self == 0) {
            return "0";
        }

        uint ui = uint(_self);
        bool positive = true;
        uint8 len = 0;
        if (_self < 0) {
            ui = uint(-_self);
            positive = false; 
            len++;
        }
        
        uint tmp = ui;
        while (tmp > 0) {
            tmp /= 10;
            len++;
        }
        
        _ret = new string(len);
        if (!positive) {
            bytes(_ret)[0] = '-';
        }
        
        uint8 i = len-1;
        while (ui > 0) {
            bytes(_ret)[i--] = byte(ui%10+0x30);
            ui /= 10;
        }
    }

    function toAddrString(uint _self) internal returns (string _ret) {      
        _ret = new string(42);
        bytes(_ret)[0] = '0';
        bytes(_ret)[1] = 'x';

        for (uint8 i=41; i>=2; --i) {
            uint8 digit = uint8(_self&0x0F);
            _self /= 16;

            if (digit < 10)
                bytes(_ret)[i] = byte(digit+0x30);
            else
                bytes(_ret)[i] = byte(digit-10+0x61);
        }
    }
    
    function toKeyValue(uint _self, string _key) internal returns (string _ret) {
        uint len = bytes(_key).length+3;

        if (_self == 0) {
            len += 1;
        } else {
            uint tmp = _self;
            while (tmp > 0) {
                tmp /= 10;
                len++;
            }
        }
        
        _ret = new string(len);
        
        uint i = 0;
        bytes(_ret)[i++] = '"';
        for (uint j=0; j < bytes(_key).length; j++) {
            bytes(_ret)[i++] = bytes(_key)[j];
        }
        bytes(_ret)[i++] = '"';
        
        bytes(_ret)[i++] = ':';

        i = len-1;
        if (_self == 0) {
            bytes(_ret)[i] = byte(0x30);
        } else {
            while (_self > 0) {
                bytes(_ret)[i--] = byte(_self%10+0x30);
                _self /= 10;
            }
        }
    }
    
    function toKeyValue(int _self, string _key) internal returns (string _ret) {
        uint ui = uint(_self);
        bool positive = true;
        uint len = bytes(_key).length+3;
        if (_self < 0) {
            ui = uint(-_self);
            positive = false; 
            len++;
        }

        if (_self == 0) {
            len += 1;
        } else {
            uint tmp = ui;
            while (tmp > 0) {
                tmp /= 10;
                len++;
            }
        }
        
        _ret = new string(len);
        
        uint i = 0;
        bytes(_ret)[i++] = '"';
        for (uint j=0; j < bytes(_key).length; j++) {
            bytes(_ret)[i++] = bytes(_key)[j];
        }
        bytes(_ret)[i++] = '"';
        
        bytes(_ret)[i++] = ':';

        if (!positive) {
            bytes(_ret)[i++] = '-';
        }
        i = len-1;
        if (_self == 0) {
            bytes(_ret)[i] = byte(0x30);
        } else {
            while (ui > 0) {
                bytes(_ret)[i--] = byte(ui%10+0x30);
                ui /= 10;
            }
        }
    }

    function toKeyValue(address _self, string _key) internal returns (string _ret) {
        uint ui = uint(_self);
        uint len = bytes(_key).length+47;
        
        _ret = new string(len);
        
        uint i = 0;
        bytes(_ret)[i++] = '"';
        for (uint j=0; j < bytes(_key).length; j++) {
            bytes(_ret)[i++] = bytes(_key)[j];
        }
        bytes(_ret)[i++] = '"';
        
        bytes(_ret)[i++] = ':';

        bytes(_ret)[i++] = '"';
        bytes(_ret)[i++] = '0';
        bytes(_ret)[i++] = 'x';

        for (i=len-2; i>=len-41; i--) {
            uint8 digit = uint8(ui&0x0F);
            ui /= 16;

            if (digit < 10)
                bytes(_ret)[i] = byte(digit+0x30);
            else
                bytes(_ret)[i] = byte(digit-10+0x61);
        }

        bytes(_ret)[len-1] = '"';
    }

    function recoveryToString(uint _self) internal returns (string _ret) {
        uint tmp = _self;
        uint len = 0;
        while (tmp > 0) {
            tmp /= 256;
            len++;
        }
        
        _ret = new string(len);
        tmp = _self;
        uint i = len-1;
        while (tmp > 0) {
            bytes(_ret)[i] = byte(tmp%256);
            tmp /= 256;
            i--;
        }
    }

    function inArray(uint _self, uint[] storage _array) internal returns (bool _ret) {
        for (uint i=0; i<_array.length; ++i) {
            if (_self==_array[i]) {
                return true;
            }
        }

        return false;
    }

    //�������㣬��ȡ�������������ֵ
    function divide(uint _divisor, uint _divident) internal returns (uint _result) {
        if (_divident == 0) {
            return;
        }
        _result = _divisor/_divident;
        uint _tmp = _divisor - _result*_divident;
        if (_result == 0){
            _tmp = _divisor;
        }
        if (2*_tmp >= _divident) {
            return _result+1;
        }
    }
}
