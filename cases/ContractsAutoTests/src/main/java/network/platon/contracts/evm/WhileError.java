package network.platon.contracts.evm;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Bool;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.util.Arrays;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.0.
 */
public class WhileError extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060c98061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063ed6f916c14602d575b600080fd5b6033604d565b604051808215151515815260200191505060405180910390f35b60008060566070565b90508060000160009054906101000a900460ff1691505090565b60005b600090508060000160009054906101000a900460ff16156091576073565b9056fea265627a7a72315820f8da7af8b887ff2407f98076d5ccb87cc00634cfba050c6d83b20f731c328b5764736f6c63430005110032";

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
