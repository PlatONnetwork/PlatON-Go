package network.platon.contracts;

import java.util.Arrays;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
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
 * <p>Generated with web3j version 0.13.0.7.
 */
public class InheritContractParentMutipleTest extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506101a6806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80631416d347146100515780636c3364ea1461009b578063853255cc146100e55780638da5cb5b146100ef575b600080fd5b610059610139565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100a3610141565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100ed610149565b005b6100f761014b565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b600033905090565b600033905090565b565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff168156fea265627a7a72315820e8959083650b4dfec8e6c89025fb538141443360592e0833afb8cbef5e8f1fd064736f6c634300050d0032";

    public static final String FUNC_GETADDRESSA = "getAddressA";

    public static final String FUNC_GETADDRESSB = "getAddressB";

    public static final String FUNC_OWNER = "owner";

    public static final String FUNC_SUM = "sum";

    protected InheritContractParentMutipleTest(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected InheritContractParentMutipleTest(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<String> getAddressA() {
        final Function function = new Function(FUNC_GETADDRESSA, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> getAddressB() {
        final Function function = new Function(FUNC_GETADDRESSB, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> owner() {
        final Function function = new Function(FUNC_OWNER, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public void sum() {
        throw new RuntimeException("cannot call constant function with void return type");
    }

    public static RemoteCall<InheritContractParentMutipleTest> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractParentMutipleTest.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<InheritContractParentMutipleTest> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractParentMutipleTest.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static InheritContractParentMutipleTest load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractParentMutipleTest(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static InheritContractParentMutipleTest load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractParentMutipleTest(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
