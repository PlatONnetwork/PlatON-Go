package network.platon.test.wasm.function;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.CoinbaseFunction;
import org.junit.Test;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;


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
            author = "liweic", showName = "wasm.CoinbaseFunctionTest验证链上函数platon_coinbase",sourcePrefix = "wasm")
    public void Coinbasefunction() {

        try {
            prepare();
            CoinbaseFunction coinbase = CoinbaseFunction.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = coinbase.getContractAddress();
            String transactionHash = coinbase.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("CoinbaseFunction issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("CoinbaseFunction deploy gasUsed:" + coinbase.getTransactionReceipt().get().getGasUsed());

            String coinbaseaddr = coinbase.get_platon_coinbase().send();
            collector.logStepPass("getPlatONConibase函数返回值:" + coinbaseaddr);

            String coinaddr = "0x0000000000000000000000000000000000000000";
            boolean iscoinbase = coinbaseaddr.equals(coinaddr);
            collector.assertEqual(!iscoinbase ,true);


        } catch (Exception e) {
            collector.logStepFail("CoinbaseFunctionTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
