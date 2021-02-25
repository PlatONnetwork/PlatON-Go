package network.platon.test.evm.versioncompatible.v0_5_0.v8_contractAndAddress;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.ContractAndAddress;
import org.junit.Before;
import org.junit.Test;


/**
 * @title 08-合约和地址
 * 1-contract合约类型不再包括 address类型的成员函数，
 * 必须显式转换成 address地址类型才能使用
 * @description:
 * @author: hudenian
 * @create: 2019/12/26
 */
public class ContractAdressTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "version_compatible.0.5.0.ContractAdressTest-合约与地址互相转换", sourcePrefix = "evm")
    public void addressAndPaybleChange() {
        try {

            ContractAndAddress contractAndAddress = ContractAndAddress.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = contractAndAddress.getNonalContractAddress().send();
            String dddressToPayable = contractAndAddress.getAddressToPayable().send();

            collector.logStepPass("ContractAdressTest force change address to payable");

            collector.assertEqual(contractAddress, dddressToPayable);


            String payableAddress = contractAndAddress.getNonalPayableAddress().send();
            String payableToAddress = contractAndAddress.getPayableToAddress().send();

            collector.logStepPass("ContractAdressTest change payable to address");

            collector.assertEqual(payableAddress, payableToAddress);

        } catch (Exception e) {
            collector.logStepFail("ContractAdressTest process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
