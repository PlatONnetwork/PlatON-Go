package network.platon.contracts.evm.v0_7_1;

import com.alaya.abi.solidity.FunctionEncoder;
import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.0.
 */
public class ConstructorInternalDeclaraction extends Contract {
    private static final String BINARY = "6080604052600a60005534801561001557600080fd5b506040516100e43803806100e48339818101604052602081101561003857600080fd5b8101908080519060200190929190505050806000819055505060858061005f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c806306661abd14602d575b600080fd5b60336049565b6040518082815260200191505060405180910390f35b6000548156fea264697066735822122008b699b37ad9c00ce7e27f22d232a3e3f38ba50d972798e06c2a2a032cd8c76464736f6c63430007010033";

    public static final String FUNC_COUNT = "count";

    protected ConstructorInternalDeclaraction(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ConstructorInternalDeclaraction(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<ConstructorInternalDeclaraction> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId, BigInteger _count) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(_count)));
        return deployRemoteCall(ConstructorInternalDeclaraction.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static RemoteCall<ConstructorInternalDeclaraction> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId, BigInteger _count) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(_count)));
        return deployRemoteCall(ConstructorInternalDeclaraction.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public RemoteCall<TransactionReceipt> count() {
        final Function function = new Function(
                FUNC_COUNT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static ConstructorInternalDeclaraction load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ConstructorInternalDeclaraction(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ConstructorInternalDeclaraction load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ConstructorInternalDeclaraction(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
