package network.platon.contracts.evm.v0_4_26;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Bool;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Bytes32;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import com.alaya.tuples.generated.Tuple2;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Callable;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.0.
 */
public class DeprecatedFunctions extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061024d806100206000396000f300608060405260043610610057576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806341c0e1b51461005c578063a3928d9914610073578063f492f3a8146100b1575b600080fd5b34801561006857600080fd5b506100716100f8565b005b34801561007f57600080fd5b5061008861012a565b604051808315151515815260200182600019166000191681526020019250505060405180910390f35b3480156100bd57600080fd5b506100de60048036038101908080351515906020019092919050505061020b565b604051808215151515815260200191505060405180910390f35b600073ca35b7d915458ef540ade6068dfe2f44e8fa733c90508073ffffffffffffffffffffffffffffffffffffffff16ff5b60008060008060008073ca35b7d915458ef540ade6068dfe2f44e8fa733c935073d25ed029c093e56bc8911a07c46545000cbf37c692508273ffffffffffffffffffffffffffffffffffffffff1684604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019150506000604051808303816000865af2915050915060405180807f77616e677a68616e6778696f6e67000000000000000000000000000000000000815250600e01905060405180910390209050818195509550505050509091565b600081151561021957600080fd5b8190509190505600a165627a7a723058209ce1c2a7d23c077c3a670d7606b0f60fd80582ff018f394d975bae29adadf3280029";

    public static final String FUNC_KILL = "kill";

    public static final String FUNC_FUNCTIONCHECK = "functionCheck";

    public static final String FUNC_THROWCHECK = "throwCheck";

    protected DeprecatedFunctions(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected DeprecatedFunctions(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> kill() {
        final Function function = new Function(
                FUNC_KILL, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<Tuple2<Boolean, byte[]>> functionCheck() {
        final Function function = new Function(FUNC_FUNCTIONCHECK, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}, new TypeReference<Bytes32>() {}));
        return new RemoteCall<Tuple2<Boolean, byte[]>>(
                new Callable<Tuple2<Boolean, byte[]>>() {
                    @Override
                    public Tuple2<Boolean, byte[]> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<Boolean, byte[]>(
                                (Boolean) results.get(0).getValue(), 
                                (byte[]) results.get(1).getValue());
                    }
                });
    }

    public RemoteCall<Boolean> throwCheck(Boolean param) {
        final Function function = new Function(FUNC_THROWCHECK, 
                Arrays.<Type>asList(new Bool(param)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public static RemoteCall<DeprecatedFunctions> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DeprecatedFunctions.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<DeprecatedFunctions> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DeprecatedFunctions.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static DeprecatedFunctions load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new DeprecatedFunctions(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static DeprecatedFunctions load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new DeprecatedFunctions(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
