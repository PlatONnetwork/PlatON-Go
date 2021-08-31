### 1. 整体消息结构

```json

{
	"blockType": 0,
	"epoch": 1,
	"consensusNo": 1,
	"nodeID": "0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f",
	"nodeAddress": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",
	"block": {
		"extraData": "0xda8301000086706c61746f6e88676f312e31362e32856c696e757800000000000119fe39ad415cf8737558d010ac90cab76ce758d429eea6f9f0eb959e9b558d4f9b9bfef8d599e5e6771095833bbfda244e75ecd7aef409f87113874c6ca16800",
		"gasLimit": "0x8fcf88",
		"gasUsed": "0x2480c",
		"hash": "0x93b7278cd08e4144a7fdd103e2e359af1a5ebe76e7b5a4740aa9dbfc5cbfbbd5",
		"logsBloom": "0x00000000000000000000000000000000000000000000002000000000000000000000000000000000100000000000004000000000000040000000000000000000000000000000000000000008000000000000200000000000000000000000000800000080020000000000000000000800000000000000000000000010000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000100000000000000000",
		"miner": "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqrdyjj2v",
		"nonce": "0x02fb713531ba733281325f7e4269ed6482d6e57c486a083499f4ab383346f8da1af319aaf8884f9f02adfadac052a7d7f84939e515cfd8d8e8bb4135cc9f902a3beebe13373374df1f5772ddc5082d6686",
		"number": "0x143ad",
		"parentHash": "0x4c5dad3243de72c85de00c5e358503296d808251cd12bb49c825df4766e00cd6",
		"receiptsRoot": "0x5ec1b4df629e040e7661aa8b25b9fcf5608ee6dd108ca1da8777ca436657b6c7",
		"size": "0x42f",
		"stateRoot": "0x8ad02aee6d4e8672d0dec13ea4f350b41ddb9a012c4ed78d696e32e069807fbb",
		"timestamp": "0x1790c8dc54a",
		"transactions": [
			{
				"blockHash": "0x93b7278cd08e4144a7fdd103e2e359af1a5ebe76e7b5a4740aa9dbfc5cbfbbd5",
				"blockNumber": "0x143ad",
				"from": "lat12ve9py9hg5m2nfxelj06nq6frdgzft62zf64xt",
				"gas": "0x89543f",
				"gasPrice": "0x3b9aca00",
				"hash": "0x9c8411715dec897a6ae653f859f8aea6940d0a144de9206659d0e97fe5c0e862",
				"input": "0xd3fc98640000000000000000000000007c7b4da6a8d60632d072a64181766b6b752c358100000000000000000000000000000000000000000000000000000000013461510000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000005c68747470733a2f2f697066732e696f2f697066732f516d653963434a6b7461547055644a7332783176414b68744c3537534b48583541586536436a31753335705a42533f66696c656e616d653d5469636b6574303030312e6a736f6e00000000",
				"nonce": "0x1",
				"to": "lat12v6d2mguvnh4wm2d65k9sf5t2t8z9urer55u08",
				"transactionIndex": "0x0",
				"value": "0x0",
				"v": "0xeb",
				"r": "0xc0385de195c584e6ad2f8eb80290dc74041b17b6f2ed4765f18d333ac21abb21",
				"s": "0x27b915d4a88f7949e712509f3c4ffbf8759c639d94a6477a4c59a0706e0477ef"
			}
		],
		"transactionsRoot": "0xa11cf45b6f62436eaa3899ebb3727c6165fa3d2202367c41777e73bfe6bd47b2"
	},
	"receipts": [
		{
			"root": "0x",
			"status": "0x1",
			"cumulativeGasUsed": "0x2480c",
			"logsBloom": "0x00000000000000000000000000000000000000000000002000000000000000000000000000000000100000000000004000000000000040000000000000000000000000000000000000000008000000000000200000000000000000000000000800000080020000000000000000000800000000000000000000000010000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000100000000000000000",
			"logs": [
				{
					"address": "lat12v6d2mguvnh4wm2d65k9sf5t2t8z9urer55u08",
					"topics": [
						"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
						"0x0000000000000000000000000000000000000000000000000000000000000000",
						"0x0000000000000000000000007c7b4da6a8d60632d072a64181766b6b752c3581",
						"0x0000000000000000000000000000000000000000000000000000000001346151"
					],
					"data": "0x",
					"blockNumber": "0x143ad",
					"transactionHash": "0x9c8411715dec897a6ae653f859f8aea6940d0a144de9206659d0e97fe5c0e862",
					"transactionIndex": "0x0",
					"blockHash": "0x93b7278cd08e4144a7fdd103e2e359af1a5ebe76e7b5a4740aa9dbfc5cbfbbd5",
					"logIndex": "0x0",
					"removed": false
				}
			],
			"transactionHash": "0x9c8411715dec897a6ae653f859f8aea6940d0a144de9206659d0e97fe5c0e862",
			"contractAddress": "lat1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq542u6a",
			"gasUsed": "0x2480c",
			"blockHash": "0x93b7278cd08e4144a7fdd103e2e359af1a5ebe76e7b5a4740aa9dbfc5cbfbbd5",
			"blockNumber": "0x143ad",
			"transactionIndex": "0x0"
		}
	],
	"exeBlockData": {
		"activeVersion": "0.14.0",
		"additionalIssuanceData":{
			"additionalNo": 0,
			"additionalBase": 1000000,
			"additionalRate": 100,
			"additionalAmount": 1000,
			"issuanceItemList": [
				{
					"address": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",
					"amount": 1000000
				}
			]
		},
		"rewardData": {
			"blockRewardAmount": 1000,
			"delegatorReward": true,
			"stakingRewardAmount": 1000,
			"candidateInfoList": [
				{
					"nodeId":"0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f",
					"minerAddress": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9"
				}
			]
		},
		"zeroSlashingItemList": [
			{
				"nodeID": "0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f",
				"slashingAmount": 100
			}
		],
		"duplicatedSignSlashingSetting": {
			"penaltyRatioByValidStakings": 10,
			"rewardRatioByPenalties": 10
		},
		"stakingSetting": {
			"operatingThreshold": 1000
		},
		"stakingFrozenItemList": [
			{
				"nodeId": "0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f",
				"nodeAddress": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",
				"frozenEpochNo": 100,
				"recovery": true
			}
		],
		"restrictingReleaseItemList": [
			{
				"destAddress": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",
				"releaseAmount": 10000,
				"lackingAmount": 10000
			}
		],
		"embedTransferTxList": [
			{
				"txHash": "0x9c8411715dec897a6ae653f859f8aea6940d0a144de9206659d0e97fe5c0e862",
				"from": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",
				"to": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",
				"amount": 1000
			}
		],
		"embedContractTxList": [
			{
				"txHash": "0x9c8411715dec897a6ae653f859f8aea6940d0a144de9206659d0e97fe5c0e862",
				"from": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",
				"contractAddress": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",
				"input": "0xd3fc98640000000000000000000000007c7b4da6a8d60632d072a64181766b6b752c358100000000000000000000000000000000000000000000000000000000013461510000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000005c68747470733a2f2f697066732e696f2f697066732f516d653963434a6b7461547055644a7332783176414b68744c3537534b48583541586536436a31753335705a42533f66696c656e616d653d5469636b6574303030312e6a736f6e00000000"
			}
		],
		"withdrawDelegationList": [
			{
				"txHash": "0x9c8411715dec897a6ae653f859f8aea6940d0a144de9206659d0e97fe5c0e862",
				"delegateAddress": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",
				"nodeId": "0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f",
				"rewardAmount": 10000
			}
		],
		"autoStakingTxMap": {
			"0x9c8411715dec897a6ae653f859f8aea6940d0a144de9206659d0e97fe5c0e862": {
				"restrictingAmount": 10000,
				"balanceAmount": 10000
			}
		},
		"epochElection": [
			"0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f"
		],
		"consensusElection": [
			"0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f"
		],
		"epochNumber": 149
	},
	"GenesisData": {
		"allocItemList": [
			{
				"address": "lat1aaczrlrzylnanv57map5lndllkf7mvtnd9h8dj",
				"amount": 0
			}
		],
		"stakingItemList": [
			{
				"nodeID": "0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f",
				"stakingAddress": "lat1fpccktpn37a94rdj9yxszxp7pt3kae05j6lr9l",
				"benefitAddress": "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqrdyjj2v",
				"nodeName": "platon.node.1",
				"amount": 1000
			}
		],
		"restrictingCreateItemList": [
			{
				"from": "lat1fpccktpn37a94rdj9yxszxp7pt3kae05j6lr9l",
				"destAddress": "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqrdyjj2v",
				"plans": [
					100
				]
			}
		],
		"initFundItemList": [
			{
				"from": "lat1fpccktpn37a94rdj9yxszxp7pt3kae05j6lr9l",
				"to": "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqrdyjj2v",
				"amount": 1000
			}
		],
		"epochElection": [
			"0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f"
		],
		"consensusElection": [
			"0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f"
		]
	},
	"statData": {
		"put":{
			"candidate":[
				{
					"nodeId": "0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f",
					"stakingAddress": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",
					"benefitAddress": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",
					"rewardPer": 100,
					"nextRewardPer": 100,
					"stakingTxIndex": 11,
					"programVersion": 1000,
					"status": 2,
					"stakingBlockNum": 1000,
					"shares": 10000,
					"released": 10000,
					"releasedHes": 10000,
					"restrictingPlan": 10000,
					"restrictingPlanHes": 10000,
					"externalId": "",
					"nodeName": "",
					"website": "",
					"details": "",
					"delegateTotal": 1000,
					"delegateTotalHes": 10000,
					"delegateRewardTotal": 1000
				}
			]
		},
		"delete": {
			"candidate":[
				{
					"nodeId": "0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f",
					"stakingBlockNum": 1000,
					"stakingTxIndex": 11
				}
			]
			
		}
	}
}

```

