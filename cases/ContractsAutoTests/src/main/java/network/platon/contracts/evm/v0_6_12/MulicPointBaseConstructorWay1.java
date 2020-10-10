package network.platon.contracts.evm.v0_6_12;

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
public class MulicPointBaseConstructorWay1 extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50600a806000819055505061010a8061002a6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80630dbe671f14603757806382ab890a146053575b600080fd5b603d60af565b6040518082815260200191505060405180910390f35b607c60048036036020811015606757600080fd5b810190808035906020019092919050505060b5565b604051808373ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390f35b60005481565b600080826000808282540192505081905550336000549150915091509156fea26469706673582212206ebb4283d4739c70cb8f3e73dbbdf711ad90d085dad2de07ea0ca06513049ae964736f6c634300060c0033";

    public static final String FUNC_A = "a";

    public static final String FUNC_UPDATE = "update";

    protected MulicPointBaseConstructorWay1(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected MulicPointBaseConstructorWay1(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<MulicPointBaseConstructorWay1> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(MulicPointBaseConstructorWay1.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<MulicPointBaseConstructorWay1> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(MulicPointBaseConstructorWay1.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public RemoteCall<TransactionReceipt> a() {
        final Function function = new Function(
                FUNC_A, 
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

    public static MulicPointBaseConstructorWay1 load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new MulicPointBaseConstructorWay1(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static MulicPointBaseConstructorWay1 load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new MulicPointBaseConstructorWay1(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
