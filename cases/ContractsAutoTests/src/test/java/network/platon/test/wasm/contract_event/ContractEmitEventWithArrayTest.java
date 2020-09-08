package network.platon.test.wasm.contract_event;

import com.platon.rlp.datatypes.Int8;
import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractEmitEventWithArray;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

import java.util.ArrayList;
import java.util.List;

/**
 * @title PLATON_EVENT 合约优化事件测试
 * http://192.168.9.66/PlatONContract/PlatONContract/blob/develop/system-design/event%E4%BC%98%E5%8C%96%E6%96%B9%E6%A1%88.md
 * @description:
 * @author: hudenian
 * @create: 2020/07/28
 */
public class ContractEmitEventWithArrayTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_event字节列表类型测试包含零个topic", sourcePrefix = "wasm")
    public void testZeroTopic() {

        try {
            prepare();
            ContractEmitEventWithArray contractEmitEventWithAddr = ContractEmitEventWithArray.deploy(web3j, transactionManager, provider, chainId).send();

            Int8[] _argsOne = new Int8[3];
            _argsOne[0] = Int8.of(0);
            _argsOne[1] = Int8.of(1);
            _argsOne[2] = Int8.of(2);


            byte[] _argsTwo = "abc".getBytes();

            byte[] _argsThree = "abcdefgh".getBytes();

            Int8[] _argsFour = new Int8[8];
            _argsFour[0] = Int8.of(3);
            _argsFour[1] = Int8.of(4);
            _argsFour[2] = Int8.of(5);
            _argsFour[4] = Int8.of(6);
            _argsFour[5] = Int8.of(7);
            _argsFour[6] = Int8.of(8);
            _argsFour[7] = Int8.of(9);

            Uint8[] _argsFive = new Uint8[8];
            _argsFive[0] = Uint8.of(3);
            _argsFive[1] = Uint8.of(4);
            _argsFive[2] = Uint8.of(5);
            _argsFive[3] = Uint8.of(6);
            _argsFive[4] = Uint8.of(7);
            _argsFive[5] = Uint8.of(8);
            _argsFive[6] = Uint8.of(9);
            _argsFive[7] = Uint8.of(10);

            List<Uint8> _argsSix = new ArrayList<>();
            _argsSix.add(Uint8.of(6));
            _argsSix.add(Uint8.of(7));
            _argsSix.add(Uint8.of(8));

            List<com.platon.rlp.datatypes.Int8> _argsSeven = new ArrayList<>();
            _argsSeven.add(Int8.of(9));
            _argsSeven.add(Int8.of(10));
            _argsSeven.add(Int8.of(11));

            String _argsEight = "lastArg";

            TransactionReceipt transactionReceipt = contractEmitEventWithAddr.zerotopic_eigthargs_event(_argsOne, _argsTwo, _argsThree, _argsFour, _argsFive, _argsSix, _argsSeven, _argsEight).send();
            List<ContractEmitEventWithArray.Transfer0EventResponse> transfer0EventResponseList = contractEmitEventWithAddr.getTransfer0Events(transactionReceipt);

            collector.assertEqual(contractEmitEventWithAddr.get_string().send(), _argsEight);
            collector.assertEqual(transfer0EventResponseList.get(0).arg1[0].value, _argsOne[0].value);
            collector.assertEqual(transfer0EventResponseList.get(0).arg2[0], _argsTwo[0]);
            collector.assertEqual(transfer0EventResponseList.get(0).arg3[0], _argsThree[0]);
            collector.assertEqual(transfer0EventResponseList.get(0).arg4[0].value, _argsFour[0].value);
            collector.assertEqual(transfer0EventResponseList.get(0).arg5[0].value, _argsFive[0].value);
            collector.assertEqual(transfer0EventResponseList.get(0).arg6.get(0), _argsSix.get(0));
            collector.assertEqual(transfer0EventResponseList.get(0).arg7.get(0), _argsSeven.get(0));
            collector.assertEqual(transfer0EventResponseList.get(0).arg8, _argsEight);

        } catch (Exception e) {
            collector.logStepFail("ContractEmitEventWithArrayTest call testZeroTopic failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_event字节列表类型测试包含1个topic", sourcePrefix = "wasm")
    public void testOneTopic() {

        try {
            prepare();
            ContractEmitEventWithArray contractEmitEventWithAddr = ContractEmitEventWithArray.deploy(web3j, transactionManager, provider, chainId).send();

            List<com.platon.rlp.datatypes.Int8> _topicOne = new ArrayList<>();
            _topicOne.add(Int8.of(9));
            _topicOne.add(Int8.of(10));
            _topicOne.add(Int8.of(11));

            Int8[] _argsOne = new Int8[3];
            _argsOne[0] = Int8.of(0);
            _argsOne[1] = Int8.of(1);
            _argsOne[2] = Int8.of(2);

            byte[] _argsTwo = "abc".getBytes();

            byte[] _argsThree = "abcdefgh".getBytes();

            Int8[] _argsFour = new Int8[8];
            _argsFour[0] = Int8.of(3);
            _argsFour[1] = Int8.of(4);
            _argsFour[2] = Int8.of(5);
            _argsFour[4] = Int8.of(6);
            _argsFour[5] = Int8.of(7);
            _argsFour[6] = Int8.of(8);
            _argsFour[7] = Int8.of(9);

            Uint8[] _argsFive = new Uint8[8];
            _argsFive[0] = Uint8.of(3);
            _argsFive[1] = Uint8.of(4);
            _argsFive[2] = Uint8.of(5);
            _argsFive[3] = Uint8.of(6);
            _argsFive[4] = Uint8.of(7);
            _argsFive[5] = Uint8.of(8);
            _argsFive[6] = Uint8.of(9);
            _argsFive[7] = Uint8.of(10);

            List<Uint8> _argsSix = new ArrayList<>();
            _argsSix.add(Uint8.of(6));
            _argsSix.add(Uint8.of(7));
            _argsSix.add(Uint8.of(8));

            String _argsSeven = "lastArg1";

            TransactionReceipt transactionReceipt = contractEmitEventWithAddr.onetopic_sevenargs_event(_topicOne, _argsOne, _argsTwo, _argsThree, _argsFour, _argsFive, _argsSix, _argsSeven).send();
            List<ContractEmitEventWithArray.Transfer1EventResponse> transfer0EventResponseList = contractEmitEventWithAddr.getTransfer1Events(transactionReceipt);

            collector.assertEqual(contractEmitEventWithAddr.get_string().send(), _argsSeven);
            collector.assertEqual(transfer0EventResponseList.get(0).arg1[0].value, _argsOne[0].value);
            collector.assertEqual(transfer0EventResponseList.get(0).arg2[0], _argsTwo[0]);
            collector.assertEqual(transfer0EventResponseList.get(0).arg3[0], _argsThree[0]);
            collector.assertEqual(transfer0EventResponseList.get(0).arg4[0].value, _argsFour[0].value);
            collector.assertEqual(transfer0EventResponseList.get(0).arg5[0].value, _argsFive[0].value);

        } catch (Exception e) {
            collector.logStepFail("ContractEmitEventWithArrayTest call testOneTopic failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_event字节列表类型测试包含2个topic", sourcePrefix = "wasm")
    public void testTwoTopic() {

        try {
            prepare();
            ContractEmitEventWithArray contractEmitEventWithAddr = ContractEmitEventWithArray.deploy(web3j, transactionManager, provider, chainId).send();

            Int8[] _topicOne = new Int8[8];
            _topicOne[0] = Int8.of(0);
            _topicOne[1] = Int8.of(1);
            _topicOne[2] = Int8.of(2);

            byte[] _topicTwo = "abc".getBytes();

            byte[] _argsOne = "abcdefgh".getBytes();

            Int8[] _argsTwo = new Int8[8];
            _argsTwo[0] = Int8.of(3);
            _argsTwo[1] = Int8.of(4);
            _argsTwo[2] = Int8.of(5);
            _argsTwo[4] = Int8.of(6);
            _argsTwo[5] = Int8.of(7);
            _argsTwo[6] = Int8.of(8);
            _argsTwo[7] = Int8.of(9);

            Uint8[] _argsThree = new Uint8[8];
            _argsThree[0] = Uint8.of(3);
            _argsThree[1] = Uint8.of(4);
            _argsThree[2] = Uint8.of(5);
            _argsThree[3] = Uint8.of(6);
            _argsThree[4] = Uint8.of(7);
            _argsThree[5] = Uint8.of(8);
            _argsThree[6] = Uint8.of(9);
            _argsThree[7] = Uint8.of(10);

            List<com.platon.rlp.datatypes.Int8> _argsFour = new ArrayList<>();
            _argsFour.add(Int8.of(9));
            _argsFour.add(Int8.of(10));
            _argsFour.add(Int8.of(11));

            List<Uint8> _argsFive = new ArrayList<>();
            _argsFive.add(Uint8.of(6));
            _argsFive.add(Uint8.of(7));
            _argsFive.add(Uint8.of(8));

            String _argsSix = "lastArg2";

            TransactionReceipt transactionReceipt = contractEmitEventWithAddr.twotopic_sixargs_event(_topicOne, _topicTwo, _argsOne, _argsTwo, _argsThree, _argsFour, _argsFive, _argsSix).send();
            List<ContractEmitEventWithArray.Transfer2EventResponse> transfer0EventResponseList = contractEmitEventWithAddr.getTransfer2Events(transactionReceipt);

            collector.assertEqual(contractEmitEventWithAddr.get_string().send(), _argsSix);
            collector.assertEqual(transfer0EventResponseList.get(0).arg1[0], _argsOne[0]);
            collector.assertEqual(transfer0EventResponseList.get(0).arg2[0], _argsTwo[0]);
            collector.assertEqual(transfer0EventResponseList.get(0).arg3[0], _argsThree[0]);
            collector.assertEqual(transfer0EventResponseList.get(0).arg4.get(0), _argsFour.get(0));
            collector.assertEqual(transfer0EventResponseList.get(0).arg5.get(0), _argsFive.get(0));
            collector.assertEqual(transfer0EventResponseList.get(0).arg6, _argsSix);

        } catch (Exception e) {
            collector.logStepFail("ContractEmitEventWithArrayTest call testTwoTopic failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_event字节列表类型测试包含3个topic", sourcePrefix = "wasm")
    public void testThreeTopic() {

        try {
            prepare();
            ContractEmitEventWithArray contractEmitEventWithAddr = ContractEmitEventWithArray.deploy(web3j, transactionManager, provider, chainId).send();

            byte[] _topicOne = "abcdefgh".getBytes();

            Int8[] _topicTwo = new Int8[8];
            _topicTwo[0] = Int8.of(3);
            _topicTwo[1] = Int8.of(4);
            _topicTwo[2] = Int8.of(5);
            _topicTwo[4] = Int8.of(6);
            _topicTwo[5] = Int8.of(7);
            _topicTwo[6] = Int8.of(8);
            _topicTwo[7] = Int8.of(9);

            Uint8[] _topicThree = new Uint8[8];
            _topicThree[0] = Uint8.of(3);
            _topicThree[1] = Uint8.of(4);
            _topicThree[2] = Uint8.of(5);
            _topicThree[3] = Uint8.of(6);
            _topicThree[4] = Uint8.of(7);
            _topicThree[5] = Uint8.of(8);
            _topicThree[6] = Uint8.of(9);
            _topicThree[7] = Uint8.of(10);

            Int8[] _argsOne = new Int8[8];
            _argsOne[0] = Int8.of(3);
            _argsOne[1] = Int8.of(4);
            _argsOne[2] = Int8.of(5);
            _argsOne[4] = Int8.of(6);
            _argsOne[5] = Int8.of(7);
            _argsOne[6] = Int8.of(8);
            _argsOne[7] = Int8.of(9);


            byte[] _argsTwo = "def".getBytes();

            List<com.platon.rlp.datatypes.Int8> _argsThree = new ArrayList<>();
            _argsThree.add(Int8.of(9));
            _argsThree.add(Int8.of(10));
            _argsThree.add(Int8.of(11));

            List<String> _argsFour = new ArrayList<>();
            _argsFour.add("a");
            _argsFour.add("b");
            _argsFour.add("c");

            String _argsFive = "lastArg3";

            TransactionReceipt transactionReceipt = contractEmitEventWithAddr.threetopic_fiveargs_event(_topicOne, _topicTwo, _topicThree, _argsOne, _argsTwo, _argsThree, _argsFour, _argsFive).send();
            List<ContractEmitEventWithArray.Transfer3EventResponse> transfer3EventResponseList = contractEmitEventWithAddr.getTransfer3Events(transactionReceipt);

            collector.assertEqual(contractEmitEventWithAddr.get_string().send(), _argsFive);
            collector.assertEqual(transfer3EventResponseList.get(0).arg1[0], _argsOne[0]);
            collector.assertEqual(transfer3EventResponseList.get(0).arg2[0], _argsTwo[0]);
            collector.assertEqual(transfer3EventResponseList.get(0).arg3.get(0), _argsThree.get(0));
            collector.assertEqual(transfer3EventResponseList.get(0).arg4.get(0), _argsFour.get(0));
            collector.assertEqual(transfer3EventResponseList.get(0).arg5, _argsFive);

        } catch (Exception e) {
            collector.logStepFail("ContractEmitEventWithArrayTest call testThreeTopic failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }
}
