package evm.evm_assembly_contract_call;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.Precompiled;
import network.platon.utils.DataChangeUtil;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.gas.ContractGasProvider;

import java.math.BigInteger;
import java.util.List;

/**
 * 添加evm合约调用系统合约场景
 *
 * @author hudenian
 * @dev 2020/02/19
 */

public class PrecompiledTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "evm_assembly_contract_call.AssemblyAddTest-evm合约调用系统合约", sourcePrefix = "evm")
    public void precompiledTest() {
        try {
            Precompiled precompiled = Precompiled.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = precompiled.getContractAddress();
            TransactionReceipt tx = precompiled.getTransactionReceipt().get();
            collector.logStepPass("PrecompiledTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + precompiled.getTransactionReceipt().get().getGasUsed());

            byte[] bytes = "hu".getBytes();


            //验证ecrecover函数
            String ecrecover = "lax132dnv620rmht2qxgfgvmkdqn0vz3vtk9lecarh";

            String hash = "e281eaa11e6e37e6f53aade5d6c5b7201ef1c66162ec42ccc3215a0c4349350d";
            byte[] a = DataChangeUtil.hexToByteArray(hash);

            BigInteger v = new BigInteger("27");

            String R = "55b60cadd4b4a3ea4fc368ef338f97e12e7328dd6e9e969a3fd8e5c10be855fe";
            byte[] b = DataChangeUtil.hexToByteArray(R);

            String S = "2b42cee2585a16ea537efcb88009c1aeac693c28b59aa6bbff0baf22730338f6";
            byte[] c = DataChangeUtil.hexToByteArray(S);

            //Address 0x01: ecrecover(hash, v, r, s)
            String resultF = precompiled.callEcrecover(a, v, b, c).send();
            collector.logStepPass("ecrecover函数返回值：" + resultF);
            collector.assertEqual(ecrecover, resultF.toLowerCase());

            //Address 0x02: sha256(data)
            byte[] resultD = precompiled.callSha256(bytes).send();
            String hexValue2 = DataChangeUtil.bytesToHex(resultD);
            collector.logStepPass("Sha256函数返回值：" + hexValue2);

            //Address 0x03: ripemd160(data)
            byte[] resultE = precompiled.callRipemd160(bytes).send();
            String hexValue3 = DataChangeUtil.bytesToHex(resultE);
            collector.logStepPass("ripemd160函数返回值：" + hexValue3);
//            collector.assertEqual(ripemd160 ,hexValue3);

            //0x04 dataCopy() test pass
            tx = precompiled.callDatacopy(bytes).send();
            collector.logStepPass("PrecompiledTest callDatacopy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            byte[] callDatacopyValueByte = precompiled.getCallDatacopyValue().send();
            collector.logStepPass("PrecompiledTest 0x04 result is:" + new String(callDatacopyValueByte));

            //0x05 callBigModExp() test pass
            byte[] base = DataChangeUtil.hexToByteArray("0000000000000000000000000000000462e4ded88953a39ce849a8a7fa163fa9");
            byte[] exponent = DataChangeUtil.hexToByteArray("1f4a3123ff1223a1b0d040057af8a9fe70baa9258e0b959273ffc5718c6d4cc7");
            byte[] modulus = DataChangeUtil.hexToByteArray("00000000000000000000000000077d29a9c710b7e616683f194f18c43b43b869");

            tx = precompiled.callBigModExp(base, exponent, modulus).send();
            collector.logStepPass("PrecompiledTest callBigModExp successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            byte[] callBigModExpValueByte = precompiled.getCallBigModExpValue().send();
            collector.logStepPass("PrecompiledTest 0x05 result is:" + DataChangeUtil.bytesToHex(callBigModExpValueByte));


            //Address 0x06: bn256Add(ax, ay, bx, by)(test pass)
            //输入参数规则：y^2 = x^3 + 3
            tx = precompiled.callBn256Add(new BigInteger("1"), new BigInteger("2"), new BigInteger("1"), new BigInteger("2")).send();
            collector.logStepPass("PrecompiledTest callBn256Add successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            List list = precompiled.getCallBn256AddValues().send();
            for (int i = 0; i < list.size(); i++) {
                collector.logStepPass("PrecompiledTest 0x06 result  " + i + " is:" + list.get(i).toString());
            }

            //Address 0x07: bn256ScalarMul(x, y, scalar)
            byte[] pointXByte = DataChangeUtil.hexToByteArray("2bd3e6d0f3b142924f5ca7b49ce5b9d54c4703d7ae5648e61d02268b1a0a9fb7");
            byte[] pointYByte = DataChangeUtil.hexToByteArray("21611ce0a6af85915e2f1d70300909ce2e49dfad4a4619c8390cae66cefdb204");
            byte[] scalarByte = DataChangeUtil.hexToByteArray("00000000000000000000000000000000000000000000000011138ce750fa15c2");
            tx = precompiled.callBn256ScalarMul(pointXByte, pointYByte, scalarByte).send();
            collector.logStepPass("PrecompiledTest callBigModExp successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            List callBn256ScalarMulList = precompiled.getCallBn256ScalarMulValues().send();
            for (int i = 0; i < callBn256ScalarMulList.size(); i++) {
                collector.logStepPass("PrecompiledTest call 0x07 result  " + i + " is:" + DataChangeUtil.bytesToHex((byte[]) callBn256ScalarMulList.get(i)));

            }

            //Address 0x08:callBn256Pairing
            String hexStr = "1c76476f4def4bb94541d57ebba1193381ffa7aa76ada664dd31c16024c43f593034dd2920f673e204fee2811c678745fc819b55d3e9d294e45c9b03a76aef41209dd15ebff5d46c4bd888e51a93cf99a7329636c63514396b4a452003a35bf704bf11ca01483bfa8b34b43561848d28905960114c8ac04049af4b6315a416782bb8324af6cfc93537a2ad1a445cfd0ca2a71acd7ac41fadbf933c2a51be344d120a2a4cf30c1bf9845f20c6fe39e07ea2cce61f0c9bb048165fe5e4de877550111e129f1cf1097710d41c4ac70fcdfa5ba2023c6ff1cbeac322de49d1b6df7c2032c61a830e3c17286de9462bf242fca2883585b93870a73853face6a6bf411198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa";
            tx = precompiled.callBn256Pairing(DataChangeUtil.hexToByteArray(hexStr)).send();
            collector.logStepPass("PrecompiledTest callBn256Pairing successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            byte[] callBn256PairingValue = precompiled.getCallBn256PairingValue().send();
            collector.logStepPass("PrecompiledTest 0x08 result is:" + DataChangeUtil.bytesToHex(callBn256PairingValue));


        } catch (Exception e) {
            collector.logStepFail("PrecompiledTest Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}


