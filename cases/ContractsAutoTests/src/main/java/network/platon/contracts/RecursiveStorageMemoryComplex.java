package network.platon.contracts;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.DynamicArray;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
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
public class RecursiveStorageMemoryComplex extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060426000800181905550600260006001018161002d91906101af565b50614200600060010160008154811061004257fe5b906000526020600020906002020160000181905550614201600060010160018154811061006b57fe5b9060005260206000209060020201600001819055506003600060010160008154811061009357fe5b9060005260206000209060020201600101816100af91906101af565b5060008090505b60038110156101145780624200000160006001016000815481106100d657fe5b906000526020600020906002020160010182815481106100f257fe5b90600052602060002090600202016000018190555080806001019150506100b6565b506004600060010160018154811061012857fe5b90600052602060002090600202016001018161014491906101af565b5060008090505b60048110156101a957806242010001600060010160018154811061016b57fe5b9060005260206000209060020201600101828154811061018757fe5b906000526020600020906002020160000181905550808060010191505061014b565b5061023c565b8154818355818111156101dc576002028160020283600052602060002091820191016101db91906101e1565b5b505050565b61021591905b80821115610211576000808201600090556001820160006102089190610218565b506002016101e7565b5090565b90565b508054600082556002029060005260206000209081019061023991906101e1565b50565b61048f8061024b6000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063a9be8c391461003b578063c04062261461009a575b600080fd5b6100436100f9565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b8381101561008657808201518184015260208101905061006b565b505050509050019250505060405180910390f35b6100a2610151565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b838110156100e55780820151818401526020810190506100ca565b505050509050019250505060405180910390f35b6060600280548060200260200160405190810160405280929190818152602001828054801561014757602002820191906000526020600020905b815481526020019060010190808311610133575b5050505050905090565b606061015b6102e4565b6000610166906102fe565b905060006101738261021f565b9050806040519080825280602002602001820160405280156101a45781602001602082028038833980820191505090505b50600290805190602001906101ba929190610373565b506101c682600061026f565b50600280548060200260200160405190810160405280929190818152602001828054801561021357602002820191906000526020600020905b8154815260200190600101908083116101ff575b50505050509250505090565b60006001905060008090505b826020015151811015610269576102588360200151828151811061024b57fe5b602002602001015161021f565b82019150808060010191505061022b565b50919050565b600082600001516002838060010194508154811061028957fe5b906000526020600020018190555060008090505b8360200151518110156102da576102cb846020015182815181106102bd57fe5b60200260200101518461026f565b9250808060010191505061029d565b5081905092915050565b604051806040016040528060008152602001606081525090565b6040518060400160405290816000820154815260200160018201805480602002602001604051908101604052809291908181526020016000905b828210156103685783829060005260206000209060020201610359906103c0565b81526020019060010190610338565b505050508152505090565b8280548282559060005260206000209081019282156103af579160200282015b828111156103ae578251825591602001919060010190610393565b5b5090506103bc9190610435565b5090565b6040518060400160405290816000820154815260200160018201805480602002602001604051908101604052809291908181526020016000905b8282101561042a578382906000526020600020906002020161041b906103c0565b815260200190600101906103fa565b505050508152505090565b61045791905b8082111561045357600081600090555060010161043b565b5090565b9056fea265627a7a723158202e8e44533423da30101996e1917c41c36796fba3b9bf24f78310b522fbbc2bb264736f6c634300050d0032";

    public static final String FUNC_GETRUNRESULT = "getRunResult";

    public static final String FUNC_RUN = "run";

    protected RecursiveStorageMemoryComplex(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected RecursiveStorageMemoryComplex(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<RecursiveStorageMemoryComplex> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(RecursiveStorageMemoryComplex.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<RecursiveStorageMemoryComplex> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(RecursiveStorageMemoryComplex.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public RemoteCall<List> getRunResult() {
        final Function function = new Function(FUNC_GETRUNRESULT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicArray<Uint256>>() {}));
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

    public RemoteCall<TransactionReceipt> run() {
        final Function function = new Function(
                FUNC_RUN, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RecursiveStorageMemoryComplex load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new RecursiveStorageMemoryComplex(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static RecursiveStorageMemoryComplex load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new RecursiveStorageMemoryComplex(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
