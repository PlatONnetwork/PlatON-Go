package network.platon.contracts.evm;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
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
public class DisallowYears extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506101d3806100206000396000f3fe60806040526004361061004a5760003560e01c80630bb2b6961461004f57806320de797e1461007a57806325b29d84146100bc578063c6d8d657146100e7578063c6f8a3b714610112575b600080fd5b34801561005b57600080fd5b5061006461013d565b6040518082815260200191505060405180910390f35b6100a66004803603602081101561009057600080fd5b8101908080359060200190929190505050610147565b6040518082815260200191505060405180910390f35b3480156100c857600080fd5b506100d1610181565b6040518082815260200191505060405180910390f35b3480156100f357600080fd5b506100fc61018b565b6040518082815260200191505060405180910390f35b34801561011e57600080fd5b50610127610194565b6040518082815260200191505060405180910390f35b6000600254905090565b60006301e13380600081905550680dd2d5fcf3bc9c000060018190555060ff600281905550680dd2d5fcf3bc9c0000600381905550919050565b6000600154905090565b60008054905090565b600060035490509056fea265627a7a72315820305b1753bfb9b5bf29faaedeedc621ea888ca9add4f05f9a62e4d819d2e2f02664736f6c63430005110032";

    public static final String FUNC_GETETHERVALUE = "getEtherValue";

    public static final String FUNC_GETHEXCOMVALUE = "getHexComValue";

    public static final String FUNC_GETHEXVALUE = "getHexValue";

    public static final String FUNC_GETTIME1 = "getTime1";

    public static final String FUNC_TESTYEAR = "testyear";

    protected DisallowYears(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected DisallowYears(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> getEtherValue() {
        final Function function = new Function(FUNC_GETETHERVALUE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getHexComValue() {
        final Function function = new Function(FUNC_GETHEXCOMVALUE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getHexValue() {
        final Function function = new Function(FUNC_GETHEXVALUE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getTime1() {
        final Function function = new Function(FUNC_GETTIME1, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> testyear(BigInteger a, BigInteger vonValue) {
        final Function function = new Function(
                FUNC_TESTYEAR, 
                Arrays.<Type>asList(new Uint256(a)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public static RemoteCall<DisallowYears> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DisallowYears.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<DisallowYears> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DisallowYears.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static DisallowYears load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new DisallowYears(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static DisallowYears load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new DisallowYears(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
