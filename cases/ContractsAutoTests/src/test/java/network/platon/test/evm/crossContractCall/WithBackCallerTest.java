package network.platon.test.evm.crossContractCall;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.WithBackCallee;
import network.platon.contracts.evm.WithBackCaller;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigDecimal;
import java.math.BigInteger;


/**
 * @title 0.5.13跨合约调用者,接收返回值
 * @description:
 * @author: hudenian
 * @create: 2019/12/28
 */
public class WithBackCallerTest extends ContractPrepareTest {

    //需要进行double的值
    private String doubleValue = "10";

    //需要前缀hello的值
    private String helloValue = "hudenian";

    @Before
    public void before() {
        this.prepare();
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "WithBackCallerTest-跨合约调用者对反回值进行编码与解码", sourcePrefix = "evm")
    public void crossContractCaller() {
        try {
            //调用者合约地址
            WithBackCaller withBackCaller = WithBackCaller.deploy(web3j, transactionManager, provider, chainId).send();
            String callerContractAddress = withBackCaller.getContractAddress();
            TransactionReceipt tx = withBackCaller.getTransactionReceipt().get();
            collector.logStepPass("WithBackCaller deploy successfully.contractAddress:" + callerContractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + withBackCaller.getTransactionReceipt().get().getGasUsed());


            //被调用者合约地址
            WithBackCallee withBackCallee = WithBackCallee.deploy(web3j, transactionManager, provider, chainId).send();
            String calleeContractAddress = withBackCallee.getContractAddress();
            TransactionReceipt tx1 = withBackCallee.getTransactionReceipt().get();
            collector.logStepPass("WithBackCallee deploy successfully.contractAddress:" + calleeContractAddress + ", hash:" + tx1.getTransactionHash());
            collector.logStepPass("deploy gas used:" + withBackCallee.getTransactionReceipt().get().getGasUsed());

            //数值类型跨合约调用
            TransactionReceipt tx2 = withBackCaller.callDoublelTest(calleeContractAddress,new BigInteger(doubleValue)).send();
            collector.logStepPass("WithBackCaller callDoublelTest successfully hash:" + tx2.getTransactionHash());
            //获取数值类型跨合约调用的结果
            String chainDoubleValue = withBackCaller.getuintResult().send().toString();
            collector.logStepPass("获取数值类型跨合约调用的结果值为:" + chainDoubleValue);
            collector.assertEqual(new BigDecimal(doubleValue).multiply(new BigDecimal("2")),new BigDecimal(chainDoubleValue));


            //字符串类型跨合约调用
            tx2 = withBackCaller.callgetNameTest(calleeContractAddress,helloValue).send();
            collector.logStepPass("WithBackCaller callDoublelTest successfully hash:" + tx2.getTransactionHash()+"gas消耗值为："+tx2.getGasUsed());
            //获取字符串类型跨合约调用
            String callerStringResult = withBackCaller.getStringResult().send().toString();
            collector.logStepPass("获取字符串类型跨合约调用的结果值为:" + callerStringResult);
            collector.assertEqual(callerStringResult,"hello"+helloValue);


            //调用被调用合约，指定足够的gas
            tx2 = withBackCaller.callgetNameTestWithGas(calleeContractAddress,helloValue,new BigInteger("90000")).send();
            collector.logStepPass("withBackCaller callgetNameTestWithGas successfully hash:" + tx2.getTransactionHash()+"gas消耗值为："+tx2.getGasUsed());
            //获取字符串类型跨合约调用
            callerStringResult = withBackCaller.getStringResult().send().toString();
            collector.logStepPass("调用被调用合约，指定足够的gas查询到的返回值为:" + callerStringResult);
            collector.assertEqual(callerStringResult,"hellogashudenian");

            //调用被调用合约，指定gas不足
            try {
                tx2 = withBackCaller.callgetNameTestWithGas(calleeContractAddress, helloValue, new BigInteger("100")).send();

                collector.logStepPass("withBackCaller callgetNameTestWithGas with less gas hash:" + tx2.getTransactionHash() + "gas消耗值为：" + tx2.getGasUsed());
                //获取字符串类型跨合约调用
                callerStringResult = withBackCaller.getStringResult().send().toString();
                collector.logStepPass("调用被调用合约，指定gas不足时查询到的返回值为:" + callerStringResult);
            }catch(Exception e1){
                collector.logStepPass("指定gas不足时，调用合约失败！");
            }


        } catch (Exception e) {
            collector.logStepFail("WithBackCallerTest process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
