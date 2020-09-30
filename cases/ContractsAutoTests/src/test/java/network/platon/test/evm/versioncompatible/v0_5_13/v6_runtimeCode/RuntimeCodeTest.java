package network.platon.test.evm.versioncompatible.v0_5_13.v6_runtimeCode;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.RuntimeCode;
import network.platon.utils.DataChangeUtil;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;


/**
 * @title type(C).runtimeCode()：提供运行时代码的访问
 * @description: 
 * @author: hudenian
 * @create: 2019/12/28
 */
public class RuntimeCodeTest extends ContractPrepareTest {

    private String creationCodeTypeBin="608060405234801561001057600080fd5b5061011d806100206000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80632096525514602d575b600080fd5b603360ab565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101560715780820151818401526020810190506058565b50505050905090810190601f168015609d5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60606040518060400160405280600581526020017f68656c6c6f00000000000000000000000000000000000000000000000000000081525090509056fea265627a7a7231582095b7b1a6ae048b3f442cbb1ef4e1f340e91048c616d80049bf7fbdb68596a6e564736f6c634300050d0032";

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "RuntimeCodeTest-type(C).runtimeCode()", sourcePrefix = "evm")
    public void testRuntimeType() {
        try {

            RuntimeCode runtimeCode = RuntimeCode.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = runtimeCode.getContractAddress();
            TransactionReceipt tx = runtimeCode.getTransactionReceipt().get();

            collector.logStepPass("RuntimeCodeTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + runtimeCode.getTransactionReceipt().get().getGasUsed());

            byte[] runtimeCodeByteArr = runtimeCode.getContractName().send();

            String hexRuntimeCodeByteArrType =DataChangeUtil.bytesToHex(runtimeCodeByteArr);

            collector.logStepPass("RuntimeCodeByteArrType的16进制字节码是："+hexRuntimeCodeByteArrType);

        } catch (Exception e) {
            collector.logStepFail("RuntimeCodeTest  process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
