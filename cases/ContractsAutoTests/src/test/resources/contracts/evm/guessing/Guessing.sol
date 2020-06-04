pragma solidity ^0.5.13;
/**
 * 竞猜合约
 * 当用户往竞猜合约发起转账操作，满足(单笔>=5LAT)，
 * 每笔交易给对应地址分配一个抽奖号码（一个地址多笔交易可分配多个抽奖号码），并记录地址转入数；
 *
 */
contract Guessing {

    uint256 public endBlock; //竞猜截止块高
    bytes32 public block_hash;//竞猜截止块高的块hash
    bool public guessingClosed = false; //竞猜是否已开奖
    uint256 public baseUnit = 5 lat; //最小转金额

    uint256 public balance; //竞猜奖池总金额
    uint256 public averageAmount; //每个中奖号码对应获奖金额
    mapping(address => uint256) public gussingerLat; //每个竞猜者发起转账的金额
    mapping(uint256 =>address payable ) public indexOfgussinger; //每个抽奖号码对应的竞猜者地址
    mapping(address => uint[]) public gussingerCodes; //每个竞猜者所有的抽奖号码集合
    mapping(address => uint256 ) public winnerMap; //中奖者对应中奖号码个数
    uint public indexKey = 1;
    address[] public winnerAddresses; //中奖者地址数组
    address public createAddress;//合约创建者
    uint public postfix = 0; //中奖尾号


    //竞猜成功通知
    event FundTransfer(address _backer, uint _amount, bool _isSuccess);

    //记录已接收的LAT通知
    event CurrentBalance(address _msgSenderAddress, uint _balance);

    /**
     * 初始化构造函数
     *
     * @param _endBlock 竞猜截止块高（达到此块高后不能再往合约转账）
     */
    constructor (uint _endBlock)public {
        createAddress = msg.sender;
        endBlock = _endBlock;
    }

    /**
     *
     * 默认函数，可以向合约直接打款
     */
    function () payable external {
        guessing(msg.value,msg.sender);
    }

    /**
     * 判断是否已经过了竞猜截止限期
     */
    modifier beforeDeadline() { if (endBlock >= block.number) _; }


    /**
     * 竞猜
     */
    function guessingWithLat() public payable{
        guessing(msg.value,msg.sender);
    }

    //竞猜内部实现
    function guessing(uint amount,address payable msgsender) beforeDeadline internal{
        //转账金额要大于5lat
        require(amount >= baseUnit);

        gussingerCodes[msgsender].push(indexKey);

        indexOfgussinger[indexKey] = msgsender;
        indexKey++;

        //竞猜人的金额累加
        gussingerLat[msgsender] += amount;

        //竞猜总额累加
        balance += amount;

        //竞猜成功通知
        emit FundTransfer(msgsender, amount, true);
    }

    /**
     * 判断是否已经过了竞猜截止限期
     */
    modifier afterDeadline() { if (block.number > endBlock ) _; }

    /**
     * 开奖操作
     * 如果当前区块超过了截止日期
     *
     * 参数A：基于开奖截止区块的hash转成uint256
     *
     * 参数B：本期内参与的总票数
     *
     * 结果C ：A 对 B取余
     *
     * 方法：A % B = C
     *
     *
     * 1.参数B,1<总票数<99，取余数C的个位数，个位数相同的为中奖；例如余数的12，取个位数，个位数是2的为中奖票
     * 2.参数B,100<总票数<9999，取余数C的两位数(个位十位)，两位数相同的为中奖；例如余数是123，取两位数，票数后两位23的为中奖票
     * 3.参数B，总票数>10000，取余数C的三位数(个位十位百位)，三位数相同的为中奖；例如余数是1234，取三位数，票数后三位带234的为中奖票
     *
     */
    function draw(bytes32 _block_hash) public afterDeadline {
        //开奖只能操作一次，且只有合约创建者可以开奖
        if(!guessingClosed && createAddress == msg.sender && indexKey > 1){
            uint256 random = uint256(keccak256(abi.encodePacked(_block_hash)));
            uint drawIndex = random%indexKey;

            if(indexKey<100){
                getwinners(drawIndex,10);
            }else if(indexKey<10000){
                getwinners(drawIndex,100);
            }else{
                getwinners(drawIndex,1000);
            }

            //每个中奖者可以分到的金额
            averageAmount = balance/winnerAddresses.length;
            address payable tempAddress;
            //向中奖者转账
            for(uint256 j=0;j<winnerAddresses.length;j++){
                //中奖者中奖票号统计
                winnerMap[winnerAddresses[j]] = winnerMap[winnerAddresses[j]]+1;
                if(winnerAddresses[j] != address("lat1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq542u6a") || winnerAddresses[j] != address("lax1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqmscn5j")){
                    tempAddress = address(uint160(winnerAddresses[j]));
                    tempAddress.transfer(averageAmount);
                }
            }

            guessingClosed = true;
            block_hash = _block_hash;
        }
    }

    //获取所有的中奖者
    function getwinners(uint drawIndex,uint times) internal{
        postfix = drawIndex%times;
        if(postfix ==0){
            for(uint256 i=1;i<indexKey;i++){
                if((i-postfix)%times == 0){
                    winnerAddresses.push(indexOfgussinger[i]);
                }
            }
        }else{
            for(uint256 i=1;i<indexKey;i++){
                if(i%times != 0 && (i-postfix)%times == 0){
                    winnerAddresses.push(indexOfgussinger[i]);
                }
            }
        }
    }


    /**
     * 查看当前合约中的余额
     */
    function getBalanceOf() view public returns (uint256){
        return address(this).balance;
    }

    /**
     * 查看一共有几个中奖号码
     */
    function getWinnerCount() view public returns (uint256){
        return winnerAddresses.length;
    }

    /**
     * 获取所有中奖人地址（可能有重复，调用方可以对此进行合并）
     */
    function getWinnerAddresses() view public returns (address[] memory){
        return winnerAddresses;
    }

    /**
     * 查询截止区块对应的hash
     */
    function getEndBlockHash() view public returns (bytes32){
        return block_hash;
    }

    /**
     * 生成指定区块的blockhash
     */
    function generateBlockHash(uint256 _blocknumber) view public returns (bytes32){
        require(_blocknumber > block.number - 257 && _blocknumber < block.number - 1);
        return blockhash(_blocknumber);
    }


    /**
     * 获取当前参与者所有幸运号码
     */
    function getMyGuessCodes() view public returns (uint[] memory){
        return gussingerCodes[msg.sender];
    }

    /**
     * 获幸运尾号
     */
    function getPostfix() view public returns (uint){
        require(guessingClosed == true);
        return postfix;
    }


}