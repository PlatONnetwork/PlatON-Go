package evm.csdccontracts;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.NewSecPledgeApplyManager;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.util.Date;


/**
 * @title 结算质押申请合约验证测试
 * @description:
 * @author: hudenian
 * @create: 2020/1/9
 */
public class NewSecPledgeApplyManagerTest extends ContractPrepareTest {


//    //模拟简单的业务数据
    private String secApply = "2-businessNo1-bizId1-4-5-6-7-8-9-10-11-12-13-14-15-16-17-18-19-20-21-22-23-24";

    //对应 secApply以“-”分隔的第一个值
    private String bizId = "2";

    @Before
    public void before() {
        this.prepare();
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "NewSecPledgeApplyManagerTest-结算复杂合约测试验证", sourcePrefix = "evm")
    public void createPledgeApplyCommonTest() {
        try {

            NewSecPledgeApplyManager newSecPledgeApplyManager = NewSecPledgeApplyManager.deploy(web3j, transactionManager, provider, chainId).send();
            String callerContractAddress = newSecPledgeApplyManager.getContractAddress();
            TransactionReceipt tx = newSecPledgeApplyManager.getTransactionReceipt().get();
            collector.logStepPass("NewSecPledgeApplyManager deploy successfully.contractAddress:" + callerContractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + newSecPledgeApplyManager.getTransactionReceipt().get().getGasUsed());

            Date start = new Date();


            tx = newSecPledgeApplyManager.createPledgeApplyCommon(secApply).send();
            //插入业务数据
            collector.logStepPass("NewSecPledgeApplyManager add successfully hash:" + tx.getTransactionHash());

            //查询质押申请信息
           String businessNo =  newSecPledgeApplyManager.select_SecPledgeApply_byId(bizId).send().toString();
           collector.logStepPass("bizId:"+bizId+"对应的business_no为："+businessNo);
            //对应 secApply以“-”分隔的第二个值
           collector.assertEqual("businessNo1",businessNo);

            //查询交易用户
            String tradeUser =  newSecPledgeApplyManager.select_tradeUser_byId(bizId).send().toString();
            collector.logStepPass("bizId:"+bizId+"对应的tradeUser为："+tradeUser);
            //对应 secApply以“-”分隔的第三个值
            collector.assertEqual("bizId1",tradeUser);

            //查询操作用户
            String tradeOperator =  newSecPledgeApplyManager.select_tradeOperator_bytId(bizId).send().toString();
            collector.logStepPass("bizId:"+bizId+"对应的tradeOperator为："+tradeOperator);
            //对应 secApply以“-”分隔的第七个值
            collector.assertEqual("7",tradeOperator);

            //查询pledgeSecurity
            String pledgeSecurity =  newSecPledgeApplyManager.select_pledgeSecurity_bytId(bizId).send().toString();
            collector.logStepPass("bizId:"+bizId+"对应的pledgeSecurity为："+pledgeSecurity);
            //对应 secApply以“-”分隔的第十四个值
            collector.assertEqual("14",pledgeSecurity);


            Date end = new Date();
            collector.logStepPass("插入到调用一共耗时："+(end.getTime()-start.getTime()));


        } catch (Exception e) {
            collector.logStepFail("NewSecPledgeApplyManagerTest process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