##### 1.1 基础消息

1. blockType 区块类型
	```
	  0 - 创世块
	  1 - 普通区块
	  2 - 共识论起始块
	  3 - 共识论选举块
	  5 - 结算周期起始块
	  6 - 结算周期结束块
	```
2. epoch 结算周期轮数，从1开始
3. consensusNo 共识周期轮数，从1开始
4. nodeID 打包的节点id，创世区块为0
5. nodeAddress 打包的节点地址，创世区块为0
6. block 区块信息，包含交易。 同rpc接口platon.getBlock(9250679, true)
7. receipts 交易回执。 同rpc接口platon.getTransactionReceipt('0x9d5fbd3e5e67368dde656b8e87145733eeea3a1a6e5858ae0cdb894e82d1de93')
8. exeBlockData 区块在执行期间，跟踪系统关注的过程数据，包括在执行显式交易时的过程数据，以及在ChainBlockReactor调用BeginBlocker()和EndBlocker()时产生的过程数据。
9. GenesisData 存放创世块中，预分配的账户地址和余额，以及内置节点质押信息，创世块的锁仓，和其他一些账户的初始化信息。
10. statData 区块执行完成后，涉及状态数据的变更，包括候选人、委托等信息。

##### 1.2 exeBlockData 定义

1. activeVersion 如果当前块有升级提案生效，则填写新版本,0.14.0
2. additionalIssuanceData 增发数据, blockType = 6 时可能存在
    ```
    {
        "additionalNo": 1,              // 增加周期轮数， 从1开始，创世块增发不包含在其中       
        "additionalBase": 1000000,      // 增发基数
        "additionalRate": 100,          // 增发比例 单位：万分之一
        "additionalAmount": 1000,       // 增发金额
        "issuanceItemList": [           // 增发分配
            {
                "address": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",  // 增发金额分配地址
                "amount": 1000000                                         // 增发金额
            }
        ]
    }
    ```
