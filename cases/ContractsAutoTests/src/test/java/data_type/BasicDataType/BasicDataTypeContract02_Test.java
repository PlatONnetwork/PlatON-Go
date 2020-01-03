//package data_type.BasicDataType;
//
//import beforetest.ContractPrepareTest;
//import network.platon.autotest.junit.annotations.DataSource;
//import network.platon.autotest.junit.enums.DataSourceType;
//import network.platon.contracts.BasicDataTypeContract2;
//import network.platon.utils.DataChangeUtil;
//import org.junit.Before;
//import org.junit.Test;
//import org.web3j.protocol.core.methods.response.TransactionReceipt;
//import org.web3j.tuples.generated.Tuple3;
//
//import java.math.BigInteger;
//
///**
// * @title 测试：合约基本数据类型字面常量
// * @description:
// * @author: qudong
// * @create: 2019/12/25 15:09
// **/
//public class BasicDataTypeContract02_Test extends ContractPrepareTest {
//
//    @Before
//    public void before() {
//       this.prepare();
//    }
//
//    @Test
//    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "02BasicDataTypeContract.合约基本数据类型字面常量")
//    public void testBasicDataTypeContract() {
//
//        BasicDataTypeContract2 basicDataTypeContract2 = null;
//        try {
//            //合约部署
//            basicDataTypeContract2 = basicDataTypeContract2.deploy(web3j, transactionManager, provider).send();
//            String contractAddress = basicDataTypeContract2.getContractAddress();
//            TransactionReceipt tx =  basicDataTypeContract2.getTransactionReceipt().get();
//            collector.logStepPass("basicDataTypeContract2 issued successfully.contractAddress:" + contractAddress
//                                    + ", hash:" + tx.getTransactionHash());
//            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
//        } catch (Exception e) {
//            collector.logStepFail("basicDataTypeContract2 deploy fail.", e.toString());
//            e.printStackTrace();
//        }
//
//        //调用合约方法
//
//        //1、验证：address类型，转账给指定地址
//        try {
//            BigInteger moneyValue1 = new BigInteger("100");
//            BigInteger moneyValue2 = new BigInteger("100");
//            int expectValue = 200;
//
//            String myContractAddress = "0x8a9B36694F1eeeb500c84A19bB34137B05162EC9";
//
//            //1)、查询地址余额
//            BigInteger contractBalance = basicDataTypeContract2.getBalance(myContractAddress).send();
//            collector.logStepPass("BasicDataTypeContract  执行 1)、查询地址余额 getBalance() successfully.contractBalance:" + contractBalance);
//
//            //2）、查询当前转账人账户余额getBalance()
//            BigInteger currentTransferUserBalance = basicDataTypeContract2.getBalance(walletAddress).send();
//            collector.logStepPass("BasicDataTypeContract 执行 2)、转账人账户余额getBalance() successfully.currentTransferUserBalance: " + currentTransferUserBalance);
//
//            //3)、给当前地址第一次转账
//            TransactionReceipt transactionReceipt = basicDataTypeContract2.goTransfer(myContractAddress,moneyValue1).send();
//            collector.logStepPass("BasicDataTypeContract 执行 3)、给当前地址转账  goTransfer() successfully hash:" + transactionReceipt.getTransactionHash());
//
//            //4)、给当前地址第二次转账
//            TransactionReceipt transactionReceipt1 = basicDataTypeContract2.goSend(myContractAddress,moneyValue2).send();
//            collector.logStepPass("BasicDataTypeContract 执行 4)、给当前地址转账  goSend() successfully hash:" + transactionReceipt1.getTransactionHash());
//
//            //5)、查询转账完毕，地址余额
//            BigInteger currentContractBalance = basicDataTypeContract2.getBalance(myContractAddress).send();
//            collector.logStepPass("BasicDataTypeContract 执行 5)、查询转账完毕，地址余额 getBalance() successfully currentContractBalance:" + currentContractBalance);
//            int actualValue = currentContractBalance.intValue() - contractBalance.intValue();
//
//            collector.assertEqual(actualValue,expectValue, "checkout execute success.");
//
//        } catch (Exception e) {
//            collector.logStepFail("BasicDataTypeContract Calling Method fail.", e.toString());
//            e.printStackTrace();
//        }
//
//        //2、验证：address类型，获取当前合约余额
//        try {
//            BigInteger expectValue = new BigInteger("0");
//            //执行getCurrentBalance()
//            BigInteger actualValue = basicDataTypeContract2.getCurrentBalance().send();
//            collector.logStepPass("BasicDataTypeContract 执行查询当前合约中余额 getCurrentBalance() successfully contractBalance:" + actualValue);
//            collector.assertEqual(actualValue,expectValue, "checkout  execute success.");
//        } catch (Exception e) {
//            collector.logStepFail("BasicDataTypeContract Calling Method fail.", e.toString());
//            e.printStackTrace();
//        }
//
//        //3、验证：字符串字面常量，进行赋值
//        try {
//            String expectValue = "hello";
//            //执行getCurrentBalance()
//            String actualValue = basicDataTypeContract2.getStr1().send();
//            collector.logStepPass("BasicDataTypeContract  字符串字面常量进行赋值 执行getStr1() successfully actualValue:" + actualValue);
//            collector.assertEqual(actualValue,expectValue, "checkout  execute success.");
//        } catch (Exception e) {
//            collector.logStepFail("BasicDataTypeContract Calling Method fail.", e.toString());
//            e.printStackTrace();
//        }
//
//        //4、验证：字符串是特殊的动态长度字节数组，获取长度
//        try {
//            BigInteger expectValue = new BigInteger("5");
//            //执行getCurrentBalance()
//            BigInteger actualValue = basicDataTypeContract2.getStr1Length().send();
//            collector.logStepPass("BasicDataTypeContract  验证字符串字面常量 执行getStr1Length() successfully actualValue:" + actualValue);
//            collector.assertEqual(actualValue,expectValue, "checkout  execute success.");
//        } catch (Exception e) {
//            collector.logStepFail("BasicDataTypeContract Calling Method fail.", e.toString());
//            e.printStackTrace();
//        }
//
//        //5、验证：字符串是特殊的动态长度字节数组，进行转换
//        try {
//            String expectValue = "aello";
//            //赋值执行setStr1()
//            String actualValue = basicDataTypeContract2.setStr1().send();
//            collector.logStepPass("BasicDataTypeContract  验证字符串字面常量  执行setStr1() successfully actualValue:" + actualValue);
//            collector.assertEqual(actualValue,expectValue, "checkout  execute success.");
//        } catch (Exception e) {
//            collector.logStepFail("BasicDataTypeContract Calling Method fail.", e.toString());
//            e.printStackTrace();
//        }
//
//
//        //6、验证：十六进制字面常量，定义赋值取值
//        try {
//            //执行getHexLiteral()
//            String expectValue = "c8";
//            byte[] byteValue = basicDataTypeContract2.getHexLiteral().send();
//            String actualValue = DataChangeUtil.bytesToHex(byteValue);
//            collector.logStepPass("BasicDataTypeContract  验证十六进制字面常量 执行getHexLiteral() successfully actualValue:" + actualValue);
//            collector.assertEqual(actualValue,expectValue, "checkout  execute success.");
//
//            //执行getHexLitera2()
//            String expectValue2 = "01f4";
//            byte[] byteValue2 = basicDataTypeContract2.getHexLitera2().send();
//            String actualValue2 = DataChangeUtil.bytesToHex(byteValue2);
//            collector.logStepPass("BasicDataTypeContract   验证十六进制字面常量  执行getHexLitera2() successfully actualValue2:" + actualValue2);
//            collector.assertEqual(actualValue2,expectValue2, "checkout  execute success.");
//
//            //执行getHexLitera3()
//            String expectValueA = "01f4";
//            String expectValueB = "01";
//            String expectValueC = "f4";
//            Tuple3<byte[], byte[], byte[]> tuple3 = basicDataTypeContract2.getHexLitera3().send();
//            String actualValueA = DataChangeUtil.bytesToHex(tuple3.getValue1());
//            String actualValueB = DataChangeUtil.bytesToHex(tuple3.getValue2());
//            String actualValueC = DataChangeUtil.bytesToHex(tuple3.getValue3());
//            collector.logStepPass("BasicDataTypeContract   验证十六进制字面常量  执行getHexLitera3() successfully  actualValueA:" + actualValueA +
//                                  ",  actualValueB:" + actualValueB + ",  actualValueC:" + actualValueC);
//
//            collector.assertEqual(actualValueA,expectValueA, "checkout  execute success.");
//            collector.assertEqual(actualValueB,expectValueB, "checkout  execute success.");
//            collector.assertEqual(actualValueC,expectValueC, "checkout  execute success.");
//
//        } catch (Exception e) {
//            collector.logStepFail("BasicDataTypeContract Calling Method fail.", e.toString());
//            e.printStackTrace();
//        }
//
//        //7、验证：枚举类型
//        try {
//            BigInteger expectValue1 = new BigInteger("1");
//            BigInteger expectValue2 = new BigInteger("3");
//            //执行getSeason1()
//            BigInteger actualValue1 = basicDataTypeContract2.getSeason1().send();
//            collector.logStepPass("BasicDataTypeContract  验证枚举类型执行 getSeason1() successfully actualValue1:" + actualValue1);
//            collector.assertEqual(actualValue1,expectValue1, "checkout  execute success.");
//
//            //执行getSeason2()
//            BigInteger actualValue2 = basicDataTypeContract2.getSeason2().send();
//            collector.logStepPass("BasicDataTypeContract 验证枚举类型执行 getSeason2() successfully actualValue2:" + actualValue2);
//            collector.assertEqual(actualValue2,expectValue2, "checkout  execute success.");
//
//        } catch (Exception e) {
//            collector.logStepFail("BasicDataTypeContract Calling Method fail.", e.toString());
//            e.printStackTrace();
//        }
//
//        //8、验证：有理数常量
//        try {
//            BigInteger expectValue1 = new BigInteger("3");
//            BigInteger expectValue2 = new BigInteger("20000000000");
//            BigInteger expectValue3 = new BigInteger("3000000000000000");
//            //执行getValue()
//            Tuple3<BigInteger, BigInteger, BigInteger> tuple3 = basicDataTypeContract2.getValue().send();
//            BigInteger actualValue1 = tuple3.getValue1();
//            BigInteger actualValue2 = tuple3.getValue2();
//            BigInteger actualValue3 = tuple3.getValue3();
//
//            collector.logStepPass("BasicDataTypeContract  验证有理数常量执行 getValue() successfully actualValue1:" + actualValue1 +
//                                  ",  actualValue2: " + actualValue2 + ",  actualValue3:" + actualValue3) ;
//            collector.assertEqual(actualValue1,expectValue1, "checkout  execute success.");
//            collector.assertEqual(actualValue2,expectValue2, "checkout  execute success.");
//            collector.assertEqual(actualValue3,expectValue3, "checkout  execute success.");
//
//        } catch (Exception e) {
//            collector.logStepFail("BasicDataTypeContract Calling Method fail.", e.toString());
//            e.printStackTrace();
//        }
//
//    }
//}
