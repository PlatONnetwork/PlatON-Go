package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint16;
import org.web3j.abi.datatypes.generated.Uint64;
import org.web3j.abi.datatypes.generated.Uint8;
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
 * <p>Generated with web3j version 0.7.5.3-SNAPSHOT.
 */
public class TypeConversionBytesToUintContract extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506101ab806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80634e9189bc146100465780636ab281811461006c5780638acc06e01461009e575b600080fd5b61004e6100c2565b604051808261ffff1661ffff16815260200191505060405180910390f35b6100746100fe565b604051808267ffffffffffffffff1667ffffffffffffffff16815260200191505060405180910390f35b6100a6610140565b604051808260ff1660ff16815260200191505060405180910390f35b6000807f6162636400000000000000000000000000000000000000000000000000000000905060008160e01c9050600081905080935050505090565b6000807f6162636400000000000000000000000000000000000000000000000000000000905060008160e01c905060008163ffffffff16905080935050505090565b6000807f6100000000000000000000000000000000000000000000000000000000000000905060008160f81c905080925050509056fea265627a7a72315820eecbae4c0e6b43a29204fdc11fcdc7674f75108c47cae7bc21c59264b017958064736f6c634300050d0032";

    public static final String FUNC_BYTESTOBIGUINT = "bytesToBigUint";

    public static final String FUNC_BYTESTOSAMEUINT = "bytesToSameUint";

    public static final String FUNC_BYTESTOSMALLUINT = "bytesToSmallUint";

    @Deprecated
    protected TypeConversionBytesToUintContract(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected TypeConversionBytesToUintContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected TypeConversionBytesToUintContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected TypeConversionBytesToUintContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<BigInteger> bytesToBigUint() {
        final Function function = new Function(FUNC_BYTESTOBIGUINT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint64>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> bytesToSameUint() {
        final Function function = new Function(FUNC_BYTESTOSAMEUINT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> bytesToSmallUint() {
        final Function function = new Function(FUNC_BYTESTOSMALLUINT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint16>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<TypeConversionBytesToUintContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(TypeConversionBytesToUintContract.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<TypeConversionBytesToUintContract> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(TypeConversionBytesToUintContract.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<TypeConversionBytesToUintContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(TypeConversionBytesToUintContract.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<TypeConversionBytesToUintContract> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(TypeConversionBytesToUintContract.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static TypeConversionBytesToUintContract load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new TypeConversionBytesToUintContract(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static TypeConversionBytesToUintContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new TypeConversionBytesToUintContract(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static TypeConversionBytesToUintContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new TypeConversionBytesToUintContract(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static TypeConversionBytesToUintContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new TypeConversionBytesToUintContract(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
