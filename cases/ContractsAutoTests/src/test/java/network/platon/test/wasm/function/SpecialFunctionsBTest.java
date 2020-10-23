package network.platon.test.wasm.function;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.SpecialFunctionsB;
import org.junit.Test;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;

/**
 *
 * @title 验证函数platon_gas,platon_gas_limit,platon_gas_price
 * @description:
 * @author: liweic
 * @create: 2020/02/07
 */
public class SpecialFunctionsBTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "wasm.SpecialFunctionsB验证gas相关函数",sourcePrefix = "wasm")
    public void SpecialfunctionsB() {

        try {
            prepare();
//            provider = new ContractGasProvider(BigInteger.valueOf(50000000004L), BigInteger.valueOf(90000000L));
            SpecialFunctionsB specialfunctionsb = SpecialFunctionsB.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = specialfunctionsb.getContractAddress();
            String transactionHash = specialfunctionsb.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("SpecialFunctionsBTest issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("SpecialFunctionsBTest deploy gasUsed:" + specialfunctionsb.getTransactionReceipt().get().getGasUsed());

            Uint64 gas =specialfunctionsb.getPlatONGas().send();
            collector.logStepPass("getPlatONGas函数返回值:" + gas.value);
            boolean result = "0".toString().equals(gas.value.toString());
            collector.assertEqual(result, false);

            Uint64 gaslimit = specialfunctionsb.getPlatONGasLimit().send();
            collector.logStepPass("getPlatONGasLimit函数返回值:" + gaslimit.value);
            int compareresult = gaslimit.value.compareTo(new BigInteger("4712388"));
            boolean resulta = compareresult >= 0;
            collector.assertEqual(resulta, true);

            String gasprice = specialfunctionsb.getPlatONGasPrice().send();
            collector.logStepPass("getPlatONGasPrice函数返回值:" + gasprice);
            boolean resultb = "0".toString().equals(gasprice.toString());
            collector.assertEqual(resultb, true);

        } catch (Exception e) {
            collector.logStepFail("SpecialFunctionsBTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}

