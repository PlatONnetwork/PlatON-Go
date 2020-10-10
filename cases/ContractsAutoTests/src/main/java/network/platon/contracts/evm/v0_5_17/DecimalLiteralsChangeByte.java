package network.platon.contracts.evm.v0_5_17;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Bytes4;
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
public class DecimalLiteralsChangeByte extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610218806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806322e20d8f146100465780633db3fb0a146100c6578063c2b21da414610122575b600080fd5b6100726004803603602081101561005c57600080fd5b810190808035906020019092919050505061017e565b60405180827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b6100ce6101bb565b60405180827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b61012a6101cd565b60405180827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b6000808290508060e01b6000806101000a81548163ffffffff021916908360e01c02179055506000809054906101000a900460e01b915050919050565b6000809054906101000a900460e01b81565b60008060009054906101000a900460e01b90509056fea265627a7a72315820a50b619ca4d56b67ed89cb5e6c5d00d15eb652233ffbc5998392a3e04cd2b4c664736f6c63430005110032";

    public static final String FUNC_B4 = "b4";

    public static final String FUNC_GETB4 = "getB4";

    public static final String FUNC_TESTCHANGE = "testChange";

    protected DecimalLiteralsChangeByte(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected DecimalLiteralsChangeByte(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<byte[]> b4() {
        final Function function = new Function(FUNC_B4, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes4>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<byte[]> getB4() {
        final Function function = new Function(FUNC_GETB4, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes4>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<TransactionReceipt> testChange(BigInteger a) {
        final Function function = new Function(
                FUNC_TESTCHANGE, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(a)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<DecimalLiteralsChangeByte> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DecimalLiteralsChangeByte.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<DecimalLiteralsChangeByte> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DecimalLiteralsChangeByte.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static DecimalLiteralsChangeByte load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new DecimalLiteralsChangeByte(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static DecimalLiteralsChangeByte load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new DecimalLiteralsChangeByte(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
