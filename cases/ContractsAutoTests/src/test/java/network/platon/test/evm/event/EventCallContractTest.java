package network.platon.test.evm.event;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import com.alibaba.fastjson.JSONObject;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.EventCallContract;
import network.platon.utils.DataChangeUtil;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;
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
            author = "albedo", showName = "event.EventCallContractTest-event关键字声明事件", sourcePrefix = "evm")
    public void testEmitEvent() {
        try {
            prepare();
            EventCallContract eventCallContract = EventCallContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = eventCallContract.getContractAddress();
            String transactionHash = eventCallContract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("EventCallContract issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + eventCallContract.getTransactionReceipt().get().getGasUsed());
            TransactionReceipt receipt = eventCallContract.emitEvent().send();
            List<EventCallContract.IncrementEventResponse> emitEventData = eventCallContract.getIncrementEvents(receipt);
            String data = emitEventData.get(0).log.getData();
            collector.assertEqual(emitEventData.get(0).who, receipt.getFrom(), "checkout declare event keyword");
        } catch (Exception e) {
            collector.logStepFail("EventCallContractTest testEmitEvent failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "indexedEvent",
            author = "albedo", showName = "event.EventCallContractTest-indexed关键字定义事件索引", sourcePrefix = "evm")
    public void testIndexedEvent() {
        try {
            prepare();
            EventCallContract eventCallContract = EventCallContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = eventCallContract.getContractAddress();
            String transactionHash = eventCallContract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("EventCallContract issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            TransactionReceipt receipt = eventCallContract.indexedEvent().send();
            List<EventCallContract.DepositEventResponse> emitEventData = eventCallContract.getDepositEvents(receipt);
            String data = emitEventData.get(0).log.getData();
            collector.assertEqual(emitEventData.get(0)._from, receipt.getFrom(), "checkout new contract param");
            collector.assertEqual(emitEventData.get(0)._value.toString(), "12", "checkout indexed keyword");
        } catch (Exception e) {
            collector.logStepFail("EventCallContractTest testIndexedEvent failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "anonymousEvent",
            author = "albedo", showName = "event.EventCallContractTest-anonymous关键字定义匿名事件", sourcePrefix = "evm")
    public void testAnonymousEvent() {
        try {
            prepare();
            EventCallContract eventCallContract = EventCallContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = eventCallContract.getContractAddress();
            String transactionHash = eventCallContract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("EventCallContract issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            TransactionReceipt receipt = eventCallContract.anonymousEvent().send();
            collector.assertEqual(DataChangeUtil.subHexData(receipt.getLogs().get(0).getData()), DataChangeUtil.subHexData("1"), "checkout anonymous keyword");
        } catch (Exception e) {
            collector.logStepFail("EventCallContractTest testAnonymousEvent failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testEmitEvents",
            author = "albedo", showName = "event.EventCallContractTest-函数多事件监听", sourcePrefix = "evm")
    public void testEmitEvents() {
        try {
            prepare();
            EventCallContract eventCallContract = EventCallContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = eventCallContract.getContractAddress();
            String transactionHash = eventCallContract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("EventCallContract issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            TransactionReceipt receipt = eventCallContract.testBool().send();
            System.out.println(JSONObject.toJSONString(receipt.getLogs()));
            List<EventCallContract.BoolEventEventResponse> bList = eventCallContract.getBoolEventEvents(receipt);
            collector.assertEqual(bList.get(0).result ,Boolean.FALSE, "checkout multi bool type event");
            collector.assertEqual(bList.get(1).result ,Boolean.TRUE, "checkout multi bool type event");
            List<EventCallContract.IncrementEventResponse> iList = eventCallContract.getIncrementEvents(receipt);
            collector.assertEqual(iList.get(0).who ,receipt.getFrom(), "checkout address type event");

            List<EventCallContract.DepositEventResponse> dList = eventCallContract.getDepositEvents(receipt);
            collector.assertEqual(dList.get(0)._from ,receipt.getFrom(), "checkout address type event");
            collector.assertEqual(dList.get(0)._value ,new BigInteger("12"), "checkout uint event");

            collector.assertEqual(DataChangeUtil.subHexData(receipt.getLogs().get(3).getData()),DataChangeUtil.subHexData("c"), "checkout address type event");

        } catch (Exception e) {
            collector.logStepFail("EventCallContractTest testAnonymousEvent failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testMultiAnonymousEvents",
            author = "albedo", showName = "event.EventCallContractTest-函数多匿名事件监听", sourcePrefix = "evm")
    public void testMultiAnonymousEvents() {
        try {
            prepare();
            EventCallContract eventCallContract = EventCallContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = eventCallContract.getContractAddress();
            String transactionHash = eventCallContract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("EventCallContract issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            TransactionReceipt receipt = eventCallContract.testMultiAnonymous().send();
            collector.assertEqual(DataChangeUtil.subHexData(receipt.getLogs().get(0).getData()),DataChangeUtil.subHexData("c"), "checkout multi anonymous event");
            collector.assertEqual(DataChangeUtil.subHexData(receipt.getLogs().get(1).getData()),DataChangeUtil.subHexData("d"), "checkout multi anonymous event");
            collector.assertEqual(DataChangeUtil.subHexData(receipt.getLogs().get(2).getData()),DataChangeUtil.subHexData("e"), "checkout multi anonymous event");
        } catch (Exception e) {
            collector.logStepFail("EventCallContractTest testMultiAnonymousEvents failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }
}
