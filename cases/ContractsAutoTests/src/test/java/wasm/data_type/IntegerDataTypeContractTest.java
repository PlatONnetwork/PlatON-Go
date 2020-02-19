package wasm.data_type;

import com.platon.rlp.Int64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.IntegerDataTypeContract_1;
import network.platon.contracts.wasm.IntegerDataTypeContract_2;
import network.platon.contracts.wasm.IntegerDataTypeContract_3;
import network.platon.contracts.wasm.IntegerDataTypeContract_4;
import org.junit.Before;
import org.junit.Test;
import org.web3j.abi.datatypes.Address;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

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
            IntegerDataTypeContract_1 contract = IntegerDataTypeContract_1.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("IntegerDataTypeContract_01 issued successfully.contractAddress:"
                    + contractAddress + ", hash:" + transactionHash
                    + " gasUsed:" + contract.getTransactionReceipt().get().getGasUsed().toString());

            // test: int8
            Int64 int8 = contract.int8().send();
            collector.logStepPass("To invoke int8 success, int8: " + int8.value);

            // test: int64
            Int64 int64 = contract.int64().send();
            collector.logStepPass("To invoke int8 success, int64: " + int64);

            // test uint8
            Byte expectByt8 = Byte.valueOf((byte) 10);
            Byte uint8 = contract.uint8t(expectByt8).send();
            collector.logStepPass("To invoke uint8 success, uint8: " + uint8.byteValue());

            // test: uint32
            Integer expectByt32 = 1000;
            Integer uint32 = contract.uint32t(expectByt32).send();
            collector.logStepPass("To invoke uint32 success, uint32: " + uint32.intValue());
            collector.assertEqual(uint32, expectByt32 * 2);

            // test: uint64
            Long expect64 = Long.valueOf(10000);
            Long uint64 = contract.uint64t(expect64).send();
            collector.logStepPass("To invoke uint64 success, uint64: " + uint64.longValue());
            collector.assertEqual(uint64, expect64 * 2);

            // test: u128
            Long expect128 = Long.valueOf(10000);
            String u128 = contract.u128t(expect128).send();
            collector.logStepPass("To invoke uint64 success, u128: " + u128);
            collector.assertEqual(u128, expect128.toString());

            // test: u256
            Long expect256 = Long.valueOf(10000);
            String u256 = contract.u256t(expect256).send();
            collector.logStepPass("To invoke u256t success, u256: " + u128);
            collector.assertEqual(u256, expect256.toString());

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
            IntegerDataTypeContract_2 contract = IntegerDataTypeContract_2.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("IntegerDataTypeContract_01 issued successfully.contractAddress:"
                    + contractAddress + ", hash:" + transactionHash
                    + " gasUsed:" + contract.getTransactionReceipt().get().getGasUsed().toString());

            //  int8
            TransactionReceipt int8Tr = contract.setInt8(Int64.of((byte) 2)).send();
            collector.logStepPass("To invoke setInt 8 success, txHash: " + int8Tr.getTransactionHash());
            Int64 getInt8 = contract.getInt8().send();
            collector.logStepPass("To invoke getInt8 8 success, getInt8: " + getInt8.getValue());

            //  int32
            TransactionReceipt int32Tr = contract.setInt32(Int64.of(100)).send();
            collector.logStepPass("To invoke setInt32 success, txHash: " + int32Tr.getTransactionHash());
            Int64 getInt32 = contract.getInt32().send();
            collector.logStepPass("To invoke getInt32 success, getInt32: " + getInt32);

            // int64
            TransactionReceipt int64Tr = contract.setInt64(Int64.of(1111111111)).send();
            collector.logStepPass("To invoke setInt64 success, txHash: " + int64Tr.getTransactionHash());
            Int64 getInt64 = contract.getInt64().send();
            collector.logStepPass("To invoke getInt64 success, getInt64: " + getInt64);

            // ======================= uint =======================
            //  uint8
            TransactionReceipt uint8Tr = contract.setUint8(Byte.valueOf((byte) 2)).send();
            collector.logStepPass("To invoke setUint8 success, txHash: " + uint8Tr.getTransactionHash());
            Byte getUint8 = contract.getUint8().send();
            collector.logStepPass("To invoke getUint8 8 success, getUint8: " + getUint8.byteValue());

            //  uint32
            TransactionReceipt uint32Tr = contract.setUint32(100).send();
            collector.logStepPass("To invoke setuUint32 success, txHash: " + uint32Tr.getTransactionHash());
            Integer getUint32 = contract.getUint32().send();
            collector.logStepPass("To invoke getUint32 success, getUint32: " + getUint32);

            // uint64
            TransactionReceipt uint64Tr = contract.setUint64(Long.valueOf("1111111111")).send();
            collector.logStepPass("To invoke setUint64 success, txHash: " + uint64Tr.getTransactionHash());
            Long getUint64 = contract.getUint64().send();
            collector.logStepPass("To invoke getUint64 success, getUint64: " + getUint64);

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
            IntegerDataTypeContract_3 contract = IntegerDataTypeContract_3.deploy(web3j, transactionManager, provider).send();
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
            TransactionReceipt charTr = contract.setChar(Int64.of(expectByte.byteValue())).send();
            collector.logStepPass("To invoke setChar success, txHash: " + charTr.getTransactionHash());
            Int64 getChar = contract.getChar().send();
            collector.logStepPass("To invoke getChar success, getChar: " + getChar.getValue());
            collector.assertEqual(getChar.getValue(), expectByte.longValue());

        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("IntegerDataTypeContract_3 and could not call contract function");
            }else{
                collector.logStepFail("IntegerDataTypeContract_3 failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.base_data_type_04",sourcePrefix = "wasm")
    public void testBaseTypeContract_04() {

        try {
            // deploy contract.
            IntegerDataTypeContract_4 contract = IntegerDataTypeContract_4.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("IntegerDataTypeContract_4 issued successfully.contractAddress:"
                    + contractAddress + ", hash:" + transactionHash
                    + " gasUsed:" + contract.getTransactionReceipt().get().getGasUsed().toString());

            // test: store address
            Address expectAddr = new Address("0x5b05e7a3e2a688c5e5cc491545a84a1efc66c1b1");
            TransactionReceipt addrTr = contract.setAddress(expectAddr.getValue()).send();
            collector.logStepPass("To invoke setAddress success, txHash: " + addrTr.getTransactionHash());
            String getAddress = contract.getAddress().send();
            collector.logStepPass("To invoke getAddress success, getAddress: " + getAddress);
            //collector.assertEqual(getAddress, expectAddr);

            // test: store u256
            String expectU256 = "100000";
            TransactionReceipt u256TR = contract.setU256(Long.valueOf(expectU256)).send();
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
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("IntegerDataTypeContract_4 and could not call contract function");
            }else{
                collector.logStepFail("IntegerDataTypeContract_4 failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.base_data_type_border",sourcePrefix = "wasm")
    public void testBaseTypeContract_border() {
        // 主要测试各类型的边界值
        try {
            // deploy contract.
            IntegerDataTypeContract_2 contract = IntegerDataTypeContract_2.deploy(web3j, transactionManager, provider).send();
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

            //  int8 -128 ~ 127
            TransactionReceipt int8Tr = contract.setInt8(Int64.of(-128)).send();
            collector.logStepPass("To invoke setInt8 success, txHash: " + int8Tr.getTransactionHash());
            Int64 getInt8 = contract.getInt8().send();
            collector.logStepPass("To invoke getInt8 success, getInt8: " + getInt8.getValue());



        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("IntegerDataTypeContract_02 and could not call contract function");
            }else{
                collector.logStepFail("IntegerDataTypeContract_02 failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }
}
