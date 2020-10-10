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
public class SumAssembly extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50600060019080600181540180825580915050600190039060005260206000200160009091909190915055600060029080600181540180825580915050600190039060005260206000200160009091909190915055600060039080600181540180825580915050600190039060005260206000200160009091909190915055600060049080600181540180825580915050600190039060005260206000200160009091909190915055600060059080600181540180825580915050600190039060005260206000200160009091909190915055610154806100f26000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063853255cc14610030575b600080fd5b61003861004e565b6040518082815260200191505060405180910390f35b600073__$bd55ce7c8d444a52551ddc65d373d70380$__6387fbcc7760006040518263ffffffff1660e01b8152600401808060200182810382528381815481526020019150805480156100c057602002820191906000526020600020905b8154815260200190600101908083116100ac575b50509250505060206040518083038186803b1580156100de57600080fd5b505af41580156100f2573d6000803e3d6000fd5b505050506040513d602081101561010857600080fd5b810190808051906020019092919050505090509056fea264697066735822122021c72b864b4a8a7d58046d81d2b4d64edf91072b01bd76a6e960cfc04969cb2b64736f6c634300060c0033\n"
            + "\n"
            + "// $bd55ce7c8d444a52551ddc65d373d70380$ -> /home/platon/.jenkins/workspace/contracts_test_alaya/cases/ContractsAutoTests/src/test/resources/contracts/evm/0.6.12/1.function/07Assembly/SumAssembly.sol:Sum";

    public static final String FUNC_SUM = "sum";

    protected SumAssembly(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected SumAssembly(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<SumAssembly> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(SumAssembly.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<SumAssembly> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(SumAssembly.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public RemoteCall<TransactionReceipt> sum() {
        final Function function = new Function(
                FUNC_SUM, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static SumAssembly load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new SumAssembly(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static SumAssembly load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new SumAssembly(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
