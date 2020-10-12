package network.platon.contracts.evm.v0_6_12;

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
public class ForError extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061014a806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063895e3ada1461003b578063ec56ae5d1461005b575b600080fd5b61004361007b565b60405180821515815260200191505060405180910390f35b6100636100a0565b60405180821515815260200191505060405180910390f35b6000806100866100c5565b90508060000160009054906101000a900460ff1691505090565b6000806100ab6100ee565b90508060000160009054906101000a900460ff1691505090565b60005b600115158160000160009054906101000a900460ff16151514156100eb576100c8565b90565b60005b600190508060000160009054906101000a900460ff1615610111576100f1565b9056fea264697066735822122019f27abdefd02ae16145788e1941c9ee6126a5079a43e3bca3b3ef04235f682564736f6c634300060c0033";

    public static final String FUNC_GETFORCONTROLRES = "getForControlRes";

    public static final String FUNC_GETFORCONTROLRES1 = "getForControlRes1";

    protected ForError(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ForError(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<Boolean> getForControlRes() {
        final Function function = new Function(FUNC_GETFORCONTROLRES, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<Boolean> getForControlRes1() {
        final Function function = new Function(FUNC_GETFORCONTROLRES1, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public static RemoteCall<ForError> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ForError.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<ForError> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ForError.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static ForError load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ForError(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ForError load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ForError(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
