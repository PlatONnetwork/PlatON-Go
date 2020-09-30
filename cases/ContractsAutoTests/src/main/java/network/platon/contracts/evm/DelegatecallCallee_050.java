package network.platon.contracts.evm;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import org.web3j.abi.EventEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
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
 * <p>Generated with web3j version 0.13.1.5.
 */
public class DelegatecallCallee_050 extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610184806100206000396000f3fe608060405260043610610057576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630c55699c1461005c578063371303c0146100875780635a3617561461009e575b600080fd5b34801561006857600080fd5b506100716100c9565b6040518082815260200191505060405180910390f35b34801561009357600080fd5b5061009c6100cf565b005b3480156100aa57600080fd5b506100b361014f565b6040518082815260200191505060405180910390f35b60005481565b60008081548092919060010191905055507fb0333e0e3a6b99318e4e2e0d7e5e5f93646f9cbf62da1587955a4092bf7df6e733600054604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a1565b6000805490509056fea165627a7a72305820b62481330f6ec2f461adc90a7dff1037c65b972cf20bd6a6e20c6c607c4377ac0029";

    public static final String FUNC_X = "x";

    public static final String FUNC_INC = "inc";

    public static final String FUNC_GETCALLEEX = "getCalleeX";

    public static final Event EVENTNAME_EVENT = new Event("EventName", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}, new TypeReference<Uint256>() {}));
    ;

    protected DelegatecallCallee_050(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected DelegatecallCallee_050(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> x() {
        final Function function = new Function(FUNC_X, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> inc() {
        final Function function = new Function(
                FUNC_INC, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> getCalleeX() {
        final Function function = new Function(FUNC_GETCALLEEX, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public List<EventNameEventResponse> getEventNameEvents(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(EVENTNAME_EVENT, transactionReceipt);
        ArrayList<EventNameEventResponse> responses = new ArrayList<EventNameEventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
            EventNameEventResponse typedResponse = new EventNameEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.seder = (String) eventValues.getNonIndexedValues().get(0).getValue();
            typedResponse.x = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<EventNameEventResponse> eventNameEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, EventNameEventResponse>() {
            @Override
            public EventNameEventResponse call(Log log) {
                EventValuesWithLog eventValues = extractEventParametersWithLog(EVENTNAME_EVENT, log);
                EventNameEventResponse typedResponse = new EventNameEventResponse();
                typedResponse.log = log;
                typedResponse.seder = (String) eventValues.getNonIndexedValues().get(0).getValue();
                typedResponse.x = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<EventNameEventResponse> eventNameEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(EVENTNAME_EVENT));
        return eventNameEventObservable(filter);
    }

    public static RemoteCall<DelegatecallCallee_050> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DelegatecallCallee_050.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<DelegatecallCallee_050> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DelegatecallCallee_050.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static DelegatecallCallee_050 load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new DelegatecallCallee_050(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static DelegatecallCallee_050 load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new DelegatecallCallee_050(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static class EventNameEventResponse {
        public Log log;

        public String seder;

        public BigInteger x;
    }
}
