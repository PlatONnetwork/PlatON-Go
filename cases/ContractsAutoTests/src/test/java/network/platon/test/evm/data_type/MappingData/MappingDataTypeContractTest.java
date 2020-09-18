package network.platon.test.evm.data_type.MappingData;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.MappingDataTypeContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 测试：映射（Mapping）定义赋值取值
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class MappingDataTypeContractTest extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "MappingDataTypeContract.映射（Mapping）定义赋值取值",sourcePrefix = "evm")
    public void testMappingContract() {

        MappingDataTypeContract mappingContractTest = null;
        try {
            //合约部署
            mappingContractTest = MappingDataTypeContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = mappingContractTest.getContractAddress();
            TransactionReceipt tx =  mappingContractTest.getTransactionReceipt().get();
            collector.logStepPass("MappingContractTest issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("MappingContractTest deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        //1、验证：数组的声明及初始化及取值(定长数组、可变数组)
        try {
            String expectValue = "Lucy";
            BigInteger index = new BigInteger("0");
            //赋值执行addName()
            TransactionReceipt transactionReceipt = mappingContractTest.addName().send();
            collector.logStepPass("MappingContractTest 执行addName() successfully.hash:" + transactionReceipt.getTransactionHash());
            //获取值getName()
            String actualValue = mappingContractTest.getName(index).send();
            collector.logStepPass("调用合约getName()方法完毕 successful actualValue:" + actualValue);
            collector.assertEqual(actualValue,expectValue, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("MappingContractTest Calling Method fail.", e.toString());
            e.printStackTrace();
        }




    }

}
