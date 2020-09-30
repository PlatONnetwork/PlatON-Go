pragma solidity 0.5.13;

/**
 * ERC20 0.5.13版本验证
 *
 *
 * @author hudenian
 * @dev 2019/12/23
 **/
contract ERC200513Token {
    string public name; // ERC20标准--代币名称
    string public symbol; // ERC20标准——代币简称
    uint8 public decimals = 18;  // ERC20标准，decimals 可以有的小数点个数，最小的代币单位。18 是建议的默认值
    uint256 public totalSupply; // ERC20标准 总供应量

    // 用mapping保存每个地址对应的余额 ERC20标准
    mapping(address => uint256) public balanceOf;
    // 存储对账号的控制 ERC20标准
    mapping(address => mapping(address => uint256)) public allowance;

    // 事件，用来通知客户端交易发生 ERC20标准
    event Transfer(address indexed from, address indexed to, uint256 value);

    // 事件，用来通知客户端代币被消费 ERC20标准
    event Burn(address indexed from, uint256 value);

    /**
     * 初始化构造
     */
    constructor(uint256 initialSupply, string memory tokenName, string memory tokenSymbol) public {
        totalSupply = initialSupply * 10 ** uint256(decimals);
        // 供应的份额，份额跟最小的代币单位有关，份额 = 币数 * 10 ** decimals。
        balanceOf[msg.sender] = totalSupply;
        // 创建者拥有所有的代币
        name = tokenName;
        // 代币名称
        symbol = tokenSymbol;
        // 代币符号
    }
    /**
     * 返回代币的名称
     */
    function getName() view public returns (string memory){
        return name;
    }

    /**
     * 返回代币的简称
     */
    function getSymbol() view public returns (string memory){
        return symbol;
    }
    /**
      * 返回代币最小分割量
      */
    function getDecimals() public view returns (uint8){
        return decimals;
    }

    function getTotalSupply() public view returns (uint256 theTotalSupply) {
        //函数声明中已经定义了返回变量theTotalSupply
        theTotalSupply = totalSupply;
        return theTotalSupply;
    }

    function getBalanceOf(address _owner) public view returns (uint256 balance) {
        //返回指定地址的通证余额
        return balanceOf[_owner];
    }
    /**
     * 代币交易转移的内部实现
     */
    function _transfer(address _from, address _to, uint _value) internal returns (bool success){
        // 确保目标地址不为0x0，因为0x0地址代表销毁
        require(_to != address("lat1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq542u6a") || _to != address("lax1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqmscn5j"));
        // 检查发送者余额
        require(balanceOf[_from] >= _value);
        // 确保转移为正数个
        require(balanceOf[_to] + _value > balanceOf[_to]);

        // 以下用来检查交易，
        uint previousBalances = balanceOf[_from] + balanceOf[_to];
        // Subtract from the sender
        balanceOf[_from] -= _value;
        // Add the same to the recipient
        balanceOf[_to] += _value;
        emit Transfer(_from, _to, _value);

        // 用assert来检查代码逻辑。
        return (balanceOf[_from] + balanceOf[_to] == previousBalances);
    }

    /**
     *  代币交易转移
     *  从自己（创建交易者）账号发送`_value`个代币到 `_to`账号
     * ERC20标准
     * @param _to 接收者地址
     * @param _value 转移数额
     */
    function transfer(address _to, uint256 _value) public returns (bool success){
        return _transfer(msg.sender, _to, _value);
    }

    /**
     * 账号之间代币交易转移
     * ERC20标准
     * @param _from 发送者地址
     * @param _to 接收者地址
     * @param _value 转移数额
     */
    function transferFrom(address _from, address _to, uint256 _value) public returns (bool success) {
        require(_value <= allowance[_from][msg.sender]);
        // Check allowance
        allowance[_from][msg.sender] -= _value;
        _transfer(_from, _to, _value);
        return true;
    }

    /**
     * 设置某个地址（合约）可以创建交易者名义花费的代币数。
     *
     * 允许发送者`_spender` 花费不多于 `_value` 个代币
     * ERC20标准
     * @param _spender The address authorized to spend
     * @param _value the max amount they can spend
     */
    function approve(address _spender, uint256 _value) public
    returns (bool success) {
        allowance[msg.sender][_spender] = _value;
        return true;
    }

    /**
     * 设置允许一个地址（合约）以我（创建交易者）的名义可最多花费的代币数。
     *-非ERC20标准
     * @param _spender 被授权的地址（合约）
     * @param _value 最大可花费代币数
     * @param _extraData 发送给合约的附加数据
     */
    // function approveAndCall(address _spender, uint256 _value, bytes memory _extraData) public returns (bool success) {
    //     tokenRecipient spender = tokenRecipient(_spender);
    //     if (approve(_spender, _value)) {
    //         // 通知合约
    //         spender.receiveApproval(msg.sender, _value, address(this), _extraData);
    //         return true;
    //     }
    // }

    /**
     *
     * 获取_spender可以从账户_owner中转出token的剩余数量
     */
    function getAllowance(address _owner, address _spender) public view returns (uint remaining){
        return allowance[_owner][_spender];
    }

    /**
     * 销毁我（创建交易者）账户中指定个代币
     *-非ERC20标准
     */
    function burn(uint256 _value) public returns (bool success) {
        require(balanceOf[msg.sender] >= _value);
        // Check if the sender has enough
        balanceOf[msg.sender] -= _value;
        // Subtract from the sender
        totalSupply -= _value;
        // Updates totalSupply
        emit Burn(msg.sender, _value);
        return true;
    }

    /**
     * 销毁用户账户中指定个代币
     *-非ERC20标准
     * Remove `_value` tokens from the system irreversibly on behalf of `_from`.
     *
     * @param _from the address of the sender
     * @param _value the amount of money to burn
     */
    function burnFrom(address _from, uint256 _value) public returns (bool success) {
        require(balanceOf[_from] >= _value);
        // Check if the targeted balance is enough
        require(_value <= allowance[_from][msg.sender]);
        // Check allowance
        balanceOf[_from] -= _value;
        // Subtract from the targeted balance
        allowance[_from][msg.sender] -= _value;
        // Subtract from the sender's allowance
        totalSupply -= _value;
        // Update totalSupply
        emit Burn(_from, _value);
        return true;
    }
}