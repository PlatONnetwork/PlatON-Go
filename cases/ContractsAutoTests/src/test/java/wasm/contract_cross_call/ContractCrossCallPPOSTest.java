//package wasm.contract_cross_call;
//
//import com.platon.rlp.datatypes.Uint64;
//import network.platon.autotest.junit.annotations.DataSource;
//import network.platon.autotest.junit.enums.DataSourceType;
//import network.platon.contracts.wasm.ContractCallPPOS;
//import network.platon.utils.DataChangeUtil;
//import org.junit.Test;
//import org.web3j.protocol.core.methods.response.TransactionReceipt;
//import wasm.beforetest.WASMContractPrepareTest;
//
//public class ContractCrossCallPPOSTest extends WASMContractPrepareTest {
//
//    // 锁仓合约
//    private String restrictingContractAddr = "0x1000000000000000000000000000000000000001";
//
//    private String stakingContractAddr = "0x1000000000000000000000000000000000000002";
//
//    private String slashingContractAddr = "0x1000000000000000000000000000000000000004";
//
//    private String govContractAddr = "0x1000000000000000000000000000000000000005";
//
//    private String DelegateRewardPoolAddr = "0x1000000000000000000000000000000000000006";
//
//    @Test
//    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
//            author = "xujiacan", showName = "wasm.contract_cross_call_ppos",sourcePrefix = "wasm")
//    public void testCrossCallPPOS4Restricting() {
//        try {
//
//            prepare();
//
//            // 测试跨合约调 ppos合约
//
//
//            // 部署
//            ContractCallPPOS ppos =  ContractCallPPOS.deploy(web3j, transactionManager, provider).send();
//
//
//            /**
//             *  锁仓
//             *
//             *  account: 0xc9E1C2B330Cf7e759F2493c5C754b34d98B07f93
//             *  plans: [{"epoch":1,"amount":1000000000000000000},{"epoch":2,"amount":1000000000000000000},{"epoch":3,"amount":1000000000000000000},{"epoch":4,"amount":1000000000000000000},{"epoch":5,"amount":1000000000000000000}]
//             *
//              */
//            String createRestrictingPlanInput = "0xf85483820fa09594c9e1c2b330cf7e759f2493c5c754b34d98b07f93b838f7ca01880de0b6b3a7640000ca02880de0b6b3a7640000ca03880de0b6b3a7640000ca04880de0b6b3a7640000ca05880de0b6b3a7640000";
//
//            TransactionReceipt createRestrictingPlanReceipt =  ppos.cross_call_ppos_send(restrictingContractAddr, createRestrictingPlanInput, Uint64.of(0), Uint64.of(60000000l)).send();
//
//            String  createRestrictingPlanDataHex = createRestrictingPlanReceipt.getLogs().get(0).getData();
//            byte[] createRestrictingPlanDataByte = DataChangeUtil.hexToByteArray(createRestrictingPlanDataHex);
//            String createRestrictingPlanDataStr = new String(createRestrictingPlanDataByte);
//            String createRestrictingPlanExpectData = "304004";
//
//            collector.logStepPass("cross_call_ppos createRestrictingPlan successfully data:" + createRestrictingPlanDataStr);
//            collector.assertEqual(createRestrictingPlanDataStr, createRestrictingPlanExpectData);
//
//
////            /**
////             *  查询 候选人详情
////             *
////             *  nodeId： ced880d4769331f47af07a8d1b79de1e40c95a37ea1890bb9d3f0da8349e1a7c0ea4cadbb9c5bf185b051061eef8e5eadca251c24e1db1d9faf0fb24cbd06f9a
////             */
////
////            String getCandidateInfoInput = "0xf84883820451b842b840ced880d4769331f47af07a8d1b79de1e40c95a37ea1890bb9d3f0da8349e1a7c0ea4cadbb9c5bf185b051061eef8e5eadca251c24e1db1d9faf0fb24cbd06f9a";
////            String getCandidateInfoHexStr =  ppos.cross_call_ppos_query(stakingContractAddr, getCandidateInfoInput, Uint64.of(0), Uint64.of(60000000l)).send();
////            byte[] getCandidateInfoByte =  DataChangeUtil.hexToByteArray(getCandidateInfoHexStr);
////            String getCandidateInfoStr = new String(getCandidateInfoByte);
////            collector.logStepPass("Str:" + getCandidateInfoStr);
//
//
//        } catch (Exception e) {
//            collector.logStepFail("Failed to call cross_call_origin_type Contract,exception msg:" , e.getMessage());
//            e.printStackTrace();
//        }
//    }
//
//
//    @Test
//    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
//            author = "xujiacan", showName = "wasm.contract_cross_call_ppos",sourcePrefix = "wasm")
//    public void testCrossCallPPOS4Staking() {
//        try {
//
//            prepare();
//
//            // 测试跨合约调 ppos合约
//
//
//            // 部署
//            ContractCallPPOS ppos =  ContractCallPPOS.deploy(web3j, transactionManager, provider).send();
//
//
//            /**
//             *  质押
//             *
//             * typ: 0
//             * benefitAddress: 0xd87E10F8efd2C32f5e88b7C279953aEF6EE58902
//             * nodeId: ced880d4769331f47af07a8d1b79de1e40c95a37ea1890bb9d3f0da8349e1a7c0ea4cadbb9c5bf185b051061eef8e5eadca251c24e1db1d9faf0fb24cbd06f9a
//             * externalId: xssssddddffffggggg
//             * nodeName: Gavin, China
//             * website: https://www.Gavin.network
//             * details: Gavin super node
//             * amount: 1000000000000000000000000
//             * rewardPer: 5000
//             * programVersion: 65536
//             * programVersionSign: 0xa6b8f5e2de519991b567b7c9af9cdf6adb331be5c4b79d26ba0dcd1fb2fbab7a276a36944765997f4d18d6135b515d9de7cceb2d7af4338f4fedb6acdc050f8001
//             * blsPubKey: e5c3fcd6ce33c06aae22113977396c295728b8c01e0bc9188d2f3ffe52ea97b465639731f3cd4956a26e3e35f96e2a10646f28a352c6453be54dda05b703c31f0fbda3abc55a75151788338917f5a60b26f92bd15cdaf0dc00779a62056f3a00
//             * blsProof: 46e9c92915ad2b423e9eea33482d2615f6a17a15a6fa3b99e3e83bc394700d14c7ace638b34be70e6903724ca217b3cf4ff85db6f38e83f06de95de2c0370916
//             */
//            String createStakingInput = "0xf901b1838203e881809594d87e10f8efd2c32f5e88b7c279953aef6ee58902b842b840ced880d4769331f47af07a8d1b79de1e40c95a37ea1890bb9d3f0da8349e1a7c0ea4cadbb9c5bf185b051061eef8e5eadca251c24e1db1d9faf0fb24cbd06f9a93927873737373646464646666666667676767678d8c476176696e2c204368696e619a9968747470733a2f2f7777772e476176696e2e6e6574776f726b9190476176696e207375706572206e6f64658b8ad3c21bcecceda1000000838213888483010000b843b841a6b8f5e2de519991b567b7c9af9cdf6adb331be5c4b79d26ba0dcd1fb2fbab7a276a36944765997f4d18d6135b515d9de7cceb2d7af4338f4fedb6acdc050f8001b862b860e5c3fcd6ce33c06aae22113977396c295728b8c01e0bc9188d2f3ffe52ea97b465639731f3cd4956a26e3e35f96e2a10646f28a352c6453be54dda05b703c31f0fbda3abc55a75151788338917f5a60b26f92bd15cdaf0dc00779a62056f3a00b842b84046e9c92915ad2b423e9eea33482d2615f6a17a15a6fa3b99e3e83bc394700d14c7ace638b34be70e6903724ca217b3cf4ff85db6f38e83f06de95de2c0370916";
//
//            TransactionReceipt createStakingReceipt =  ppos.cross_call_ppos_send(stakingContractAddr, createStakingInput, Uint64.of(0), Uint64.of(60000000l)).send();
//
//            String  createStakingDataHex = createStakingReceipt.getLogs().get(0).getData();
//            byte[] createStakingDataByte = DataChangeUtil.hexToByteArray(createStakingDataHex);
//            String createStakingDataStr = new String(createStakingDataByte);
//            String createStakingExpectData = "301111";
//
//            collector.logStepPass("cross_call_ppos createStaking successfully data:" + createStakingDataStr);
//            collector.assertEqual(createStakingDataStr, createStakingExpectData);
//
//
//            /**
//             *  查询 候选人详情
//             *
//             *  nodeId： ced880d4769331f47af07a8d1b79de1e40c95a37ea1890bb9d3f0da8349e1a7c0ea4cadbb9c5bf185b051061eef8e5eadca251c24e1db1d9faf0fb24cbd06f9a
//             */
//
//            String getCandidateInfoInput = "0xf84883820451b842b840ced880d4769331f47af07a8d1b79de1e40c95a37ea1890bb9d3f0da8349e1a7c0ea4cadbb9c5bf185b051061eef8e5eadca251c24e1db1d9faf0fb24cbd06f9a";
//            String getCandidateInfoHexStr =  ppos.cross_call_ppos_query(stakingContractAddr, getCandidateInfoInput, Uint64.of(0), Uint64.of(60000000l)).send();
//            byte[] getCandidateInfoByte =  DataChangeUtil.hexToByteArray(getCandidateInfoHexStr);
//            String getCandidateInfoStr = new String(getCandidateInfoByte);
//            collector.logStepPass("Str:" + getCandidateInfoStr);
//
//
//        } catch (Exception e) {
//            collector.logStepFail("Failed to call cross_call_origin_type Contract,exception msg:" , e.getMessage());
//            e.printStackTrace();
//        }
//    }
//
//
//}
