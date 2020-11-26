package network.platon.contracts.evm.v0_6_12;

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
public class BubbleSort extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061033f806100206000396000f3fe6080604052600436106100295760003560e01c8063970f17bd1461002e578063f6dd00aa1461009a575b600080fd5b34801561003a57600080fd5b5061004361015c565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b8381101561008657808201518184015260208101905061006b565b505050509050019250505060405180910390f35b61015a600480360360408110156100b057600080fd5b81019080803590602001906401000000008111156100cd57600080fd5b8201836020820111156100df57600080fd5b8035906020019184602083028401116401000000008311171561010157600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f820116905080830192505050505050509192919290803590602001909291905050506101b4565b005b606060008054806020026020016040519081016040528092919081815260200182805480156101aa57602002820191906000526020600020905b815481526020019060010190808311610196575b5050505050905090565b60005b600182038110156102835760005b816001840303811015610275578360018201815181106101e157fe5b60200260200101518482815181106101f557fe5b6020026020010151131561026857600084828151811061021157fe5b6020026020010151905084600183018151811061022a57fe5b602002602001015185838151811061023e57fe5b6020026020010181815250508085600184018151811061025a57fe5b602002602001018181525050505b80806001019150506101c5565b5080806001019150506101b7565b50816000908051906020019061029a92919061029f565b505050565b8280548282559060005260206000209081019282156102db579160200282015b828111156102da5782518255916020019190600101906102bf565b5b5090506102e891906102ec565b5090565b5b808211156103055760008160009055506001016102ed565b509056fea2646970667358221220a3d65d818c1a0df7fb11fdc00da3d3524caa43571b919ce39050ccdb214ea92564736f6c634300060c0033";

    public static final String FUNC_BUBBLEARRAYS = "BubbleArrays";

    public static final String FUNC_GET_ARR = "get_arr";

    protected BubbleSort(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected BubbleSort(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> BubbleArrays(List<BigInteger> arr, BigInteger n, BigInteger vonValue) {
        final Function function = new Function(
                FUNC_BUBBLEARRAYS, 
                Arrays.<Type>asList(new DynamicArray<Int256>(
                Int256.class,
                        com.alaya.abi.solidity.Utils.typeMap(arr, Int256.class)),
                new com.alaya.abi.solidity.datatypes.generated.Uint256(n)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public RemoteCall<List> get_arr() {
        final Function function = new Function(FUNC_GET_ARR, 
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

    public static RemoteCall<BubbleSort> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(BubbleSort.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<BubbleSort> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(BubbleSort.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static BubbleSort load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new BubbleSort(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static BubbleSort load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new BubbleSort(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
