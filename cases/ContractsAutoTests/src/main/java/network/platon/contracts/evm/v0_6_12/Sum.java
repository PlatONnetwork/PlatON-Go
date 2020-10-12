package network.platon.contracts.evm.v0_6_12;

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
import java.util.List;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.1.
 */
public class Sum extends Contract {
    private static final String BINARY = "61016c610026600b82828239805160001a60731461001957fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600436106100355760003560e01c806387fbcc771461003a575b600080fd5b6100f06004803603602081101561005057600080fd5b810190808035906020019064010000000081111561006d57600080fd5b82018360208201111561007f57600080fd5b803590602001918460208302840111640100000000831117156100a157600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f820116905080830192505050505050509192919290505050610106565b6040518082815260200191505060405180910390f35b600080600090505b825181101561013057602081026020840101518201915080600101905061010e565b5091905056fea26469706673582212205348bed70e29de2214caae24a6872babdeea9203111542652c7a7b3417391faf64736f6c634300060c0033";

    public static final String FUNC_SUMUSINGINLINEASSEMBLY = "sumUsingInlineAssembly";

    protected Sum(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected Sum(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> sumUsingInlineAssembly(List<BigInteger> _data) {
        final Function function = new Function(FUNC_SUMUSINGINLINEASSEMBLY, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.DynamicArray<Uint256>(
                Uint256.class,
                        com.alaya.abi.solidity.Utils.typeMap(_data, Uint256.class))),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<Sum> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Sum.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<Sum> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Sum.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static Sum load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new Sum(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static Sum load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new Sum(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
