package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.FunctionEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
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
public class ConstructorPublicDeclaraction extends Contract {
    private static final String BINARY = "60806040526000805534801561001457600080fd5b506040516020806101b18339810180604052602081101561003457600080fd5b810190808051906020019092919050505080600081905550506101558061005c6000396000f3fe60806040526004361061004c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806382ab890a14610051578063a87d942c146100d3575b600080fd5b34801561005d57600080fd5b5061008a6004803603602081101561007457600080fd5b81019080803590602001909291905050506100fe565b604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390f35b3480156100df57600080fd5b506100e8610120565b6040518082815260200191505060405180910390f35b6000808260008082825401925050819055503360005481915091509150915091565b6000805490509056fea165627a7a7230582009a531d64713af46bc52bf05b3bed5791766a59d980ac13ec78a043aaa23ffd10029";

    public static final String FUNC_UPDATE = "update";

    public static final String FUNC_GETCOUNT = "getCount";

    @Deprecated
    protected ConstructorPublicDeclaraction(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected ConstructorPublicDeclaraction(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected ConstructorPublicDeclaraction(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected ConstructorPublicDeclaraction(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> update(BigInteger amount) {
        final Function function = new Function(
                FUNC_UPDATE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(amount)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> getCount() {
        final Function function = new Function(FUNC_GETCOUNT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<ConstructorPublicDeclaraction> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, BigInteger _count) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_count)));
        return deployRemoteCall(ConstructorPublicDeclaraction.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor);
    }

    public static RemoteCall<ConstructorPublicDeclaraction> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, BigInteger _count) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_count)));
        return deployRemoteCall(ConstructorPublicDeclaraction.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor);
    }

    @Deprecated
    public static RemoteCall<ConstructorPublicDeclaraction> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit, BigInteger _count) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_count)));
        return deployRemoteCall(ConstructorPublicDeclaraction.class, web3j, credentials, gasPrice, gasLimit, BINARY, encodedConstructor);
    }

    @Deprecated
    public static RemoteCall<ConstructorPublicDeclaraction> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit, BigInteger _count) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_count)));
        return deployRemoteCall(ConstructorPublicDeclaraction.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, encodedConstructor);
    }

    @Deprecated
    public static ConstructorPublicDeclaraction load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new ConstructorPublicDeclaraction(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static ConstructorPublicDeclaraction load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new ConstructorPublicDeclaraction(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static ConstructorPublicDeclaraction load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new ConstructorPublicDeclaraction(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static ConstructorPublicDeclaraction load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new ConstructorPublicDeclaraction(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