3. rewardData 分配奖励，包括出块奖励，质押奖励
    ```
    {
        "blockRewardAmount": 1000,        // 出块奖励，每个块都有
        "delegatorReward": true,          // 出块奖励中，是否分配给委托人的奖励，每个块都有
        "stakingRewardAmount": 1000,      // 一结算周期内所有101节点的质押奖励，blockType = 6 时存在
        "candidateInfoList": [            // 备选节点信息，blockType = 6 时存在
            {
                "nodeId":"0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f",   // 备选节点ID
                "minerAddress": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9"  // 备选节点的矿工地址（收益地址）
            }
        ]
    }
    ```
4. zeroSlashingItemList 零出块惩罚节点明细 blockType = 3 时可能存在

    ```
    [
        {
            "nodeID": "0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f",   // 备选节点ID
            "slashingAmount": 100  // 0出块处罚金(先从已生效的自有质押扣除，不够在从已生效的锁仓质押扣除)
        }
    ]
    ```
5. duplicatedSignSlashingSetting 双签惩罚明细设置，双签交易时存在
    ```
    {
        "penaltyRatioByValidStakings": 10,  // 罚金 = 有效质押 * PenaltyRatioByValidStakings / 10000
        "rewardRatioByPenalties": 10        // 给举报人的赏金=罚金 * RewardRatioByPenalties / 100
    }
    ```	
