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
import org.web3j.abi.datatypes.Utf8String;
import org.web3j.abi.datatypes.generated.StaticArray3;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple2;
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
public class PramaAndReturns extends Contract {
    private static final String BINARY = "60806040526040518060400160405280600d81526020017f576861742773207570206d616e000000000000000000000000000000000000008152506000908051906020019061004f929190610062565b5034801561005c57600080fd5b50610107565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100a357805160ff19168380011785556100d1565b828001600101855582156100d1579182015b828111156100d05782518255916020019190600101906100b5565b5b5090506100de91906100e2565b5090565b61010491905b808211156101005760008160009055506001016100e8565b5090565b90565b610563806101166000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c806386b714e21161005b57806386b714e2146102905780639e2eea06146102ae578063e93314ab146102fa578063f8adff321461031857610088565b80630965e1451461008d5780631aa72bf81461012e57806354d410c8146101b15780637f0ffe3114610258575b600080fd5b6100f0600480360360608110156100a357600080fd5b8101908080606001906003806020026040519081016040528092919082600360200280828437600081840152601f19601f820116905080830192505050505050919291929050505061035a565b6040518082600360200280838360005b8381101561011b578082015181840152602081019050610100565b5050505090500191505060405180910390f35b610136610382565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561017657808201518184015260208101905061015b565b50505050905090810190601f1680156101a35780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101b9610424565b604051808060200180602001838103835285818151815260200191508051906020019060200280838360005b838110156102005780820151818401526020810190506101e5565b50505050905001838103825284818151815260200191508051906020019060200280838360005b83811015610242578082015181840152602081019050610227565b5050505090500194505050505060405180910390f35b61028e6004803603604081101561026e57600080fd5b8101908080359060200190929190803590602001909291905050506104d9565b005b6102986104e4565b6040518082815260200191505060405180910390f35b6102e4600480360360408110156102c457600080fd5b8101908080359060200190929190803590602001909291905050506104ea565b6040518082815260200191505060405180910390f35b6103026104f5565b6040518082815260200191505060405180910390f35b6103446004803603602081101561032e57600080fd5b81019080803590602001909291905050506104ff565b6040518082815260200191505060405180910390f35b61036261050c565b60038260026003811061037157fe5b602002018181525050819050919050565b606060008054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561041a5780601f106103ef5761010080835404028352916020019161041a565b820191906000526020600020905b8154815290600101906020018083116103fd57829003601f168201915b5050505050905090565b606080606060036040519080825280602002602001820160405280156104595781602001602082028038833980820191505090505b50905060018160008151811061046b57fe5b60200260200101818152505060028160018151811061048657fe5b6020026020010181815250506003816002815181106104a157fe5b6020026020010181815250506060819050600a826000815181106104c157fe5b60200260200101818152505081819350935050509091565b816001819055505050565b60015481565b600082905092915050565b6000600154905090565b6000819050809050919050565b604051806060016040528060039060208202803883398082019150509050509056fea265627a7a723158208555080edb3c6d74ae21b0506fec977e2b2a47d18af2086f25a516bb1f18cc9a64736f6c634300050d0032";

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

    public RemoteCall<BigInteger> InputParam(BigInteger a) {
        final Function function = new Function(FUNC_INPUTPARAM, 
                Arrays.<Type>asList(new Uint256(a)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<List> IuputArray(List<BigInteger> y) {
        final Function function = new Function(FUNC_IUPUTARRAY, 
                Arrays.<Type>asList(new StaticArray3<Uint256>(
                Uint256.class,
                        org.web3j.abi.Utils.typeMap(y, Uint256.class))),
                Arrays.<TypeReference<?>>asList(new TypeReference<StaticArray3<Uint256>>() {}));
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

    public RemoteCall<TransactionReceipt> NoOutput(BigInteger a, BigInteger b) {
        final Function function = new Function(
                FUNC_NOOUTPUT, 
                Arrays.<Type>asList(new Uint256(a),
                new Uint256(b)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> OmitParam(BigInteger y, BigInteger param1) {
        final Function function = new Function(FUNC_OMITPARAM, 
                Arrays.<Type>asList(new Uint256(y),
                new Uint256(param1)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<Tuple2<List<BigInteger>, List<BigInteger>>> OuputArrays() {
        final Function function = new Function(FUNC_OUPUTARRAYS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicArray<Uint256>>() {}, new TypeReference<DynamicArray<Uint256>>() {}));
        return new RemoteCall<Tuple2<List<BigInteger>, List<BigInteger>>>(
                new Callable<Tuple2<List<BigInteger>, List<BigInteger>>>() {
                    @Override
                    public Tuple2<List<BigInteger>, List<BigInteger>> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<List<BigInteger>, List<BigInteger>>(
                                convertToNative((List<Uint256>) results.get(0).getValue()), 
                                convertToNative((List<Uint256>) results.get(1).getValue()));
                    }
                });
    }

    public RemoteCall<String> OuputString() {
        final Function function = new Function(FUNC_OUPUTSTRING, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> getS() {
        final Function function = new Function(FUNC_GETS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> s() {
        final Function function = new Function(FUNC_S, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
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
