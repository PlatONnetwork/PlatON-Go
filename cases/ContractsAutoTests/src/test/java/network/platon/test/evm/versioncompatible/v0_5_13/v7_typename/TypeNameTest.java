package network.platon.test.evm.versioncompatible.v0_5_13.v7_typename;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.TypeName;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;


/**
 * @title type(C).name()：提供对合约名称的访问
 * @description: 
 * @author: hudenian
 * @create: 2019/12/28
 */
public class TypeNameTest extends ContractPrepareTest {

    private String creationCodeTypeBin="6080604052348015600f57600080fd5b506004361060285760003560e01c80632096525514602d575b600080fd5b603360ab565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101560715780820151818401526020810190506058565b50505050905090810190601f168015609d5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60606040518060400160405280600581526020017f68656c6c6f00000000000000000000000000000000000000000000000000000081525090509056fea265627a7a7231582000701fa96aee5985c61ad55fed3b64e0df39f0356f468a390a28c47b0159281d64736f6c634300050d0032";

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "TypeNameTest-type(C).name()", sourcePrefix = "evm")
    public void testRuntimeType() {
        try {

            TypeName typeName = TypeName.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = typeName.getContractAddress();
            TransactionReceipt tx = typeName.getTransactionReceipt().get();

            collector.logStepPass("RuntimeCodeTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + typeName.getTransactionReceipt().get().getGasUsed());

            String contractNameStr = typeName.name().send().toString();


            collector.logStepPass("获取到的合约名是："+contractNameStr);

            collector.assertEqual("TypeName",contractNameStr);

        } catch (Exception e) {
            collector.logStepFail("TypeNameTest  process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
