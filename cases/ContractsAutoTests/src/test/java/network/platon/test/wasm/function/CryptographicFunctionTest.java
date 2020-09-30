package network.platon.test.wasm.function;

import com.platon.rlp.datatypes.WasmAddress;
import com.platon.sdk.utlis.Bech32;
import com.platon.sdk.utlis.NetworkParameters;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.CryptographicFunction;
import network.platon.utils.DataChangeUtil;
import org.junit.Test;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;


/**
 *
 * @title 内置函数
 * platon_ecrecover
 * platon_ripemd160
 * platon_sha256
 * @description:
 * @author: hudenian
 * @create: 2020/03/18
 */
public class CryptographicFunctionTest extends WASMContractPrepareTest {

    // 公钥
    String publicKey = "0x60320b8a71bc314404ef7d194ad8cac0bee1e331";
    // 被签名数据的哈希结果值 Hash.sha3("abc");
    String hexHash = "4e03657aea45a94fc7d47ba826c8d667c0d1e6e33a64a036ec44f58fa12d6c45";
    // 签名后的数据
    String hexSignature = "f4128988cbe7df8315440adde412a8955f7f5ff9a5468a791433727f82717a6753bd71882079522207060b681fbd3f5623ee7ed66e33fc8e581f442acbcf6ab800";



    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.CryptographicFunctionTest-验证内置函数",sourcePrefix = "wasm")
    public void cryptographicTest() {

        try {
            prepare();
            CryptographicFunction cryptographicFunction = CryptographicFunction.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = cryptographicFunction.getContractAddress();
            String transactionHash = cryptographicFunction.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("CryptographicFunctionTest issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("CryptographicFunctionTest deploy gasUsed:" + cryptographicFunction.getTransactionReceipt().get().getGasUsed());



            //platon_ecrecover
            byte[] hash = DataChangeUtil.hexToByteArray(hexHash);
            byte[] signature = DataChangeUtil.hexToByteArray(hexSignature);

            WasmAddress wasmAddress1 = cryptographicFunction.call_platon_ecrecover(hash,signature).send();
            String bech32Address = Bech32.addressEncode(NetworkParameters.TestNetParams.getHrp(), publicKey);
            collector.logStepPass("call_platon_ecrecover函数返回值:" + wasmAddress1.getAddress());
            collector.assertEqual(bech32Address,wasmAddress1.getAddress());

            //platon_ripemd160
            byte[] wasmAddress = cryptographicFunction.call_platon_ripemd160(DataChangeUtil.stringToBytes32("hellow")).send();
            collector.logStepPass("call_platon_ripemd160函数返回值:" + DataChangeUtil.bytesToHex(wasmAddress));

            //platon_sha256
            byte[] resultByte = cryptographicFunction.call_platon_sha256(DataChangeUtil.stringToBytes32("hellow")).send();
            collector.logStepPass("call_platon_sha256函数返回值:" + DataChangeUtil.bytesToHex(resultByte));


        } catch (Exception e) {
            collector.logStepFail("CryptographicFunctionTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}