6. stakingSetting 质押合约设置，解除委托时存在
    ```
    {
        "operatingThreshold": 10            //委托要求的最小数量；当减持委托时，委托数量少于该值，则全部减持。
    }
    ```	
7. stakingFrozenItemList 质押冻结信息。 在解除质押、 零出块惩罚、 双签惩罚、 锁定恢复（blockType = 6）时可能存在。
    ```
    [
        {
            "nodeId": "0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f",
            "nodeAddress": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",
            "frozenEpochNo": 100,
            "recovery": true
        }
    ]
    ```
8. restrictingReleaseItemList 锁仓释放明细, 创世块锁仓释放（blockType = 6）、 普通锁仓释放（blockType = 6）、
    ```
    [
        {
            "destAddress": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",  //释放地址
            "releaseAmount": 10000,                                       //释放金额
            "lackingAmount": 10000                                        //欠释放金额，该释放但是因为质押无法按时释放的资金
        }
    ]
    ```
9. embedTransferTxList 一个显式交易引起的内置转账交易, 合约内部转账、 存在金额合约自毁时存在，需要结合交易回执中状态处理，如果交易失败则忽略。
    ```
    [
        {
            "txHash": "0x9c8411715dec897a6ae653f859f8aea6940d0a144de9206659d0e97fe5c0e862",
            "from": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",
            "to": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",
            "amount": 1000
        }
    ]
    ```
10. embedContractTxList 一个显式交易引起的内置合约交易，需要结合交易回执中状态处理，如果交易失败则忽略。
    ```
    [
        {
            "txHash": "0x9c8411715dec897a6ae653f859f8aea6940d0a144de9206659d0e97fe5c0e862",
            "from": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",
            "contractAddress": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",
            "input": "0xd3fc98640000000000000000000000007c7b4da6a8d60632d072a64181766b6b752c358100000000000000000000000000000000000000000000000000000000013461510000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000005c68747470733a2f2f697066732e696f2f697066732f516d653963434a6b7461547055644a7332783176414b68744c3537534b48583541586536436a31753335705a42533f66696c656e616d653d5469636b6574303030312e6a736f6e00000000"
        }
    ]
    ```
11. withdrawDelegationList 解除委托时，存在奖励时
    ```
    [
        {
            "txHash": "0x9c8411715dec897a6ae653f859f8aea6940d0a144de9206659d0e97fe5c0e862",    // 委托用户撤销节点的全部委托的交易HASH
            "delegateAddress": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",            // 委托用户地址
            "nodeId": "0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f",  // 委托用户委托的节点ID
            "rewardAmount": 10000    // 委托用户从此节点获取的全部委托奖励
        }
    ]
    ```
12. autoStakingTxMap 混合质押时存在
    ```
    {
        "0x9c8411715dec897a6ae653f859f8aea6940d0a144de9206659d0e97fe5c0e862": {   // 交易hash
            "restrictingAmount": 10000,    // 质押中锁仓金额
            "balanceAmount": 10000         // 质押中自有金额
        }
    }
    ```
13. epochElection 结算周期验证人， blockType = 6 时存在
    ```
    [
        "0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f"  
    ]
    ```
14. consensusElection 共识周期验证人， blockType = 3 时存在
    ```
    [
        "0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f"  
    ]
     ```
15. epochNumber 当前增发周期内可能共识周期数， blockType = 6 时存在

##### 1.3 GenesisData 定义 blockType = 0 时存在

