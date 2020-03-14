pragma solidity ^0.5.0;

library SafeMath {
    function add(uint a, uint b) internal pure returns (uint c) {
        c = a + b;
        require(c >= a, "SafeMath add wrong value");
        return c;
    }
    function sub(uint a, uint b) internal pure returns (uint) {
        require(b <= a, "SafeMath sub wrong value");
        return a - b;
    }
}

contract Destructable {
    address payable private owner;

    constructor() public {
        owner = msg.sender;
    }

    modifier isDestructable {
        require(msg.sender == owner, "only owner");
        require(address(this).balance == 0, "balance is not zero");
        _;
    }

    function destruct() public isDestructable {
        selfdestruct(owner);
    }
}

contract AtomicSwap is Destructable {
    using SafeMath for uint;

    enum State { Empty, Initiated, Redeemed, Refunded }

    struct Swap {
        bytes32 hashedSecret;
        bytes32 secret;
        address payable initiator;
        address payable participant;
        uint refundTimestamp;
        uint value;
        uint payoff;
        State state;
    }

    event Initiated(
        bytes32 indexed _hashedSecret,
        address indexed _participant,
        address _initiator,
        uint _refundTimestamp,
        uint _value,
        uint _payoff
    );

    event Added(
        bytes32 indexed _hashedSecret,
        address _sender,
        uint _value
    );

    event Redeemed(
        bytes32 indexed _hashedSecret,
        bytes32 _secret
    );

    event Refunded(
        bytes32 indexed _hashedSecret
    );

    mapping(bytes32 => Swap) public swaps;

    modifier isRefundable(bytes32 _hashedSecret) {
        require(block.timestamp >= swaps[_hashedSecret].refundTimestamp, "refundTimestamp has not passed");
        _;
    }

    modifier isRedeemable(bytes32 _hashedSecret, bytes32 _secret) {
        require(block.timestamp < swaps[_hashedSecret].refundTimestamp, "refundTimestamp has already passed");
        require(sha256(abi.encodePacked(sha256(abi.encodePacked(_secret)))) == _hashedSecret, "secret is not correct");
        _;
    }

    modifier isInitiated(bytes32 _hashedSecret) {
        require(swaps[_hashedSecret].state == State.Initiated, "swap for this hash is empty or already spent");
        _;
    }

    modifier isInitiatable(bytes32 _hashedSecret, uint _refundTimestamp) {
        require(swaps[_hashedSecret].state == State.Empty, "swap for this hash is already initiated");
        require(_refundTimestamp > block.timestamp, "refundTimestamp has already passed");
        _;
    }

    modifier isAddable(bytes32 _hashedSecret) {
        require(block.timestamp <= swaps[_hashedSecret].refundTimestamp, "refundTime has already come");
        _;
    }

    function initiate (bytes32 _hashedSecret, address payable _participant, uint _refundTimestamp, uint _payoff)
    public payable isInitiatable(_hashedSecret, _refundTimestamp)
    {
        swaps[_hashedSecret].value = msg.value.sub(_payoff);
        swaps[_hashedSecret].hashedSecret = _hashedSecret;
        swaps[_hashedSecret].initiator = msg.sender;
        swaps[_hashedSecret].participant = _participant;
        swaps[_hashedSecret].refundTimestamp = _refundTimestamp;
        swaps[_hashedSecret].payoff = _payoff;
        swaps[_hashedSecret].state = State.Initiated;

        emit Initiated(
            _hashedSecret,
            swaps[_hashedSecret].participant,
            msg.sender,
            swaps[_hashedSecret].refundTimestamp,
            swaps[_hashedSecret].value,
            swaps[_hashedSecret].payoff
        );
    }

    function add (bytes32 _hashedSecret)
    public payable isInitiated(_hashedSecret) isAddable(_hashedSecret)
    {
        swaps[_hashedSecret].value = swaps[_hashedSecret].value.add(msg.value);

        emit Added(
            _hashedSecret,
            msg.sender,
            swaps[_hashedSecret].value
        );
    }

    function redeem(bytes32 _hashedSecret, bytes32 _secret)
    public isInitiated(_hashedSecret) isRedeemable(_hashedSecret, _secret)
    {
        swaps[_hashedSecret].secret = _secret;
        swaps[_hashedSecret].state = State.Redeemed;

        emit Redeemed(
            _hashedSecret,
            _secret
        );

        swaps[_hashedSecret].participant.transfer(swaps[_hashedSecret].value);
        if (swaps[_hashedSecret].payoff > 0) {
            msg.sender.transfer(swaps[_hashedSecret].payoff);
        }

        delete swaps[_hashedSecret];
    }

    function refund(bytes32 _hashedSecret)
    public isInitiated(_hashedSecret) isRefundable(_hashedSecret)
    {
        swaps[_hashedSecret].state = State.Refunded;

        emit Refunded(
            _hashedSecret
        );

        swaps[_hashedSecret].initiator.transfer(swaps[_hashedSecret].value.add(swaps[_hashedSecret].payoff));

        delete swaps[_hashedSecret];
    }
}