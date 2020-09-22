//package network.platon.test.evm.innerContract;
//
//import network.platon.test.evm.beforetest.ContractPrepareTest;
//import network.platon.autotest.junit.annotations.DataSource;
//import network.platon.autotest.junit.enums.DataSourceType;
//import network.platon.contracts.evm.Guessing;
//import network.platon.contracts.evm.PlatonUnit;
//import org.junit.Before;
//import org.junit.Test;
//import org.web3j.protocol.core.DefaultBlockParameterName;
//import org.web3j.protocol.core.methods.response.TransactionReceipt;
//import org.web3j.tx.Transfer;
//import org.web3j.utils.Convert;
//
//import java.math.BigDecimal;
//import java.math.BigInteger;
//
///**
// * @title Platon金额单位测试
// * @description:
// * @author: hudenian
// * @create: 2020/03/05 16:42
// *
// **/
//
//public class PlatonUnitTest extends ContractPrepareTest {
//    private String endBlock = "100000000"; //设置竞猜截止块高
//
//    @Before
//    public void before() {
//        this.prepare();
//    }
//
//    @Test
//    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
//            author = "hudenian", showName = "PlatonUnitTest-Platon金额单位测试", sourcePrefix = "network.platon.test.evm")
//    public void unitTest() {
//
//        try {
//
//            PlatonUnit platonUnit = PlatonUnit.deploy(web3j, transactionManager, provider, chainId).send();
//
//            String contractAddress = platonUnit.getContractAddress();
//            TransactionReceipt tx = platonUnit.getTransactionReceipt().get();
//            collector.logStepPass("PlatonUnitTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
//            collector.logStepPass("PlatonUnitTest deploy gasUsed:" + platonUnit.getTransactionReceipt().get().getGasUsed());
//
//            //发起转账
//            Transfer transfer = new Transfer(web3j, transactionManager);
//            TransactionReceipt transactionReceipt = transfer.sendFunds(contractAddress, new BigDecimal("1"), Convert.Unit.GLAT, new BigInteger("1000000000"), new BigInteger("4712388")).send();
//            transactionReceipt = transfer.sendFunds(contractAddress, new BigDecimal("1"), Convert.Unit.MLAT, new BigInteger("1000000000"), new BigInteger("4712388")).send();
//            transactionReceipt = transfer.sendFunds(contractAddress, new BigDecimal("1"), Convert.Unit.KLAT, new BigInteger("1000000000"), new BigInteger("4712388")).send();
//            transactionReceipt = transfer.sendFunds(contractAddress, new BigDecimal("1"), Convert.Unit.LAT, new BigInteger("1000000000"), new BigInteger("4712388")).send();
//            transactionReceipt = transfer.sendFunds(contractAddress, new BigDecimal("1"), Convert.Unit.FINNEY, new BigInteger("1000000000"), new BigInteger("4712388")).send();
//            transactionReceipt = transfer.sendFunds(contractAddress, new BigDecimal("1"), Convert.Unit.SZABO, new BigInteger("1000000000"), new BigInteger("4712388")).send();
//            transactionReceipt = transfer.sendFunds(contractAddress, new BigDecimal("1"), Convert.Unit.GVON, new BigInteger("1000000000"), new BigInteger("4712388")).send();
//            transactionReceipt = transfer.sendFunds(contractAddress, new BigDecimal("1"), Convert.Unit.MVON, new BigInteger("1000000000"), new BigInteger("4712388")).send();
//            transactionReceipt = transfer.sendFunds(contractAddress, new BigDecimal("1"), Convert.Unit.KVON, new BigInteger("1000000000"), new BigInteger("4712388")).send();
//            transactionReceipt = transfer.sendFunds(contractAddress, new BigDecimal("1"), Convert.Unit.VON, new BigInteger("1000000000"), new BigInteger("4712388")).send();
//
//            //查询合约余额
//            String contractBalance = platonUnit.getBalance().send().toString();
//            collector.logStepPass("合约中的余额为："+contractBalance);
//            collector.assertEqual(contractBalance,"1001001001001001001001001001");
//
//            //
//
//
//        } catch (Exception e) {
//            collector.logStepFail("PlatonUnitTest Calling Method fail.", e.toString());
//            e.printStackTrace();
//        }
//    }
//}
