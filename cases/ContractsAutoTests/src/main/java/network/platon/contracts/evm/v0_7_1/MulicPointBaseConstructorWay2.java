package network.platon.contracts.evm.v0_7_1;

import com.alaya.abi.solidity.FunctionEncoder;
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
 * <p>Generated with web3j version 0.13.2.1.
 */
public class MulicPointBaseConstructorWay2 extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506040516101a43803806101a48339818101604052602081101561003357600080fd5b81019080805190602001909291905050508081028060008190555050506101458061005f6000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80630dbe671f1461004657806382ab890a14610064578063d46300fd146100c3575b600080fd5b61004e6100e1565b6040518082815260200191505060405180910390f35b6100906004803603602081101561007a57600080fd5b81019080803590602001909291905050506100e7565b604051808373ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390f35b6100cb610106565b6040518082815260200191505060405180910390f35b60005481565b6000808260008082825401925050819055503360005491509150915091565b6000805490509056fea26469706673582212203be9d5c9a4fc3e48c136fe29cb6f0221a5af893d6a61f13cf1543a158accc88064736f6c63430007010033";

    public static final String FUNC_A = "a";

    public static final String FUNC_GETA = "getA";

    public static final String FUNC_UPDATE = "update";

    protected MulicPointBaseConstructorWay2(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected MulicPointBaseConstructorWay2(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<MulicPointBaseConstructorWay2> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId, BigInteger _y) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new Uint256(_y)));
        return deployRemoteCall(MulicPointBaseConstructorWay2.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static RemoteCall<MulicPointBaseConstructorWay2> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId, BigInteger _y) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new Uint256(_y)));
        return deployRemoteCall(MulicPointBaseConstructorWay2.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public RemoteCall<BigInteger> a() {
        final Function function = new Function(FUNC_A, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getA() {
        final Function function = new Function(FUNC_GETA, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> update(BigInteger amount) {
        final Function function = new Function(
                FUNC_UPDATE, 
                Arrays.<Type>asList(new Uint256(amount)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static MulicPointBaseConstructorWay2 load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new MulicPointBaseConstructorWay2(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static MulicPointBaseConstructorWay2 load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new MulicPointBaseConstructorWay2(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
