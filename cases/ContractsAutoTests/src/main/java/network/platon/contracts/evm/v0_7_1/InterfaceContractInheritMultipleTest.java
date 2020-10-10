package network.platon.contracts.evm.v0_7_1;

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
public class InterfaceContractInheritMultipleTest extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610119806100206000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c806399ecedf6146037578063cad0899b146080575b600080fd5b606a60048036036040811015604b57600080fd5b81019080803590602001909291908035906020019092919050505060c9565b6040518082815260200191505060405180910390f35b60b360048036036040811015609457600080fd5b81019080803590602001909291908035906020019092919050505060d6565b6040518082815260200191505060405180910390f35b6000818303905092915050565b600081830190509291505056fea26469706673582212204357f547d17b11ba0b8c2f2b9d1ed0b82ce4589b71f2916cf55e9a7aecd00fbb64736f6c63430007010033";

    public static final String FUNC_REDUCE = "reduce";

    public static final String FUNC_SUM = "sum";

    protected InterfaceContractInheritMultipleTest(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected InterfaceContractInheritMultipleTest(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> reduce(BigInteger c, BigInteger d) {
        final Function function = new Function(
                FUNC_REDUCE, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(c), 
                new com.alaya.abi.solidity.datatypes.generated.Uint256(d)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> sum(BigInteger a, BigInteger b) {
        final Function function = new Function(
                FUNC_SUM, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(a), 
                new com.alaya.abi.solidity.datatypes.generated.Uint256(b)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<InterfaceContractInheritMultipleTest> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InterfaceContractInheritMultipleTest.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<InterfaceContractInheritMultipleTest> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InterfaceContractInheritMultipleTest.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static InterfaceContractInheritMultipleTest load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new InterfaceContractInheritMultipleTest(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static InterfaceContractInheritMultipleTest load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new InterfaceContractInheritMultipleTest(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
