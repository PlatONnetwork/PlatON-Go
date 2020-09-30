package network.platon.test.evm.oop.Interface;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.InterfaceContractParentTest;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 测试：在5.0以后版本，接口的函数只能声明外部类型(external)，否则会编译失败
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class InterfaceContractALimitTest extends ContractPrepareTest {

    private String param1,param2,sumParam;

    @Before
    public void before() {
       this.prepare();
        param1 = driverService.param.get("param1");
        param2 = driverService.param.get("param2");
        sumParam = driverService.param.get("sumParam");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "InterfaceContract.接口的函数只能声明外部类型",sourcePrefix = "evm")
    public void testInterfaceContractLimit() {

        InterfaceContractParentTest interfaceContractParentTest= null;
        try {
            //合约部署
            interfaceContractParentTest = InterfaceContractParentTest.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = interfaceContractParentTest.getContractAddress();
            TransactionReceipt tx = interfaceContractParentTest.getTransactionReceipt().get();

            collector.logStepPass("InterfaceContractParentTest issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());

            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());

        } catch (Exception e) {
            collector.logStepFail("InterfaceContractParentTest deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        try {
            BigInteger sumBigInt = interfaceContractParentTest.sumExternal(new BigInteger(param1),new BigInteger(param2)).send();
            collector.logStepPass("调用合约方法完毕 successful.sumBigInt:" + sumBigInt);
            collector.assertEqual(sumBigInt, new BigInteger(sumParam), "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("InterfaceContractParentTest Calling Method fail.", e.toString());
            e.printStackTrace();
        }



    }

}
