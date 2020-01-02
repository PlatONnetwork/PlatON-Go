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
public class InheritContractOverload extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506101a9806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80639450268b146100515780639944cc71146100a7578063a898fd70146100c5578063cad0899b146100e3575b600080fd5b6100916004803603606081101561006757600080fd5b8101908080359060200190929190803590602001909291908035906020019092919050505061012f565b6040518082815260200191505060405180910390f35b6100af61013f565b6040518082815260200191505060405180910390f35b6100cd610152565b6040518082815260200191505060405180910390f35b610119600480360360408110156100f957600080fd5b810190808035906020019092919080359060200190929190505050610167565b6040518082815260200191505060405180910390f35b6000818385010190509392505050565b600061014d60016002610167565b905090565b600061016260016002600361012f565b905090565b600081830190509291505056fea265627a7a72315820c0eb8b9c622d2758d02de2c08c8d7ed29ca96a784cd3b1df18208394bc8926c464736f6c634300050d0032";

    public static final String FUNC_GETDATA1 = "getData1";

    public static final String FUNC_GETDATA2 = "getData2";

    public static final String FUNC_SUM = "sum";

    @Deprecated
    protected InheritContractOverload(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected InheritContractOverload(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected InheritContractOverload(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected InheritContractOverload(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<BigInteger> getData1() {
        final Function function = new Function(FUNC_GETDATA1, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getData2() {
        final Function function = new Function(FUNC_GETDATA2, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> sum(BigInteger a, BigInteger b, BigInteger c) {
        final Function function = new Function(FUNC_SUM, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(a), 
                new org.web3j.abi.datatypes.generated.Uint256(b), 
                new org.web3j.abi.datatypes.generated.Uint256(c)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> sum(BigInteger a, BigInteger b) {
        final Function function = new Function(FUNC_SUM, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(a), 
                new org.web3j.abi.datatypes.generated.Uint256(b)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<InheritContractOverload> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(InheritContractOverload.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<InheritContractOverload> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(InheritContractOverload.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<InheritContractOverload> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(InheritContractOverload.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<InheritContractOverload> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(InheritContractOverload.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static InheritContractOverload load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new InheritContractOverload(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static InheritContractOverload load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new InheritContractOverload(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static InheritContractOverload load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new InheritContractOverload(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static InheritContractOverload load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new InheritContractOverload(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
