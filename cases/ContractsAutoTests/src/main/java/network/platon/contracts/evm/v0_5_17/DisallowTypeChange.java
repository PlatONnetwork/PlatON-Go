package network.platon.contracts.evm.v0_5_17;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint32;
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
 * <p>Generated with web3j version 0.13.2.1.
 */
public class DisallowTypeChange extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610165806100206000396000f3fe6080604052600436106100345760003560e01c80630b7f166514610039578063420343a414610070578063a56dfe4a1461007a575b600080fd5b34801561004557600080fd5b5061004e6100b1565b604051808263ffffffff1663ffffffff16815260200191505060405180910390f35b6100786100ca565b005b34801561008657600080fd5b5061008f61011b565b604051808263ffffffff1663ffffffff16815260200191505060405180910390f35b60008060009054906101000a900463ffffffff16905090565b60007faaaa000000000000000000000000000000000000000000000000000000000000905060008190508060e01c6000806101000a81548163ffffffff021916908363ffffffff1602179055505050565b6000809054906101000a900463ffffffff168156fea265627a7a72315820e8091c954fc1e9e837a99bdf576c1c9a4aaa7d9ce48af4698c5b39fec85a07ab64736f6c63430005110032";

    public static final String FUNC_GETY = "getY";

    public static final String FUNC_TESTCHANGE = "testChange";

    public static final String FUNC_Y = "y";

    protected DisallowTypeChange(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected DisallowTypeChange(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> getY() {
        final Function function = new Function(FUNC_GETY, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint32>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> testChange(BigInteger vonValue) {
        final Function function = new Function(
                FUNC_TESTCHANGE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public RemoteCall<BigInteger> y() {
        final Function function = new Function(FUNC_Y, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint32>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<DisallowTypeChange> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DisallowTypeChange.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<DisallowTypeChange> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DisallowTypeChange.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static DisallowTypeChange load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new DisallowTypeChange(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static DisallowTypeChange load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new DisallowTypeChange(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
