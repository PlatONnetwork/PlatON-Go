package network.platon.test.evm.versioncompatible.v0_5_0.v10_deprecatedFunction;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.AssertContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/*
 * @title  assert() 代替 0.4.25版本函数 throw
 * @description:
 * @author: hudenian
 * @create: 2020/1/8
 */
public class AssertContractTest extends ContractPrepareTest {

    //减数
    protected String first;

    //被减数
    protected String second;

    @Before
    public void before() {
        this.prepare();
        first = driverService.param.get("first");
        second = driverService.param.get("second");
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "AssertContractTest-assert函数用法", sourcePrefix = "evm")
    public void assertTest() {
        try {

            AssertContract assertContract = AssertContract.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = assertContract.getContractAddress();
            TransactionReceipt tx = assertContract.getTransactionReceipt().get();

            collector.logStepPass("AssertContract deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + assertContract.getTransactionReceipt().get().getGasUsed());

            tx = assertContract.toSenderAmount(new BigInteger(first),new BigInteger(second)).send();

            if(Integer.valueOf(first).intValue()>Integer.valueOf(second)){
                String chainResult = assertContract.getResult().send().toString();
                collector.assertEqual(Integer.valueOf(chainResult).intValue(),Integer.valueOf(first).intValue()+Integer.valueOf(second).intValue());
            }

        } catch (Exception e) {
            if(Integer.valueOf(first).intValue()-Integer.valueOf(second).intValue()<0){
                collector.logStepPass("assert processed");
            }else{
                collector.logStepFail("AssertContractTest assertTest process fail.", e.toString());
                e.printStackTrace();
            }
        }
    }
}
