package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
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
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.0.
 */
public class AddressFunctions extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061025c806100206000396000f3fe60806040526004361061003f5760003560e01c80631a69523014610044578063ecbde5e614610088578063ee082a0c146100b3578063f8b2cb4f146100f7575b600080fd5b6100866004803603602081101561005a57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061015c565b005b34801561009457600080fd5b5061009d6101a6565b6040518082815260200191505060405180910390f35b6100f5600480360360208110156100c957600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506101c5565b005b34801561010357600080fd5b506101466004803603602081101561011a57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610206565b6040518082815260200191505060405180910390f35b8073ffffffffffffffffffffffffffffffffffffffff166108fc349081150290604051600060405180830381858888f193505050501580156101a2573d6000803e3d6000fd5b5050565b60003073ffffffffffffffffffffffffffffffffffffffff1631905090565b8073ffffffffffffffffffffffffffffffffffffffff166108fc678ac7230489e800009081150290604051600060405180830381858888f193505050505050565b60008173ffffffffffffffffffffffffffffffffffffffff1631905091905056fea265627a7a723158200379d83f8bae39e3bea0958218de396284158740bc1a41134670aed62b3d002664736f6c634300050d0032";

    public static final String FUNC_GETBALANCE = "getBalance";

    public static final String FUNC_GETBALANCEOF = "getBalanceOf";

    public static final String FUNC_TRANSFER = "transfer";

    public static final String FUNC_TRANSFER4 = "transfer4";

    @Deprecated
    protected AddressFunctions(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected AddressFunctions(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected AddressFunctions(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected AddressFunctions(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> getBalance(String addr) {
        final Function function = new Function(
                FUNC_GETBALANCE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(addr)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getBalanceOf() {
        final Function function = new Function(
                FUNC_GETBALANCEOF, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> transfer(String addr, BigInteger weiValue) {
        final Function function = new Function(
                FUNC_TRANSFER, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(addr)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, weiValue);
    }

    public RemoteCall<TransactionReceipt> transfer4(String addr, BigInteger weiValue) {
        final Function function = new Function(
                FUNC_TRANSFER4, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(addr)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, weiValue);
    }

    public static RemoteCall<AddressFunctions> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(AddressFunctions.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<AddressFunctions> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(AddressFunctions.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<AddressFunctions> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(AddressFunctions.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<AddressFunctions> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(AddressFunctions.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static AddressFunctions load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new AddressFunctions(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static AddressFunctions load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new AddressFunctions(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static AddressFunctions load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new AddressFunctions(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static AddressFunctions load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new AddressFunctions(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