1. allocItemList 初始分配的金额，创世文件中 alloc 和 innerAcc中地址及金额。 
	```
	[
		{
			"address": "lat1aaczrlrzylnanv57map5lndllkf7mvtnd9h8dj",
			"amount": 0
		}
	]
	```
2. stakingItemList 初始质押的信息。  
	```
	[
		{
			"nodeID": "0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f",   // 初始验证人节点id  
			"stakingAddress": "lat1fpccktpn37a94rdj9yxszxp7pt3kae05j6lr9l",   // 初始验证人质押地址，即cdfAccount地址（计算士基金）
			"benefitAddress": "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqrdyjj2v",   // 初始验证人收益地址，即激励池地址
			"nodeName": "platon.node.1",
			"amount": 1.5e23
		}
	]
	``` 
3. restrictingCreateItemList 锁仓信息。 第一年不在列表中，直接转给激励池了
	```
	[
		{
			"from": "lat1fpccktpn37a94rdj9yxszxp7pt3kae05j6lr9l",            // 锁仓发起地址，即cdfAccount地址（计算士基金）
			"destAddress": "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqrdyjj2v",     // 锁仓目标地址，即激励池地址
			"plans": [
				5.5965742e25
			]
		}
	]
	```
4. initFundItemList 计算士基金第一年补贴给激励池的余额。
	```
	[
		{
			"from": "lat1fpccktpn37a94rdj9yxszxp7pt3kae05j6lr9l",
			"to": "lat1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqrdyjj2v",
			"amount": 6.2215742e25
		}
	]
	```
5. epochElection 初始的结算周期验证人
	```
	[
		"0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f"  
	]
	```
6. consensusElection 初始的共识周期验证人
	```
	[
		"0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f"  
	]
	```

##### 1.4 statData 定义

1. put 代表对象新增或更新
2. delete 代表对象被删除
3. candidate 定义
```
{
	"nodeId": "0xd2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f",   // 节点id
	"stakingAddress": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",  // 质押地址
	"benefitAddress": "lat153qkj3uk04yyagkgmplx7rhsagpy5n3k9gkwn9",  // 收益地址
	"rewardPer": 100,                     // 当前结算周期奖励分成比例，采用BasePoint 1BP=0.01%
	"nextRewardPer": 100,                 // 下一个结算周期奖励分成比例，采用BasePoint 1BP=0.01%   
	"stakingTxIndex": 11,                 // 发起质押时的交易索引
	"programVersion": 1000,               // 被质押节点的PlatON进程的真实版本号(获取版本号的接口由治理提供)
	"status": 2,                          // 候选人的状态(状态是根据uint32的32bit来放置的，可同时存在多个状态，值为多个同时存在的状态值相加【0: 节点可用 (32个bit全为0)； 1: 节点不可用 (只有最后一bit为1)； 2： 节点零出块需要锁定但无需解除质押(只有倒数第二bit为1)； 4： 节点的von不足最低质押门槛(只有倒数第三bit为1)； 8：节点被举报双签(只有倒数第四bit为1)); 16: 节点零出块需要锁定并解除质押(倒数第五位bit为1); 32: 节点主动发起撤销(只有倒数第六位bit为1)】
	"stakingBlockNum": 1000,              // 发起质押时的区块高度
	"shares": 10000,                      // 当前候选人总共质押加被委托的von数目,
	"released": 10000,                    // 发起质押账户的自由金额的锁定期质押的von
	"releasedHes": 10000,                 // 发起质押账户的自由金额的犹豫期质押的von
	"restrictingPlan": 10000,             // 发起质押账户的锁仓金额的锁定期质押的von
	"restrictingPlanHes": 10000,          // 发起质押账户的锁仓金额的犹豫期质押的von
	"externalId": "",                     // 外部Id(有长度限制，给第三方拉取节点描述的Id)
	"nodeName": "",                       // 被质押节点的名称(有长度限制，表示该节点的名称)
	"website": "",                        // 节点的第三方主页(有长度限制，表示该节点的主页)
	"details": "",	                      // 节点的描述(有长度限制，表示该节点的描述)
	"delegateTotal": 1000,                // 节点被委托的生效总数量
	"delegateTotalHes": 10000,            // 节点被委托的未生效的总数量
	"delegateRewardTotal": 1000           // 节点当前已发放的总委托奖励
}
```

