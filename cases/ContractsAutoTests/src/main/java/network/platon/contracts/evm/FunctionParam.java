package network.platon.contracts.evm;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
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
public class FunctionParam extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060d08061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c806326121ff014603757806392d0d153146053575b600080fd5b603d606f565b6040518082815260200191505060405180910390f35b6059607e565b6040518082815260200191505060405180910390f35b60006079607e6087565b905090565b60006007905090565b600060948263ffffffff16565b905091905056fea265627a7a72315820b8b25a17e403c6580f3f280b04695bdf34adda7e7f87e551d259ec8d4057608864736f6c634300050d0032";

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

    public RemoteCall<BigInteger> t() {
        final Function function = new Function(FUNC_T, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
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
