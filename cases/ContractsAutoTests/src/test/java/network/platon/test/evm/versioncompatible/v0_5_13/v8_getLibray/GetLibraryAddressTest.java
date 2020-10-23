package network.platon.test.evm.versioncompatible.v0_5_13.v8_getLibray;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.GetLibraryAddress;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.RawTransactionManager;
import org.web3j.tx.gas.ContractGasProvider;

import java.math.BigInteger;


/**
 * @title address(LibraryName)：获取链接库的地址
 * @description:
 * @author: hudenian
 * @create: 2019/12/28
 */
public class GetLibraryAddressTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "GetLibraryAddressTest-获取链接库的地址", sourcePrefix = "evm")
    public void getLibraryAddress() {

        ContractGasProvider provider = new ContractGasProvider(new BigInteger("50000000000"), new BigInteger("3000000"));
        RawTransactionManager transactionManager = new RawTransactionManager(web3j, credentials, chainId);

        try {

            GetLibraryAddress getLibraryAddress = GetLibraryAddress.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = getLibraryAddress.getContractAddress();
            TransactionReceipt tx = getLibraryAddress.getTransactionReceipt().get();

            collector.logStepPass("GetLibraryAddressTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + getLibraryAddress.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceipt = getLibraryAddress.setUserLibAddress().send();

            collector.logStepPass("StringmappingSupportTest testMapping successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            String userLibAddress = getLibraryAddress.getUserLibAddress().send().toString();

            collector.logStepPass("获取到的library库的地址是："+userLibAddress);

        } catch (Exception e) {
            collector.logStepFail("getLibraryAddress  process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
