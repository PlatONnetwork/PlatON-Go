package network.platon.contracts.evm.v0_6_12;

import com.alaya.abi.solidity.FunctionEncoder;
import com.alaya.abi.solidity.datatypes.Type;
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
public class InheritContractBase extends Contract {
    private static final String BINARY = "608060405260008055348015601357600080fd5b506040516098380380609883398181016040526020811015603357600080fd5b81019080805190602001909291905050508060008190555050603f8060596000396000f3fe6080604052600080fdfea2646970667358221220d610e1bbf67c3cf40869b9a86f7ccc166c194502baf92f2c37bfde3e9830463e64736f6c634300060c0033";

    protected InheritContractBase(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected InheritContractBase(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<InheritContractBase> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId, BigInteger x) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(x)));
        return deployRemoteCall(InheritContractBase.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static RemoteCall<InheritContractBase> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId, BigInteger x) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(x)));
        return deployRemoteCall(InheritContractBase.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static InheritContractBase load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractBase(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static InheritContractBase load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractBase(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
