package event;

import beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.EventCallContract;
import network.platon.utils.DataChangeUtil;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

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
            collector.assertEqual(DataChangeUtil.subHexData(data), DataChangeUtil.subHexData(receipt.getFrom()), "checkout declare event keyword");
        } catch (Exception e) {
            collector.logStepFail("EventCallContractTest testEmitEvent failure,exception msg:" , e.getMessage());
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
            collector.assertEqual(DataChangeUtil.subHexData(emitEventData.get(0).log.getTopics().get(1)), DataChangeUtil.subHexData(receipt.getFrom()), "checkout new contract param");
            collector.assertEqual(DataChangeUtil.subHexData(data), DataChangeUtil.subHexData("c"), "checkout indexed keyword");
        } catch (Exception e) {
            collector.logStepFail("EventCallContractTest testIndexedEvent failure,exception msg:" , e.getMessage());
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
            collector.assertEqual(DataChangeUtil.subHexData(receipt.getLogs().get(0).getData()), DataChangeUtil.subHexData("1") ,"checkout anonymous keyword");
        } catch (Exception e) {
            collector.logStepFail("EventCallContractTest testAnonymousEvent failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

}
