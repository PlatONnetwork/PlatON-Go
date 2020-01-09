pragma solidity ^0.4.12;

library Errnos {

    enum FilterType {
        FILTER_TYPE_START,
        FILTER_TYPE_ROLE,
        //other filter start
        //other filter end
        FILTER_TYPE_END//type in (FILTER_TYPE_START,FILTER_TYPE_END) is valid.
    }
    
    enum State {
        STATE_INVALID,
        STATE_VALID,
        STATE_DELETED
    }

    enum Base_Errnos {
        NO_ERROR,//don't add BASE_ERRNO_OFFSET
        BAD_PARAMETER,
        ID_EMPTY,
        ID_OVERFLOW,
        ID_CONFLICTED,
        ID_NONEXIST,
        NAME_EMPTY,
        VERSION_EMPTY,
        STATE_INVALID,
        ADDRESS_INVALID,
        NAME_VERSION_CONFLICTED
    }
    
    enum Filter_Errnos {
        FILTER_TYPE_INVALID
    }
}