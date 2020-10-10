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
public class InheritContractOverload extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506101aa806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80635873f056146100515780639450268b1461006f578063a46cf4b3146100c5578063cad0899b146100e3575b600080fd5b61005961012f565b6040518082815260200191505060405180910390f35b6100af6004803603606081101561008557600080fd5b81019080803590602001909291908035906020019092919080359060200190929190505050610144565b6040518082815260200191505060405180910390f35b6100cd610154565b6040518082815260200191505060405180910390f35b610119600480360360408110156100f957600080fd5b810190808035906020019092919080359060200190929190505050610167565b6040518082815260200191505060405180910390f35b600061013f600160026003610144565b905090565b6000818385010190509392505050565b600061016260016002610167565b905090565b600081830190509291505056fea2646970667358221220c323a35ca79886008f12216ff452b88f123cf03b20b4b0917086041e471652f664736f6c63430007010033";

    public static final String FUNC_GETDATAA = "getDataA";

    public static final String FUNC_GETDATAB = "getDataB";

    public static final String FUNC_SUM = "sum";

    protected InheritContractOverload(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected InheritContractOverload(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getDataA() {
        final Function function = new Function(
                FUNC_GETDATAA, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getDataB() {
        final Function function = new Function(
                FUNC_GETDATAB, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> sum(BigInteger a, BigInteger b, BigInteger c) {
        final Function function = new Function(
                FUNC_SUM, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(a), 
                new com.alaya.abi.solidity.datatypes.generated.Uint256(b), 
                new com.alaya.abi.solidity.datatypes.generated.Uint256(c)), 
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

    public static RemoteCall<InheritContractOverload> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractOverload.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<InheritContractOverload> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractOverload.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static InheritContractOverload load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractOverload(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static InheritContractOverload load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractOverload(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
