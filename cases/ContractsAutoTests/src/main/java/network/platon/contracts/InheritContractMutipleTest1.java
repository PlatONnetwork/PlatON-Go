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
public class InheritContractMutipleTest1 extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060ba8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c8063430fe9c11460375780636f45ba21146053575b600080fd5b603d606f565b6040518082815260200191505060405180910390f35b60596078565b6040518082815260200191505060405180910390f35b60006001905090565b60006080606f565b90509056fea265627a7a7231582053247a3d8978a5a867d2d03b1d4152b8ba13842bcbfe1e1573a26068a0b6fa8464736f6c634300050d0032";

    public static final String FUNC_CALLGETDATE1 = "callGetDate1";

    public static final String FUNC_GETDATE = "getDate";

    @Deprecated
    protected InheritContractMutipleTest1(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected InheritContractMutipleTest1(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected InheritContractMutipleTest1(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected InheritContractMutipleTest1(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<BigInteger> callGetDate1() {
        final Function function = new Function(FUNC_CALLGETDATE1, 
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

    public static RemoteCall<InheritContractMutipleTest1> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(InheritContractMutipleTest1.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<InheritContractMutipleTest1> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(InheritContractMutipleTest1.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<InheritContractMutipleTest1> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(InheritContractMutipleTest1.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<InheritContractMutipleTest1> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(InheritContractMutipleTest1.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static InheritContractMutipleTest1 load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new InheritContractMutipleTest1(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static InheritContractMutipleTest1 load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new InheritContractMutipleTest1(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static InheritContractMutipleTest1 load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new InheritContractMutipleTest1(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static InheritContractMutipleTest1 load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new InheritContractMutipleTest1(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
