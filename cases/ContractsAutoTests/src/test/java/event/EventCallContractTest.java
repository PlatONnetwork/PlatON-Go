package event;

import beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.AbiEncoder;
import network.platon.contracts.EventCallContract;
import org.apache.commons.lang.StringUtils;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.List;

/**
 * @title 事件验证测试
 * @description:
 * @author: albedo
 * @create: 2019/12/28
 */
public class EventCallContractTest extends ContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "emitEvent",
            author = "albedo", showName = "event.EventCallContractTest-event关键字声明事件")
    public void testEmitEvent() {
        try {
            prepare();
            EventCallContract eventCallContract = EventCallContract.deploy(web3j, transactionManager, provider).send();
            String contractAddress = eventCallContract.getContractAddress();
            String transactionHash = eventCallContract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("EventCallContract issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            TransactionReceipt receipt = eventCallContract.emitEvent().send();
            List<EventCallContract.IncrementEventResponse> emitEventData = eventCallContract.getIncrementEvents(receipt);
            String data = emitEventData.get(0).log.getData();
            collector.assertEqual(subHexData(data), subHexData(receipt.getFrom()), "checkout declare event keyword");
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "indexedEvent",
            author = "albedo", showName = "event.EventCallContractTest-indexed关键字定义事件索引")
    public void testIndexedEvent() {
        try {
            prepare();
            EventCallContract eventCallContract = EventCallContract.deploy(web3j, transactionManager, provider).send();
            String contractAddress = eventCallContract.getContractAddress();
            String transactionHash = eventCallContract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("EventCallContract issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            TransactionReceipt receipt = eventCallContract.indexedEvent().send();
            List<EventCallContract.DepositEventResponse> emitEventData = eventCallContract.getDepositEvents(receipt);
            String data = emitEventData.get(0).log.getData();
            collector.assertEqual(subHexData(emitEventData.get(0).log.getTopics().get(1)), subHexData(receipt.getFrom()), "checkout new contract param");
            collector.assertEqual(subHexData(data), subHexData("c"), "checkout indexed keyword");
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "anonymousEvent",
            author = "albedo", showName = "event.EventCallContractTest-anonymous关键字定义匿名事件")
    public void testAnonymousEvent() {
        try {
            prepare();
            EventCallContract eventCallContract = EventCallContract.deploy(web3j, transactionManager, provider).send();
            String contractAddress = eventCallContract.getContractAddress();
            String transactionHash = eventCallContract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("EventCallContract issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            TransactionReceipt receipt = eventCallContract.anonymousEvent().send();
            List<EventCallContract.AnonymousEventResponse> data=eventCallContract.getAnonymousEvents(receipt);
            collector.assertEqual(data.get(0)._id, new BigInteger("1") ,"checkout anonymous keyword");
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "anonymousEvent",
            author = "albedo", showName = "event.EventCallContractTest-anonymous关键字定义匿名事件")
    public void testAbiEncoderEvent() {
        try {
            prepare();
            AbiEncoder abiEncoder = AbiEncoder.deploy(web3j, transactionManager, provider).send();
            String contractAddress = abiEncoder.getContractAddress();
            String transactionHash = abiEncoder.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("EventCallContract issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            TransactionReceipt receipt = abiEncoder.test().send();
            List<AbiEncoder.EEventResponse> emitEventData = abiEncoder.getEEvents(receipt);
            List<BigInteger> data = emitEventData.get(0).multi;
            List<BigInteger> except=new ArrayList<>(5);
            for (int i = 0; i < 5; i++) {
                except.add(new BigInteger(Integer.toString(i)));
            }
            collector.assertEqual(data, except, "checkout declare event keyword");
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    private String subHexData(String hexStr) {
        if (StringUtils.isBlank(hexStr)) {
            throw new IllegalArgumentException("string is blank");
        }
        if (StringUtils.startsWith(hexStr, "0x")) {
            hexStr = StringUtils.substringAfter(hexStr, "0x");
        }
        byte[] addi = hexStr.getBytes();
        for (int i = 0; i < addi.length; i++) {
            if (addi[i] != 0) {
                hexStr = StringUtils.substring(hexStr, i - 1);
                break;
            }
        }
        return hexStr;
    }
}
