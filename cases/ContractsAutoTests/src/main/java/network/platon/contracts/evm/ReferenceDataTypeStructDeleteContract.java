package network.platon.contracts.evm;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Bool;
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
public class ReferenceDataTypeStructDeleteContract extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506005600081905550600a600160020181905550600180600301600060018152602001908152602001600020819055506002600160030160006002815260200190815260200160002081905550600360016000016000018190555060018060000160010160006001815260200190815260200160002060006101000a81548160ff02191690831515021790555060018060000160010160006002815260200190815260200160002060006101000a81548160ff02191690831515021790555060016000808201600080820160009055505060028201600090555050600080905561019f806100ff6000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c806311977c5c1461005c5780631268893e1461007a5780635ff76c8a1461009857806379e44a38146100b6578063d587919c146100d4575b600080fd5b6100646100f6565b6040518082815260200191505060405180910390f35b6100826100ff565b6040518082815260200191505060405180910390f35b6100a061011e565b6040518082815260200191505060405180910390f35b6100be61012e565b6040518082815260200191505060405180910390f35b6100dc61013b565b604051808215151515815260200191505060405180910390f35b60008054905090565b6000600160030160006001815260200190815260200160002054905090565b6000600160000160000154905090565b6000600160020154905090565b6000600160000160010160006001815260200190815260200160002060009054906101000a900460ff1690509056fea265627a7a72315820c32f99a185e87e67d75f232ff928f3d4dad2d15ff91172d83ac1bae11fcd1ba264736f6c63430005110032";

    public static final String FUNC_GETNESTEDMAPPING = "getNestedMapping";

    public static final String FUNC_GETNESTEDVALUE = "getNestedValue";

    public static final String FUNC_GETTODELETEINT = "getToDeleteInt";

    public static final String FUNC_GETTOPMAPPING = "getTopMapping";

    public static final String FUNC_GETTOPVALUE = "getTopValue";

    protected ReferenceDataTypeStructDeleteContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ReferenceDataTypeStructDeleteContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<ReferenceDataTypeStructDeleteContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ReferenceDataTypeStructDeleteContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<ReferenceDataTypeStructDeleteContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ReferenceDataTypeStructDeleteContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public RemoteCall<Boolean> getNestedMapping() {
        final Function function = new Function(FUNC_GETNESTEDMAPPING, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<BigInteger> getNestedValue() {
        final Function function = new Function(FUNC_GETNESTEDVALUE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getToDeleteInt() {
        final Function function = new Function(FUNC_GETTODELETEINT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getTopMapping() {
        final Function function = new Function(FUNC_GETTOPMAPPING, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getTopValue() {
        final Function function = new Function(FUNC_GETTOPVALUE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static ReferenceDataTypeStructDeleteContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ReferenceDataTypeStructDeleteContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ReferenceDataTypeStructDeleteContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ReferenceDataTypeStructDeleteContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
