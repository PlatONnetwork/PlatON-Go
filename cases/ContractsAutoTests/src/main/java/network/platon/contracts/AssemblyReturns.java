package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.Bool;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Bytes2;
import org.web3j.abi.datatypes.generated.Bytes3;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.tuples.generated.Tuple5;
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
 * <p>Generated with web3j version 0.9.1.2-SNAPSHOT.
 */
public class AssemblyReturns extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506101c5806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c806326121ff014610030575b600080fd5b61003861011c565b60405180868152602001857dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001847cffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167cffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001831515151581526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019550505050505060405180910390f35b6000806000806000600294507fabcd00000000000000000000000000000000000000000000000000000000000093507f61626300000000000000000000000000000000000000000000000000000000009250600191507312121212121212121212121212121212121212129050909192939456fea265627a7a723158209e50878a732b0ff6c09591adf5828c6e5f8798d6756ac18f2589ca7ad1ad17fe64736f6c634300050d0032";

    public static final String FUNC_F = "f";

    @Deprecated
    protected AssemblyReturns(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected AssemblyReturns(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected AssemblyReturns(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected AssemblyReturns(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<Tuple5<BigInteger, byte[], byte[], Boolean, String>> f() {
        final Function function = new Function(FUNC_F, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Bytes2>() {}, new TypeReference<Bytes3>() {}, new TypeReference<Bool>() {}, new TypeReference<Address>() {}));
        return new RemoteCall<Tuple5<BigInteger, byte[], byte[], Boolean, String>>(
                new Callable<Tuple5<BigInteger, byte[], byte[], Boolean, String>>() {
                    @Override
                    public Tuple5<BigInteger, byte[], byte[], Boolean, String> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple5<BigInteger, byte[], byte[], Boolean, String>(
                                (BigInteger) results.get(0).getValue(), 
                                (byte[]) results.get(1).getValue(), 
                                (byte[]) results.get(2).getValue(), 
                                (Boolean) results.get(3).getValue(), 
                                (String) results.get(4).getValue());
                    }
                });
    }

    public static RemoteCall<AssemblyReturns> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(AssemblyReturns.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<AssemblyReturns> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(AssemblyReturns.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<AssemblyReturns> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(AssemblyReturns.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<AssemblyReturns> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(AssemblyReturns.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static AssemblyReturns load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new AssemblyReturns(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static AssemblyReturns load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new AssemblyReturns(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static AssemblyReturns load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new AssemblyReturns(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static AssemblyReturns load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new AssemblyReturns(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
