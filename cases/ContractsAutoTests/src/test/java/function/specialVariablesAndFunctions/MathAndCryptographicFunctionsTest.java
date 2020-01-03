package function.specialVariablesAndFunctions;

import beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.MathAndCryptographicFunctions;
import network.platon.utils.DataChangeUtil;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 验证数学和加密函数
 * @description:
 * @author: liweic
 * @create: 2019/12/30 20:01
 **/

public class MathAndCryptographicFunctionsTest extends ContractPrepareTest {

    @Before
    public void before() {

        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.MathAndCryptographicFunctionsTest-数学和加密函数测试")
    public void MathAndCryptographicfunction() {
        try {
            MathAndCryptographicFunctions mathAndCryptographicFunctions = MathAndCryptographicFunctions.deploy(web3j, transactionManager, provider).send();

            String contractAddress = mathAndCryptographicFunctions.getContractAddress();
            TransactionReceipt tx = mathAndCryptographicFunctions.getTransactionReceipt().get();
            collector.logStepPass("MathAndCryptographicFunctionsTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //验证addmod函数
            BigInteger resultA = mathAndCryptographicFunctions.callAddMod().send();
            collector.logStepPass("addmod函数返回值：" + resultA);
            collector.assertEqual("2",resultA.toString());

            //验证mulmod函数
            BigInteger resultB = mathAndCryptographicFunctions.callMulMod().send();
            collector.logStepPass("mulmod函数返回值：" + resultB);
            collector.assertEqual("0",resultB.toString());

            //验证keccak256函数
            byte[] resultC = mathAndCryptographicFunctions.callKeccak256().send();
            String hexValue = DataChangeUtil.bytesToHex(resultC);
            collector.logStepPass("keccak256函数返回值：" + hexValue);
            collector.assertEqual("e1629b9dda060bb30c7908346f6af189c16773fa148d3366701fbaa35d54f3c8",hexValue);

            //验证sha256函数
            byte[] resultD = mathAndCryptographicFunctions.callSha256().send();
            String hexValue2 = DataChangeUtil.bytesToHex(resultD);
            collector.logStepPass("sha256函数返回值：" + hexValue2);
            collector.assertEqual("b5d4045c3f466fa91fe2cc6abe79232a1a57cdf104f7a26e716e0a1e2789df78",hexValue2);

            //验证ripemd160函数
            byte[] resultE = mathAndCryptographicFunctions.callRipemd160().send();
            String hexValue3 = DataChangeUtil.bytesToHex(resultE);
            collector.logStepPass("ripemd160函数返回值：" + hexValue3);
            collector.assertEqual("df62d400e51d3582d53c2d89cfeb6e10d32a3ca6000000000000000000000000",hexValue3);

            //验证ecrecover函数
            String hash = "e281eaa11e6e37e6f53aade5d6c5b7201ef1c66162ec42ccc3215a0c4349350d";
            byte[] a = DataChangeUtil.hexToByteArray(hash);

            BigInteger v = new BigInteger("27");

            String R = "55b60cadd4b4a3ea4fc368ef338f97e12e7328dd6e9e969a3fd8e5c10be855fe";
            byte[] b = DataChangeUtil.hexToByteArray(R);

            String S = "2b42cee2585a16ea537efcb88009c1aeac693c28b59aa6bbff0baf22730338f6";
            byte[] c = DataChangeUtil.hexToByteArray(S);

            String resultF = mathAndCryptographicFunctions.callEcrecover(a, v, b, c).send();
            collector.logStepPass("ecrecover函数返回值：" + resultF);
            collector.assertEqual("0x8a9B36694F1eeeb500c84A19bB34137B05162EC5".toLowerCase(),resultF.toLowerCase());

        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
