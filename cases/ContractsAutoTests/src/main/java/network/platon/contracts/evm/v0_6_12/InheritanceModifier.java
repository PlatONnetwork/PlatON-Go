package network.platon.contracts.evm.v0_6_12;

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
public class InheritanceModifier extends Contract {
    private static final String BINARY = "6080604052600a60005534801561001557600080fd5b5060d3806100246000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80632e2619a3146037578063d46300fd146062575b600080fd5b606060048036036020811015604b57600080fd5b8101908080359060200190929190505050607e565b005b60686094565b6040518082815260200191505060405180910390f35b80600054811060905760016000819055505b5050565b6000805490509056fea26469706673582212209eeb513dc31107c5d7d8dfe126337814ac74d0067e6e5d57da4f9b9cef6613ea64736f6c634300060c0033";

    public static final String FUNC_GETA = "getA";

    public static final String FUNC_INHERITANCE = "inheritance";

    protected InheritanceModifier(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected InheritanceModifier(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getA() {
        final Function function = new Function(
                FUNC_GETA, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> inheritance(BigInteger c) {
        final Function function = new Function(
                FUNC_INHERITANCE, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(c)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<InheritanceModifier> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritanceModifier.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<InheritanceModifier> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritanceModifier.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static InheritanceModifier load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new InheritanceModifier(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static InheritanceModifier load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new InheritanceModifier(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
