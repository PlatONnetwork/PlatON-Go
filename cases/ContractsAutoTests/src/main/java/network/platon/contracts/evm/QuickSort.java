package network.platon.contracts.evm;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.DynamicArray;
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
 * <p>Generated with web3j version 0.13.1.5.
 */
public class QuickSort extends Contract {
    private static final String BINARY = "6060604052341561000f57600080fd5b5b5b5b610475806100216000396000f30060606040526000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680631703539c1461005457806371e5ee5f146100c0578063cc80f6f3146100f7575b600080fd5b341561005f57600080fd5b6100be600480803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843782019150505050505091908035906020019091908035906020019091905050610162565b005b34156100cb57600080fd5b6100e1600480803590602001909190505061018a565b6040518082815260200191505060405180910390f35b341561010257600080fd5b61010a6101af565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b8381101561014e5780820151818401525b602081019050610132565b505050509050019250505060405180910390f35b61016d83838361020e565b82600090805190602001906101839291906103c3565b505b505050565b60008181548110151561019957fe5b906000526020600020900160005b915090505481565b6101b7610410565b600080548060200260200160405190810160405280929190818152602001828054801561020357602002820191906000526020600020905b8154815260200190600101908083116101ef575b505050505090505b90565b60008183101561024d57610223848484610254565b905060008114151561023e5761023d84846001840361020e565b5b61024c84600183018461020e565b5b5b50505050565b600080600080868681518110151561026857fe5b9060200190602002015192508591508490505b8082141515610398575b80821080156102aa575082878281518110151561029e57fe5b90602001906020020151135b156102ba57600181039050610285565b808210156102fd5786818151811015156102d057fe5b9060200190602002015187838151811015156102e857fe5b90602001906020020181815250506001820191505b5b8082108015610323575082878381518110151561031757fe5b90602001906020020151125b15610333576001820191506102fe565b8082101561037657868281518110151561034957fe5b90602001906020020151878281518110151561036157fe5b90602001906020020181815250506001810390505b82878381518110151561038557fe5b906020019060200201818152505061027b565b8287838151811015156103a757fe5b90602001906020020181815250508193505b5050509392505050565b8280548282559060005260206000209081019282156103ff579160200282015b828111156103fe5782518255916020019190600101906103e3565b5b50905061040c9190610424565b5090565b602060405190810160405280600081525090565b61044691905b8082111561044257600081600090555060010161042a565b5090565b905600a165627a7a723058201a5f1caa593b4db19f51a2c1190a53af9536cf6dffd3ccc5e8d206e6cd4cface0029";

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
                        org.web3j.abi.Utils.typeMap(_arr, Int256.class)),
                new org.web3j.abi.datatypes.generated.Uint256(low), 
                new org.web3j.abi.datatypes.generated.Uint256(high)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> arr(BigInteger param0) {
        final Function function = new Function(FUNC_ARR, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param0)), 
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
