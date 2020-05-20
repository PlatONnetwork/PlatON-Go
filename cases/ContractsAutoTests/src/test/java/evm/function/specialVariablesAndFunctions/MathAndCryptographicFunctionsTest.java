package evm.function.specialVariablesAndFunctions;

import com.platon.sdk.utlis.Bech32;
import com.platon.sdk.utlis.NetworkParameters;
import evm.beforetest.ContractPrepareTest;
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

    private String addmod;
    private String mulmod;
    private String keccak256;
    private String sha256;
    private String ripemd160;
    private String ecrecover;

    @Before
    public void before() {
        this.prepare();
        addmod = driverService.param.get("addmod");
        mulmod = driverService.param.get("mulmod");
        keccak256 = driverService.param.get("keccak256");
        sha256 = driverService.param.get("sha256");
        ripemd160 = driverService.param.get("ripemd160");
        ecrecover = driverService.param.get("ecrecover");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.MathAndCryptographicFunctionsTest-数学和加密函数测试", sourcePrefix = "evm")
    public void MathAndCryptographicfunction() {
        try {
            MathAndCryptographicFunctions mathAndCryptographicFunctions = MathAndCryptographicFunctions.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = mathAndCryptographicFunctions.getContractAddress();
            TransactionReceipt tx = mathAndCryptographicFunctions.getTransactionReceipt().get();
            collector.logStepPass("MathAndCryptographicFunctionsTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("MathAndCryptographicFunctionsTest deploy gasUsed:" + mathAndCryptographicFunctions.getTransactionReceipt().get().getGasUsed());

            //验证addmod函数
            BigInteger resultA = mathAndCryptographicFunctions.callAddMod().send();
            collector.logStepPass("addmod函数返回值：" + resultA);
            collector.assertEqual(addmod ,resultA.toString());

            //验证mulmod函数
            BigInteger resultB = mathAndCryptographicFunctions.callMulMod().send();
            collector.logStepPass("mulmod函数返回值：" + resultB);
            collector.assertEqual(mulmod ,resultB.toString());

            //验证keccak256函数
            byte[] resultC = mathAndCryptographicFunctions.callKeccak256().send();
            String hexValue = DataChangeUtil.bytesToHex(resultC);
            collector.logStepPass("keccak256函数返回值：" + hexValue);
            collector.assertEqual(keccak256 ,hexValue);

            //验证sha256函数
            byte[] resultD = mathAndCryptographicFunctions.callSha256().send();
            String hexValue2 = DataChangeUtil.bytesToHex(resultD);
            collector.logStepPass("sha256函数返回值：" + hexValue2);
            collector.assertEqual(sha256 ,hexValue2);

            //验证ripemd160函数
            byte[] resultE = mathAndCryptographicFunctions.callRipemd160().send();
            String hexValue3 = DataChangeUtil.bytesToHex(resultE);
            collector.logStepPass("ripemd160函数返回值：" + hexValue3);
            collector.assertEqual(ripemd160 ,hexValue3);

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
            String bech32Address = Bech32.addressEncode(NetworkParameters.TestNetParams.getHrp(), ecrecover);
            collector.assertEqual(bech32Address ,resultF.toLowerCase());

        } catch (Exception e) {
            collector.logStepFail("MathAndCryptographicfunctionsContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}
