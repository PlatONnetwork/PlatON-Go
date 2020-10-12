package network.platon.contracts.evm.v0_7_1;

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
 * <p>Generated with web3j version 0.13.2.1.
 */
public class IfError extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610132806100206000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c806349b9d30f146037578063c77beeb6146055575b600080fd5b603d6073565b60405180821515815260200191505060405180910390f35b605b6098565b60405180821515815260200191505060405180910390f35b600080607e600160bd565b90508060000160009054906101000a900460ff1691505090565b60008060a3600160d6565b90508060000160009054906101000a900460ff1691505090565b6000811560cc576000905060d1565b600090505b919050565b6000811560e5576000905060f7565b8160f1576000905060f6565b600090505b5b91905056fea2646970667358221220205fe28f4931e078d195fddd901b10d8da1b464c3c692b7301dce36313b08f0064736f6c63430007010033";

    public static final String FUNC_GETIFCONTROLRES = "getIfControlRes";

    public static final String FUNC_GETIFCONTROLRES1 = "getIfControlRes1";

    protected IfError(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected IfError(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<Boolean> getIfControlRes() {
        final Function function = new Function(FUNC_GETIFCONTROLRES, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<Boolean> getIfControlRes1() {
        final Function function = new Function(FUNC_GETIFCONTROLRES1, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public static RemoteCall<IfError> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(IfError.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<IfError> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(IfError.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static IfError load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new IfError(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static IfError load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new IfError(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
