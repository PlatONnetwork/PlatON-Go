package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import org.web3j.abi.FunctionEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
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
public class NoParenthesesConstructor extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506040516020806100f88339810180604052602081101561003057600080fd5b81019080805190602001909291905050506001806000819055505050609e8061005a6000396000f3fe608060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630dbe671f146044575b600080fd5b348015604f57600080fd5b506056606c565b6040518082815260200191505060405180910390f35b6000548156fea165627a7a723058203b932c211fa73d548ce6470091ff4fa4e7c835108bf437abb0bdf4bf8c99437f0029";

    public static final String FUNC_A = "a";

    protected NoParenthesesConstructor(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected NoParenthesesConstructor(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> a() {
        final Function function = new Function(FUNC_A, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<NoParenthesesConstructor> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId, BigInteger _b) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_b)));
        return deployRemoteCall(NoParenthesesConstructor.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static RemoteCall<NoParenthesesConstructor> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId, BigInteger _b) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_b)));
        return deployRemoteCall(NoParenthesesConstructor.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static NoParenthesesConstructor load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new NoParenthesesConstructor(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static NoParenthesesConstructor load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new NoParenthesesConstructor(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
