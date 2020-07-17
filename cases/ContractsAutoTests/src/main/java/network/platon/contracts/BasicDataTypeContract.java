package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Int8;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.abi.datatypes.generated.Uint8;
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
public class BasicDataTypeContract extends Contract {
    private static final String BINARY = "60806040527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60005560018055600280556001600360006101000a81548160ff021916908360ff16021790555060ff600360016101000a81548160ff021916908360ff1602179055506001600360026101000a81548161ffff021916908361ffff16021790555061ffff600360046101000a81548161ffff021916908361ffff1602179055506001600360066101000a81548160ff021916908360000b60ff160217905550607f600360076101000a81548160ff021916908360000b60ff1602179055507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600360086101000a81548160ff021916908360000b60ff1602179055507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80600360096101000a81548160ff021916908360000b60ff16021790555060016003600a6101000a81548160ff02191690831515021790555060006003600b6101000a81548160ff0219169083151502179055507f61000000000000000000000000000000000000000000000000000000000000006003600c6101000a81548160ff02191690837f01000000000000000000000000000000000000000000000000000000000000009004021790555060017f0100000000000000000000000000000000000000000000000000000000000000026003600d6101000a81548160ff02191690837f0100000000000000000000000000000000000000000000000000000000000000900402179055507f61620000000000000000000000000000000000000000000000000000000000006003600e6101000a81548161ffff02191690837e01000000000000000000000000000000000000000000000000000000000000900402179055507f6162630000000000000000000000000000000000000000000000000000000000600360106101000a81548162ffffff02191690837d010000000000000000000000000000000000000000000000000000000000900402179055506040805190810160405280600181526020017f61000000000000000000000000000000000000000000000000000000000000008152506004908051906020019061034f9291906103fa565b506040805190810160405280600281526020017f61620000000000000000000000000000000000000000000000000000000000008152506005908051906020019061039b9291906103fa565b506040805190810160405280600381526020017f6162630000000000000000000000000000000000000000000000000000000000815250600690805190602001906103e79291906103fa565b503480156103f457600080fd5b5061049f565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061043b57805160ff1916838001178555610469565b82800160010185558215610469579182015b8281111561046857825182559160200191906001019061044d565b5b509050610476919061047a565b5090565b61049c91905b80821115610498576000816000905550600101610480565b5090565b90565b6101be806104ae6000396000f3fe608060405234801561001057600080fd5b5060043610610069576000357c01000000000000000000000000000000000000000000000000000000009004806307da3eae1461006e5780633d9ceb371461008c5780635c5c8419146100d7578063d29f1598146100f5575b600080fd5b610076610140565b6040518082815260200191505060405180910390f35b6100bb600480360360208110156100a257600080fd5b81019080803560000b906020019092919050505061015e565b604051808260000b60000b815260200191505060405180910390f35b6100df61016b565b6040518082815260200191505060405180910390f35b6101246004803603602081101561010b57600080fd5b81019080803560ff169060200190929190505050610185565b604051808260ff1660ff16815260200191505060405180910390f35b60006006805460018160011615610100020316600290049050905090565b6000600182019050919050565b60006003600e9054906101000a905050600260ff16905090565b600060018201905091905056fea165627a7a7230582033a530bdde5427d3ace252376bef55d9c267c5f1a7e70f823ec996d7930c96360029";

    public static final String FUNC_GETBYTESLENGTH = "getBytesLength";

    public static final String FUNC_ADDINTOVERFLOW = "addIntOverflow";

    public static final String FUNC_GETBYTES1LENGTH = "getBytes1Length";

    public static final String FUNC_ADDUINTOVERFLOW = "addUintOverflow";

    protected BasicDataTypeContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected BasicDataTypeContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> getBytesLength() {
        final Function function = new Function(FUNC_GETBYTESLENGTH, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> addIntOverflow(BigInteger a) {
        final Function function = new Function(FUNC_ADDINTOVERFLOW, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Int8(a)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Int8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getBytes1Length() {
        final Function function = new Function(FUNC_GETBYTES1LENGTH, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> addUintOverflow(BigInteger a) {
        final Function function = new Function(FUNC_ADDUINTOVERFLOW, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint8(a)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<BasicDataTypeContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(BasicDataTypeContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<BasicDataTypeContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(BasicDataTypeContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static BasicDataTypeContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new BasicDataTypeContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static BasicDataTypeContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new BasicDataTypeContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
