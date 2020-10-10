package network.platon.contracts.evm.v0_5_17;

import com.alaya.abi.solidity.FunctionEncoder;
import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.Arrays;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.0.
 */
public class NoParenthesesConstructor extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506040516100e13803806100e18339818101604052602081101561003357600080fd5b8101908080519060200190929190505050600180600081905550505060848061005d6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80630dbe671f14602d575b600080fd5b60336049565b6040518082815260200191505060405180910390f35b6000548156fea265627a7a723158207548c59930fb3f925f3189b697207ac8f4b41af0143f6b40227a5c0200b39e7b64736f6c63430005110032";

    public static final String FUNC_A = "a";

    protected NoParenthesesConstructor(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected NoParenthesesConstructor(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<NoParenthesesConstructor> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId, BigInteger _b) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new Uint256(_b)));
        return deployRemoteCall(NoParenthesesConstructor.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static RemoteCall<NoParenthesesConstructor> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId, BigInteger _b) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new Uint256(_b)));
        return deployRemoteCall(NoParenthesesConstructor.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public RemoteCall<BigInteger> a() {
        final Function function = new Function(FUNC_A, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static NoParenthesesConstructor load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new NoParenthesesConstructor(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static NoParenthesesConstructor load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new NoParenthesesConstructor(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
