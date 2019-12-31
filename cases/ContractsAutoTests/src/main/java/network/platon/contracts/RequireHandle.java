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
public class RequireHandle extends Contract {
    private static final String BINARY = "6080604052734b0897b0513fdc7c541b6d9d7e929c4e5364d2db600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555034801561006557600080fd5b5060007314723a09acff6d2a60dcdf7aa4aff308fddc160c90508073ffffffffffffffffffffffffffffffffffffffff166108fc600a9081150290604051600060405180830381858888f193505050501580156100c6573d6000803e3d6000fd5b505061055a806100d76000396000f3fe6080604052600436106100705760003560e01c80635995caa71161004e5780635995caa7146101195780639ba1cf1114610130578063bfa34d5114610147578063fd3d30051461015e57610070565b80631e94691f146100e15780631f4c7d9c146100eb578063326f7d6714610102575b34801561007c57600080fd5b5060007314723a09acff6d2a60dcdf7aa4aff308fddc160c90508073ffffffffffffffffffffffffffffffffffffffff166108fc600a9081150290604051600060405180830381858888f193505050501580156100dd573d6000803e3d6000fd5b5050005b6100e9610175565b005b3480156100f757600080fd5b506101006101d9565b005b34801561010e57600080fd5b506101176101e9565b005b34801561012557600080fd5b5061012e61025d565b005b34801561013c57600080fd5b506101456102c7565b005b34801561015357600080fd5b5061015c61034e565b005b34801561016a57600080fd5b506101736103b2565b005b60007314723a09acff6d2a60dcdf7aa4aff308fddc160c90508073ffffffffffffffffffffffffffffffffffffffff166108fc600a9081150290604051600060405180830381858888f193505050501580156101d5573d6000803e3d6000fd5b5050565b60006001106101e757600080fd5b565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc680579a814e10a7400009081150290604051600060405180830381858888f1935050505015801561025a573d6000803e3d6000fd5b50565b60405161026990610460565b604051809103906000f080158015610285573d6000803e3d6000fd5b506000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16634d4310976103206040518263ffffffff1660e01b8152600401600060405180830381600088803b15801561033357600080fd5b5087f1158015610347573d6000803e3d6000fd5b5050505050565b60007314723a09acff6d2a60dcdf7aa4aff308fddc160c90508073ffffffffffffffffffffffffffffffffffffffff166108fc600a9081150290604051600060405180830381858888f193505050501580156103ae573d6000803e3d6000fd5b5050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663370158ea600a610320906040518363ffffffff1660e01b81526004016020604051808303818589803b15801561042057600080fd5b5088f1158015610434573d6000803e3d6000fd5b5050505050506040513d602081101561044c57600080fd5b810190808051906020019092919050505050565b60b98061046d8339019056fe6080604052348015600f57600080fd5b50609b8061001e6000396000f3fe60806040526004361060265760003560e01c8063370158ea14602b5780634d431097146047575b600080fd5b6031605b565b6040518082815260200191505060405180910390f35b348015605257600080fd5b5060596064565b005b6000602a905090565b56fea265627a7a72315820cbd2b14c9aa5642f69f1d77873065de8c094e0850300124092a7a0de5024099964736f6c634300050d0032a265627a7a723158208975a60e4c4bcbbc9021154023707e8e6c1b34bf54f0d6d30c94a52b4e91400164736f6c634300050d0032";

    public static final String FUNC_FUNCTIONCALLECECPTION = "functionCallEcecption";

    public static final String FUNC_NEWCONTRACTEXCEPTION = "newContractException";

    public static final String FUNC_NONPAYABLERECEIVEETHEXCEPTION = "nonPayableReceiveEthException";

    public static final String FUNC_OUTFUNCTIONCALLEXCEPTION = "outFunctionCallException";

    public static final String FUNC_PARAMEXCEPTION = "paramException";

    public static final String FUNC_PUBLICGETTERRECEIVEETHEXCEPTION = "publicGetterReceiveEthException";

    public static final String FUNC_TRANSFERCALLEXCEPTION = "transferCallException";

    @Deprecated
    protected RequireHandle(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected RequireHandle(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected RequireHandle(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected RequireHandle(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public static RemoteCall<RequireHandle> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(RequireHandle.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    public static RemoteCall<RequireHandle> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(RequireHandle.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<RequireHandle> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(RequireHandle.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<RequireHandle> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(RequireHandle.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    public RemoteCall<TransactionReceipt> functionCallEcecption() {
        final Function function = new Function(
                FUNC_FUNCTIONCALLECECPTION, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> newContractException() {
        final Function function = new Function(
                FUNC_NEWCONTRACTEXCEPTION, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> nonPayableReceiveEthException() {
        final Function function = new Function(
                FUNC_NONPAYABLERECEIVEETHEXCEPTION, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> outFunctionCallException() {
        final Function function = new Function(
                FUNC_OUTFUNCTIONCALLEXCEPTION, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> paramException() {
        final Function function = new Function(
                FUNC_PARAMEXCEPTION, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> publicGetterReceiveEthException() {
        final Function function = new Function(
                FUNC_PUBLICGETTERRECEIVEETHEXCEPTION, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> transferCallException(BigInteger weiValue) {
        final Function function = new Function(
                FUNC_TRANSFERCALLEXCEPTION, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, weiValue);
    }

    @Deprecated
    public static RequireHandle load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new RequireHandle(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static RequireHandle load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new RequireHandle(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static RequireHandle load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new RequireHandle(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static RequireHandle load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new RequireHandle(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
