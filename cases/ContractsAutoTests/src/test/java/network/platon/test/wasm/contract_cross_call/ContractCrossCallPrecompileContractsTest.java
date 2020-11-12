package network.platon.test.wasm.contract_cross_call;

import com.platon.rlp.datatypes.Uint64;
import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractCallPrecompile;
import network.platon.utils.DataChangeUtil;
import org.junit.Test;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

public class ContractCrossCallPrecompileContractsTest extends WASMContractPrepareTest {


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "xujiacan", showName = "wasm.contract_cross_call_precompile",sourcePrefix = "wasm")
    public void testCrossCallPreCompile() {
        try {

            prepare();

            // 测试跨合约调 ecrecover 预编译合约
            //
            // 入参的规则为: bytes32(dataHash) + bytes32(v) + bytes32(r) + bytes32(s)
            //
            // dataHash: "0xe281eaa11e6e37e6f53aade5d6c5b7201ef1c66162ec42ccc3215a0c4349350d", this hash is not txHash
            // V = 27
            // R = "0x55b60cadd4b4a3ea4fc368ef338f97e12e7328dd6e9e969a3fd8e5c10be855fe"
            // S = "0x2b42cee2585a16ea537efcb88009c1aeac693c28b59aa6bbff0baf22730338f6"
            //
            // 求出最终的hex应为 0000000000000000000000008a9b36694f1eeeb500c84a19bb34137b05162ec5
            // 对应的address: "0x8a9B36694F1eeeb500c84A19bB34137B05162EC5"
            String msghStr = "0xe281eaa11e6e37e6f53aade5d6c5b7201ef1c66162ec42ccc3215a0c4349350d";

            byte[] msgh = DataChangeUtil.hexToByteArray(msghStr);


            Uint8 v = Uint8.of("27");

            String rStr = "0x55b60cadd4b4a3ea4fc368ef338f97e12e7328dd6e9e969a3fd8e5c10be855fe";
            byte[] r = DataChangeUtil.hexToByteArray(rStr);

            String sStr = "0x2b42cee2585a16ea537efcb88009c1aeac693c28b59aa6bbff0baf22730338f6";
            byte[] s = DataChangeUtil.hexToByteArray(sStr);

            String expectAddrStr = "0000000000000000000000008a9b36694f1eeeb500c84a19bb34137b05162ec5";


            ContractCallPrecompile precompile =  ContractCallPrecompile.deploy(web3j, transactionManager, provider, chainId).send();
            collector.logStepPass("gas used after deploy cross_call_precompile contract:" + precompile.getTransactionReceipt().get().getGasUsed());
            collector.logStepPass("cross_call_precompile deployed sucessfully, contractAddress:" + precompile.getContractAddress() + ", txHash:" + precompile.getTransactionReceipt().get().getTransactionHash());

            String addr =  precompile.cross_call_ecrecover(msgh, v, r, s, Uint64.of(0), Uint64.of(60000000l)).send();

            collector.logStepPass("cross_call_precompile cross_call_ecrecover successfully addr:" + addr);
            collector.assertEqual(addr, expectAddrStr);

            // 测试跨合约调 sha256hash 预编译合约
            String sha3Str = "0x414243"; // hex(ABC)
            String sha3ExpectHash = "b5d4045c3f466fa91fe2cc6abe79232a1a57cdf104f7a26e716e0a1e2789df78";
            String sha3Hash = precompile.cross_call_sha256hash(sha3Str,Uint64.of(0), Uint64.of(60000000l)).send();
            collector.logStepPass("cross_call_precompile cross_call_sha256hash successfully sha3Hash:" + sha3Hash);
            collector.assertEqual(sha3Hash, sha3ExpectHash);

            // 测试跨合约调 ripemd160hash 预编译合约
            String ripemd160Str = "0x414243"; // hex(ABC)
            String ripemd160ExpectHash = "000000000000000000000000df62d400e51d3582d53c2d89cfeb6e10d32a3ca6"; // 这一点注意, sol中返回的是被经处理成: df62d400e51d3582d53c2d89cfeb6e10d32a3ca6000000000000000000000000
            String ripemd160Hash = precompile.cross_call_ripemd160hash(ripemd160Str,Uint64.of(0), Uint64.of(60000000l)).send();
            collector.logStepPass("cross_call_precompile cross_call_ripemd160hash successfully ripemd160Hash:" + ripemd160Hash);
            collector.assertEqual(ripemd160Hash, ripemd160ExpectHash);

            // 测试跨合约调 dataCopy 预编译合约
            String dataCopyStr = "414243"; // hex(ABC)
            String dataCopyHash = precompile.cross_call_dataCopy(sha3Str, Uint64.of(0), Uint64.of(60000000l)).send();
            collector.logStepPass("cross_call_precompile cross_call_dataCopy successfully dataCopyHash:" + dataCopyHash);
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
            String zExpectHash = "0000000000000000000000000000000000000000000000000000000000000003";

            String baseStr = "0000000000000000000000000000000000000000000000000000000000000020";
            byte[] base = DataChangeUtil.hexToByteArray(baseStr);

            String expStr = "0000000000000000000000000000000000000000000000000000000000000003";
            byte[] exp = DataChangeUtil.hexToByteArray(expStr);

            String modStr = "0000000000000000000000000000000000000000000000000000000000000005";
            byte[] mod = DataChangeUtil.hexToByteArray(modStr);


            String zHash = precompile.cross_call_bigModExp(base, exp, mod,Uint64.of(0), Uint64.of(60000000l)).send();
            collector.logStepPass("cross_call_precompile cross_call_bigModExp successfully zHash:" + zHash);
            collector.assertEqual(zHash, zExpectHash);


            // 测试跨合约调 bn256Add 预编译合约
            //
            // 主要是对 bn256 做 + 操作
            // 其中入参的  ax, ay, bx, by 其中 (ax, ay) 为椭圆曲线 bn256 上的一个点A， (bx, by) 是椭圆曲线上的另外一个点B
            // 该函数就是 求两个点的 G点
            //
            // input = ax + ay + bx + by
            //
            // 对于点A坐标取值:
            // ax:  10744596414106452074759370245733544594153395043370666422502510773307029471145 其十六进制为: 0x17c139df0efee0f766bc0204762b774362e4ded88953a39ce849a8a7fa163fa9
            // ay:  848677436511517736191562425154572367705380862894644942948681172815252343932  其十六进制为: 0x01e0559bacb160664764a357af8a9fe70baa9258e0b959273ffc5718c6d4cc7c
            //
            // 对于点B坐标取值:
            // bx: 1624070059937464756887933993293429854168590106605707304006200119738501412969  其十六进制为: 0x039730ea8dff1254c0fee9c0ea777d29a9c710b7e616683f194f18c43b43b869
            // by: 3269329550605213075043232856820720631601935657990457502777101397807070461336  其十六进制为: 0x073a5ffcc6fc7a28c30723d6e58ce577356982d65b833a5a5c15bf9024b43d98
            //
            // 则，内置合约 bn256Add 的input入参为: ax+ay+bx+by = 0x17c139df0efee0f766bc0204762b774362e4ded88953a39ce849a8a7fa163fa901e0559bacb160664764a357af8a9fe70baa9258e0b959273ffc5718c6d4cc7c039730ea8dff1254c0fee9c0ea777d29a9c710b7e616683f194f18c43b43b869073a5ffcc6fc7a28c30723d6e58ce577356982d65b833a5a5c15bf9024b43d98
            String cExpect = "15bf2bb17880144b5d1cd2b1f46eff9d617bffd1ca57c37fb5a49bd84e53cf66049c797f9ce0d17083deb32b5e36f2ea2a212ee036598dd7624c168993d1355f";


            String axStr = "0x17c139df0efee0f766bc0204762b774362e4ded88953a39ce849a8a7fa163fa9";
            byte[] ax = DataChangeUtil.hexToByteArray(axStr);

            String ayStr = "0x01e0559bacb160664764a357af8a9fe70baa9258e0b959273ffc5718c6d4cc7c";
            byte[] ay = DataChangeUtil.hexToByteArray(ayStr);

            String bxStr = "0x039730ea8dff1254c0fee9c0ea777d29a9c710b7e616683f194f18c43b43b869";
            byte[] bx = DataChangeUtil.hexToByteArray(bxStr);

            String byStr = "0x073a5ffcc6fc7a28c30723d6e58ce577356982d65b833a5a5c15bf9024b43d98";
            byte[] by = DataChangeUtil.hexToByteArray(byStr);

            String cCoordinate = precompile.cross_call_bn256Add(ax, ay, bx, by, Uint64.of(0), Uint64.of(60000000l)).send();
            collector.logStepPass("cross_call_precompile cross_call_bn256Add successfully c ponit coordinate:" + cCoordinate);
            collector.assertEqual(cCoordinate, cExpect);



            // 测试跨合约调 bn256ScalarMul 预编译合约
            //
            // 主要是对 bn256 做 × 操作
            // 其中入参的 ax, ay, 和 系数 scalar 其中 (ax, ay) 为椭圆曲线 bn256 上的一个点A, scalar 为椭圆曲线的质数域中的 N×G， 即 A×NG == B 得到椭圆曲线上的另外一个点B的坐标(bx, by)
            //
            // input = ax + ay + scalar
            //
            // 对于点A坐标取值:
            // ax: 19823850254741169819033785099293761935467223354323761392354670518001715552183  其十六进制为: 0x2bd3e6d0f3b142924f5ca7b49ce5b9d54c4703d7ae5648e61d02268b1a0a9fb7
            // ay: 15097907474011103550430959168661954736283086276546887690628027914974507414020   其十六进制为: 0x21611ce0a6af85915e2f1d70300909ce2e49dfad4a4619c8390cae66cefdb204
            //
            // 系数 scalar 取值: 1230482048326178242  其十六进制为: 0x00000000000000000000000000000000000000000000000011138ce750fa15c2
            //
            // 我们将会得到B的坐标为:  (返回的B坐标的 bx+by的十六进制为: 0x070a8d6a982153cae4be29d434e8faef8a47b274a053f5a4ee2a6c9c13c31e5c031b8ce914eba3a9ffb989f9cdd5b0f01943074bf4f0f315690ec3cec6981afc)
            // bx: 3184834430741071145030522771540763108892281233703148152311693391954704539228  其十六进制为: 0x070a8d6a982153cae4be29d434e8faef8a47b274a053f5a4ee2a6c9c13c31e5c
            // by: 1405615944858121891163559530323310827496899969303520166098610312148921359100  其十六进制为: 0x031b8ce914eba3a9ffb989f9cdd5b0f01943074bf4f0f315690ec3cec6981afc
            //
            // 则，内置合约 bn256ScalarMul 的input入参为: ax+ay+scalar = 0x2bd3e6d0f3b142924f5ca7b49ce5b9d54c4703d7ae5648e61d02268b1a0a9fb721611ce0a6af85915e2f1d70300909ce2e49dfad4a4619c8390cae66cefdb20400000000000000000000000000000000000000000000000011138ce750fa15c2
            String bExpect = "070a8d6a982153cae4be29d434e8faef8a47b274a053f5a4ee2a6c9c13c31e5c031b8ce914eba3a9ffb989f9cdd5b0f01943074bf4f0f315690ec3cec6981afc";

            String xStr = "0x2bd3e6d0f3b142924f5ca7b49ce5b9d54c4703d7ae5648e61d02268b1a0a9fb7";
            byte[] x = DataChangeUtil.hexToByteArray(xStr);

            String yStr = "0x21611ce0a6af85915e2f1d70300909ce2e49dfad4a4619c8390cae66cefdb204";
            byte[] y = DataChangeUtil.hexToByteArray(yStr);

            String scalarStr = "0x00000000000000000000000000000000000000000000000011138ce750fa15c2";
            byte[] scalar = DataChangeUtil.hexToByteArray(scalarStr);

            String bCoordinate = precompile.cross_call_bn256ScalarMul(x, y, scalar, Uint64.of(0), Uint64.of(60000000l)).send();
            collector.logStepPass("cross_call_precompile cross_call_bn256ScalarMul successfully b ponit coordinate:" + bCoordinate);
            collector.assertEqual(bCoordinate, bExpect);


            // 测试跨合约调 bn256ScalarMul 预编译合约
            //
            // 主要是对 bn256 做 pairing 操作, 也就是 配对操作
            // 通过没两组 (x, y) 分别代表 去 曲线点和 扭曲点 来一一匹配，看看是否同属于 一条曲线上的
            //
            //入参数据如下:
            //
            // x1: 12873740738727497448187997291915224677121726020054032516825496230827252793177   其对应的十六进制为:  0x1c76476f4def4bb94541d57ebba1193381ffa7aa76ada664dd31c16024c43f59
            // y1: 21804419174137094775122804775419507726154084057848719988004616848382402162497   其对应的十六进制为:  0x3034dd2920f673e204fee2811c678745fc819b55d3e9d294e45c9b03a76aef41
            // x2: 14752851163271972921165116810778899752274893127848647655434033030151679466487   其对应的十六进制为:  0x209dd15ebff5d46c4bd888e51a93cf99a7329636c63514396b4a452003a35bf7
            // y2: 2146841959437886920191033516947821737903543682424168472444605468016078231160   其对应的十六进制为:  0x04bf11ca01483bfa8b34b43561848d28905960114c8ac04049af4b6315a41678
            // x3: 19774899457345372253936887903062884289284519982717033379297427576421785416781   其对应的十六进制为:  0x2bb8324af6cfc93537a2ad1a445cfd0ca2a71acd7ac41fadbf933c2a51be344d
            // y3: 8159591693044959083845993640644415462154314071906244874217244895511876957520   其对应的十六进制为:  0x120a2a4cf30c1bf9845f20c6fe39e07ea2cce61f0c9bb048165fe5e4de877550
            // x4: 7742452358972543465462254569134860944739929848367563713587808717088650354556   其对应的十六进制为:  0x111e129f1cf1097710d41c4ac70fcdfa5ba2023c6ff1cbeac322de49d1b6df7c
            // y4: 14563720768440487558151020426243236708567496944263114635856508834497000371217   其对应的十六进制为:  0x2032c61a830e3c17286de9462bf242fca2883585b93870a73853face6a6bf411
            // x5: 11559732032986387107991004021392285783925812861821192530917403151452391805634   其对应的十六进制为:  0x198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c2
            // y5: 10857046999023057135944570762232829481370756359578518086990519993285655852781   其对应的十六进制为:  0x1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed
            // x6: 4082367875863433681332203403145435568316851327593401208105741076214120093531   其对应的十六进制为:  0x090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b
            //
            // 规则 (x1, y1) 和 (x2, y2) 做匹配 以此类推
            //

            String pairingStr = "0x1c76476f4def4bb94541d57ebba1193381ffa7aa76ada664dd31c16024c43f593034dd2920f673e204fee2811c678745fc819b55d3e9d294e45c9b03a76aef41209dd15ebff5d46c4bd888e51a93cf99a7329636c63514396b4a452003a35bf704bf11ca01483bfa8b34b43561848d28905960114c8ac04049af4b6315a416782bb8324af6cfc93537a2ad1a445cfd0ca2a71acd7ac41fadbf933c2a51be344d120a2a4cf30c1bf9845f20c6fe39e07ea2cce61f0c9bb048165fe5e4de877550111e129f1cf1097710d41c4ac70fcdfa5ba2023c6ff1cbeac322de49d1b6df7c2032c61a830e3c17286de9462bf242fca2883585b93870a73853face6a6bf411198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa";
            String expextFlag = "0000000000000000000000000000000000000000000000000000000000000001"; // 1 表示全部配对成功, 0 表示有配对失败的存在
            String retFlag = precompile.cross_call_bn256Pairing(pairingStr,Uint64.of(0), Uint64.of(60000000l)).send();
            collector.logStepPass("cross_call_precompile cross_call_bn256Pairing successfully ret flag:" + retFlag);
            collector.assertEqual(retFlag, expextFlag);





        } catch (Exception e) {
            collector.logStepFail("Failed to call cross_call_origin_type Contract,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }


}
