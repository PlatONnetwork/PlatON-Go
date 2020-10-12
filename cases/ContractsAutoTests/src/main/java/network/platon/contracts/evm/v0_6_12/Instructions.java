package network.platon.contracts.evm.v0_6_12;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
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
 * <p>Generated with web3j version 0.13.2.1.
 */
public class Instructions extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610203806100206000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c8063048a5fed1461005c578063165c4a161461007a5780633408e470146100c65780635a0db89e146100e4578063eb8ac92114610130575b600080fd5b61006461017c565b6040518082815260200191505060405180910390f35b6100b06004803603604081101561009057600080fd5b810190808035906020019092919080359060200190929190505050610189565b6040518082815260200191505060405180910390f35b6100ce610196565b6040518082815260200191505060405180910390f35b61011a600480360360408110156100fa57600080fd5b8101908080359060200190929190803590602001909291905050506101a3565b6040518082815260200191505060405180910390f35b6101666004803603604081101561014657600080fd5b8101908080359060200190929190803590602001909291905050506101b7565b6040518082815260200191505060405180910390f35b6000804790508091505090565b6000818302905092915050565b6000804690508091505090565b60006101af8383610189565b905092915050565b60006101c5600260036101a3565b90509291505056fea2646970667358221220fec7dd0754fcd3813c7f8e02ef653f4728dd727b6326880bffc70cfa0bff345c64736f6c634300060c0033";

    public static final String FUNC_GETCHAINID = "getChainId";

    public static final String FUNC_GETSELFBALANCE = "getSelfBalance";

    public static final String FUNC_MULTIPLY = "multiply";

    public static final String FUNC_TEST = "test";

    public static final String FUNC_TEST_MUL = "test_mul";

    protected Instructions(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected Instructions(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> getChainId() {
        final Function function = new Function(FUNC_GETCHAINID, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getSelfBalance() {
        final Function function = new Function(FUNC_GETSELFBALANCE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> multiply(BigInteger x, BigInteger y) {
        final Function function = new Function(
                FUNC_MULTIPLY, 
                Arrays.<Type>asList(new Uint256(x),
                new Uint256(y)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> test(BigInteger x, BigInteger y) {
        final Function function = new Function(
                FUNC_TEST, 
                Arrays.<Type>asList(new Uint256(x),
                new Uint256(y)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> test_mul(BigInteger x, BigInteger y) {
        final Function function = new Function(
                FUNC_TEST_MUL, 
                Arrays.<Type>asList(new Uint256(x),
                new Uint256(y)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<Instructions> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Instructions.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<Instructions> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Instructions.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static Instructions load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new Instructions(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static Instructions load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new Instructions(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
