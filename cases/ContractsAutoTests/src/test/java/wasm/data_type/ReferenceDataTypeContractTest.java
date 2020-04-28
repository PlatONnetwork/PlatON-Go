package wasm.data_type;

import com.platon.rlp.datatypes.Uint64;
import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ReferenceDataTypeContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @author zjsunzone
 *
 * This class is used to test data type of reference.
 */
public class ReferenceDataTypeContractTest extends WASMContractPrepareTest {

    @Before
    public void before(){
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.address_map_type",sourcePrefix = "wasm")
    public void testAddressMapContract() {

        try {
            // deploy contract.
            ReferenceDataTypeContract contract = ReferenceDataTypeContract.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ReferenceDataTypeContract deploy successfully.contractAddress:"
                    + contractAddress + ", hash:" + transactionHash
                    + " gasUsed:" + contract.getTransactionReceipt().get().getGasUsed().toString());

            // test: map
            String expectKey1 = "name";
            String expectValue11 = "0xc4482dd68fbaa5f4da143145198672bd17245ff2";
            TransactionReceipt mapTr = contract.setAddressMap(expectKey1, expectValue11).send();
            collector.logStepPass("To invoke setAddressMap success, txHash1: " + mapTr.getTransactionHash());

            String expectKey2 = "name2";
            String expectValue2 = "0xc4482dd68fbaa5f4da143145198672bd17245ff1";
            TransactionReceipt mapTr2 = contract.setAddressMap(expectKey2, expectValue2).send();
            collector.logStepPass("To invoke setAddressMap success, txHash2: " + mapTr2.getTransactionHash());

            String actValue1 = contract.getAddrFromMap(expectKey1).send();
            String actValue2 = contract.getAddrFromMap(expectKey2).send();
            collector.logStepPass("To invoke getAddrFromMap success, value1: " + actValue1 + " value2:" + actValue2);

            Uint8 mapSize = contract.sizeOfAddrMap().send();
            collector.logStepPass("To invoke sizeOfAddrMap success, mapSize: " + mapSize.getValue().toString());
            collector.assertEqual(mapSize.getValue().intValue(), 2);
            collector.assertEqual(prependHexPrefix(actValue1).toUpperCase(), prependHexPrefix(expectValue11).toUpperCase());
            collector.assertEqual(prependHexPrefix(actValue2).toUpperCase(), prependHexPrefix(expectValue2).toUpperCase());


        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("ReferenceDataTypeContract and could not call contract function");
            }else{
                collector.logStepFail("ReferenceDataTypeContract failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.u256_map_type",sourcePrefix = "wasm")
    public void testU256MapContract() {

        try {
            // deploy contract.
            ReferenceDataTypeContract contract = ReferenceDataTypeContract.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ReferenceDataTypeContract deploy successfully.contractAddress:"
                    + contractAddress + ", hash:" + transactionHash
                    + " gasUsed:" + contract.getTransactionReceipt().get().getGasUsed().toString());

            // test: u256
            String expectKey1 = "name";
            String expectValue11 = "100000";
            TransactionReceipt tr1 = contract.setU256Map(expectKey1, Uint64.of(expectValue11)).send();
            collector.logStepPass("To invoke setU256Map success, txHash1: " + tr1.getTransactionHash());

            String expectKey2 = "name2";
            String expectValue2 = "200000";
            TransactionReceipt mapTr2 = contract.setU256Map(expectKey2, Uint64.of(expectValue2)).send();
            collector.logStepPass("To invoke setU256Map success, txHash2: " + mapTr2.getTransactionHash());

            String actValue1 = contract.getU256FromMap(expectKey1).send();
            String actValue2 = contract.getU256FromMap(expectKey2).send();
            collector.logStepPass("To invoke getU256FromMap success, value1: " + actValue1 + " value2:" + actValue2);

            Uint8 mapSize = contract.sizeOfU256Map().send();
            collector.logStepPass("To invoke sizeOfU256Map success, mapSize: " + mapSize.getValue().toString());
            collector.assertEqual(mapSize.getValue().intValue(), 2);
            collector.assertEqual(actValue1.toUpperCase(), expectValue11.toUpperCase());
            collector.assertEqual(actValue2.toUpperCase(), expectValue2.toUpperCase());


        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("ReferenceDataTypeContract and could not call contract function");
            }else{
                collector.logStepFail("ReferenceDataTypeContract failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.h256_map_type",sourcePrefix = "wasm")
    public void testH256MapContract() {

        try {
            // deploy contract.
            ReferenceDataTypeContract contract = ReferenceDataTypeContract.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ReferenceDataTypeContract deploy successfully.contractAddress:"
                    + contractAddress + ", hash:" + transactionHash
                    + " gasUsed:" + contract.getTransactionReceipt().get().getGasUsed().toString());

            // test: u256
            String expectKey1 = "name";
            String expectValue11 = "0xc648eb226f98cbb05835c9fad2cbaa419c602782ed4fbd4a6b5c6789e8292a85";
            TransactionReceipt tr1 = contract.setH256Map(expectKey1, expectValue11).send();
            collector.logStepPass("To invoke setH256Map success, txHash1: " + tr1.getTransactionHash());

            String expectKey2 = "name2";
            String expectValue2 = "0xc648eb226f98cbb05835c9fad2cbaa419c602782ed4fbd4a6b5c6789e8292a86";
            TransactionReceipt mapTr2 = contract.setH256Map(expectKey2, expectValue2).send();
            collector.logStepPass("To invoke setH256Map success, txHash2: " + mapTr2.getTransactionHash());

            String actValue1 = contract.getH256FromMap(expectKey1).send();
            String actValue2 = contract.getH256FromMap(expectKey2).send();
            collector.logStepPass("To invoke getH256FromMap success, value1: " + actValue1 + " value2:" + actValue2);

            Uint8 mapSize = contract.sizeOfH256Map().send();
            collector.logStepPass("To invoke sizeOfH256Map success, mapSize: " + mapSize.getValue().toString());
            collector.assertEqual(mapSize.getValue().intValue(), 2);
            //collector.assertEqual(actValue1.toUpperCase(), expectValue11.toUpperCase());
            //collector.assertEqual(actValue2.toUpperCase(), expectValue2.toUpperCase());


        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("ReferenceDataTypeContract and could not call contract function");
            }else{
                collector.logStepFail("ReferenceDataTypeContract failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }
}
