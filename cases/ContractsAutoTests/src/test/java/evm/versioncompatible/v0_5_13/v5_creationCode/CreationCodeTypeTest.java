package evm.versioncompatible.v0_5_13.v5_creationCode;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.CreationCode;
import network.platon.utils.DataChangeUtil;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;


/**
 * @title type(C).creationCode()：提供对合约的创建
 * @description:
 * @author: hudenian
 * @create: 2019/12/28
 */
public class CreationCodeTypeTest extends ContractPrepareTest {

    private String creationCodeTypeBin="608060405234801561001057600080fd5b5061011d806100206000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80632096525514602d575b600080fd5b603360ab565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101560715780820151818401526020810190506058565b50505050905090810190601f168015609d5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60606040518060400160405280600581526020017f68656c6c6f00000000000000000000000000000000000000000000000000000081525090509056fea265627a7a723158202347b0e035ea94f30c71aead790f4db2a69c4287904ff219b796729285b3640a64736f6c634300050d0032";

    @Before
    public void before() {
        this.prepare();
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "OverloadTest-type(C).creationCode()", sourcePrefix = "evm")
    public void testType() {
        try {

            CreationCode creationCode = CreationCode.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = creationCode.getContractAddress();
            TransactionReceipt tx = creationCode.getTransactionReceipt().get();

            collector.logStepPass("CreationCode deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + creationCode.getTransactionReceipt().get().getGasUsed());

            byte[] createCodeByteArr = creationCode.getContractName().send();

            String hexCreationCodeType =DataChangeUtil.bytesToHex(createCodeByteArr);

            collector.logStepPass("CreationCodeType的16进制字节码是："+hexCreationCodeType);


        } catch (Exception e) {
            collector.logStepFail("CreationCodeTypeTest  process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
