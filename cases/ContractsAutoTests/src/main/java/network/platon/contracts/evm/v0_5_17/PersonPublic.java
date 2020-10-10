package network.platon.contracts.evm.v0_5_17;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.Utf8String;
import com.alaya.abi.solidity.datatypes.generated.Int256;
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
public class PersonPublic extends Contract {
    private static final String BINARY = "6080604052601d60018190555060aa6002819055506040518060400160405280600981526020017f4c75636b7920646f6700000000000000000000000000000000000000000000008152506003908051906020019061005f9291906100b1565b506040518060400160405280600a81526020017f323031312d30312d303100000000000000000000000000000000000000000000815250600090805190602001906100ab9291906100b1565b50610156565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100f257805160ff1916838001178555610120565b82800160010185558215610120579182015b8281111561011f578251825591602001919060010190610104565b5b50905061012d9190610131565b5090565b61015391905b8082111561014f576000816000905550600101610137565b5090565b90565b610187806101656000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c8063262a9dff14610046578063beb0067e14610064578063f377bd5b146100e7575b600080fd5b61004e610105565b6040518082815260200191505060405180910390f35b61006c61010f565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100ac578082015181840152602081019050610091565b50505050905090810190601f1680156100d95780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6100ef61014c565b6040518082815260200191505060405180910390f35b6000600154905090565b60606040518060400160405280600a81526020017f323032302d31322d313500000000000000000000000000000000000000000000815250905090565b6001548156fea265627a7a72315820e6c8db8a323d164dccb50ba1af9c44d18b01ab3dfa77e6899a39803232f2e18a64736f6c63430005110032";

    public static final String FUNC__AGE = "_age";

    public static final String FUNC_AGE = "age";

    public static final String FUNC_BIRTHDAY = "birthDay";

    protected PersonPublic(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected PersonPublic(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> _age() {
        final Function function = new Function(FUNC__AGE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Int256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> age() {
        final Function function = new Function(FUNC_AGE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Int256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> birthDay() {
        final Function function = new Function(FUNC_BIRTHDAY, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public static RemoteCall<PersonPublic> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(PersonPublic.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<PersonPublic> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(PersonPublic.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static PersonPublic load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new PersonPublic(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static PersonPublic load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new PersonPublic(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
