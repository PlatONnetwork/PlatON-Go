package network.platon.contracts.evm.v0_7_1;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.DynamicBytes;
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
 * <p>Generated with web3j version 0.13.2.1.
 */
public class PlatonInner extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506103a1806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063eb51cc911461003b578063f40ae8d9146100be575b600080fd5b610043610199565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610083578082015181840152602081019050610068565b50505050905090810190601f1680156100b05780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610197600480360360408110156100d457600080fd5b81019080803590602001906401000000008111156100f157600080fd5b82018360208201111561010357600080fd5b8035906020019184600183028401116401000000008311171561012557600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061023b565b005b606060008054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156102315780601f1061020657610100808354040283529160200191610231565b820191906000526020600020905b81548152906001019060200180831161021457829003601f168201915b5050505050905090565b6000825190506000606060008084602088016000885af161025857fe5b3d91508167ffffffffffffffff8111801561027257600080fd5b506040519080825280601f01601f1916602001820160405280156102a55781602001600182028036833780820191505090505b5090503d6000602083013e80600090805190602001906102c69291906102ce565b505050505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061030f57805160ff191683800117855561033d565b8280016001018555821561033d579182015b8281111561033c578251825591602001919060010190610321565b5b50905061034a919061034e565b5090565b5b8082111561036757600081600090555060010161034f565b509056fea2646970667358221220283460bd7fba1dfcee18504bae45379c653ddb6c34a6ed19ad75fb9dec23c90664736f6c63430007010033";

    public static final String FUNC_ASSEMBLYCALLPPOS = "assemblyCallppos";

    public static final String FUNC_GETRETURNVALUE = "getReturnValue";

    protected PlatonInner(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected PlatonInner(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> assemblyCallppos(byte[] data, String addr) {
        final Function function = new Function(
                FUNC_ASSEMBLYCALLPPOS, 
                Arrays.<Type>asList(new DynamicBytes(data),
                new com.alaya.abi.solidity.datatypes.Address(addr)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<byte[]> getReturnValue() {
        final Function function = new Function(FUNC_GETRETURNVALUE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicBytes>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public static RemoteCall<PlatonInner> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(PlatonInner.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<PlatonInner> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(PlatonInner.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static PlatonInner load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new PlatonInner(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static PlatonInner load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new PlatonInner(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
