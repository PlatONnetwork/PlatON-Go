package network.platon.contracts.evm;

import java.util.Arrays;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Bool;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.GasProvider;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.1.5.
 */
public class WhileError extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060c98061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063ed6f916c14602d575b600080fd5b6033604d565b604051808215151515815260200191505060405180910390f35b60008060566070565b90508060000160009054906101000a900460ff1691505090565b60005b600090508060000160009054906101000a900460ff16156091576073565b9056fea265627a7a7231582027abdddc06996992d8773bd854a31ab3eeb4107ef4bcc3317a82f5de8d7a895f64736f6c634300050d0032";

    public static final String FUNC_GETWHILECONTROLRES = "getWhileControlRes";

    protected WhileError(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected WhileError(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<Boolean> getWhileControlRes() {
        final Function function = new Function(FUNC_GETWHILECONTROLRES, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public static RemoteCall<WhileError> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(WhileError.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<WhileError> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(WhileError.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static WhileError load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new WhileError(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static WhileError load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new WhileError(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
