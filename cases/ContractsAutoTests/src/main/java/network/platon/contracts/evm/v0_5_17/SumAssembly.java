package network.platon.contracts.evm.v0_5_17;

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
public class SumAssembly extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50600060019080600181540180825580915050906001820390600052602060002001600090919290919091505550600060029080600181540180825580915050906001820390600052602060002001600090919290919091505550600060039080600181540180825580915050906001820390600052602060002001600090919290919091505550600060049080600181540180825580915050906001820390600052602060002001600090919290919091505550600060059080600181540180825580915050906001820390600052602060002001600090919290919091505550610153806101016000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063853255cc14610030575b600080fd5b61003861004e565b6040518082815260200191505060405180910390f35b600073__$0b57193f600960d8b32f465d99d292bef2$__6387fbcc7760006040518263ffffffff1660e01b8152600401808060200182810382528381815481526020019150805480156100c057602002820191906000526020600020905b8154815260200190600101908083116100ac575b50509250505060206040518083038186803b1580156100de57600080fd5b505af41580156100f2573d6000803e3d6000fd5b505050506040513d602081101561010857600080fd5b810190808051906020019092919050505090509056fea265627a7a72315820c53b8a5b9972a12b2d7eb12d5af89ca51b1d0585ddd3d5a9fc65d508598b713f64736f6c63430005110032\n"
            + "\n"
            + "// $0b57193f600960d8b32f465d99d292bef2$ -> /home/platon/.jenkins/workspace/contracts_test_alaya/cases/ContractsAutoTests/src/test/resources/contracts/evm/0.5.17/1.function/07Assembly/SumAssembly.sol:Sum";

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

    public RemoteCall<BigInteger> sum() {
        final Function function = new Function(FUNC_SUM, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static SumAssembly load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new SumAssembly(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static SumAssembly load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new SumAssembly(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
