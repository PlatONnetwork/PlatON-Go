//package oop.inherit;
//
//import beforetest.ContractPrepareTest;
//import network.platon.autotest.junit.annotations.DataSource;
//import network.platon.autotest.junit.enums.DataSourceType;
//import network.platon.contracts.InheritContractMutipleTest1;
//import network.platon.contracts.InheritContractMutipleTest2;
//import org.junit.Before;
//import org.junit.Test;
//import org.web3j.protocol.core.methods.response.TransactionReceipt;
//
//import java.math.BigInteger;
//
///**
// * @title 测试：多重合约继承重名问题，遵循最远继承原则
// * @description:
// * @author: qudong
// * @create: 2019/12/25 15:09
// **/
//public class InheritContractMutiple01_Test extends ContractPrepareTest {
//
//
//    @Before
//    public void before() {
//       this.prepare();
//    }
//
//    @Test
//    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "01InheritContractMultipleTest.多重合约继承重名问题(遵循最远继承原则)")
//    public void testInheritContractMutipleTest1() {
//
//        InheritContractMutipleTest1 inheritContractMutipleTest1 = null;
//        InheritContractMutipleTest2 inheritContractMutipleTest2 = null;
//        try {
//            //合约部署(InheritContractMutipleTest1)
//            inheritContractMutipleTest1 = InheritContractMutipleTest1.deploy(web3j, transactionManager, provider).send();
//            String contractAddress = inheritContractMutipleTest1.getContractAddress();
//            TransactionReceipt tx =  inheritContractMutipleTest1.getTransactionReceipt().get();
//            collector.logStepPass("InheritContractMutipleTest1 issued successfully.contractAddress:" + contractAddress
//                                    + ", hash:" + tx.getTransactionHash());
//            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
//
//
//            //合约部署(InheritContractMutipleTest1)
//            inheritContractMutipleTest2 = InheritContractMutipleTest2.deploy(web3j, transactionManager, provider).send();
//            String contractAddress2 = inheritContractMutipleTest2.getContractAddress();
//            TransactionReceipt tx2 =  inheritContractMutipleTest2.getTransactionReceipt().get();
//            collector.logStepPass("InheritContractMutipleTest2 issued successfully.contractAddress:" + contractAddress2
//                    + ", hash:" + tx2.getTransactionHash());
//            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx2.getBlockNumber());
//
//        } catch (Exception e) {
//            collector.logStepFail("InheritContractMutipleTest deploy fail.", e.toString());
//            e.printStackTrace();
//        }
//
//        //调用合约方法
//        //1、执行callGetDate1()
//        try {
//            BigInteger expectBookResult = new BigInteger("1");
//            BigInteger actualBigInteger = inheritContractMutipleTest1.callGetDate1().send();
//            collector.logStepPass("调用合约callGetDate1()方法完毕 successful actualValue:" + actualBigInteger);
//            collector.assertEqual(actualBigInteger,expectBookResult, "checkout  execute success.");
//        } catch (Exception e) {
//            collector.logStepFail("InterfaceContractStructTest1 Calling Method fail.", e.toString());
//            e.printStackTrace();
//        }
//
//        //2、执行callGetDate2()
//        try {
//            BigInteger expectBookResult = new BigInteger("2");
//            BigInteger actualBigInteger = inheritContractMutipleTest2.callGetDate2().send();
//            collector.logStepPass("调用合约callGetDate2()方法完毕 successful actualValue:" + actualBigInteger);
//            collector.assertEqual(actualBigInteger,expectBookResult, "checkout  execute success.");
//        } catch (Exception e) {
//            collector.logStepFail("InterfaceContractStructTest2 Calling Method fail.", e.toString());
//            e.printStackTrace();
//        }
//    }
//
//}
