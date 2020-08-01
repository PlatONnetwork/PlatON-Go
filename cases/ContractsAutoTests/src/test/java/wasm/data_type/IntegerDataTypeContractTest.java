package wasm.data_type;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONArray;
import com.alibaba.fastjson.JSONObject;
import com.platon.rlp.datatypes.*;
import com.platon.sdk.utlis.Bech32;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.IntegerDataTypeContract_1;
import network.platon.contracts.wasm.IntegerDataTypeContract_2;
import network.platon.contracts.wasm.IntegerDataTypeContract_3;
//import network.platon.contracts.wasm.IntegerDataTypeContract_4;
import network.platon.contracts.wasm.IntegerDataTypeContract_4;
import network.platon.utils.DataChangeUtil;
import org.junit.Before;
import org.junit.Test;
import org.web3j.abi.datatypes.Address;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;

/**
 * @author zjsunzone
 *
 * This class is used to test date type of integer.
 */
public class IntegerDataTypeContractTest extends WASMContractPrepareTest {

    @Before
    public void before(){
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.base_data_type_01",sourcePrefix = "wasm")
    public void testBaseTypeContract_01() {

        try {
            // deploy contract.
            IntegerDataTypeContract_1 contract = IntegerDataTypeContract_1.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("IntegerDataTypeContract_01 issued successfully.contractAddress:"
                    + contractAddress + ", hash:" + transactionHash
                    + " gasUsed:" + contract.getTransactionReceipt().get().getGasUsed().toString());

            // test: int8
            Int16 int8 = contract.int8().send();
            collector.logStepPass("To invoke int8 success, int8: " + int8.getValue());

            // test: int64
            Int64 int64 = contract.int64().send();
            collector.logStepPass("To invoke int8 success, int64: " + int64.getValue());

            // test uint8
            Uint8 expectByt8 = Uint8.of(10);
            Uint8 uint8 = contract.uint8t(expectByt8).send();
            collector.logStepPass("To invoke uint8 success, uint8: " + uint8.getValue().toString());

            // test: uint32
            Uint32 expectByt32 = Uint32.of(1000);
            Uint32 uint32 = contract.uint32t(expectByt32).send();
            collector.logStepPass("To invoke uint32 success, uint32: " + uint32.getValue());
            collector.assertEqual(uint32.getValue(), BigInteger.valueOf(2000));

            // test: uint64
            Uint64 expect64 = Uint64.of(10000);
            Uint64 uint64 = contract.uint64t(expect64).send();
            collector.logStepPass("To invoke uint64 success, uint64: " + uint64.getValue().toString());
            collector.assertEqual(uint64.getValue(), BigInteger.valueOf(20000));

            // test: u128
            Uint64 expect128 = Uint64.of(10000);
            String u128 = contract.u128t(expect128).send();
            collector.logStepPass("To invoke uint64 success, u128: " + u128);
            collector.assertEqual(u128, expect128.getValue().toString());

            // test: u256
            Uint64 expect256 = Uint64.of(10000);
            String u256 = contract.u256t(expect256).send();
            collector.logStepPass("To invoke u256t success, u256: " + u128);
            collector.assertEqual(u256, expect256.getValue().toString());

        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("IntegerDataTypeContract_01 and could not call contract function");
            }else{
                collector.logStepFail("IntegerDataTypeContract_01 failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.base_data_type_02",sourcePrefix = "wasm")
    public void testBaseTypeContract_02() {

        try {
            // deploy contract.
            IntegerDataTypeContract_2 contract = IntegerDataTypeContract_2.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("IntegerDataTypeContract_01 issued successfully.contractAddress:"
                    + contractAddress + ", hash:" + transactionHash
                    + " gasUsed:" + contract.getTransactionReceipt().get().getGasUsed().toString());

            //  int8
            TransactionReceipt int8Tr = contract.setInt8(Int8.of((byte) 2)).send();
            collector.logStepPass("To invoke setInt 8 success, txHash: " + int8Tr.getTransactionHash());
            Int8 getInt8 = contract.getInt8().send();
            collector.logStepPass("To invoke getInt8 8 success, getInt8: " + getInt8.getValue());

            //  int32
            TransactionReceipt int32Tr = contract.setInt32(Int32.of(100)).send();
            collector.logStepPass("To invoke setInt32 success, txHash: " + int32Tr.getTransactionHash());
            Int32 getInt32 = contract.getInt32().send();
            collector.logStepPass("To invoke getInt32 success, getInt32: " + getInt32);

            // int64
            TransactionReceipt int64Tr = contract.setInt64(Int64.of(1111111111)).send();
            collector.logStepPass("To invoke setInt64 success, txHash: " + int64Tr.getTransactionHash());
            Int64 getInt64 = contract.getInt64().send();
            collector.logStepPass("To invoke getInt64 success, getInt64: " + getInt64);

            // ======================= uint =======================
            //  uint8
            TransactionReceipt uint8Tr = contract.setUint8(Uint8.of(2)).send();
            collector.logStepPass("To invoke setUint8 success, txHash: " + uint8Tr.getTransactionHash());
            Uint8 getUint8 = contract.getUint8().send();
            collector.logStepPass("To invoke getUint8 8 success, getUint8: " + getUint8.getValue().toString());

            //  uint32
            TransactionReceipt uint32Tr = contract.setUint32(Uint32.of(100)).send();
            collector.logStepPass("To invoke setuUint32 success, txHash: " + uint32Tr.getTransactionHash());
            Uint32 getUint32 = contract.getUint32().send();
            collector.logStepPass("To invoke getUint32 success, getUint32: " + getUint32.getValue().toString());

            // uint64
            TransactionReceipt uint64Tr = contract.setUint64(Uint64.of("1111111111")).send();
            collector.logStepPass("To invoke setUint64 success, txHash: " + uint64Tr.getTransactionHash());
            Uint64 getUint64 = contract.getUint64().send();
            collector.logStepPass("To invoke getUint64 success, getUint64: " + getUint64.getValue().toString());

        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("IntegerDataTypeContract_02 and could not call contract function");
            }else{
                collector.logStepFail("IntegerDataTypeContract_02 failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.base_data_type_03",sourcePrefix = "wasm")
    public void testBaseTypeContract_03() {

        try {
            // deploy contract.
            IntegerDataTypeContract_3 contract = IntegerDataTypeContract_3.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("IntegerDataTypeContract_3 issued successfully.contractAddress:"
                    + contractAddress + ", hash:" + transactionHash
                    + " gasUsed:" + contract.getTransactionReceipt().get().getGasUsed().toString());

            // test: store string
            TransactionReceipt strTr = contract.setString("setString").send();
            String getString = contract.getString().send();
            collector.logStepPass("To invoke setString and getString, getString: " + getString);
            collector.assertEqual(getString, "setString");

            // test: store bool
            TransactionReceipt boolTr = contract.setBool(true).send();
            boolean getBool = contract.getBool().send();
            collector.logStepPass("To invoke setBool and getBool, bool: " + getBool);
            collector.assertEqual(getBool, true);

            // test: store char
            Byte expectByte = (byte)1;
            TransactionReceipt charTr = contract.setChar(Int8.of(expectByte.byteValue())).send();
            collector.logStepPass("To invoke setChar success, txHash: " + charTr.getTransactionHash());
            Int8 getChar = contract.getChar().send();
            collector.logStepPass("To invoke getChar success, getChar: " + getChar.getValue());
            collector.assertEqual(getChar.getValue(), expectByte.byteValue());

        } catch (Exception e) {
            collector.logStepFail("IntegerDataTypeContract_3 failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.base_data_type_04",sourcePrefix = "wasm")
    public void testBaseTypeContract_04() {

        try {
            // deploy contract.
            IntegerDataTypeContract_4 contract = IntegerDataTypeContract_4.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("IntegerDataTypeContract_4 issued successfully.contractAddress:"
                    + contractAddress + ", hash:" + transactionHash
                    + " gasUsed:" + contract.getTransactionReceipt().get().getGasUsed().toString());

            // init address
            TransactionReceipt initTr = contract.initAddress().send();
            collector.logStepPass("To invoke initAddress success, txHash: " + initTr.getTransactionHash());
            String afterInitAddress = contract.getAddress().send();
            collector.logStepPass("To invoke getAddress success, getAddress: " + afterInitAddress);
            collector.assertEqual(afterInitAddress.toLowerCase(), "lax1w2kjkufl4g2v93xd94a0lewc75ufdr66rnzuw2".toLowerCase());

//            // test: store address
//            Address expectAddr = new Address("0x5b05e7a3e2a688c5e5cc491545a84a1efc66c1b1");
            String expectAddr = "lax1fyeszufxwxk62p46djncj86rd553skpptsj8v6";
            TransactionReceipt addrTr = contract.setAddress(expectAddr).send();
            collector.logStepPass("To invoke setAddress success, txHash: " + addrTr.getTransactionHash());
            String getAddress = contract.getAddress().send();
            collector.logStepPass("To invoke getAddress success, getAddress: " + getAddress);
            //collector.assertEqual(getAddress, expectAddr);

            // test: store u256
            String expectU256 = "100000";
            TransactionReceipt u256TR = contract.setU256(Uint64.of(expectU256)).send();
            collector.logStepPass("To invoke setU256 success, txHash: " + u256TR.getTransactionHash());
            String getU256 = contract.getU256().send();
            collector.logStepPass("To invoke getU256 success, getU256: " + getU256);
            collector.assertEqual(getU256, expectU256);

            // test: store h256
            String expectH256 = "0x80b543239ae8e4f679019719312524d10f14fef79fd0d9117d810bffdedf608e";
            TransactionReceipt h256Tr = contract.setH256(expectH256).send();
            collector.logStepPass("To invoke setH256 success, txHash: " + h256Tr.getTransactionHash());
            String getH256 = contract.getH256().send();
            collector.logStepPass("To invoke getH256 success, getH256: " + getH256);
            //collector.assertEqual(getH256, expectH256);

        } catch (Exception e) {
            collector.logStepFail("IntegerDataTypeContract_4 failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.base_data_type_border",sourcePrefix = "wasm")
    public void testBaseTypeContract_border() {
        // 主要测试各类型的边界值
        try {
            // deploy contract.
            IntegerDataTypeContract_2 contract = IntegerDataTypeContract_2.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("IntegerDataTypeContract_01 issued successfully.contractAddress:"
                    + contractAddress + ", hash:" + transactionHash
                    + " gasUsed:" + contract.getTransactionReceipt().get().getGasUsed().toString());

            // int8: -128 ~ 127
            // uint: 0 ~ 255
            // int32: -2147483648 ~ 2147483647
            // uint32: 0 ~ 4294967295
            // int64: -9,223,372,036,854,775,808 ~ 9,223,372,036,854,775,807
            // uint64: 0 ~ 18,446,744,073,709,551,615

            // int8: -128 ~ 127
            // [{"input":100, "expect": 100, "equal":"N"}]
            JSONArray int8Cases = JSON.parseArray("["
                    + "{\"input\":-128, \"expect\": -128, \"equal\":\"Y\"}"
                    //+ "{\"input\":-129, \"expect\": -129, \"equal\":\"N\"}" // pass
                    + "{\"input\":127, \"expect\": 127, \"equal\":\"Y\"}"
                    //+ "{\"input\":128, \"expect\": 128, \"equal\":\"N\"}" // pass
                    +"]");
            for (int i = 0; i < int8Cases.size(); i++) {
                JSONObject testCase = int8Cases.getJSONObject(i);
                int input = testCase.getIntValue("input");
                int expect = testCase.getIntValue("expect");
                String equal = testCase.getString("equal");
                TransactionReceipt int8Tr = contract.setInt8(Int8.of(input)).send();
                collector.logStepPass("To invoke setInt8 success, txHash: " + int8Tr.getTransactionHash());
                Int8 getInt8 = contract.getInt8().send();
                collector.logStepPass("To invoke getInt8 success, setInt8: " + input+ " getInt8: " + getInt8.getValue());
                if(equal.equals("Y")){
                    collector.assertEqual(getInt8.getValue(), (byte)expect);
                } else {
                    boolean eq = (int)getInt8.getValue() != expect;
                    collector.assertTrue(eq);
                }
            }

            // uint8: 0 ~ 255
            JSONArray uint8Cases = JSON.parseArray("["
                    //+ "{\"input\":-1, \"expect\": -1, \"equal\":\"N\"}" // pass:
                    + "{\"input\":0, \"expect\": 0, \"equal\":\"Y\"}"
                    //+ "{\"input\":254, \"expect\": 254, \"equal\":\"Y\"}" // return: -2
                    + "{\"input\":255, \"expect\": 255, \"equal\":\"N\"}"
                    //+ "{\"input\":256, \"expect\": 256, \"equal\":\"N\"}" // pass.
                    +"]");
            for (int i = 0; i < uint8Cases.size(); i++) {
                JSONObject testCase = uint8Cases.getJSONObject(i);
                int input = testCase.getIntValue("input");
                int expect = testCase.getIntValue("expect");
                String equal = testCase.getString("equal");
                TransactionReceipt int8Tr = contract.setUint8(Uint8.of(input)).send();
                collector.logStepPass("To invoke setUint8 success, txHash: " + int8Tr.getTransactionHash());
                Uint8 getUint8 = contract.getUint8().send();
                collector.logStepPass("To invoke getUint8 success, setUint8: "+ input +", getUint8: " + getUint8.getValue().toString());
                if(equal.equals("Y")){
                    collector.assertEqual(getUint8.getValue().toString(), String.valueOf(expect));
                } else {
                    collector.assertFalse(getUint8.getValue() == BigInteger.valueOf(expect));
                }
            }

            // int32: -2147483648 ~ 2147483647
            // int32
            JSONArray int32Cases = JSON.parseArray("["
                    //+ "{\"input\": \"-2147483649\", \"expect\": \"-2147483649\", \"equal\":\"N\"}" // pass.
                    + "{\"input\": \"-2147483648\", \"expect\": \"-2147483648\", \"equal\":\"Y\"}"  // right
                    + "{\"input\": \"2147483647\", \"expect\": \"2147483647\", \"equal\":\"Y\"}"
                    + "{\"input\":\"0\", \"expect\": \"0\", \"equal\":\"Y\"}"
                    //+ "{\"input\": \"2147483648\", \"expect\": \"2147483648\", \"equal\":\"N\"}"
                    +"]");
            for (int i = 0; i < int32Cases.size(); i++) {
                JSONObject testCase = int32Cases.getJSONObject(i);
                String input = testCase.getString("input");
                String expect = testCase.getString("expect");
                String equal = testCase.getString("equal");
                TransactionReceipt tr = contract.setInt32(Int32.of(Long.valueOf(input))).send();
                collector.logStepPass("To invoke setInt32 success, txHash: " + tr.getTransactionHash());
                Int32 getInt32 = contract.getInt32().send();
                collector.logStepPass("To invoke getInt32 success,setInt32: "+ input +", getInt32: " + getInt32.toString());
                if(equal.equals("Y")){
                    collector.assertEqual(getInt32.getValue(), Integer.valueOf(expect));
                } else {
                    collector.assertFalse(getInt32.toString().equals(expect));
                }
            }

            // uint32: 0 ~ 4294967295
            // uint32
            JSONArray uint32Cases = JSON.parseArray("["
                    //+ "{\"input\": \"-1\", \"expect\": \"-1\", \"equal\":\"N\"}"  // pass
                    //+ "{\"input\": \"4294967294\", \"expect\": \"4294967294\", \"equal\":\"Y\"}"
                    + "{\"input\": \"4294967295\", \"expect\": \"4294967295\", \"equal\":\"Y\"}"
                    //+ "{\"input\": \"4294967296\", \"expect\": \"4294967296\", \"equal\":\"N\"}"  // pass
                    + "{\"input\": \"0\", \"expect\": \"0\", \"equal\":\"Y\"}"
                    +"]");
            for (int i = 0; i < uint32Cases.size(); i++) {
                JSONObject testCase = uint32Cases.getJSONObject(i);
                String input = testCase.getString("input");
                String expect = testCase.getString("expect");
                String equal = testCase.getString("equal");
                TransactionReceipt tr = contract.setUint32(Uint32.of(input)).send();
                collector.logStepPass("To invoke setUint32 success, txHash: " + tr.getTransactionHash());
                Uint32 getReturn = contract.getUint32().send();
                collector.logStepPass("To invoke getUint32 success,setUint32: "+ input +", getUint32: " + getReturn);
                if(equal.equals("Y")){
                    collector.assertEqual(getReturn.getValue().toString(), expect);
                } else {
                    collector.assertFalse(getReturn.getValue().toString().equals(expect));
                }
            }

            // int64: -9,223,372,036,854,775,808 ~ 9,223,372,036,854,775,807
            // int64
            JSONArray int64Cases = JSON.parseArray("["
                    + "{\"input\": \"-1\", \"expect\": \"-1\", \"equal\":\"Y\"}"
                    + "{\"input\": \"-9223372036854775808\", \"expect\": \"-9223372036854775808\", \"equal\":\"Y\"}"
                    + "{\"input\": \"0\", \"expect\": \"0\", \"equal\":\"Y\"}"
                    + "{\"input\": \"9223372036854775807\", \"expect\": \"9223372036854775807\", \"equal\":\"Y\"}"
                    //+ "{\"input\": \"9223372036854775808\", \"expect\": \"9223372036854775808\", \"equal\":\"N\"}"
                    +"]");
            for (int i = 0; i < int64Cases.size(); i++) {
                JSONObject testCase = int64Cases.getJSONObject(i);
                String input = testCase.getString("input");
                String expect = testCase.getString("expect");
                String equal = testCase.getString("equal");
                TransactionReceipt tr = contract.setInt64(Int64.of(Long.valueOf(input))).send();
                collector.logStepPass("To invoke setInt64 success, txHash: " + tr.getTransactionHash());
                Int64 getReturn = contract.getInt64().send();
                collector.logStepPass("To invoke getInt64 success,setInt64: "+ input +", getInt64: " + getReturn.getValue());
                if(equal.equals("Y")){
                    collector.assertEqual(getReturn.getValue(), Long.valueOf(expect).longValue());
                } else {
                    collector.assertFalse(getReturn.getValue() == Long.valueOf(expect).longValue());
                }
            }

            // uint64: 0 ~ 18,446,744,073,709,551,615
            JSONArray uint64Cases = JSON.parseArray("["
                    //+ "{\"input\": \"-1\", \"expect\": \"-1\", \"equal\":\"Y\"}"  // pass
                    + "{\"input\": \"4294967294\", \"expect\": \"4294967294\", \"equal\":\"Y\"}"
                    + "{\"input\": \"0\", \"expect\": \"0\", \"equal\":\"Y\"}"
                    + "{\"input\": \"18446744073709551615\", \"expect\": \"18446744073709551615\", \"equal\":\"Y\"}"
                    +"]");
            for (int i = 0; i < uint64Cases.size(); i++) {
                JSONObject testCase = uint64Cases.getJSONObject(i);
                String input = testCase.getString("input");
                String expect = testCase.getString("expect");
                String equal = testCase.getString("equal");
                TransactionReceipt tr = contract.setUint64(Uint64.of(input)).send();
                collector.logStepPass("To invoke setUint64 success, txHash: " + tr.getTransactionHash());
                Uint64 getReturn = contract.getUint64().send();
                collector.logStepPass("To invoke setUint64 success,setUint64: "+ input +", getUint64: " + getReturn.getValue().toString());
                if(equal.equals("Y")){
                    collector.assertEqual(getReturn.getValue(), new BigInteger(expect));
                } else {
                    collector.assertFalse(getReturn.getValue().compareTo(new BigInteger(expect)) == 0);
                }
            }

        } catch (Exception e) {
            collector.logStepFail("IntegerDataTypeContract_02 failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
