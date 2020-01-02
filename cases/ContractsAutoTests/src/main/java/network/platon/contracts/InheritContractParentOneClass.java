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
public class InheritContractParentOneClass extends Contract {
    private static final String BINARY = "6080604052348015600f57600080fd5b5060878061001e6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063430fe9c114602d575b600080fd5b60336049565b6040518082815260200191505060405180910390f35b6000600190509056fea265627a7a72315820874e0944ac66cfdc955eb65560745a66e12e87085465119ddd5c26a6e0df5c6764736f6c634300050d0032";

    public static final String FUNC_GETDATE = "getDate";

    @Deprecated
    protected InheritContractParentOneClass(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected InheritContractParentOneClass(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected InheritContractParentOneClass(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected InheritContractParentOneClass(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<BigInteger> getDate() {
        final Function function = new Function(FUNC_GETDATE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<InheritContractParentOneClass> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(InheritContractParentOneClass.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<InheritContractParentOneClass> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(InheritContractParentOneClass.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<InheritContractParentOneClass> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(InheritContractParentOneClass.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<InheritContractParentOneClass> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(InheritContractParentOneClass.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static InheritContractParentOneClass load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new InheritContractParentOneClass(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static InheritContractParentOneClass load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new InheritContractParentOneClass(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static InheritContractParentOneClass load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new InheritContractParentOneClass(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static InheritContractParentOneClass load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new InheritContractParentOneClass(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
