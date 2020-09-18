package network.platon.test.evm.versioncompatible.v0_4_25;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import com.alibaba.fastjson.JSONObject;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.StorageLocation;
import org.junit.Test;
/**
 * @title 数据位置测试
 * 1. 0.4.25版本验证使用constrictor关键字定义构造函数，但是不强制声明可见性(默认为public可见性）
 * 2. 0.4.25版本同一继承层次结构中允许多次指定基类构造函数参数验证:
 * (1) 允许合约直接声明构造函数 ———— is Base(7)
 *（2）子合约构造函数继承父合约构造函数———— constructor(uint _y) Base(_y * _y)
 * 两种引用构造函数方式共存时，合约优先选择（2）方式
 * @description:
 * @author: albedo
 * @create: 2019/12/28
 */
public class StorageLocationTest extends ContractPrepareTest {

    byte[] test=new byte[]{1,2,3,4};

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testStorageLocaltionCheck",
            author = "albedo", showName = "network.platon.test.evm.StorageLocationTest-参数变量数组类型不须显式声明", sourcePrefix = "evm")
    public void testStorageLocationCheck() {
        try {
            prepare();
            StorageLocation storageLocation = StorageLocation.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = storageLocation.getContractAddress();
            String transactionHash = storageLocation.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("StorageLocation issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + storageLocation.getTransactionReceipt().get().getGasUsed());
            byte[] result = storageLocation.storageLocaltionCheck(test).send();
            collector.assertEqual(JSONObject.toJSONString(result), JSONObject.toJSONString(test), "checkout array storage location result");
        } catch (Exception e) {
            collector.logStepFail("StorageLocationTest testStorageLocationCheck failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testTransfer",
            author = "albedo", showName = "StorageLocationTest-external函数数组类型参数不需显式声明 calldata", sourcePrefix = "evm")
    public void testTransfer() {
        try {
            prepare();
            StorageLocation storageLocation = StorageLocation.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = storageLocation.getContractAddress();
            String transactionHash = storageLocation.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("StorageLocation issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + storageLocation.getTransactionReceipt().get().getGasUsed());
            byte[] result = storageLocation.transfer(test).send();
            collector.assertEqual(JSONObject.toJSONString(result), JSONObject.toJSONString(test), "checkout external declare function array location result");
        } catch (Exception e) {
            collector.logStepFail("StorageLocationTest testTransfer failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

}
