package network.platon.test.evm.function.specialVariablesAndFunctions;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.ABIFunctions;
import network.platon.utils.DataChangeUtil;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

/**
 * @title 验证ABI编解码函数
 * @description:
 * @author: liweic
 * @create: 2019/12/30 19:01
 **/

public class ABIFunctionsTest extends ContractPrepareTest {

    private String encodeWithSignature;
    private String encode;
    private String encodePacked;

    @Before
    public void before() {
        this.prepare();
        encodeWithSignature = driverService.param.get("encodeWithSignature");
        encode = driverService.param.get("encode");
        encodePacked = driverService.param.get("encodePacked");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.ABIFunctionsTest-ABI函数测试", sourcePrefix = "evm")
    public void ABIfunction() {
        try {
            ABIFunctions abiFunctions = ABIFunctions.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = abiFunctions.getContractAddress();
            TransactionReceipt tx = abiFunctions.getTransactionReceipt().get();
            collector.logStepPass("ABIFunctionsTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("ABIFunctionsTest deploy gasUsed:" + abiFunctions.getTransactionReceipt().get().getGasUsed());

            //验证abi.encodeWithSignature函数
            byte[] resultA = abiFunctions.getEncodeWithSignature().send();
            String hexValue = DataChangeUtil.bytesToHex(resultA);
            collector.logStepPass("getEncodeWithSignature函数返回值：" + hexValue);
            collector.assertEqual(encodeWithSignature ,hexValue);

            //验证abi.encode函数
            byte[] resultB = abiFunctions.getEncode().send();
            String hexValue2 = DataChangeUtil.bytesToHex(resultB);
            collector.logStepPass("getEncode函数返回值：" + hexValue2);
            collector.assertEqual(encode ,hexValue2);

            //验证abi.encodePacked函数
            byte[] resultC = abiFunctions.getEncodePacked().send();
            String hexValue3 = DataChangeUtil.bytesToHex(resultC);
            collector.logStepPass("getEncodePacked函数返回值：" + hexValue3);
            collector.assertEqual(encodePacked ,hexValue3);

        } catch (Exception e) {
            collector.logStepFail("ABIFuctionsContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}
