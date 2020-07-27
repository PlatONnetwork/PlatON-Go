package evm.data_type.TypeConversion;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.TypeConversionBytesToUintContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 测试：测试不同类型转换(字节转换整型)
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class TypeConversionBytesToUintContractTest extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "TypeConversionContract.字节转换整型",sourcePrefix = "evm")
    public void testTypeConversionContract() {

        TypeConversionBytesToUintContract typeConversionBytesToUintContract = null;
        try {
            //合约部署
            typeConversionBytesToUintContract = TypeConversionBytesToUintContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = typeConversionBytesToUintContract.getContractAddress();
            TransactionReceipt tx =  typeConversionBytesToUintContract.getTransactionReceipt().get();
            collector.logStepPass("typeConversion issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("typeConversion deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        try {
            BigInteger expectValue = new BigInteger("1633837924");
            //1、执行字节转换大位整型 bytesToBigUint()
            BigInteger actualValue = typeConversionBytesToUintContract.bytesToBigUint().send();
            collector.logStepPass("typeConversion 执行【字节转换大位整型 bytesToBigUint()】 successfully.actualValue:" + actualValue);
            collector.assertEqual(actualValue,expectValue, "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("typeConversion Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        try {
            BigInteger expectValue = new BigInteger("97");
            //2、字节转换相同位数整数 bytesToSameUint()
            BigInteger actualValue = typeConversionBytesToUintContract.bytesToSameUint().send();
            collector.logStepPass("typeConversion 执行【字节转换相同位数整数 bytesToSameUint()】 successfully.actualValue:" + actualValue);
            collector.assertEqual(actualValue,expectValue, "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("typeConversion Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        try {
            BigInteger expectValue = new BigInteger("25444");
            //3、字节转换小位整型 bytesToSmallUint()
            BigInteger actualValue = typeConversionBytesToUintContract.bytesToSmallUint().send();
            collector.logStepPass("typeConversion 执行【字节转换小位整型 bytesToSmallUint()】 successfully.actualValue:" + actualValue);
            collector.assertEqual(actualValue,expectValue, "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("typeConversion Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}
