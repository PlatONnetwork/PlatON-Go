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
public class SuicideAndSelfdestruct extends Contract {
    private static final String BINARY = "60806040526000805534801561001457600080fd5b5033600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101b5806100656000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806341c0e1b514610046578063a87d942c14610050578063d09de08a1461006e575b600080fd5b61004e610078565b005b61005861010b565b6040518082815260200191505060405180910390f35b610076610114565b005b3373ffffffffffffffffffffffffffffffffffffffff16600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16141561010957600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16ff5b565b60008054905090565b6000600a90503373ffffffffffffffffffffffffffffffffffffffff16600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16141561017c5780600054016000819055505b5056fea2646970667358221220add4a856904c91ad83f663fb7aaf0b72c74c60629cd9ec01dff117b35946d1ef64736f6c63430007010033";

    public static final String FUNC_GETCOUNT = "getCount";

    public static final String FUNC_INCREMENT = "increment";

    public static final String FUNC_KILL = "kill";

    protected SuicideAndSelfdestruct(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected SuicideAndSelfdestruct(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<SuicideAndSelfdestruct> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(SuicideAndSelfdestruct.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<SuicideAndSelfdestruct> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(SuicideAndSelfdestruct.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public RemoteCall<TransactionReceipt> getCount() {
        final Function function = new Function(
                FUNC_GETCOUNT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> increment() {
        final Function function = new Function(
                FUNC_INCREMENT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> kill() {
        final Function function = new Function(
                FUNC_KILL, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static SuicideAndSelfdestruct load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new SuicideAndSelfdestruct(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static SuicideAndSelfdestruct load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new SuicideAndSelfdestruct(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
