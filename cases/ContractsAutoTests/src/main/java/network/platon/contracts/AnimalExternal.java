package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.Utf8String;
import org.web3j.abi.datatypes.generated.Int256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.GasProvider;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.0.7.
 */
public class AnimalExternal extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50601d60018190555060aa6002819055506040518060400160405280600981526020017f4c75636b7920646f6700000000000000000000000000000000000000000000008152506003908051906020019061006c9291906100be565b506040518060400160405280600a81526020017f323031312d30312d303100000000000000000000000000000000000000000000815250600090805190602001906100b89291906100be565b50610163565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100ff57805160ff191683800117855561012d565b8280016001018555821561012d579182015b8281111561012c578251825591602001919060010190610111565b5b50905061013a919061013e565b5090565b61016091905b8082111561015c576000816000905550600101610144565b5090565b90565b6101ec806101726000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c8063262a9dff14610046578063beb0067e14610064578063f377bd5b146100e7575b600080fd5b61004e610105565b6040518082815260200191505060405180910390f35b61006c61010f565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100ac578082015181840152602081019050610091565b50505050905090810190601f1680156100d95780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6100ef6101b1565b6040518082815260200191505060405180910390f35b6000600154905090565b606060008054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156101a75780601f1061017c576101008083540402835291602001916101a7565b820191906000526020600020905b81548152906001019060200180831161018a57829003601f168201915b5050505050905090565b6001548156fea265627a7a7231582021f78e5ef44f25bf15434f51d057dde0b711813f35dac01fb6bc0e5278aebf4764736f6c634300050d0032";

    public static final String FUNC__AGE = "_age";

    public static final String FUNC_AGE = "age";

    public static final String FUNC_BIRTHDAY = "birthDay";

    protected AnimalExternal(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected AnimalExternal(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<AnimalExternal> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AnimalExternal.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<AnimalExternal> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AnimalExternal.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
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

    public static AnimalExternal load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new AnimalExternal(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static AnimalExternal load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new AnimalExternal(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
