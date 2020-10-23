package network.platon.contracts.evm;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import org.web3j.abi.EventEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Bool;
import org.web3j.abi.datatypes.Event;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.DefaultBlockParameter;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.request.PlatonFilter;
import org.web3j.protocol.core.methods.response.Log;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.GasProvider;
import rx.Observable;
import rx.functions.Func1;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.1.5.
 */
public class LibraryUsingFor extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610183806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063f207564e14610030575b600080fd5b61005c6004803603602081101561004657600080fd5b8101908080359060200190929190505050610076565b604051808215151515815260200191505060405180910390f35b60008073f78d03c62b6701ff2b7506b69481cc87fd2f41ad637ae1e0589091846040518363ffffffff1660e01b8152600401808381526020018281526020019250505060206040518083038186803b1580156100d157600080fd5b505af41580156100e5573d6000803e3d6000fd5b505050506040513d60208110156100fb57600080fd5b810190808051906020019092919050505090507f0b3bdb70bcb1393d4319be3261bd6ab95e2ea1665e718029d24cecca39e84ccc81604051808215151515815260200191505060405180910390a191905056fea265627a7a72315820654b60685a929796a83c80bafca34b25c1b300665142db0f8623579c8582e48764736f6c634300050d0032";

    public static final String FUNC_REGISTER = "register";

    public static final Event RESULT_EVENT = new Event("Result", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
    ;

    protected LibraryUsingFor(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected LibraryUsingFor(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public List<ResultEventResponse> getResultEvents(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(RESULT_EVENT, transactionReceipt);
        ArrayList<ResultEventResponse> responses = new ArrayList<ResultEventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
            ResultEventResponse typedResponse = new ResultEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.result = (Boolean) eventValues.getNonIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<ResultEventResponse> resultEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, ResultEventResponse>() {
            @Override
            public ResultEventResponse call(Log log) {
                EventValuesWithLog eventValues = extractEventParametersWithLog(RESULT_EVENT, log);
                ResultEventResponse typedResponse = new ResultEventResponse();
                typedResponse.log = log;
                typedResponse.result = (Boolean) eventValues.getNonIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<ResultEventResponse> resultEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(RESULT_EVENT));
        return resultEventObservable(filter);
    }

    public RemoteCall<TransactionReceipt> register(BigInteger value) {
        final Function function = new Function(
                FUNC_REGISTER, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(value)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<LibraryUsingFor> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(LibraryUsingFor.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<LibraryUsingFor> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(LibraryUsingFor.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static LibraryUsingFor load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new LibraryUsingFor(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static LibraryUsingFor load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new LibraryUsingFor(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static class ResultEventResponse {
        public Log log;

        public Boolean result;
    }
}
