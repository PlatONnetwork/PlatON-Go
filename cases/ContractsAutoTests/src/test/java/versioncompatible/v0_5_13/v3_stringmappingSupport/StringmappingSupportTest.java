//package versioncompatible.v0_5_13.v3_stringmappingSupport;
//
//import complexcontracts.BaseTest;
//import network.platon.autotest.junit.annotations.DataSource;
//import network.platon.autotest.junit.enums.DataSourceType;
//import network.platon.contracts.PersonPublic;
//import network.platon.contracts.StringmappingSupport;
//import network.platon.utils.DataChangeUtil;
//import org.junit.Test;
//import org.web3j.protocol.core.methods.response.TransactionReceipt;
//import org.web3j.tx.RawTransactionManager;
//import org.web3j.tx.gas.ContractGasProvider;
//
//import java.math.BigInteger;
//
//
///**
// * @title   添加对具有string或bytes键类型的mapping的获取器的支持
// * @description: 此功能部署执行后内存溢出，后面重新验证测试
// * @author: hudenian
// * @create: 2019/12/27
// */
//public class StringmappingSupportTest extends BaseTest {
//
//    private String strKey="hdn";
//
//    private String strValue="hudenian_value";
//
//
//    @Test
//    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
//            author = "hudenian", showName = "version_compatible.0.5.13.StringmappingSupportTest-添加对具有string键类型的mapping的获取器的支持")
//    public void testStringMapping() {
//
//        ContractGasProvider provider = new ContractGasProvider(new BigInteger("50000000000"), new BigInteger("3000000"));
//        RawTransactionManager transactionManager = new RawTransactionManager(web3j, credentials, chainId);
//
//        try {
//
//            StringmappingSupport stringmappingSupport = StringmappingSupport.deploy(web3j, transactionManager, provider).send();
//
//            String contractAddress = stringmappingSupport.getContractAddress();
//            TransactionReceipt tx = stringmappingSupport.getTransactionReceipt().get();
//
//            collector.logStepPass("StringmappingSupportTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
//
//            TransactionReceipt transactionReceipt = stringmappingSupport.setStringmapValue(strKey,strValue).send();
//
//            collector.logStepPass("StringmappingSupportTest testMapping successful.transactionHash:" + transactionReceipt.getTransactionHash());
//            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());
//
//            String chainValue = stringmappingSupport.getStringmapValue(strKey).send().toString();
//
//            collector.assertEqual(strValue,chainValue);
//
//        } catch (Exception e) {
//            e.printStackTrace();
//        }
//    }
//
//    @Test
//    @DataSource(type = DataSourceType.EXCEL, file = "test1.xls", sheetName = "Sheet1",
//            author = "hudenian", showName = "version_compatible.0.5.13.StringmappingSupportTest-添加对具有bytes键类型的mapping的获取器的支持")
//    public void testBytesMapping() {
//
//        ContractGasProvider provider = new ContractGasProvider(new BigInteger("50000000000"), new BigInteger("3000000"));
//        RawTransactionManager transactionManager = new RawTransactionManager(web3j, credentials, chainId);
//
//        try {
//
//            StringmappingSupport stringmappingSupport = StringmappingSupport.deploy(web3j, transactionManager, provider).send();
//
//            String contractAddress = stringmappingSupport.getContractAddress();
//            TransactionReceipt tx = stringmappingSupport.getTransactionReceipt().get();
//
//            collector.logStepPass("StringmappingSupportTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
//
//            byte[] byte32Key = DataChangeUtil.stringToBytes32(strKey);
//            TransactionReceipt transactionReceipt = stringmappingSupport.setByte32mapValue(byte32Key,strValue).send();
//
//            collector.logStepPass("StringmappingSupportTest testBytesMapping successful.transactionHash:" + transactionReceipt.getTransactionHash());
//            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());
//
//            String chainValue1 = stringmappingSupport.getByte32mapValue(byte32Key).send().toString();
//
//            collector.assertEqual(strValue,chainValue1);
//
//        } catch (Exception e) {
//            e.printStackTrace();
//        }
//    }
//
//}
