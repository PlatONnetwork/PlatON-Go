package network.platon.contracts.evm.v0_6_12;

import com.alaya.abi.solidity.FunctionEncoder;
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
public class ErrorParamConstructor extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060405161020d38038061020d8339818101604052602081101561003357600080fd5b8101908080519060200190929190505050600a806000819055505080600181905550506101a8806100656000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c80630dbe671f1461005c5780634df7e3d01461007a57806382ab890a14610098578063a1c51915146100f7578063d46300fd14610115575b600080fd5b610064610133565b6040518082815260200191505060405180910390f35b610082610139565b6040518082815260200191505060405180910390f35b6100c4600480360360208110156100ae57600080fd5b810190808035906020019092919050505061013f565b604051808373ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390f35b6100ff61015f565b6040518082815260200191505060405180910390f35b61011d610169565b6040518082815260200191505060405180910390f35b60005481565b60015481565b600080826001600082825401925050819055503360015491509150915091565b6000600154905090565b6000805490509056fea264697066735822122065af30e171b8cb03c0b37c94fd74012c0f758cc6c32678912aea7f8747dbc23a64736f6c634300060c0033";

    public static final String FUNC_A = "a";

    public static final String FUNC_B = "b";

    public static final String FUNC_GETA = "getA";

    public static final String FUNC_GETB = "getB";

    public static final String FUNC_UPDATE = "update";

    protected ErrorParamConstructor(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ErrorParamConstructor(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<ErrorParamConstructor> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId, BigInteger _b) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(_b)));
        return deployRemoteCall(ErrorParamConstructor.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static RemoteCall<ErrorParamConstructor> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId, BigInteger _b) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(_b)));
        return deployRemoteCall(ErrorParamConstructor.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public RemoteCall<TransactionReceipt> a() {
        final Function function = new Function(
                FUNC_A, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> b() {
        final Function function = new Function(
                FUNC_B, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getA() {
        final Function function = new Function(
                FUNC_GETA, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getB() {
        final Function function = new Function(
                FUNC_GETB, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> update(BigInteger amount) {
        final Function function = new Function(
                FUNC_UPDATE, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(amount)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static ErrorParamConstructor load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ErrorParamConstructor(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ErrorParamConstructor load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ErrorParamConstructor(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
