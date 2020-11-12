package network.platon.contracts.evm;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.FunctionEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.tuples.generated.Tuple2;
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
 * <p>Generated with web3j version 0.13.1.5.
 */
public class ConstructorPublicVisibility extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506040516020806102c38339810180604052810190808051906020019092919050505080600181905550506102798061004a6000396000f300608060405260043610610078576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806321687d5c1461007d57806335f646c0146100a857806383f2ac99146100d357806388b5383c1461011e578063bfe9ef4f14610150578063cc80f6f31461017b575b600080fd5b34801561008957600080fd5b506100926101a6565b6040518082815260200191505060405180910390f35b3480156100b457600080fd5b506100bd6101b0565b6040518082815260200191505060405180910390f35b3480156100df57600080fd5b5061010860048036038101908080359060200190929190803590602001909291905050506101ba565b6040518082815260200191505060405180910390f35b34801561012a57600080fd5b506101336101c7565b604051808381526020018281526020019250505060405180910390f35b34801561015c57600080fd5b50610165610226565b6040518082815260200191505060405180910390f35b34801561018757600080fd5b50610190610245565b6040518082815260200191505060405180910390f35b6000600154905090565b6000600154905090565b6000818301905092915050565b600080600080600080600080600060019650600095506001945086156101ec57600193505b5b85156101fc57600192506101ed565b600091505b600282101561021b57600190508180600101925050610201565b505050505050509091565b6000607b60008190555060005460015401600181905550600154905090565b6000349050905600a165627a7a723058205bfb2953c6c4670ded5194fa79c4309775438d379cf216d6c304cf5ca3b214bd0029";

    public static final String FUNC_CONSTANTCHECK = "constantCheck";

    public static final String FUNC_GETOUTI = "getOutI";

    public static final String FUNC_NAMEDRETURN = "namedReturn";

    public static final String FUNC_GRAMMARCHECK = "grammarCheck";

    public static final String FUNC_ABSTRACTFUNCTION = "abstractFunction";

    public static final String FUNC_SHOW = "show";

    protected ConstructorPublicVisibility(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ConstructorPublicVisibility(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> constantCheck() {
        final Function function = new Function(FUNC_CONSTANTCHECK, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getOutI() {
        final Function function = new Function(FUNC_GETOUTI, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> namedReturn(BigInteger a, BigInteger b) {
        final Function function = new Function(FUNC_NAMEDRETURN, 
                Arrays.<Type>asList(new Uint256(a),
                new Uint256(b)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<Tuple2<BigInteger, BigInteger>> grammarCheck() {
        final Function function = new Function(FUNC_GRAMMARCHECK, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}));
        return new RemoteCall<Tuple2<BigInteger, BigInteger>>(
                new Callable<Tuple2<BigInteger, BigInteger>>() {
                    @Override
                    public Tuple2<BigInteger, BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<BigInteger, BigInteger>(
                                (BigInteger) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue());
                    }
                });
    }

    public RemoteCall<BigInteger> abstractFunction() {
        final Function function = new Function(FUNC_ABSTRACTFUNCTION, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> show() {
        final Function function = new Function(FUNC_SHOW, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<ConstructorPublicVisibility> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId, BigInteger _y) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new Uint256(_y)));
        return deployRemoteCall(ConstructorPublicVisibility.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static RemoteCall<ConstructorPublicVisibility> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId, BigInteger _y) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new Uint256(_y)));
        return deployRemoteCall(ConstructorPublicVisibility.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static ConstructorPublicVisibility load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ConstructorPublicVisibility(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ConstructorPublicVisibility load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ConstructorPublicVisibility(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
