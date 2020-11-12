package network.platon.test.wasm.contract_create;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InitWithObjectParams;
import org.junit.Test;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 创建合约init带一个入参对象测试
 * @description:
 * @author: hudenian
 * @create: 2020/02/07
 */
public class InitWithObjectParamsTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_create创建合约init带一个入参",sourcePrefix = "wasm")
    public void testNewContract() {

        String body = "myBody";
        String end = "myEnd";
        try {
            prepare();
            InitWithObjectParams.My_message myMessage = new InitWithObjectParams.My_message();
            myMessage.body = body;
            myMessage.end = end;
            InitWithObjectParams initWithObjectParams = InitWithObjectParams.deploy(web3j, transactionManager, provider, chainId,myMessage).send();
            String contractAddress = initWithObjectParams.getContractAddress();
            String transactionHash = initWithObjectParams.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("InitWithObjectParams issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + initWithObjectParams.getTransactionReceipt().get().getGasUsed());

            InitWithObjectParams.My_message[] messages = initWithObjectParams.get_message("").send();
            collector.assertEqual(messages[0].body,body);

        } catch (Exception e) {
            collector.logStepFail("InitWithObjectParamsTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
