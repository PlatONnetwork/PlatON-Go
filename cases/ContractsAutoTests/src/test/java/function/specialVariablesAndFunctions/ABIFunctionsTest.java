package function.specialVariablesAndFunctions;

import beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.ABIFunctions;
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

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.ABIFunctionsTest-ABI函数测试")
    public void ABIfunction() {
        try {
            ABIFunctions abiFunctions = ABIFunctions.deploy(web3j, transactionManager, provider).send();

            String contractAddress = abiFunctions.getContractAddress();
            TransactionReceipt tx = abiFunctions.getTransactionReceipt().get();
            collector.logStepPass("ABIFunctionsTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //验证abi.encodeWithSignature函数
            byte[] resultA = abiFunctions.getEncodeWithSignature().send();
            String hexValue = DataChangeUtil.bytesToHex(resultA);
            collector.logStepPass("getEncodeWithSignature函数返回值：" + hexValue);
            collector.assertEqual("60fe47b10000000000000000000000000000000000000000000000000000000000000001",hexValue);

            //验证abi.encode函数
            byte[] resultB = abiFunctions.getEncode().send();
            String hexValue2 = DataChangeUtil.bytesToHex(resultB);
            collector.logStepPass("getEncode函数返回值：" + hexValue2);
            collector.assertEqual("0000000000000000000000000000000000000000000000000000000000000001",hexValue2);

            //验证abi.encodePacked函数
            byte[] resultC = abiFunctions.getEncodePacked().send();
            String hexValue3 = DataChangeUtil.bytesToHex(resultC);
            collector.logStepPass("getEncodePacked函数返回值：" + hexValue3);
            collector.assertEqual("31",hexValue3);

        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
