package evm.data_type.MappingData;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.MappingArrayDataTypeContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple4;

import java.math.BigInteger;

/**
 * @title 测试：验证mapping数组类型
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class MappingArrayDataTypeContractTest extends ContractPrepareTest {

    private String keyValue,a1Value,a2Value,b1Value,b2Value;

    @Before
    public void before() {
       this.prepare();
        keyValue = driverService.param.get("keyValue");
        a1Value = driverService.param.get("a1Value");
        a2Value = driverService.param.get("a2Value");
        b1Value = driverService.param.get("b1Value");
        b2Value = driverService.param.get("b2Value");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "MappingDataTypeContract.mapping数组类型",sourcePrefix = "evm")
    public void testMappingContract() {

        MappingArrayDataTypeContract mappingArrayDataTypeContract = null;
        try {
            //合约部署
            mappingArrayDataTypeContract = MappingArrayDataTypeContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = mappingArrayDataTypeContract.getContractAddress();
            TransactionReceipt tx =  mappingArrayDataTypeContract.getTransactionReceipt().get();
            collector.logStepPass("mappingArrayDataTypeContract issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("mappingArrayDataTypeContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        //1、执行mapping数组赋值
        try {
            BigInteger keyValueBig = new BigInteger(keyValue);
            BigInteger a1ValueBig = new BigInteger(a1Value);
            BigInteger a2ValueBig = new BigInteger(a2Value);
            BigInteger b1ValueBig = new BigInteger(b1Value);
            BigInteger b2ValueBig = new BigInteger(b2Value);
            TransactionReceipt transactionReceipt = mappingArrayDataTypeContract.set(keyValueBig,a1ValueBig,a2ValueBig,b1ValueBig,b2ValueBig).send();
            collector.logStepPass("mappingArrayDataTypeContract 【执行mapping数组赋值 set()】 successfully.hash:" + transactionReceipt.getTransactionHash());
            //获取值getValueByKey()
            Tuple4<BigInteger, BigInteger, BigInteger, BigInteger> tuple4 = mappingArrayDataTypeContract.getValueByKey(keyValueBig).send();
            BigInteger actualValue1 = tuple4.getValue1();
            BigInteger actualValue2 = tuple4.getValue2();
            BigInteger actualValue3 = tuple4.getValue3();
            BigInteger actualValue4 = tuple4.getValue4();
            collector.logStepPass("mappingArrayDataTypeContract 【执行获取值getValueByKey()】 successful actualValue:" + tuple4.toString());
            collector.assertEqual(actualValue1,a1ValueBig, "checkout  execute success.");
            collector.assertEqual(actualValue2,a2ValueBig, "checkout  execute success.");
            collector.assertEqual(actualValue3,b1ValueBig, "checkout  execute success.");
            collector.assertEqual(actualValue4,b2ValueBig, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("mappingArrayDataTypeContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }

}
