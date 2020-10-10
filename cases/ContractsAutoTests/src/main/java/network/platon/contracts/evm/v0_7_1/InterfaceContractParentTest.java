package network.platon.contracts.evm.v0_7_1;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.0.
 */
public class InterfaceContractParentTest extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060b98061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063e288b0b514602d575b600080fd5b606060048036036040811015604157600080fd5b8101908080359060200190929190803590602001909291905050506076565b6040518082815260200191505060405180910390f35b600081830190509291505056fea2646970667358221220c52dd2a37532ff326362e8d896b304396e1662ec67dc55de66f5e7e472907e2964736f6c63430007010033";

    public static final String FUNC_SUMEXTERNAL = "sumExternal";

    protected InterfaceContractParentTest(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected InterfaceContractParentTest(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> sumExternal(BigInteger a, BigInteger b) {
        final Function function = new Function(
                FUNC_SUMEXTERNAL, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Int256(a), 
                new com.alaya.abi.solidity.datatypes.generated.Int256(b)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<InterfaceContractParentTest> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InterfaceContractParentTest.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<InterfaceContractParentTest> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InterfaceContractParentTest.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static InterfaceContractParentTest load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new InterfaceContractParentTest(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static InterfaceContractParentTest load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new InterfaceContractParentTest(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
