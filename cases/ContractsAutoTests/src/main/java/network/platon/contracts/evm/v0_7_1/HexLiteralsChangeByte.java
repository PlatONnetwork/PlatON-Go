package network.platon.contracts.evm.v0_7_1;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
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
public class HexLiteralsChangeByte extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506101c8806100206000396000f3fe6080604052600436106100345760003560e01c80630b7f166514610039578063420343a414610086578063ee495002146100c6575b600080fd5b34801561004557600080fd5b5061004e610113565b60405180827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b61008e610129565b60405180827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b3480156100d257600080fd5b506100db610181565b60405180827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b60008060009054906101000a900460f81b905090565b6000807f6162000000000000000000000000000000000000000000000000000000000000905060f160f81b6000806101000a81548160ff021916908360f81c021790555060008054906101000a900460f81b91505090565b60008054906101000a900460f81b8156fea2646970667358221220e732d437603895dfd7a82bfaaa9b4eea38e04b2265776912c8d400f08299589264736f6c63430007010033";

    public static final String FUNC_B1 = "b1";

    public static final String FUNC_GETY = "getY";

    public static final String FUNC_TESTCHANGE = "testChange";

    protected HexLiteralsChangeByte(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected HexLiteralsChangeByte(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> b1() {
        final Function function = new Function(
                FUNC_B1, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getY() {
        final Function function = new Function(
                FUNC_GETY, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> testChange() {
        final Function function = new Function(
                FUNC_TESTCHANGE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<HexLiteralsChangeByte> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(HexLiteralsChangeByte.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<HexLiteralsChangeByte> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(HexLiteralsChangeByte.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static HexLiteralsChangeByte load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new HexLiteralsChangeByte(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static HexLiteralsChangeByte load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new HexLiteralsChangeByte(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
