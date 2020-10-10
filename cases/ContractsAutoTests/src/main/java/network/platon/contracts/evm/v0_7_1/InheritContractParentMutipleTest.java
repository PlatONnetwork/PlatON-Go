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
public class InheritContractParentMutipleTest extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610165806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80631416d347146100515780636c3364ea14610085578063853255cc146100b95780638da5cb5b146100c3575b600080fd5b6100596100f7565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61008d6100ff565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100c1610107565b005b6100cb610109565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b600033905090565b600033905090565b565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff168156fea26469706673582212203998131d6105b7bd0f66d78ebb744e370651fb4fa2f6fc730bda9189be5407b364736f6c63430007010033";

    public static final String FUNC_GETADDRESSA = "getAddressA";

    public static final String FUNC_GETADDRESSB = "getAddressB";

    public static final String FUNC_OWNER = "owner";

    public static final String FUNC_SUM = "sum";

    protected InheritContractParentMutipleTest(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected InheritContractParentMutipleTest(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getAddressA() {
        final Function function = new Function(
                FUNC_GETADDRESSA, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getAddressB() {
        final Function function = new Function(
                FUNC_GETADDRESSB, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> owner() {
        final Function function = new Function(
                FUNC_OWNER, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> sum() {
        final Function function = new Function(
                FUNC_SUM, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<InheritContractParentMutipleTest> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractParentMutipleTest.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<InheritContractParentMutipleTest> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractParentMutipleTest.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static InheritContractParentMutipleTest load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractParentMutipleTest(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static InheritContractParentMutipleTest load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractParentMutipleTest(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
