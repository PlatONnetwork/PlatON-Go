package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.tuples.generated.Tuple3;
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
public class SameNameConstructorInternalVisibility extends Contract {
    private static final String BINARY = "6060604052341561000f57600080fd5b60d48061001d6000396000f300606060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063bb8220ea146044575b600080fd5b3415604e57600080fd5b60546078565b60405180848152602001838152602001828152602001935050505060405180910390f35b6000806000806000806000600180905080945050600190508383828060ff169050965096509650505050509091925600a165627a7a72305820aa0f83409d03beeb7be2c4378f747e8401f3f56acccd4ba8ed887bde2a938c430029";

    public static final String FUNC_DISCARDVARIABLE = "discardVariable";

    @Deprecated
    protected SameNameConstructorInternalVisibility(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected SameNameConstructorInternalVisibility(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected SameNameConstructorInternalVisibility(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected SameNameConstructorInternalVisibility(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<Tuple3<BigInteger, BigInteger, BigInteger>> discardVariable() {
        final Function function = new Function(FUNC_DISCARDVARIABLE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}));
        return new RemoteCall<Tuple3<BigInteger, BigInteger, BigInteger>>(
                new Callable<Tuple3<BigInteger, BigInteger, BigInteger>>() {
                    @Override
                    public Tuple3<BigInteger, BigInteger, BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple3<BigInteger, BigInteger, BigInteger>(
                                (BigInteger) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue(), 
                                (BigInteger) results.get(2).getValue());
                    }
                });
    }

    public static RemoteCall<SameNameConstructorInternalVisibility> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(SameNameConstructorInternalVisibility.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<SameNameConstructorInternalVisibility> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(SameNameConstructorInternalVisibility.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<SameNameConstructorInternalVisibility> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(SameNameConstructorInternalVisibility.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<SameNameConstructorInternalVisibility> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(SameNameConstructorInternalVisibility.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static SameNameConstructorInternalVisibility load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new SameNameConstructorInternalVisibility(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static SameNameConstructorInternalVisibility load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new SameNameConstructorInternalVisibility(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static SameNameConstructorInternalVisibility load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new SameNameConstructorInternalVisibility(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static SameNameConstructorInternalVisibility load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new SameNameConstructorInternalVisibility(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
