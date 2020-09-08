package network.platon.test.wasm.function;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.OriginFunction;
import org.junit.Test;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 *
 * @title 验证函数platon_origin
 * @description:
 * @author: liweic
 * @create: 2020/02/11
 */
public class OriginFunctionTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "wasm.OriginFunctionTest验证链上函数platon_origin",sourcePrefix = "wasm")
    public void Originfunction() {

        try {
            prepare();
            OriginFunction origin = OriginFunction.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = origin.getContractAddress();
            String transactionHash = origin.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("OriginFunction issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("OriginFunction deploy gasUsed:" + origin.getTransactionReceipt().get().getGasUsed());

            String originaddr = origin.get_platon_origin().send();
            collector.logStepPass("getPlatONOrigin函数返回值:" + originaddr);
            collector.assertEqual(originaddr, walletAddress.toLowerCase());


        } catch (Exception e) {
            collector.logStepFail("OriginFunctionTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}


