## 系统地址及交易格式定义
系统预设内置地址以0x100000000000000000000000000000000000000开始，根据需求在此地址上自加。目前系统用到的地址定义如下：
0x100000000000000000000000000000000000000	储备池账户地址
0x100000000000000000000000000000000000001	候选人合约地址
0x100000000000000000000000000000000000002	票池合约地址

前端调用内置合约接口，需要把调用参数组成rlp格式放到交易的data字段。格式如下：

data = rlp(type [8]byte, funcname string, parma1 []byte, parma2 []byte, ...)

| **参数名** | **类型** | **大小** | **参数说明** |
| --- | --- | --- | --- |
| type | bytes | 8byte | 交易类型(供前端扩展使用，底层不使用该字段) |
| funcname | string | 不定 | 内置合约定义的具体函数名称 |
| param1 | 不定 | 不定 | 函数参数 |
| ...... | ...... | ...... | ...... |


## 候选池合约接口定义

### **`CandidateDeposit`**
节点候选人申请/增加质押，质押金额为交易的value值
入参：
* `nodeId`: [64]byte 节点ID(公钥)
* `owner`: [20]byte 质押金退款地址
* `fee`: uint32 出块奖励佣金比，以10000为基数(eg：5%，则fee=500)
* `host`: string 节点IP
* `port`: string 节点端口号

出参：
* `Ret`: bool 操作结果
* `ErrMsg`: string 错误信息

### **`CandidateApplyWithdraw`**
节点质押金退回申请，申请成功后节点将被重新排序，权限校验from==owner。
入参：
* `nodeId`: [64]byte 节点ID(公钥)
* `withdraw`: uint256 退款金额 (单位：ADP)

出参：
* `Ret`: bool 操作结果
* `ErrMsg`: string 错误信息

### **`CandidateWithdrawInfos`**
获取节点申请的退款记录列表
入参：
* `nodeId`: [64]byte 节点ID(公钥)

出参：
* `Ret`: bool 操作结果
* `ErrMsg`: string 错误信息
* `[]`:列表
	* 'Balance': uint256 退款金额 (单位：ADP)
	* `LockNumber`: uint256 退款申请所在块高
	* `LockBlockCycle`: uint256 退款金额锁定周期

### **`CandidateWithdraw`**
节点质押金提取，调用成功后会提取所有已申请退回的质押金到owner账户。
入参：
* `nodeId`: [64]byte 节点ID(公钥)

出参：
* `Ret`: bool 操作结果
* `ErrMsg`: string 错误信息

### **`SetCandidateExtra`**
设置节点附加信息，供前端扩展使用。
入参：
* `nodeId`: [64]byte 节点ID(公钥)
* `extra`: string 附加信息

出参：
* `Ret`: bool 操作结果
* `ErrMsg`: string 错误信息

### **`CandidateDetails`**
获取候选人信息。
入参：
* `nodeId`: [64]byte 节点ID(公钥)

出参：
* `Ret`: bool 操作结果
* `ErrMsg`: string 错误信息
* `Deposit`: uint256 质押金额 (单位：ADP)
* `BlockNumber`: uint256 质押金更新的最新块高
* `Owner`: [20]byte 质押金退款地址
* `Fee`: uint32 出块奖励佣金比，以10000为基数(eg：5%，则fee=500)
* `Host`: string 节点IP
* `Port`: string 节点端口号

### **`CandidateList`**
获取所有入围节点的信息列表
入参：
* 无

出参：
* `Ret`: bool 操作结果
* `ErrMsg`: string 错误信息
* `[]`:列表
	* `Deposit`: uint256 质押金额 (单位：ADP)
	* `BlockNumber`: uint256 质押金更新的最新块高
	* `Owner`: [20]byte 质押金退款地址
	* `Fee`: uint32 出块奖励佣金比，以10000为基数(eg：5%，则fee=500)
	* `Host`: string 节点IP
	* `Port`: string 节点端口号

### **`VerifiersList`**
获取参与当前共识的验证人列表
入参：
* 无

出参：
* `Ret`: bool 操作结果
* `ErrMsg`: string 错误信息
* `[]`:列表
	* `Deposit`: uint256 质押金额 (单位：ADP)
	* `BlockNumber`: uint256 质押金更新的最新块高
	* `Owner`: [20]byte 质押金退款地址
	* `Fee`: uint32 出块奖励佣金比，以10000为基数(eg：5%，则fee=500)
	* `Host`: string 节点IP
	* `Port`: string 节点端口号


##### 注：所有返回值均为json格式

##例子
候选人质押申请，参数：
质押金：100*10^18ADP
* `nodeId`: "0xe152be5f5f0167250592a12a197ab19b215c5295d5eb0bb1133673dc8607530db1bfa5415b2ec5e94113f2fce0c4a60e697d5d703a29609b197b836b020446c7"
* `owner`: "0x4fed1fc4144c223ae3c1553be203cdfcbd38c581"
* `fee`: 500
* `host`: "127.0.0.1"
* `port`: "7890"

rlp[0x00000000000000f1, "CandidateDeposit", 0xe152be5f5f0167250592a12a197ab19b215c5295d5eb0bb1133673dc8607530db1bfa5415b2ec5e94113f2fce0c4a60e697d5d703a29609b197b836b020446c7, 0x4fed1fc4144c223ae3c1553be203cdfcbd38c581, 500, "127.0.0.1", "7890"]

序列化：
0xf8898800000000000000f19043616e6469646174654465706f736974b840e152be5f5f0167250592a12a197ab19b215c5295d5eb0bb1133673dc8607530db1bfa5415b2ec5e94113f2fce0c4a60e697d5d703a29609b197b836b020446c7944fed1fc4144c223ae3c1553be203cdfcbd38c5818800000000000001f4893132372e302e302e318437383930

交易：
{
from:eth.accounts[0],
to:"0x1000000000000000000000000000000000000001",
value:web3.toWei(100, "ether"),
data:"0xf8898800000000000000f19043616e6469646174654465706f736974b840e152be5f5f0167250592a12a197ab19b215c5295d5eb0bb1133673dc8607530db1bfa5415b2ec5e94113f2fce0c4a60e697d5d703a29609b197b836b020446c7944fed1fc4144c223ae3c1553be203cdfcbd38c5818800000000000001f4893132372e302e302e318437383930"
}





























