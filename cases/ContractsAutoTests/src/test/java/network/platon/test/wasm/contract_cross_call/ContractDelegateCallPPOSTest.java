package wasm.contract_cross_call;

import com.google.gson.Gson;
import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractDelegateCallPPOS;
import network.platon.utils.DataChangeUtil;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

public class ContractDelegateCallPPOSTest extends WASMContractPrepareTest {

    // 锁仓合约
    private String restrictingContractAddr = "lax1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqp3yp7hw";

    private String stakingContractAddr = "lax1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzlh5ge3";

    private String slashingContractAddr = "lax1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqyrchd9x";

    private String govContractAddr = "lax1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq97wrcc5";

    private String delegateRewardPoolAddr = "lax1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqxsakwkt";


    private Gson gson = new Gson();

    // {"Code":305001,"Ret": xxx}
    class pposResult {
        public  int Code;
        public  Object Ret;
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "xujiacan", showName = "wasm.contract_delegate_call_ppos_RestrictingPlan",sourcePrefix = "wasm")
    public void testDelegateCallPPOS4Restricting() {
        try {

            prepare();

            // 测试跨合约调 ppos合约


            // 部署
            ContractDelegateCallPPOS ppos =  ContractDelegateCallPPOS.deploy(web3j, transactionManager, provider, chainId).send();
            collector.logStepPass("gas used after deploy delegate_call_ppos:" + ppos.getTransactionReceipt().get().getGasUsed());

            /**
             *  锁仓
             *
             *  account: 0xc9E1C2B330Cf7e759F2493c5C754b34d98B07f93
             *  plans: [{"epoch":1,"amount":1000000000000000000},{"epoch":2,"amount":1000000000000000000},{"epoch":3,"amount":1000000000000000000},{"epoch":4,"amount":1000000000000000000},{"epoch":5,"amount":1000000000000000000}]
             *
             */
            String createRestrictingPlanInput = "0xf85483820fa09594c9e1c2b330cf7e759f2493c5c754b34d98b07f93b838f7ca01880de0b6b3a7640000ca02880de0b6b3a7640000ca03880de0b6b3a7640000ca04880de0b6b3a7640000ca05880de0b6b3a7640000";

            TransactionReceipt createRestrictingPlanReceipt =  ppos.delegate_call_ppos_send(restrictingContractAddr, createRestrictingPlanInput, Uint64.of(60000000l)).send();

            String  createRestrictingPlanDataHex = createRestrictingPlanReceipt.getLogs().get(0).getData();
            String createRestrictingPlanDataStr = DataChangeUtil.decodeSystemContractRlp(createRestrictingPlanDataHex, chainId);
            String createRestrictingPlanExpectData = "0";

            collector.logStepPass("delegate_call_ppos createRestrictingPlan successfully txHash:" + createRestrictingPlanReceipt.getTransactionHash());
            collector.assertEqual(createRestrictingPlanDataStr, createRestrictingPlanExpectData);


            /**
             *  查询 账户的锁仓计划
             *
             *  account： 0xc9E1C2B330Cf7e759F2493c5C754b34d98B07f93
             *
             *
             */

            String getRestrictingInfoInput = "0xda838210049594c9e1c2b330cf7e759f2493c5c754b34d98b07f93";
            String getRestrictingInfoHexStr =  ppos.delegate_call_ppos_query(restrictingContractAddr, getRestrictingInfoInput, Uint64.of(60000000l)).send();
            byte[] getRestrictingInfoByte =  DataChangeUtil.hexToByteArray(getRestrictingInfoHexStr);
            String getRestrictingInfoStr = new String(getRestrictingInfoByte);
            collector.logStepPass("获取锁仓计划:" + getRestrictingInfoStr);
            ContractCrossCallPPOSTest.pposResult res =  gson.fromJson(getRestrictingInfoStr, ContractCrossCallPPOSTest.pposResult.class);
            collector.assertEqual(res.Code, 0, "查询账户的锁仓计划 result == expect res: {\"Code\":0,\"Ret\": xxxx }");


        } catch (Exception e) {
            collector.logStepFail("Failed to delegateCall delegate_call_ppos Contract,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "xujiacan", showName = "wasm.contract_delegate_call_ppos_Staking",sourcePrefix = "wasm")
    public void testDelegateCallPPOS4Staking() {
        try {

            prepare();

            // 测试跨合约调 ppos合约


            // 部署
            ContractDelegateCallPPOS ppos =  ContractDelegateCallPPOS.deploy(web3j, transactionManager, provider, chainId).send();
            collector.logStepPass("gas used after deploy delegate_call_ppos:" + ppos.getTransactionReceipt().get().getGasUsed());

            /**
             *  质押
             *
             * typ: 0
             * benefitAddress: 0xd87E10F8efd2C32f5e88b7C279953aEF6EE58902
             * nodeId: ced880d4769331f47af07a8d1b79de1e40c95a37ea1890bb9d3f0da8349e1a7c0ea4cadbb9c5bf185b051061eef8e5eadca251c24e1db1d9faf0fb24cbd06f9a
             * externalId: xssssddddffffggggg
             * nodeName: Gavin, China
             * website: https://www.Gavin.network
             * details: Gavin super node
             * amount: 1000000000000000000000000
             * rewardPer: 5000
             * programVersion: 65536
             * programVersionSign: 0xa6b8f5e2de519991b567b7c9af9cdf6adb331be5c4b79d26ba0dcd1fb2fbab7a276a36944765997f4d18d6135b515d9de7cceb2d7af4338f4fedb6acdc050f8001
             * blsPubKey: e5c3fcd6ce33c06aae22113977396c295728b8c01e0bc9188d2f3ffe52ea97b465639731f3cd4956a26e3e35f96e2a10646f28a352c6453be54dda05b703c31f0fbda3abc55a75151788338917f5a60b26f92bd15cdaf0dc00779a62056f3a00
             * blsProof: 46e9c92915ad2b423e9eea33482d2615f6a17a15a6fa3b99e3e83bc394700d14c7ace638b34be70e6903724ca217b3cf4ff85db6f38e83f06de95de2c0370916
             */
            String createStakingInput = "0xf901b1838203e881809594d87e10f8efd2c32f5e88b7c279953aef6ee58902b842b840ced880d4769331f47af07a8d1b79de1e40c95a37ea1890bb9d3f0da8349e1a7c0ea4cadbb9c5bf185b051061eef8e5eadca251c24e1db1d9faf0fb24cbd06f9a93927873737373646464646666666667676767678d8c476176696e2c204368696e619a9968747470733a2f2f7777772e476176696e2e6e6574776f726b9190476176696e207375706572206e6f64658b8ad3c21bcecceda1000000838213888483010000b843b841a6b8f5e2de519991b567b7c9af9cdf6adb331be5c4b79d26ba0dcd1fb2fbab7a276a36944765997f4d18d6135b515d9de7cceb2d7af4338f4fedb6acdc050f8001b862b860e5c3fcd6ce33c06aae22113977396c295728b8c01e0bc9188d2f3ffe52ea97b465639731f3cd4956a26e3e35f96e2a10646f28a352c6453be54dda05b703c31f0fbda3abc55a75151788338917f5a60b26f92bd15cdaf0dc00779a62056f3a00b842b84046e9c92915ad2b423e9eea33482d2615f6a17a15a6fa3b99e3e83bc394700d14c7ace638b34be70e6903724ca217b3cf4ff85db6f38e83f06de95de2c0370916";

            TransactionReceipt createStakingReceipt =  ppos.delegate_call_ppos_send(stakingContractAddr, createStakingInput, Uint64.of(60000000l)).send();

            String  createStakingDataHex = createStakingReceipt.getLogs().get(0).getData();
            String createStakingDataStr = DataChangeUtil.decodeSystemContractRlp(createStakingDataHex, chainId);
            String createStakingExpectData = "301005";

            collector.logStepPass("delegate_call_ppos createStaking successfully txHash:" + createStakingReceipt.getTransactionHash());
            collector.assertEqual(createStakingDataStr, createStakingExpectData);


            /**
             *  查询 候选人详情
             *
             *  nodeId： ced880d4769331f47af07a8d1b79de1e40c95a37ea1890bb9d3f0da8349e1a7c0ea4cadbb9c5bf185b051061eef8e5eadca251c24e1db1d9faf0fb24cbd06f9a
             *
             *
             */

            String getCandidateInfoInput = "0xf84883820451b842b840ced880d4769331f47af07a8d1b79de1e40c95a37ea1890bb9d3f0da8349e1a7c0ea4cadbb9c5bf185b051061eef8e5eadca251c24e1db1d9faf0fb24cbd06f9a";
            String getCandidateInfoHexStr =  ppos.delegate_call_ppos_query(stakingContractAddr, getCandidateInfoInput, Uint64.of(60000000l)).send();
            byte[] getCandidateInfoByte =  DataChangeUtil.hexToByteArray(getCandidateInfoHexStr);
            String getCandidateInfoStr = new String(getCandidateInfoByte);
            ContractCrossCallPPOSTest.pposResult res =  gson.fromJson(getCandidateInfoStr, ContractCrossCallPPOSTest.pposResult.class);
            collector.assertEqual(res.Code, 301204, "查询候选人详情 result == expect res: {\"Code\":301204,\"Ret\":\"Query candidate info failed:Candidate info is not found\"}");

        } catch (Exception e) {
            collector.logStepFail("Failed to delegateCall delegate_call_ppos Contract,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }



    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "xujiacan", showName = "wasm.contract_delegate_call_ppos_Slashing",sourcePrefix = "wasm")
    public void testDelegateCallPPOS4Slashing() {
        try {

            prepare();

            // 测试跨合约调 ppos合约


            // 部署
            ContractDelegateCallPPOS ppos =  ContractDelegateCallPPOS.deploy(web3j, transactionManager, provider, chainId).send();
            collector.logStepPass("gas used after deploy delegate_call_ppos:" + ppos.getTransactionReceipt().get().getGasUsed());

            /**
             *  举报
             *
             *  dupType： 1
             *  data： {
             *           "prepareA": {
             *            "epoch": 1,
             *            "viewNumber": 1,
             *            "blockHash": "0xbb6d4b83af8667929a9cb4918bbf790a97bb136775353765388d0add3437cde6",
             *            "blockNumber": 1,
             *            "blockIndex": 1,
             *            "blockData": "0x45b20c5ba595be254943aa57cc80562e84f1fb3bafbf4a414e30570c93a39579",
             *            "validateNode": {
             *             "index": 0,
             *             "address": "0x195667cdefcad94c521bdff0bf85079761e0f8f3",
             *             "nodeId": "51c0559c065400151377d71acd7a17282a7c8abcfefdb11992dcecafde15e100b8e31e1a5e74834a04792d016f166c80b9923423fe280570e8131debf591d483",
             *             "blsPubKey": "752fe419bbdc2d2222009e450f2932657bbc2370028d396ba556a49439fe1cc11903354dcb6dac552a124e0b3db0d90edcd334d7aabda0c3f1ade12ca22372f876212ac456d549dbbd04d2c8c8fb3e33760215e114b4d60313c142f7b8bbfd87"
             *            },
             *            "signature": "0x36015fee15253487e8125b86505377d8540b1a95d1a6b13f714baa55b12bd06ec7d5755a98230cdc88858470afa8cb0000000000000000000000000000000000"
             *           },
             *           "prepareB": {
             *            "epoch": 1,
             *            "viewNumber": 1,
             *            "blockHash": "0xf46c45f7ebb4a999dd030b9f799198b785654293dbe41aa7e909223af0e8c4ba",
             *            "blockNumber": 1,
             *            "blockIndex": 1,
             *            "blockData": "0xd630e96d127f55319392f20d4fd917e3e7cba19ad366c031b9dff05e056d9420",
             *            "validateNode": {
             *             "index": 0,
             *             "address": "0x195667cdefcad94c521bdff0bf85079761e0f8f3",
             *             "nodeId": "51c0559c065400151377d71acd7a17282a7c8abcfefdb11992dcecafde15e100b8e31e1a5e74834a04792d016f166c80b9923423fe280570e8131debf591d483",
             *             "blsPubKey": "752fe419bbdc2d2222009e450f2932657bbc2370028d396ba556a49439fe1cc11903354dcb6dac552a124e0b3db0d90edcd334d7aabda0c3f1ade12ca22372f876212ac456d549dbbd04d2c8c8fb3e33760215e114b4d60313c142f7b8bbfd87"
             *            },
             *            "signature": "0x783892b9b766f9f4c2a1d45b1fd53ca9ea56a82e38a998939edc17bc7fd756267d3c145c03bc6c1412302cf590645d8200000000000000000000000000000000"
             *           }
             *          }
             */
            String reportDuplicateSignInput = "0xf907e683820bb801b907deb907db7b0a20202020202020202020227072657061726541223a207b0a20202020202020202020202265706f6368223a20312c0a202020202020202020202022766965774e756d626572223a20312c0a202020202020202020202022626c6f636b48617368223a2022307862623664346238336166383636373932396139636234393138626266373930613937626231333637373533353337363533383864306164643334333763646536222c0a202020202020202020202022626c6f636b4e756d626572223a20312c0a202020202020202020202022626c6f636b496e646578223a20312c0a202020202020202020202022626c6f636b44617461223a2022307834356232306335626135393562653235343934336161353763633830353632653834663166623362616662663461343134653330353730633933613339353739222c0a20202020202020202020202276616c69646174654e6f6465223a207b0a20202020202020202020202022696e646578223a20302c0a2020202020202020202020202261646472657373223a2022307831393536363763646566636164393463353231626466663062663835303739373631653066386633222c0a202020202020202020202020226e6f64654964223a20223531633035353963303635343030313531333737643731616364376131373238326137633861626366656664623131393932646365636166646531356531303062386533316531613565373438333461303437393264303136663136366338306239393233343233666532383035373065383133316465626635393164343833222c0a20202020202020202020202022626c735075624b6579223a2022373532666534313962626463326432323232303039653435306632393332363537626263323337303032386433393662613535366134393433396665316363313139303333353464636236646163353532613132346530623364623064393065646364333334643761616264613063336631616465313263613232333732663837363231326163343536643534396462626430346432633863386662336533333736303231356531313462346436303331336331343266376238626266643837220a20202020202020202020207d2c0a2020202020202020202020227369676e6174757265223a202230783336303135666565313532353334383765383132356238363530353337376438353430623161393564316136623133663731346261613535623132626430366563376435373535613938323330636463383838353834373061666138636230303030303030303030303030303030303030303030303030303030303030303030220a202020202020202020207d2c0a20202020202020202020227072657061726542223a207b0a20202020202020202020202265706f6368223a20312c0a202020202020202020202022766965774e756d626572223a20312c0a202020202020202020202022626c6f636b48617368223a2022307866343663343566376562623461393939646430333062396637393931393862373835363534323933646265343161613765393039323233616630653863346261222c0a202020202020202020202022626c6f636b4e756d626572223a20312c0a202020202020202020202022626c6f636b496e646578223a20312c0a202020202020202020202022626c6f636b44617461223a2022307864363330653936643132376635353331393339326632306434666439313765336537636261313961643336366330333162396466663035653035366439343230222c0a20202020202020202020202276616c69646174654e6f6465223a207b0a20202020202020202020202022696e646578223a20302c0a2020202020202020202020202261646472657373223a2022307831393536363763646566636164393463353231626466663062663835303739373631653066386633222c0a202020202020202020202020226e6f64654964223a20223531633035353963303635343030313531333737643731616364376131373238326137633861626366656664623131393932646365636166646531356531303062386533316531613565373438333461303437393264303136663136366338306239393233343233666532383035373065383133316465626635393164343833222c0a20202020202020202020202022626c735075624b6579223a2022373532666534313962626463326432323232303039653435306632393332363537626263323337303032386433393662613535366134393433396665316363313139303333353464636236646163353532613132346530623364623064393065646364333334643761616264613063336631616465313263613232333732663837363231326163343536643534396462626430346432633863386662336533333736303231356531313462346436303331336331343266376238626266643837220a20202020202020202020207d2c0a2020202020202020202020227369676e6174757265223a202230783738333839326239623736366639663463326131643435623166643533636139656135366138326533386139393839333965646331376263376664373536323637643363313435633033626336633134313233303263663539303634356438323030303030303030303030303030303030303030303030303030303030303030220a202020202020202020207d0a2020202020202020207d";

            TransactionReceipt reportDuplicateSignReceipt =  ppos.delegate_call_ppos_send(slashingContractAddr, reportDuplicateSignInput, Uint64.of(60000000l)).send();

            String reportDuplicateSignDataHex = reportDuplicateSignReceipt.getLogs().get(0).getData();
            String reportDuplicateSignDataStr = DataChangeUtil.decodeSystemContractRlp(reportDuplicateSignDataHex, chainId);
            String reportDuplicateSignExpectData = "0";
            boolean actual = reportDuplicateSignDataStr.equals(reportDuplicateSignExpectData);

            collector.logStepPass("delegate_call_ppos reportDuplicateSign successfully txHash:" + reportDuplicateSignReceipt.getTransactionHash());
            collector.assertEqual(actual, false);


            /**
             *  查询 节点是否有多签过
             *
             * dupType: 1
             * addr: 0x9e3e0f0f366b26b965f3aa3ed67603fb480b1257
             * blockNumber: 1
             *
             * 0
             */

            String checkDuplicateSignInput = "0xf183820bb801abaa6c6178316e636c713772656b64766e746a65306e34676c64766173726c6479716b796a686b7667727172";
            String checkDuplicateSignHexStr =  ppos.delegate_call_ppos_query(slashingContractAddr, checkDuplicateSignInput, Uint64.of(60000000l)).send();
            byte[] checkDuplicateSignByte =  DataChangeUtil.hexToByteArray(checkDuplicateSignHexStr);
            String checkDuplicateSignStr = new String(checkDuplicateSignByte);
            ContractCrossCallPPOSTest.pposResult res =  gson.fromJson(checkDuplicateSignStr, ContractCrossCallPPOSTest.pposResult.class);
            if(res != null){
                collector.assertEqual(res.Code, 0, "查询节点是否有多签过 result == expect res: {\"Code\":0,\"Ret\":\"\"}");
            }

        } catch (Exception e) {
            collector.logStepFail("Failed to delegateCall delegate_call_ppos Contract,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "xujiacan", showName = "wasm.contract_delegate_call_ppos_Governance",sourcePrefix = "wasm")
    public void testDelegateCallPPOS4Governance() {
        try {

            prepare();

            // 测试跨合约调 ppos合约


            // 部署
            ContractDelegateCallPPOS ppos =  ContractDelegateCallPPOS.deploy(web3j, transactionManager, provider, chainId).send();
            collector.logStepPass("gas used after deploy delegate_call_ppos:" + ppos.getTransactionReceipt().get().getGasUsed());

            /**
             *  提交文本提案
             *
             *  verifierID: ced880d4769331f47af07a8d1b79de1e40c95a37ea1890bb9d3f0da8349e1a7c0ea4cadbb9c5bf185b051061eef8e5eadca251c24e1db1d9faf0fb24cbd06f9a
             *  PIPID: textUrl
             */
            String submitTextInput = "0xf851838207d0b842b840ced880d4769331f47af07a8d1b79de1e40c95a37ea1890bb9d3f0da8349e1a7c0ea4cadbb9c5bf185b051061eef8e5eadca251c24e1db1d9faf0fb24cbd06f9a88877465787455726c";

            TransactionReceipt submitTextReceipt =  ppos.delegate_call_ppos_send(govContractAddr, submitTextInput, Uint64.of(60000000l)).send();

            String submitTextDataHex = submitTextReceipt.getLogs().get(0).getData();
            String submitTextDataStr =DataChangeUtil.decodeSystemContractRlp(submitTextDataHex, chainId);
            String submitTextExpectData = "302022";

            collector.logStepPass("delegate_call_ppos submitText successfully txHash:" + submitTextReceipt.getTransactionHash());
            collector.assertEqual(submitTextDataStr, submitTextExpectData);


            /**
             *  查询 提案
             *
             * proposalID: 0x12c171900f010b17e969702efa044d077e86808212c171900f010b17e969702e
             *
             *
             */

            String getProposalInput = "0xe683820834a1a012c171900f010b17e969702efa044d077e86808212c171900f010b17e969702e";
            String getProposalHexStr =  ppos.delegate_call_ppos_query(govContractAddr, getProposalInput, Uint64.of(60000000l)).send();
            byte[] getProposalByte =  DataChangeUtil.hexToByteArray(getProposalHexStr);
            String getProposalStr = new String(getProposalByte);
            ContractCrossCallPPOSTest.pposResult res =  gson.fromJson(getProposalStr, ContractCrossCallPPOSTest.pposResult.class);
            collector.assertEqual(res.Code, 302006, "查询提案 result == expect res: {\"Code\":302006,\"Ret\":\"proposal not found\"}");

        } catch (Exception e) {
            collector.logStepFail("Failed to delegateCall delegate_call_ppos Contract,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }



    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "xujiacan", showName = "wasm.contract_delegate_call_ppos_DelegateReward",sourcePrefix = "wasm")
    public void testDelegateCallPPOS4DelegateReward() {
        try {

            prepare();

            // 测试跨合约调 ppos合约


            // 部署
            ContractDelegateCallPPOS ppos =  ContractDelegateCallPPOS.deploy(web3j, transactionManager, provider, chainId).send();
            collector.logStepPass("gas used after deploy delegate_call_ppos:" + ppos.getTransactionReceipt().get().getGasUsed());

            /**
             *  提取委托奖励
             *
             */
            String withdrawDelegateRewardInput = "0xc483821388";

            TransactionReceipt withdrawDelegateRewardReceipt =  ppos.delegate_call_ppos_send(delegateRewardPoolAddr, withdrawDelegateRewardInput, Uint64.of(60000000l)).send();

            String withdrawDelegateRewardDataHex = withdrawDelegateRewardReceipt.getLogs().get(0).getData();
            String withdrawDelegateRewardDataStr = DataChangeUtil.decodeSystemContractRlp(withdrawDelegateRewardDataHex, chainId);
            String withdrawDelegateRewardExpectData = "305001";

            collector.logStepPass("delegate_call_ppos withdrawDelegateReward successfully txHash:" + withdrawDelegateRewardReceipt.getTransactionHash());
            collector.assertEqual(withdrawDelegateRewardDataStr, withdrawDelegateRewardExpectData);


            /**
             *  查询 账户在各节点未提取委托奖励
             *
             " Addr": "0x12c171900f010b17e969702efa044d077e868082"
             *
             * NodeIDs": [
             * "db18af9be2af9dff2347c3d06db4b1bada0598d099a210275251b68fa7b5a863d47fcdd382cc4b3ea01e5b55e9dd0bdbce654133b7f58928ce74629d5e68b974",
             * "1f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee2840e429"]
             *
             *
             *
             */
            String getDelegateRewardInput = "0xf8a2838213ec959412c171900f010b17e969702efa044d077e868082b886f884b840db18af9be2af9dff2347c3d06db4b1bada0598d099a210275251b68fa7b5a863d47fcdd382cc4b3ea01e5b55e9dd0bdbce654133b7f58928ce74629d5e68b974b8401f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee2840e429";
            String getDelegateRewardHexStr =  ppos.delegate_call_ppos_query(delegateRewardPoolAddr, getDelegateRewardInput, Uint64.of(60000000l)).send();
            byte[] getDelegateRewardByte =  DataChangeUtil.hexToByteArray(getDelegateRewardHexStr);
            String getDelegateRewardStr = new String(getDelegateRewardByte);
            ContractCrossCallPPOSTest.pposResult res =  gson.fromJson(getDelegateRewardStr, ContractCrossCallPPOSTest.pposResult.class);
            collector.assertEqual(res.Code, 305001, "查询候选人详情 result == expect res: {\"Code\":305001,\"Ret\":\"delegation info not found\"}");

        } catch (Exception e) {
            collector.logStepFail("Failed to delegateCall delegate_call_ppos Contract,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }


}
