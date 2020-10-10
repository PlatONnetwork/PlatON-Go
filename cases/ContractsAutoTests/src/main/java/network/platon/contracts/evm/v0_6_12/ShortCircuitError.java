package network.platon.contracts.evm.v0_6_12;

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
public class ShortCircuitError extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610134806100206000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c806304c09ce91460375780630c204dbc146055575b600080fd5b603d6073565b60405180821515815260200191505060405180910390f35b605b6096565b60405180821515815260200191505060405180910390f35b600080607c60b9565b90508060000160009054906101000a900460ff1691505090565b600080609f60dd565b90508060000160009054906101000a900460ff1691505090565b6000600190508060000160009054906101000a900460ff168060d9575060015b5090565b60008060000160009054906101000a900460ff16801560fa575060005b509056fea2646970667358221220a633e07a453ededf00d48bca657d8421bb35b0bc732ff4245eb9f7a41603c24d64736f6c634300060c0033";

    public static final String FUNC_GETF = "getF";

    public static final String FUNC_GETG = "getG";

    protected ShortCircuitError(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ShortCircuitError(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getF() {
        final Function function = new Function(
                FUNC_GETF, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getG() {
        final Function function = new Function(
                FUNC_GETG, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<ShortCircuitError> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ShortCircuitError.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<ShortCircuitError> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ShortCircuitError.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static ShortCircuitError load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ShortCircuitError(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ShortCircuitError load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ShortCircuitError(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
