package network.platon.contracts.evm.v0_4_26;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.DynamicArray;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Int256;
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
import java.util.List;
import java.util.concurrent.Callable;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.1.
 */
public class QuickSort extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610476806100206000396000f300608060405260043610610057576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680631703539c1461005c57806371e5ee5f146100d6578063cc80f6f314610117575b600080fd5b34801561006857600080fd5b506100d4600480360381019080803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843782019150505050505091929192908035906020019092919080359060200190929190505050610183565b005b3480156100e257600080fd5b50610101600480360381019080803590602001909291905050506101aa565b6040518082815260200191505060405180910390f35b34801561012357600080fd5b5061012c6101cd565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b8381101561016f578082015181840152602081019050610154565b505050509050019250505060405180910390f35b61018e838383610225565b82600090805190602001906101a49291906103d8565b50505050565b6000818154811015156101b957fe5b906000526020600020016000915090505481565b6060600080548060200260200160405190810160405280929190818152602001828054801561021b57602002820191906000526020600020905b815481526020019060010190808311610207575b5050505050905090565b6000818310156102645761023a84848461026a565b905060008114151561025557610254848460018403610225565b5b610263846001830184610225565b5b50505050565b600080600080868681518110151561027e57fe5b9060200190602002015192508591508490505b80821415156103ae575b80821080156102c057508287828151811015156102b457fe5b90602001906020020151135b156102d05760018103905061029b565b808210156103135786818151811015156102e657fe5b9060200190602002015187838151811015156102fe57fe5b90602001906020020181815250506001820191505b5b8082108015610339575082878381518110151561032d57fe5b90602001906020020151125b1561034957600182019150610314565b8082101561038c57868281518110151561035f57fe5b90602001906020020151878281518110151561037757fe5b90602001906020020181815250506001810390505b82878381518110151561039b57fe5b9060200190602002018181525050610291565b8287838151811015156103bd57fe5b90602001906020020181815250508193505050509392505050565b828054828255906000526020600020908101928215610414579160200282015b828111156104135782518255916020019190600101906103f8565b5b5090506104219190610425565b5090565b61044791905b8082111561044357600081600090555060010161042b565b5090565b905600a165627a7a72305820441576b1cf4874aa9409b5eb4820c676afe9e0e5102182c119c2f4d76a90d46d0029";

    public static final String FUNC_SORT = "sort";

    public static final String FUNC_ARR = "arr";

    public static final String FUNC_SHOW = "show";

    protected QuickSort(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected QuickSort(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> sort(List<BigInteger> _arr, BigInteger low, BigInteger high) {
        final Function function = new Function(
                FUNC_SORT, 
                Arrays.<Type>asList(new DynamicArray<Int256>(
                Int256.class,
                        com.alaya.abi.solidity.Utils.typeMap(_arr, Int256.class)),
                new com.alaya.abi.solidity.datatypes.generated.Uint256(low), 
                new com.alaya.abi.solidity.datatypes.generated.Uint256(high)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> arr(BigInteger param0) {
        final Function function = new Function(FUNC_ARR, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(param0)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Int256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<List> show() {
        final Function function = new Function(FUNC_SHOW, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicArray<Int256>>() {}));
        return new RemoteCall<List>(
                new Callable<List>() {
                    @Override
                    @SuppressWarnings("unchecked")
                    public List call() throws Exception {
                        List<Type> result = (List<Type>) executeCallSingleValueReturn(function, List.class);
                        return convertToNative(result);
                    }
                });
    }

    public static RemoteCall<QuickSort> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(QuickSort.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<QuickSort> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(QuickSort.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static QuickSort load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new QuickSort(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static QuickSort load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new QuickSort(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
