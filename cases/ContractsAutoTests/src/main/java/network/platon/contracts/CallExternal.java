package network.platon.contracts;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import org.web3j.abi.EventEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Event;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
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
 * <p>Generated with web3j version 0.13.0.7.
 */
public class CallExternal extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610221806100206000396000f3fe60806040526004361061001e5760003560e01c8063de29278914610023575b600080fd5b61002b610041565b6040518082815260200191505060405180910390f35b6000806040516100509061012d565b604051809103906000f08015801561006c573d6000803e3d6000fd5b5090508073ffffffffffffffffffffffffffffffffffffffff1663569c5f6d6040518163ffffffff1660e01b815260040160206040518083038186803b1580156100b557600080fd5b505afa1580156100c9573d6000803e3d6000fd5b505050506040513d60208110156100df57600080fd5b810190808051906020019092919050505091507f0a9f1213b326cb97c7a18f80791661027e1cf7a53125f3d6729d0ae093bd8ad2826040518082815260200191505060405180910390a15090565b60b38061013a8339019056fe6080604052348015600f57600080fd5b5060958061001e6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063569c5f6d14602d575b600080fd5b60336049565b6040518082815260200191505060405180910390f35b60008060019050600060029050808201925050509056fea265627a7a723158202e49189591614e050054613f168605019d1cb2b825ed23422a61e58bb317acd664736f6c634300050d0032a265627a7a723158208ae646dd80702eb4d9c97ed20c0257b8630a03f0aeee1196f164ae1abb0004df64736f6c634300050d0032";

    public static final String FUNC_GETRESULT = "getResult";

    public static final Event EXTERNALCVALUE_EVENT = new Event("ExternalCValue", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
    ;

    protected CallExternal(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected CallExternal(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public List<ExternalCValueEventResponse> getExternalCValueEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(EXTERNALCVALUE_EVENT, transactionReceipt);
        ArrayList<ExternalCValueEventResponse> responses = new ArrayList<ExternalCValueEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            ExternalCValueEventResponse typedResponse = new ExternalCValueEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.value = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<ExternalCValueEventResponse> externalCValueEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, ExternalCValueEventResponse>() {
            @Override
            public ExternalCValueEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(EXTERNALCVALUE_EVENT, log);
                ExternalCValueEventResponse typedResponse = new ExternalCValueEventResponse();
                typedResponse.log = log;
                typedResponse.value = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<ExternalCValueEventResponse> externalCValueEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(EXTERNALCVALUE_EVENT));
        return externalCValueEventObservable(filter);
    }

    public RemoteCall<TransactionReceipt> getResult(BigInteger vonValue) {
        final Function function = new Function(
                FUNC_GETRESULT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public static RemoteCall<CallExternal> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(CallExternal.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<CallExternal> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(CallExternal.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static CallExternal load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new CallExternal(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static CallExternal load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new CallExternal(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static class ExternalCValueEventResponse {
        public Log log;

        public BigInteger value;
    }
}
