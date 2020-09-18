package network.platon.test.evm.oop.inherit;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.InheritContractOverloadChild;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 测试：合约函数重载(Overload)复杂情况
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class InheritContractOverloadComplexTest extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "InheritContract.合约函数重载(Overload)复杂情况",sourcePrefix = "evm")
    public void testInheritContractMutipleTest1() {

        InheritContractOverloadChild inheritContractOverloadChild = null;
        try {
            //合约部署
            inheritContractOverloadChild = InheritContractOverloadChild.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = inheritContractOverloadChild.getContractAddress();
            TransactionReceipt tx =  inheritContractOverloadChild.getTransactionReceipt().get();
            collector.logStepPass("InheritContractOverload issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("InheritContractOverload deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        //1、执行Base父类重载方法单个参数
        try {
            BigInteger param1 = new BigInteger("3");
            TransactionReceipt transactionReceipt = inheritContractOverloadChild.initBase(param1).send();
            collector.logStepPass("执行【Base父类赋值重载方法单个参数 initBase()】 successful，hash:" + transactionReceipt.getTransactionHash());
            //执行查询
            BigInteger actualValue = inheritContractOverloadChild.getX().send();
            collector.logStepPass("查询赋值结果 getX() successful，actualValue:" + actualValue);
            collector.assertEqual(actualValue,param1, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("InheritContractOverload Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //2、执行Base父类重载方法多个参数
        try {
            BigInteger param1 = new BigInteger("3");
            BigInteger param2 = new BigInteger("4");
            TransactionReceipt transactionReceipt = inheritContractOverloadChild.initBase(param1,param2).send();
            collector.logStepPass("执行【Base父类赋值重载方法多个参数 initBase()】 successful，hash:" + transactionReceipt.getTransactionHash());
            //执行查询
            BigInteger actualValueX = inheritContractOverloadChild.getX().send();
            BigInteger actualValueY = inheritContractOverloadChild.getY().send();
            collector.logStepPass("查询赋值结果 getX() successful，actualValueX:" + actualValueX);
            collector.logStepPass("查询赋值结果 getY() successful，actualValueY:" + actualValueY);
            collector.assertEqual(actualValueX,param1, "checkout  execute success.");
            collector.assertEqual(actualValueY,param2, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("InheritContractOverload Calling Method fail.", e.toString());
            e.printStackTrace();
        }


        //3、执行BaseBase父类重载方法单个参数
        try {
            BigInteger param1 = new BigInteger("3");
            BigInteger expectResult = new BigInteger("4");
            TransactionReceipt transactionReceipt = inheritContractOverloadChild.initBaseBase(param1).send();
            collector.logStepPass("执行【BaseBase父类重载方法单个参数 initBase()】 successful，hash:" + transactionReceipt.getTransactionHash());
            //执行查询
            BigInteger actualValueX = inheritContractOverloadChild.getX().send();
            collector.logStepPass("查询赋值结果 getX() successful，actualValueX:" + actualValueX);
            collector.assertEqual(actualValueX,expectResult, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("InheritContractOverload Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //4、执行BaseBase父类重载方法多个参数
        try {
            BigInteger param1 = new BigInteger("3");
            BigInteger param2 = new BigInteger("4");
            BigInteger expectResultY = new BigInteger("3");
            BigInteger expectResultX = new BigInteger("4");
            TransactionReceipt transactionReceipt = inheritContractOverloadChild.initBaseBase(param1,param2).send();
            collector.logStepPass("执行【BaseBase父类重载方法多个参数 initBase()】 successful，hash:" + transactionReceipt.getTransactionHash());
            //执行查询
            BigInteger actualValueX = inheritContractOverloadChild.getX().send();
            BigInteger actualValueY = inheritContractOverloadChild.getY().send();
            collector.logStepPass("查询赋值结果 getX() successful，actualValueX:" + actualValueX);
            collector.logStepPass("查询赋值结果 getY() successful，actualValueY:" + actualValueY);
            collector.assertEqual(actualValueX,expectResultX, "checkout  execute success.");
            collector.assertEqual(actualValueY,expectResultY, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("InheritContractOverload Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }

}
