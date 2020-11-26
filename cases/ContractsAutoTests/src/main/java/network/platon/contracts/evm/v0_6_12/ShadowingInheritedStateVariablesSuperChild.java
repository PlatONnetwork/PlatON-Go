package network.platon.contracts.evm.v0_6_12;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.Utf8String;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.util.Arrays;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.1.
 */
public class ShadowingInheritedStateVariablesSuperChild extends Contract {
    private static final String BINARY = "60806040526040518060600160405280602e815260200161040a602e913960009080519060200190610032929190610074565b5034801561003f57600080fd5b506040518060600160405280603381526020016103d7603391396000908051906020019061006e929190610074565b50610111565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100b557805160ff19168380011785556100e3565b828001600101855582156100e3579182015b828111156100e25782518255916020019190600101906100c7565b5b5090506100f091906100f4565b5090565b5b8082111561010d5760008160009055506001016100f5565b5090565b6102b7806101206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806306fdde031461003b57806317d7de7c146100be575b600080fd5b610043610141565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610083578082015181840152602081019050610068565b50505050905090810190601f1680156100b05780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6100c66101df565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101065780820151818401526020810190506100eb565b50505050905090810190601f1680156101335780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60008054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156101d75780601f106101ac576101008083540402835291602001916101d7565b820191906000526020600020905b8154815290600101906020018083116101ba57829003601f168201915b505050505081565b606060008054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156102775780601f1061024c57610100808354040283529160200191610277565b820191906000526020600020905b81548152906001019060200180831161025a57829003601f168201915b505050505090509056fea2646970667358221220b257e0b8ce68afd1a112930f290cd55f2d13b665dac12104723daf03e7160b1e64736f6c634300060c0033436f6e747261637420536861646f77696e67496e6865726974656453746174655661726961626c657353757065724368696c64436f6e747261637420536861646f77696e67496e6865726974656453746174655661726961626c65735375706572";

    public static final String FUNC_GETNAME = "getName";

    public static final String FUNC_NAME = "name";

    protected ShadowingInheritedStateVariablesSuperChild(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ShadowingInheritedStateVariablesSuperChild(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<ShadowingInheritedStateVariablesSuperChild> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ShadowingInheritedStateVariablesSuperChild.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<ShadowingInheritedStateVariablesSuperChild> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ShadowingInheritedStateVariablesSuperChild.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public RemoteCall<String> getName() {
        final Function function = new Function(FUNC_GETNAME, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> name() {
        final Function function = new Function(FUNC_NAME, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public static ShadowingInheritedStateVariablesSuperChild load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ShadowingInheritedStateVariablesSuperChild(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ShadowingInheritedStateVariablesSuperChild load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ShadowingInheritedStateVariablesSuperChild(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
