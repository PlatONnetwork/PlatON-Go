pragma solidity ^0.5.13;
/**
 * 竞猜合约
 * 每次下注合约接收用户转账，满足条件下(单笔>=5LAT，超过5LAT也只分配一个彩票)，
 * 每笔交易给对应地址分配一个彩票（一个地址多笔交易可分配多个彩票），并记录地址转入数；
 *
 */
contract Guessing {

    uint256 public endBlock; //竞猜截止块高
    bool public guessingClosed = false; //竞猜是否已开奖
    uint256 public baseUnit = 5 lat; //最小转金额（为了公平起见，保证只5个lat可以取到一个选票。）

    uint256 public balance; //竞猜总金额
    mapping(address => uint256) public gussingerLat; //每个竞猜者对应的金额
    mapping(uint256 =>address payable ) public IndexOfgussinger; //每个竞猜者对应的下标（5个lat就给他分配一个随机数）
    uint public indexKey = 0;
    address[] public winnerAddresses; //中奖者地址
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
    constructor (uint _endBlock)public {
        createAddress = msg.sender;
        endBlock = _endBlock;
    }

    /**
     * 默认函数
     *
     * 默认函数，可以向合约直接打款
     */
    function () payable external {
        require(msg.value >= baseUnit);

        uint amount = msg.value;
        IndexOfgussinger[indexKey] = msg.sender;
        indexKey++;

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
    modifier beforeDeadline() { if (endBlock >= block.number) _; }

    /**
     * 判断转账金额要大于5lat
     */
    modifier checkAmount() { if (msg.value >= baseUnit) _; }


    /**
     * 竞猜(带上金额)
     */
    function guessingWithLat() public beforeDeadline checkAmount payable{

        uint amount = msg.value;

        IndexOfgussinger[indexKey] = msg.sender;
        indexKey++;

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
     * 参数A：基于开奖截止区块的hash转成uint256
     *
     * 参数B：本期内参与的总票数
     *
     * 结果C ：A 对 B取余
     *
     * 方法：A % B = C
     *
     *
     * 1.参数B为个位数和十位数，取余数C的个位数，个位数相同的为中奖
     *   例如余数的12，取个位数，个位数是2的为中奖票
     *
     * 2.参数B为百位数和千位数，取余数C的两位数(个位十位)，两位数相同的为中奖；
     *   例如余数是123，取两位数，票数后两位23的为中奖票
     *
     * 3.参数B为万位数和十万位数，取余数C的三位数(个位十位百位)，三位数相同的为中奖；
     *   例如余数是1234，取三位数，票数后三位带234的为中奖票
     *
     */
    function draw() public afterDeadline {
        //只有合约创建者可以开奖
        if(!guessingClosed && createAddress == msg.sender){

            uint256 random = uint256(keccak256(abi.encodePacked(blockhash(endBlock))));
            uint drawIndex = random%indexKey;
            uint postfix;

            if(indexKey<100){
                postfix = drawIndex%10;
                if(postfix ==0){
                    for(uint256 i=0;i<indexKey;i++){
                        if((i-postfix)%10 == 0){
                            winnerAddresses.push(IndexOfgussinger[i]);
                        }
                    }
                }else{
                    for(uint256 i=0;i<indexKey;i++){
                        if(i%10 != 0 && (i-postfix)%10 == 0){
                            winnerAddresses.push(IndexOfgussinger[i]);
                        }
                    }
                }
            }else if(indexKey<10000){
                postfix = drawIndex%100;
                if(postfix ==0){
                    for(uint256 i=0;i<indexKey;i++){
                        if((i-postfix)%100 == 0){
                            winnerAddresses.push(IndexOfgussinger[i]);
                        }
                    }
                }else{
                    for(uint256 i=0;i<indexKey;i++){
                        if(i%100 != 0 && (i-postfix)%100 == 0){
                            winnerAddresses.push(IndexOfgussinger[i]);
                        }
                    }
                }
            }else{
                postfix = drawIndex%1000;
                if(postfix ==0){
                    for(uint256 i=0;i<indexKey;i++){
                        if((i-postfix)%1000 == 0){
                            winnerAddresses.push(IndexOfgussinger[i]);
                        }
                    }
                }else{
                    for(uint256 i=0;i<indexKey;i++){
                        if(i%1000 != 0 && (i-postfix)%1000 == 0){
                            winnerAddresses.push(IndexOfgussinger[i]);
                        }
                    }
                }
            }

            //向中奖者转账
            for(uint256 j=0;j<winnerAddresses.length;j++){
                if(winnerAddresses[j] != address(0x0)){
                    address(uint160(winnerAddresses[j])).transfer(balance);
                }
            }

            guessingClosed = true;
        }
    }

    /**
     * 查看当前合约中的余额
     */
    function getBalanceOf() view public returns (uint256){
        return address(this).balance;
    }

    /**
     * 查看一共有几个中奖者
     */
    function getWinnerCount() view public returns (uint256){
        return winnerAddresses.length;
    }

}