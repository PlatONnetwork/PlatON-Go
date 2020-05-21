package evm.data_type.BasicDataType;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.BasicDataTypeDeleteContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple2;
import org.web3j.tuples.generated.Tuple6;

import java.math.BigInteger;

/**
 * @title 测试：数据类型操作符delete
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class BasicDataTypeDeleteContractTest extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "BasicDataTypeDelete.数据类型操作符delete",sourcePrefix = "evm")
    public void testBasicDataTypeContract() {
        BasicDataTypeDeleteContract basicDataTypeDeleteContract = null;
        try {
            //合约部署
            basicDataTypeDeleteContract = BasicDataTypeDeleteContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = basicDataTypeDeleteContract.getContractAddress();
            TransactionReceipt tx =  basicDataTypeDeleteContract.getTransactionReceipt().get();
            collector.logStepPass("BasicDataTypeDelete issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("BasicDataTypeDelete deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        //1、delete基本数据类型
        try {
            //初始化赋值
            TransactionReceipt initTx = basicDataTypeDeleteContract.initBasicData().send();
            collector.logStepPass("BasicDataTypeDelete 【执行初始化赋值操作 initBasicData（）】successfully hash:" + initTx.getTransactionHash());
            //执行delete基本数据类型
            TransactionReceipt transactionReceipt = basicDataTypeDeleteContract.deleteBasicData().send();
            collector.logStepPass("BasicDataTypeDelete 【执行delete基本数据类型 deleteBasicData（）】successfully hash:" + transactionReceipt.getTransactionHash());
            //获取基本数据类型数值
            Tuple6<Boolean, BigInteger, String, byte[], String, BigInteger> tuple6 = basicDataTypeDeleteContract.getBasicData().send();
            collector.logStepPass("BasicDataTypeDelete 【执行获取基本数据类型数值】successfully value:" + tuple6.toString());
            Boolean boolActual = tuple6.getValue1();
            BigInteger uintActual = tuple6.getValue2();
            String addrStr = tuple6.getValue3();
            byte[]  bytesActual = tuple6.getValue4();
            String str = tuple6.getValue5();
            BigInteger intActual = tuple6.getValue6();
            collector.assertEqual(boolActual,false, "checkout delete bool execute success.");
            collector.assertEqual(uintActual,new BigInteger("0"), "checkout delete uint execute success.");
            collector.assertEqual(addrStr,"lax1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqmscn5j", "checkout delete address execute success.");
            //collector.assertEqual(bytesActual,new byte[]("0x0"), "checkout delete bytes execute success.");
            collector.assertEqual(str,"", "checkout delete string execute success.");
            collector.assertEqual(intActual,new BigInteger("0"), "checkout delete int execute success.");
        } catch (Exception e) {
            collector.logStepFail("BasicDataTypeDelete Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //2、delete结构体
        try {
            TransactionReceipt transactionReceipt = basicDataTypeDeleteContract.deleteStruct().send();
            collector.logStepPass("BasicDataTypeDelete 【执行delete结构体 deleteStruct（）】successfully hash:" + transactionReceipt.getTransactionHash());
            //获取结构体数值
            Tuple2<BigInteger, String> tuple2 =basicDataTypeDeleteContract.getStruct().send();
            collector.logStepPass("BasicDataTypeDelete 【执行获取结构体数值】successfully value:" + tuple2.toString());
            BigInteger uintActual = tuple2.getValue1();
            String addrStr = tuple2.getValue2();
            collector.assertEqual(uintActual,new BigInteger("0"), "checkout delete struct execute success.");
            collector.assertEqual(addrStr,"", "checkout delete struct execute success.");
        } catch (Exception e) {
            collector.logStepFail("BasicDataTypeDelete Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //3、delete数组
        try {
            TransactionReceipt transactionReceipt = basicDataTypeDeleteContract.deleteArray().send();
            collector.logStepPass("BasicDataTypeDelete 【执行delete数组 deleteArray（）】successfully hash:" + transactionReceipt.getTransactionHash());
            //获取数组长度
            BigInteger arrayLength =basicDataTypeDeleteContract.getArrayLength().send();
            collector.logStepPass("BasicDataTypeDelete 【获取数组长度】successfully value:" + arrayLength);
            collector.assertEqual(arrayLength,new BigInteger("0"), "checkout delete struct execute success.");
        } catch (Exception e) {
            collector.logStepFail("BasicDataTypeDelete Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //3、delete枚举
        try {
            TransactionReceipt transactionReceipt = basicDataTypeDeleteContract.deleteEnum().send();
            collector.logStepPass("BasicDataTypeDelete 【执行delete枚举 deleteEnum（）】successfully hash:" + transactionReceipt.getTransactionHash());
            //获取枚举值
            BigInteger enumValue =basicDataTypeDeleteContract.getEnum().send();
            collector.logStepPass("BasicDataTypeDelete 【获取枚举数值】successfully value:" + enumValue);
            collector.assertEqual(enumValue,new BigInteger("0"), "checkout delete struct execute success.");
        } catch (Exception e) {
            collector.logStepFail("BasicDataTypeDelete Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //4、deleteMapping
        try {
            TransactionReceipt transactionReceipt = basicDataTypeDeleteContract.deleteMapping().send();
            collector.logStepPass("BasicDataTypeDelete 【执行delete映射 deleteMapping（）】successfully hash:" + transactionReceipt.getTransactionHash());
            //获取映射值
            BigInteger mappingValue =basicDataTypeDeleteContract.getMapping().send();
            collector.logStepPass("BasicDataTypeDelete 【获取映射数值】successfully value:" + mappingValue);
            collector.assertEqual(mappingValue,new BigInteger("0"), "checkout delete struct execute success.");
        } catch (Exception e) {
            collector.logStepFail("BasicDataTypeDelete Calling Method fail.", e.toString());
            e.printStackTrace();
        }


    }
}
