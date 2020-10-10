package network.platon.test.evm.event;

import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import org.junit.Test;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.List;

/**
 * @title 事件类型测试
 * @description:
 * @author: albedo
 * @create: 2020/01/06
 */
public class EventTypeContractTest extends ContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testOneDimensionalArray",
            author = "albedo", showName = "event.EventTypeContractTest-一维数组类型", sourcePrefix = "evm")
    public void testOneDimensionalArray() {
        try {
            prepare();
            EventTypeContract eventTypeContract = EventTypeContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = eventTypeContract.getContractAddress();
            String transactionHash = eventTypeContract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("EventTypeContract issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + eventTypeContract.getTransactionReceipt().get().getGasUsed());
            TransactionReceipt receipt = eventTypeContract.testOneDimensionalArray().send();
            List<EventTypeContract.OneDimensionalArrayEventEventResponse> one = eventTypeContract.getOneDimensionalArrayEventEvents(receipt);
            List<BigInteger> data = one.get(0).array;
            List<Uint256> except = new ArrayList<>(5);
            for (int i = 0; i < 5; i++) {
                except.add(new Uint256(new BigInteger(Integer.toString(i))));
            }
            collector.assertEqual(data, except, "checkout one dimensional array type declare event");
        } catch (Exception e) {
            collector.logStepFail("EventTypeContractTest testOneDimensionalArray failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testTwoDimensionalArray",
            author = "albedo", showName = "event.EventTypeContractTest-二维数组类型", sourcePrefix = "evm")
    public void testTwoDimensionalArray() {
        try {
            prepare();
            EventTypeContract eventCallContract = EventTypeContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = eventCallContract.getContractAddress();
            String transactionHash = eventCallContract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("EventTypeContract issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + eventCallContract.getTransactionReceipt().get().getGasUsed());
            TransactionReceipt receipt = eventCallContract.testTwoDimensionalArray().send();
            try {
                eventCallContract.getTwoDimensionalArrayEventEvents(receipt);
            } catch (UnsupportedOperationException e) {
                collector.assertEqual(e.getCause().getMessage(),"com.alaya.abi.solidity.datatypes.generated.StaticArray2<com.alaya.abi.solidity.datatypes.generated.Uint256>");
            }
        } catch (Exception e) {
            collector.logStepFail("EventTypeContractTest testTwoDimensionalArray failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testStr",
            author = "albedo", showName = "event.EventTypeContractTest-字符串", sourcePrefix = "evm")
    public void testStr() {
        try {
            prepare();
            EventTypeContract eventCallContract = EventTypeContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = eventCallContract.getContractAddress();
            String transactionHash = eventCallContract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("EventTypeContract issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + eventCallContract.getTransactionReceipt().get().getGasUsed());
            TransactionReceipt receipt = eventCallContract.testStr().send();
            List<EventTypeContract.StringEventEventResponse> str = eventCallContract.getStringEventEvents(receipt);
            String s = str.get(0).str;
            collector.assertEqual(s, "1234567890097865432112345678900987654321123456789009764354666663242444444444475831546856", "checkout string type declare event");
        } catch (Exception e) {
            collector.logStepFail("EventTypeContractTest testStr failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testEnum",
            author = "albedo", showName = "event.EventTypeContractTest-枚举", sourcePrefix = "evm")
    public void testEnum() {
        try {
            prepare();
            EventTypeContract eventCallContract = EventTypeContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = eventCallContract.getContractAddress();
            String transactionHash = eventCallContract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("EventTypeContract issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + eventCallContract.getTransactionReceipt().get().getGasUsed());
            TransactionReceipt receipt = eventCallContract.testEnum().send();
            List<EventTypeContract.EnumEventEventResponse> str = eventCallContract.getEnumEventEvents(receipt);
            BigInteger s = str.get(0).choices;
            collector.assertEqual(s, new BigInteger("0"), "checkout string type declare event");
        } catch (Exception e) {
            collector.logStepFail("EventTypeContractTest testStr failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
