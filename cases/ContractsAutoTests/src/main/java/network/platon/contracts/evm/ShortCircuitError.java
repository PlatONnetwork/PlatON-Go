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
public class ShortCircuitError extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610148806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806304c09ce91461003b5780630c204dbc1461005d575b600080fd5b61004361007f565b604051808215151515815260200191505060405180910390f35b6100656100a4565b604051808215151515815260200191505060405180910390f35b60008061008a6100c9565b90508060000160009054906101000a900460ff1691505090565b6000806100af6100ee565b90508060000160009054906101000a900460ff1691505090565b6000600190508060000160009054906101000a900460ff16806100ea575060015b5090565b60008090508060000160009054906101000a900460ff16801561010f575060005b509056fea265627a7a723158204165eb40f3d49ef0b3e0d9d27b1827055285d19bd5e2316169fc792c3ee1190964736f6c634300050d0032";

    public static final String FUNC_GETF = "getF";

    public static final String FUNC_GETG = "getG";

    protected ShortCircuitError(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ShortCircuitError(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<Boolean> getF() {
        final Function function = new Function(FUNC_GETF, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<Boolean> getG() {
        final Function function = new Function(FUNC_GETG, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
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
