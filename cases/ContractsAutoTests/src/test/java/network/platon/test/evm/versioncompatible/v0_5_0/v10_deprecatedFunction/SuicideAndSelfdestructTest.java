package network.platon.test.evm.versioncompatible.v0_5_0.v10_deprecatedFunction;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.SuicideAndSelfdestruct;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;


/**
 * @title  suicide()已弃用, 请使用 selfdestruct()
 * @description:
 * @author: hudenian
 * @create: 2019/12/27
 */
public class SuicideAndSelfdestructTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "SuicideAndSelfdestructTest-suicide", sourcePrefix = "evm")
    public void selfKill() {
        try {

            SuicideAndSelfdestruct suicideAndSelfdestruct = SuicideAndSelfdestruct.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = suicideAndSelfdestruct.getContractAddress();
            TransactionReceipt tx = suicideAndSelfdestruct.getTransactionReceipt().get();

            collector.logStepPass("SuicideAndSelfdestructTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + suicideAndSelfdestruct.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceipt = suicideAndSelfdestruct.increment().send();

            collector.logStepPass("SuicideAndSelfdestructTest increment successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            String count = suicideAndSelfdestruct.getCount().send().toString();

            collector.logStepPass("链上的count值为："+count);

            //调用自杀函数
            TransactionReceipt transactionReceipt1 = suicideAndSelfdestruct.kill().send();

            collector.logStepPass("SuicideAndSelfdestructTest kill successful.transactionHash:" + transactionReceipt1.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt1.getBlockNumber());


            BigInteger count1 = suicideAndSelfdestruct.getCount().send();

            collector.logStepPass("自杀后链上的count值为："+count1);

        } catch (Exception e) {
            if(e.getMessage().startsWith("Empty")){
                collector.logStepPass("自杀后查询链上的count值为 Empty");
            }else {
                collector.logStepFail("SuicideAndSelfdestructTest keccake256 process fail.", e.toString());
            }
        }
    }


    /**
     * @title do...while结果值
     * @description:
     * @author: hudenian
     * @create: 2019/12/27
     */
    public static String dowhile(int x){
        int y = x+10;
        int z = x+9;
        do{
            x+=1;
            if(x>z) continue;
        }while (x<y);
        return  String.valueOf(x);
    }

    /**
     * @title for 循环后的结果值
     * @description:
     * @author: hudenian
     * @create: 2019/12/27
     */
    public static String forsum(int x){
        int forSum = 0;
        for(int i=0;i<x;i++){
            forSum = forSum +i;
        }
        return  String.valueOf(forSum);
    }

}
