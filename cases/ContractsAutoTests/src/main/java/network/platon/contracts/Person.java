package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.Utf8String;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
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
public class Person extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060b4600181905550601460008190555033600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040805190810160405280600681526020017f71637869616f0000000000000000000000000000000000000000000000000000815250600390805190602001906100ad9291906100b3565b50610158565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100f457805160ff1916838001178555610122565b82800160010185558215610122579182015b82811115610121578251825591602001919060010190610106565b5b50905061012f9190610133565b5090565b61015591905b80821115610151576000816000905550600101610139565b5090565b90565b6104b1806101676000396000f30060806040526004361061008e576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806306fdde03146100935780630ef267431461012357806326121ff01461014e578063262a9dff1461016557806341c0e1b514610190578063741a3944146101a75780638da5cb5b146101d4578063d5dcf1271461022b575b600080fd5b34801561009f57600080fd5b506100a8610258565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100e85780820151818401526020810190506100cd565b50505050905090810190601f1680156101155780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561012f57600080fd5b506101386102fa565b6040518082815260200191505060405180910390f35b34801561015a57600080fd5b50610163610304565b005b34801561017157600080fd5b5061017a610310565b6040518082815260200191505060405180910390f35b34801561019c57600080fd5b506101a5610319565b005b3480156101b357600080fd5b506101d2600480360381019080803590602001909291905050506103ac565b005b3480156101e057600080fd5b506101e96103b6565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561023757600080fd5b50610256600480360381019080803590602001909291905050506103e0565b005b606060038054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156102f05780601f106102c5576101008083540402835291602001916102f0565b820191906000526020600020905b8154815290600101906020018083116102d357829003601f168201915b5050505050905090565b6000600154905090565b61030e60036103ea565b565b60008054905090565b3373ffffffffffffffffffffffffffffffffffffffff16600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614156103aa57600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16ff5b565b8060018190555050565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b8060008190555050565b7f5100000000000000000000000000000000000000000000000000000000000000816000815460018160011615610100020316600290048110151561042b57fe5b81546001161561044a5790600052602060002090602091828204019190065b601f036101000a81548160ff021916907f010000000000000000000000000000000000000000000000000000000000000084040217905550505600a165627a7a7230582029b75740ddec30d09eed2cd01ce9ad0212919c37be04f49e0c388793246cb2b00029";

    public static final String FUNC_NAME = "name";

    public static final String FUNC_HEIGHT = "height";

    public static final String FUNC_F = "f";

    public static final String FUNC_AGE = "age";

    public static final String FUNC_KILL = "kill";

    public static final String FUNC_SETHEIGHT = "setHeight";

    public static final String FUNC_OWNER = "owner";

    public static final String FUNC_SETAGE = "setAge";

    protected Person(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected Person(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<String> name() {
        final Function function = new Function(FUNC_NAME, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> height() {
        final Function function = new Function(FUNC_HEIGHT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> f() {
        final Function function = new Function(
                FUNC_F, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> age() {
        final Function function = new Function(FUNC_AGE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public void kill() {
        throw new RuntimeException("cannot call constant function with void return type");
    }

    public RemoteCall<TransactionReceipt> setHeight(BigInteger height) {
        final Function function = new Function(
                FUNC_SETHEIGHT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(height)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> owner() {
        final Function function = new Function(FUNC_OWNER, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> setAge(BigInteger age) {
        final Function function = new Function(
                FUNC_SETAGE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(age)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<Person> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Person.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<Person> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Person.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static Person load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new Person(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static Person load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new Person(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
