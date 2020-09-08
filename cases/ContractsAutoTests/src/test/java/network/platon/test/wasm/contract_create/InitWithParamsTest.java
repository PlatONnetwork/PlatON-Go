package wasm.contract_create;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InitWithParams;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 创建合约init带一个入参测试
 * @description:
 * @author: hudenian
 * @create: 2020/02/07
 */
public class InitWithParamsTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_create创建合约init带一个入参",sourcePrefix = "wasm")
    public void testNewContract() {

        String initName = "hudenian";
        String addName = "hudenian1";
        try {
            prepare();
            InitWithParams initWithParams = InitWithParams.deploy(web3j, transactionManager, provider, chainId,initName).send();
            String contractAddress = initWithParams.getContractAddress();
            String transactionHash = initWithParams.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("initWithParams issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + initWithParams.getTransactionReceipt().get().getGasUsed());

            InitWithParams.Person person =  new InitWithParams.Person();
            person.name=addName;
            TransactionReceipt transactionReceipt = initWithParams.add_person(person).send();
            collector.logStepPass("initWithParams add_vector successfully hash:" + transactionReceipt.getTransactionHash());

            InitWithParams.Person[] peoples = initWithParams.get_person().send();
            collector.assertEqual(peoples[1].name,addName);

        } catch (Exception e) {
            collector.logStepFail("initWithParamsTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
