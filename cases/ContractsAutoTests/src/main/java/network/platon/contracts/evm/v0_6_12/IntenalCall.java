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
public class IntenalCall extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060d18061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c8063569c5f6d146037578063de292789146053575b600080fd5b603d606f565b6040518082815260200191505060405180910390f35b60596086565b6040518082815260200191505060405180910390f35b600080600190506000600290508082019250505090565b600080600390506093606f565b81029150509056fea264697066735822122014e2175660a4a8ae5030cc2efda5e9607590cc9679f3edb7277b5450910d2fad64736f6c634300060c0033";

    public static final String FUNC_GETRESULT = "getResult";

    public static final String FUNC_GETSUM = "getSum";

    protected IntenalCall(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected IntenalCall(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getResult() {
        final Function function = new Function(
                FUNC_GETRESULT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getSum() {
        final Function function = new Function(
                FUNC_GETSUM, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<IntenalCall> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(IntenalCall.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<IntenalCall> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(IntenalCall.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static IntenalCall load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new IntenalCall(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static IntenalCall load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new IntenalCall(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
