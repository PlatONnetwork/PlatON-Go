package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.GasProvider;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.0.
 */
public class InheritContractMutipleTest2 extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060ba8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c8063430fe9c1146037578063e19ced32146053575b600080fd5b603d606f565b6040518082815260200191505060405180910390f35b60596078565b6040518082815260200191505060405180910390f35b60006002905090565b60006080606f565b90509056fea265627a7a7231582014bf87a9c177faf3208f04185111e734f2fbb2e9ec4f79bbb4bc9c907d0cae2664736f6c634300050d0032";

    public static final String FUNC_CALLGETDATE2 = "callGetDate2";

    public static final String FUNC_GETDATE = "getDate";

    @Deprecated
    protected InheritContractMutipleTest2(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected InheritContractMutipleTest2(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected InheritContractMutipleTest2(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected InheritContractMutipleTest2(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<BigInteger> callGetDate2() {
        final Function function = new Function(FUNC_CALLGETDATE2, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getDate() {
        final Function function = new Function(FUNC_GETDATE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<InheritContractMutipleTest2> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(InheritContractMutipleTest2.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<InheritContractMutipleTest2> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(InheritContractMutipleTest2.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<InheritContractMutipleTest2> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(InheritContractMutipleTest2.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<InheritContractMutipleTest2> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(InheritContractMutipleTest2.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static InheritContractMutipleTest2 load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new InheritContractMutipleTest2(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static InheritContractMutipleTest2 load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new InheritContractMutipleTest2(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static InheritContractMutipleTest2 load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new InheritContractMutipleTest2(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static InheritContractMutipleTest2 load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new InheritContractMutipleTest2(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
