package wasm.function;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.CoinbaseFunction;
import org.junit.Test;
import wasm.beforetest.WASMContractPrepareTest;


/**
 *
 * @title 验证函数platon_coinbase
 * @description:
 * @author: liweic
 * @create: 2020/02/11
 */
public class CoinbaseFunctionTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "wasm.SpecialFunctionsA验证链上函数platon_coinbase",sourcePrefix = "wasm")
    public void Coinbasefunction() {

        try {
            prepare();
            CoinbaseFunction coinbase = CoinbaseFunction.deploy(web3j, transactionManager, provider).send();
            String contractAddress = coinbase.getContractAddress();
            String transactionHash = coinbase.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("CallerFunctionTest issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            String coinbaseaddr = coinbase.get_platon_coinbase().send();
            collector.logStepPass("getPlatONConibase函数返回值:" + coinbaseaddr);
            collector.assertEqual(coinbaseaddr, "1000000000000000000000000000000000000003");


        } catch (Exception e) {
            collector.logStepFail("CoinbaseFunctionTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
