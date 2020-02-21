package wasm.contract_cross_call;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractCallPrecompile;
import org.junit.Test;
import wasm.beforetest.WASMContractPrepareTest;

public class ContractCrossCallPrecompileContractsTest extends WASMContractPrepareTest {


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "xujiacan", showName = "wasm.contract_cross_call_precompile",sourcePrefix = "wasm")
    public void testCrossCallPreCompile() {
        try {

            prepare();

            // 测试跨合约调 ecrecover 预编译合约
            //
            // uint256[4] memory input;
            // input[0] = uint256(msgh);
            // input[1] = v;
            // input[2] = uint256(r);
            // input[3] = uint256(s);
            //
            // dataHash: "0xe281eaa11e6e37e6f53aade5d6c5b7201ef1c66162ec42ccc3215a0c4349350d", this hash is not txHash
            //V = 27
            //R = "0x55b60cadd4b4a3ea4fc368ef338f97e12e7328dd6e9e969a3fd8e5c10be855fe"
            //S = "0x2b42cee2585a16ea537efcb88009c1aeac693c28b59aa6bbff0baf22730338f6"
            //address: "0x8a9B36694F1eeeb500c84A19bB34137B05162EC5"
            String input = "0xe281eaa11e6e37e6f53aade5d6c5b7201ef1c66162ec42ccc3215a0c4349350d000000000000000000000000000000000000000000000000000000000000001b55b60cadd4b4a3ea4fc368ef338f97e12e7328dd6e9e969a3fd8e5c10be855fe2b42cee2585a16ea537efcb88009c1aeac693c28b59aa6bbff0baf22730338f6";

            ContractCallPrecompile precompile =  ContractCallPrecompile.deploy(web3j, transactionManager, provider).send();
            String addr =  precompile.cross_call_ecrecover(input, Uint64.of(0), Uint64.of(60000000l)).send();
            collector.logStepPass("cross_call_origin_type cross_call_ecrecover successfully addr:" + addr);

            // 测试跨合约调 sha256hash 预编译合约
            String sha3Str = "0x414243"; // hex(ABC)
            String sha3ExpectHash = "b5d4045c3f466fa91fe2cc6abe79232a1a57cdf104f7a26e716e0a1e2789df78";
            String sha3Hash = precompile.cross_call_sha256hash(sha3Str,Uint64.of(0), Uint64.of(60000000l)).send();
            collector.logStepPass("cross_call_origin_type cross_call_sha256hash successfully sha3Hash:" + sha3Hash);
            collector.assertEqual(sha3Hash, sha3ExpectHash);

            // 测试跨合约调 ripemd160hash 预编译合约
            String ripemd160Str = "0x414243"; // hex(ABC)
            String ripemd160ExpectHash = "000000000000000000000000df62d400e51d3582d53c2d89cfeb6e10d32a3ca6"; // 这一点注意, sol中返回的是被经处理成: df62d400e51d3582d53c2d89cfeb6e10d32a3ca6000000000000000000000000
            String ripemd160Hash = precompile.cross_call_ripemd160hash(ripemd160Str,Uint64.of(0), Uint64.of(60000000l)).send();
            collector.logStepPass("cross_call_origin_type cross_call_ripemd160hash successfully ripemd160Hash:" + ripemd160Hash);
            collector.assertEqual(ripemd160Hash, ripemd160ExpectHash);

            // 测试跨合约调 dataCopy 预编译合约
            String dataCopyStr = "414243"; // hex(ABC)
            String dataCopyHash = precompile.cross_call_dataCopy(sha3Str, Uint64.of(0), Uint64.of(60000000l)).send();
            collector.logStepPass("cross_call_origin_type cross_call_dataCopy successfully dataCopyHash:" + dataCopyHash);
            collector.assertEqual(dataCopyHash, dataCopyStr);

            // 测试跨合约调 bigModExp 预编译合约
            //
            // 入参的input拼接规则, input = [32]byte(baseLen) + [32]byte(expLen) + [32]byte(modLen) + [baseLen]byte(base) + [expLen]byte(exp) + [modLen]byte32(mod)
            // 其中 base 为基数， exp 为指数， mod 为模数
            // 求出 z, 其中 (公式: z = 基数 ** 指数 mod | 模数 |)
            // 如: z = 32 ** 3 mod | 5 |; z = 3
            //
            // baseLen = "0000000000000000000000000000000000000000000000000000000000000020"
            // expLen = "0000000000000000000000000000000000000000000000000000000000000020"
            // modLen = "0000000000000000000000000000000000000000000000000000000000000020"
            // base = "0000000000000000000000000000000000000000000000000000000000000020"
            // exp = "0000000000000000000000000000000000000000000000000000000000000003"
            // mod = "0000000000000000000000000000000000000000000000000000000000000005"
            // 求得的 z 应该为: 0000000000000000000000000000000000000000000000000000000000000003
            //
            String bigModExpStr = "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000005";
            String zExpectHash = "0000000000000000000000000000000000000000000000000000000000000003";
            String zHash = precompile.cross_call_bigModExp(bigModExpStr,Uint64.of(0), Uint64.of(60000000l)).send();
            collector.logStepPass("cross_call_origin_type cross_call_bigModExp successfully zHash:" + zHash);
            collector.assertEqual(zHash, zExpectHash);

        } catch (Exception e) {
            collector.logStepFail("Failed to call cross_call_origin_type Contract,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }


}
