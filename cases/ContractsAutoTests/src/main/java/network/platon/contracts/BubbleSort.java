package network.platon.contracts;

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
 * <p>Generated with web3j version 0.13.0.7.
 */
public class BubbleSort extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061034c806100206000396000f3fe6080604052600436106100295760003560e01c8063970f17bd1461002e578063f6dd00aa1461009a575b600080fd5b34801561003a57600080fd5b5061004361015c565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b8381101561008657808201518184015260208101905061006b565b505050509050019250505060405180910390f35b61015a600480360360408110156100b057600080fd5b81019080803590602001906401000000008111156100cd57600080fd5b8201836020820111156100df57600080fd5b8035906020019184602083028401116401000000008311171561010157600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f820116905080830192505050505050509192919290803590602001909291905050506101b4565b005b606060008054806020026020016040519081016040528092919081815260200182805480156101aa57602002820191906000526020600020905b815481526020019060010190808311610196575b5050505050905090565b60008090505b600182038110156102895760008090505b81600184030381101561027b578360018201815181106101e757fe5b60200260200101518482815181106101fb57fe5b6020026020010151131561026e57600084828151811061021757fe5b6020026020010151905084600183018151811061023057fe5b602002602001015185838151811061024457fe5b6020026020010181815250508085600184018151811061026057fe5b602002602001018181525050505b80806001019150506101cb565b5080806001019150506101ba565b5081600090805190602001906102a09291906102a5565b505050565b8280548282559060005260206000209081019282156102e1579160200282015b828111156102e05782518255916020019190600101906102c5565b5b5090506102ee91906102f2565b5090565b61031491905b808211156103105760008160009055506001016102f8565b5090565b9056fea265627a7a723158205d0df03514919524fc77dec9d704c4711ddd724525cab8142b88d41a06e6bb6464736f6c634300050d0032";

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
                Arrays.<Type>asList(new org.web3j.abi.datatypes.DynamicArray<org.web3j.abi.datatypes.generated.Int256>(
                        org.web3j.abi.Utils.typeMap(arr, org.web3j.abi.datatypes.generated.Int256.class)), 
                new org.web3j.abi.datatypes.generated.Uint256(n)), 
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
