package network.platon.contracts.evm.v0_5_17;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint8;
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
public class InterfaceContractEnumTest extends Contract {
    private static final String BINARY = "60806040526001600060016101000a81548160ff0219169083600281111561002357fe5b021790555034801561003457600080fd5b5061011d806100446000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c806367cb61b6146041578063694ebe4f14606a578063843f7258146072575b600080fd5b6047608e565b60405180826002811115605657fe5b60ff16815260200191505060405180910390f35b607060a4565b005b607860c8565b6040518082815260200191505060405180910390f35b60008060009054906101000a900460ff16905090565b60026000806101000a81548160ff0219169083600281111560c157fe5b0217905550565b60008060019054906101000a900460ff16600281111560e357fe5b90509056fea265627a7a723158203864c08c42316112f604119875355f05af4c2241ea8408ba1f1c25b07cbbf8e964736f6c63430005110032";

    public static final String FUNC_GETCHOICE = "getChoice";

    public static final String FUNC_GETDEFAULTCHOICE = "getDefaultChoice";

    public static final String FUNC_SETLARGE = "setLarge";

    protected InterfaceContractEnumTest(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected InterfaceContractEnumTest(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> getChoice() {
        final Function function = new Function(FUNC_GETCHOICE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> getDefaultChoice() {
        final Function function = new Function(
                FUNC_GETDEFAULTCHOICE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> setLarge() {
        final Function function = new Function(
                FUNC_SETLARGE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<InterfaceContractEnumTest> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InterfaceContractEnumTest.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<InterfaceContractEnumTest> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InterfaceContractEnumTest.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static InterfaceContractEnumTest load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new InterfaceContractEnumTest(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static InterfaceContractEnumTest load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new InterfaceContractEnumTest(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
