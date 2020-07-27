package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Int256;
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
public class InterfaceContractStructTest extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610264806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c8063079d318b14610046578063d629546c14610064578063e36afde31461006e575b600080fd5b61004e61008c565b6040518082815260200191505060405180910390f35b61006c6100ae565b005b61007661017e565b6040518082815260200191505060405180910390f35b600080600080015414156100a357600090506100ab565b600080015490505b90565b6040518060600160405280600181526020016040518060400160405280600681526020017f506c61744f4e000000000000000000000000000000000000000000000000000081525081526020016040518060400160405280600f81526020017f506c61744f4e20446573637269626500000000000000000000000000000000008152508152506000808201518160000155602082015181600101908051906020019061015b92919061018a565b50604082015181600201908051906020019061017892919061018a565b50905050565b60008060000154905090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106101cb57805160ff19168380011785556101f9565b828001600101855582156101f9579182015b828111156101f85782518255916020019190600101906101dd565b5b509050610206919061020a565b5090565b61022c91905b80821115610228576000816000905550600101610210565b5090565b9056fea265627a7a72315820a95f16e2a57ba406793117febd6897a756bac6b708645605392201ac6f49f7be64736f6c634300050d0032";

    public static final String FUNC_GETBOOKID = "getBookID";

    public static final String FUNC_GETDEFAULTBOOKID = "getDefaultBookID";

    public static final String FUNC_SETBOOK = "setBook";

    protected InterfaceContractStructTest(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected InterfaceContractStructTest(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> getBookID() {
        final Function function = new Function(FUNC_GETBOOKID, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Int256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> getDefaultBookID() {
        final Function function = new Function(
                FUNC_GETDEFAULTBOOKID, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> setBook() {
        final Function function = new Function(
                FUNC_SETBOOK, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<InterfaceContractStructTest> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InterfaceContractStructTest.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<InterfaceContractStructTest> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InterfaceContractStructTest.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static InterfaceContractStructTest load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new InterfaceContractStructTest(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static InterfaceContractStructTest load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new InterfaceContractStructTest(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
