package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import org.web3j.abi.FunctionEncoder;
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
public class InheritContractSub2 extends Contract {
    private static final String BINARY = "60806040526000805534801561001457600080fd5b506040516100f83803806100f88339818101604052602081101561003757600080fd5b810190808051906020019092919050505060028060008190555050600190506000548101600181905550506088806100706000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80633bc5de3014602d575b600080fd5b60336049565b6040518082815260200191505060405180910390f35b600060015490509056fea265627a7a72315820e07df224cb231fd280847c33fd2ea70d2c2a24eac13a774a9a981776224b775164736f6c634300050d0032";

    public static final String FUNC_GETDATA = "getData";

    @Deprecated
    protected InheritContractSub2(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected InheritContractSub2(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected InheritContractSub2(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected InheritContractSub2(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public static RemoteCall<InheritContractSub2> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, BigInteger y) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(y)));
        return deployRemoteCall(InheritContractSub2.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor);
    }

    public static RemoteCall<InheritContractSub2> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, BigInteger y) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(y)));
        return deployRemoteCall(InheritContractSub2.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor);
    }

    @Deprecated
    public static RemoteCall<InheritContractSub2> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit, BigInteger y) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(y)));
        return deployRemoteCall(InheritContractSub2.class, web3j, credentials, gasPrice, gasLimit, BINARY, encodedConstructor);
    }

    @Deprecated
    public static RemoteCall<InheritContractSub2> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit, BigInteger y) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(y)));
        return deployRemoteCall(InheritContractSub2.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, encodedConstructor);
    }

    public RemoteCall<BigInteger> getData() {
        final Function function = new Function(FUNC_GETDATA, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    @Deprecated
    public static InheritContractSub2 load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new InheritContractSub2(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static InheritContractSub2 load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new InheritContractSub2(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static InheritContractSub2 load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new InheritContractSub2(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static InheritContractSub2 load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new InheritContractSub2(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
