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
import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.0.
 */
public class PramaAndReturns extends Contract {
    private static final String BINARY = "60806040526040518060400160405280600d81526020017f576861742773207570206d616e000000000000000000000000000000000000008152506000908051906020019061004f929190610062565b5034801561005c57600080fd5b506100ff565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100a357805160ff19168380011785556100d1565b828001600101855582156100d1579182015b828111156100d05782518255916020019190600101906100b5565b5b5090506100de91906100e2565b5090565b5b808211156100fb5760008160009055506001016100e3565b5090565b6105788061010e6000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c806386b714e21161005b57806386b714e2146102905780639e2eea06146102ae578063e93314ab146102fa578063f8adff321461031857610088565b80630965e1451461008d5780631aa72bf81461012e57806354d410c8146101b15780637f0ffe3114610258575b600080fd5b6100f0600480360360608110156100a357600080fd5b8101908080606001906003806020026040519081016040528092919082600360200280828437600081840152601f19601f820116905080830192505050505050919291929050505061035a565b6040518082600360200280838360005b8381101561011b578082015181840152602081019050610100565b5050505090500191505060405180910390f35b610136610382565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561017657808201518184015260208101905061015b565b50505050905090810190601f1680156101a35780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101b9610424565b604051808060200180602001838103835285818151815260200191508051906020019060200280838360005b838110156102005780820151818401526020810190506101e5565b50505050905001838103825284818151815260200191508051906020019060200280838360005b83811015610242578082015181840152602081019050610227565b5050505090500194505050505060405180910390f35b61028e6004803603604081101561026e57600080fd5b8101908080359060200190929190803590602001909291905050506104f0565b005b6102986104fb565b6040518082815260200191505060405180910390f35b6102e4600480360360408110156102c457600080fd5b810190808035906020019092919080359060200190929190505050610501565b6040518082815260200191505060405180910390f35b61030261050c565b6040518082815260200191505060405180910390f35b6103446004803603602081101561032e57600080fd5b8101908080359060200190929190505050610516565b6040518082815260200191505060405180910390f35b610362610520565b60038260026003811061037157fe5b602002018181525050819050919050565b606060008054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561041a5780601f106103ef5761010080835404028352916020019161041a565b820191906000526020600020905b8154815290600101906020018083116103fd57829003601f168201915b5050505050905090565b6060806060600367ffffffffffffffff8111801561044157600080fd5b506040519080825280602002602001820160405280156104705781602001602082028036833780820191505090505b50905060018160008151811061048257fe5b60200260200101818152505060028160018151811061049d57fe5b6020026020010181815250506003816002815181106104b857fe5b6020026020010181815250506060819050600a826000815181106104d857fe5b60200260200101818152505081819350935050509091565b816001819055505050565b60015481565b600082905092915050565b6000600154905090565b6000819050919050565b604051806060016040528060039060208202803683378082019150509050509056fea2646970667358221220ecc088b1e90ba499501aac66c6586642f911caebbef15d42fa29d8a7ff2ddca964736f6c63430007010033";

    public static final String FUNC_INPUTPARAM = "InputParam";

    public static final String FUNC_IUPUTARRAY = "IuputArray";

    public static final String FUNC_NOOUTPUT = "NoOutput";

    public static final String FUNC_OMITPARAM = "OmitParam";

    public static final String FUNC_OUPUTARRAYS = "OuputArrays";

    public static final String FUNC_OUPUTSTRING = "OuputString";

    public static final String FUNC_GETS = "getS";

    public static final String FUNC_S = "s";

    protected PramaAndReturns(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected PramaAndReturns(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> InputParam(BigInteger a) {
        final Function function = new Function(
                FUNC_INPUTPARAM, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(a)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> IuputArray(List<BigInteger> y) {
        final Function function = new Function(
                FUNC_IUPUTARRAY, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.StaticArray3<com.alaya.abi.solidity.datatypes.generated.Uint256>(
                com.alaya.abi.solidity.datatypes.generated.Uint256.class,
                        com.alaya.abi.solidity.Utils.typeMap(y, com.alaya.abi.solidity.datatypes.generated.Uint256.class))), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> NoOutput(BigInteger a, BigInteger b) {
        final Function function = new Function(
                FUNC_NOOUTPUT, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(a), 
                new com.alaya.abi.solidity.datatypes.generated.Uint256(b)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> OmitParam(BigInteger y, BigInteger param1) {
        final Function function = new Function(
                FUNC_OMITPARAM, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(y), 
                new com.alaya.abi.solidity.datatypes.generated.Uint256(param1)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> OuputArrays() {
        final Function function = new Function(
                FUNC_OUPUTARRAYS, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> OuputString() {
        final Function function = new Function(
                FUNC_OUPUTSTRING, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getS() {
        final Function function = new Function(
                FUNC_GETS, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> s() {
        final Function function = new Function(
                FUNC_S, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<PramaAndReturns> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(PramaAndReturns.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<PramaAndReturns> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(PramaAndReturns.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static PramaAndReturns load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new PramaAndReturns(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static PramaAndReturns load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new PramaAndReturns(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
