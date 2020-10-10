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
public class FunctionParam extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060d18061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c806326121ff014603757806392d0d153146053575b600080fd5b603d606f565b6040518082815260200191505060405180910390f35b6059607e565b6040518082815260200191505060405180910390f35b60006079607e6087565b905090565b60006007905090565b600060948263ffffffff16565b905091905056fea2646970667358221220151129aac5a586dc804a6238031b2c2936ec63cdbe9ee279ae50662d2cc83c5764736f6c634300060c0033";

    public static final String FUNC_F = "f";

    public static final String FUNC_T = "t";

    protected FunctionParam(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected FunctionParam(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> f() {
        final Function function = new Function(
                FUNC_F, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> t() {
        final Function function = new Function(
                FUNC_T, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<FunctionParam> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(FunctionParam.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<FunctionParam> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(FunctionParam.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static FunctionParam load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new FunctionParam(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static FunctionParam load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new FunctionParam(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
