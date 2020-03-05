pragma solidity ^0.5.13;
/**
 * 竞猜合约
 * 1.创建一个竞猜合约，设置竞猜截止块高
 * 2.每个竞猜者都可以发起竞猜，竞猜金额>=5 LAT，每5个LAT可以获取一个中签号码（金额越多能拿到的中签号也越多）
 * 3.记录每个中签号码对应的钱包地址（中签号码生成规则为每个调用者竞猜的顺序号）
 * 4.当截止时间一到，竞猜开关关闭，由合约创建者进行开奖，并将合约中的所有奖金转至其账户下（开奖规则，将所所有开奖区块交易hash的进行keccak256后转成uint256再对所有的中签号码进行取余操作）
 * 注：在我们测试网络把ether修改成lat
 */
contract Guessing {

    uint256 public endBlock; //竞猜截止块高
    bool public guessingClosed = false; //竞猜是否已开奖
    uint256 public baseUnit = 5 lat; //最小转金额（为了公平起见，保证只5个lat可以取到一个选票。）【在我们测试网络把ether修改成lat】
    

    uint256 public balance; //竞猜总金额
    mapping(address => uint256) public gussingerLat; //每个竞猜者对应的金额
    mapping(uint256 =>address payable ) public IndexOfgussinger; //每个竞猜者对应的下标（5个lat就给他分配一个随机数）
    // mapping(string =>uint256 ) public indexMap; //自增序列下标
    uint public indexKey = 0;
    address payable public winnerAddress; //中奖者地址
    address public createAddress;//合约创建者

    //竞猜成功通知
    event FundTransfer(address _backer, uint _amount, bool _isSuccess);

     //记录已接收的LAT通知
    event CurrentBalance(address _msgSenderAddress, uint _balance);

    //校验地址是否为空
    modifier validAddress(address _address) {
        require(_address != address(0x0));
        _;
    }

    /**
     * 初始化构造函数
     *
     * @param _endBlock 竞猜截止块高（达到此块高后不能再往合约转账）
     */
    constructor (uint _endBlock)public payable {
        createAddress = msg.sender;
        endBlock = _endBlock;
        // indexMap[indexKey] = 0; //自增索引下标
    }

    /**
     * 默认函数
     *
     * 默认函数，可以向合约直接打款
     */
    function () payable external {

        require(msg.value > 0);

        //竞猜总额累加
        balance = add(balance,msg.value);

        uint amount = msg.value;
        uint num = amount/baseUnit; //判断可以获取几抽奖码
        for(uint i=0;i<num;i++){
            IndexOfgussinger[indexKey] = msg.sender;
            indexKey++;
        }

        //竞猜人的金额累加(可以投多次)
        gussingerLat[msg.sender] += amount;


        //竞猜成功通知
        emit FundTransfer(msg.sender, amount, true);
    }


    /**
     * 判断是否已经过了竞猜截止限期
     */
    modifier beforeDeadline() { if (endBlock >= block.number) _; }

    /**
     * 判断转账金额要大于5lat
     */
    modifier checkAmount() { if (msg.value >= baseUnit) _; }


    /**
     * 竞猜(带上金额)
     */
    function guessingWithLat() public beforeDeadline checkAmount payable{
        require(msg.value > 0);
        uint amount = msg.value;
        uint num = amount/baseUnit; //判断可以获取几抽奖码
        for(uint i=0;i<num;i++){
            IndexOfgussinger[indexKey] = msg.sender;
            indexKey++;
        }

        //竞猜人的金额累加(可以投多次)
        gussingerLat[msg.sender] += amount;

        //竞猜总额累加
        balance += amount;

        //竞猜成功通知
        emit FundTransfer(msg.sender, amount, true);
    }

    /**
     * 判断是否已经过了竞猜截止限期
     */
    modifier afterDeadline() { if (block.number > endBlock ) _; }

    /**
     * 开奖操作
     * 如果当前区块超过了截止日期
     * 
     */
    function draw() public afterDeadline {
        //只有合约创建者可以开奖
        if(!guessingClosed && createAddress == msg.sender){
            uint256 random = uint256(keccak256(abi.encodePacked(blockhash(endBlock))));
            uint drawIndex = random%indexKey;

            //取到中奖者
            winnerAddress = IndexOfgussinger[drawIndex];

            //向中奖者转账
            winnerAddress.transfer(balance);
            emit FundTransfer(winnerAddress, balance, false);

            guessingClosed = true;
        }
    }

    function getBalanceOf() view public returns (uint256){
        return address(this).balance;
    }

    //累加函数
    function add(uint256 a, uint256 b) internal pure returns (uint256) {
        uint256 c = a + b;
        assert(c >= a);
        return c;
    }

}