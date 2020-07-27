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
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.0.7.
 */
public class ErrorParamConstructor extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506040516020806102708339810180604052602081101561003057600080fd5b8101908080519060200190929190505050600a8060008190555050806001819055505061020e806100626000396000f3fe60806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630dbe671f146100725780634df7e3d01461009d57806382ab890a146100c8578063a1c519151461014a578063d46300fd14610175575b600080fd5b34801561007e57600080fd5b506100876101a0565b6040518082815260200191505060405180910390f35b3480156100a957600080fd5b506100b26101a6565b6040518082815260200191505060405180910390f35b3480156100d457600080fd5b50610101600480360360208110156100eb57600080fd5b81019080803590602001909291905050506101ac565b604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390f35b34801561015657600080fd5b5061015f6101cf565b6040518082815260200191505060405180910390f35b34801561018157600080fd5b5061018a6101d9565b6040518082815260200191505060405180910390f35b60005481565b60015481565b600080826001600082825401925050819055503360015481915091509150915091565b6000600154905090565b6000805490509056fea165627a7a723058208280a4d3318c5ddce1f4c63b76536ed7648a15a9de2ccfcdfe42e30573e9f5b20029";

    public static final String FUNC_A = "a";

    public static final String FUNC_B = "b";

    public static final String FUNC_UPDATE = "update";

    public static final String FUNC_GETB = "getB";

    public static final String FUNC_GETA = "getA";

    protected ErrorParamConstructor(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ErrorParamConstructor(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> a() {
        final Function function = new Function(FUNC_A, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> b() {
        final Function function = new Function(FUNC_B, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> update(BigInteger amount) {
        final Function function = new Function(
                FUNC_UPDATE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(amount)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> getB() {
        final Function function = new Function(FUNC_GETB, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getA() {
        final Function function = new Function(FUNC_GETA, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<ErrorParamConstructor> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId, BigInteger _b) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_b)));
        return deployRemoteCall(ErrorParamConstructor.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static RemoteCall<ErrorParamConstructor> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId, BigInteger _b) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_b)));
        return deployRemoteCall(ErrorParamConstructor.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static ErrorParamConstructor load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ErrorParamConstructor(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ErrorParamConstructor load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ErrorParamConstructor(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
