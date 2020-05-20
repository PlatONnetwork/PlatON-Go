package evm.versioncompatible.v0_5_13.v3_stringmappingSupport;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.StringmappingSupport;
import network.platon.utils.DataChangeUtil;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.RawTransactionManager;
import org.web3j.tx.gas.ContractGasProvider;

import java.math.BigInteger;


/**
 * @title   添加对具有string或bytes键类型的mapping的获取器的支持
 * @description: 此功能部署执行后内存溢出，后面重新验证测试
 * @author: hudenian
 * @create: 2019/12/27
 */
public class StringmappingSupportTest extends ContractPrepareTest {

    private String strKey="hdn";

    private String strValue="hudenian_value";

    @Before
    public void before() {
        this.prepare();
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "StringmappingSupportTest-添加对具有string键类型的mapping的获取器的支持", sourcePrefix = "evm")
    public void testStringMapping() {

        ContractGasProvider provider = new ContractGasProvider(new BigInteger("50000000000"), new BigInteger("3000000"));
        RawTransactionManager transactionManager = new RawTransactionManager(web3j, credentials, chainId);

        try {

            StringmappingSupport stringmappingSupport = StringmappingSupport.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = stringmappingSupport.getContractAddress();
            TransactionReceipt tx = stringmappingSupport.getTransactionReceipt().get();

            collector.logStepPass("StringmappingSupportTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + stringmappingSupport.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceipt = stringmappingSupport.setStringmapValue(strKey,strValue).send();

            collector.logStepPass("StringmappingSupportTest testMapping successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            String chainValue = stringmappingSupport.getStringmapValue(strKey).send().toString();

            collector.assertEqual(strValue,chainValue);

        } catch (Exception e) {
            collector.logStepFail("StringmappingSupportTest testStringMapping process fail.", e.toString());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test1.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "StringmappingSupportTest-添加对具有bytes键类型的mapping的获取器的支持", sourcePrefix = "evm")
    public void testBytesMapping() {

        ContractGasProvider provider = new ContractGasProvider(new BigInteger("50000000000"), new BigInteger("3000000"));
        RawTransactionManager transactionManager = new RawTransactionManager(web3j, credentials, chainId);

        try {

            StringmappingSupport stringmappingSupport = StringmappingSupport.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = stringmappingSupport.getContractAddress();
            TransactionReceipt tx = stringmappingSupport.getTransactionReceipt().get();

            collector.logStepPass("StringmappingSupportTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + stringmappingSupport.getTransactionReceipt().get().getGasUsed());

            byte[] byte9Key = DataChangeUtil.stringToBytesN(strKey,9);
            TransactionReceipt transactionReceipt = stringmappingSupport.setByte32mapValue(byte9Key,strValue).send();

            collector.logStepPass("StringmappingSupportTest testBytesMapping successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            String chainValue1 = stringmappingSupport.getByte32mapValue(byte9Key).send().toString();

            collector.assertEqual(strValue,chainValue1);

        } catch (Exception e) {
            collector.logStepFail("StringmappingSupportTest testBytesMapping process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
