package network.platon.contracts;

import java.math.BigInteger;
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
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.8-SNAPSHOT.
 */
public class ChainFunction extends Contract {
    private static final String BINARY = "6080604052336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550346001819055506000600260006101000a81548160ff0219169083151502179055506101c7806100756000396000f3fe60806040526004361061004c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680637eed92c0146100515780639f9232f4146100ce575b600080fd5b34801561005d57600080fd5b5061008c6004803603602081101561007457600080fd5b81019080803515159060200190929190505050610155565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156100da57600080fd5b50610113600480360360408110156100f157600080fd5b8101908080351515906020019092919080359060200190929190505050610172565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b60006001151582151514151561016a57600080fd5b339050919050565b60006001151583151514151561018457fe5b600982101561019257600080fd5b3390509291505056fea165627a7a723058201a87536874081d756ebc8e1ac473bf0d7af3a13a43c599f78b11c7ba9a8d93850029";

    public static final String FUNC_DECEASEDWITHMODIFY = "deceasedWithModify";

    public static final String FUNC_DECEASED = "deceased";

    @Deprecated
    protected ChainFunction(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected ChainFunction(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected ChainFunction(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected ChainFunction(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<String> deceasedWithModify(Boolean _isDeceased) {
        final Function function = new Function(FUNC_DECEASEDWITHMODIFY, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Bool(_isDeceased)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> deceased(Boolean isDeceased, BigInteger less9) {
        final Function function = new Function(FUNC_DECEASED, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Bool(isDeceased), 
                new org.web3j.abi.datatypes.generated.Uint256(less9)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public static RemoteCall<ChainFunction> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, BigInteger initialWeiValue) {
        return deployRemoteCall(ChainFunction.class, web3j, credentials, contractGasProvider, BINARY, "", initialWeiValue);
    }

    public static RemoteCall<ChainFunction> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, BigInteger initialWeiValue) {
        return deployRemoteCall(ChainFunction.class, web3j, transactionManager, contractGasProvider, BINARY, "", initialWeiValue);
    }

    @Deprecated
    public static RemoteCall<ChainFunction> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit, BigInteger initialWeiValue) {
        return deployRemoteCall(ChainFunction.class, web3j, credentials, gasPrice, gasLimit, BINARY, "", initialWeiValue);
    }

    @Deprecated
    public static RemoteCall<ChainFunction> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit, BigInteger initialWeiValue) {
        return deployRemoteCall(ChainFunction.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "", initialWeiValue);
    }

    @Deprecated
    public static ChainFunction load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new ChainFunction(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static ChainFunction load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new ChainFunction(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static ChainFunction load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new ChainFunction(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static ChainFunction load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new ChainFunction(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
