package network.platon.contracts.evm.v0_7_1;

import com.alaya.abi.solidity.EventEncoder;
import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Address;
import com.alaya.abi.solidity.datatypes.Event;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.DefaultBlockParameter;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.protocol.core.methods.request.PlatonFilter;
import com.alaya.protocol.core.methods.response.Log;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import rx.Observable;
import rx.functions.Func1;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.1.
 */
public class CallerTwo extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506103a7806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80630c55699c14610046578063371303c0146100645780635a3617561461006e575b600080fd5b61004e61008c565b6040518082815260200191505060405180910390f35b61006c610092565b005b61007661020b565b6040518082815260200191505060405180910390f35b60005481565b60006040516100a090610214565b604051809103906000f0801580156100bc573d6000803e3d6000fd5b5090508073ffffffffffffffffffffffffffffffffffffffff167f371303c051bff726100ad13871cababf50c20dd920fca137e519f98f089a74b4604051602001808281526020019150506040516020818303038152906040526040518082805190602001908083835b602083106101495780518252602082019150602081019050602083039250610126565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381855af49150503d80600081146101a9576040519150601f19603f3d011682016040523d82523d6000602084013e6101ae565b606091505b5050507fb0333e0e3a6b99318e4e2e0d7e5e5f93646f9cbf62da1587955a4092bf7df6e733600054604051808373ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a150565b60008054905090565b610150806102228339019056fe608060405234801561001057600080fd5b50610130806100206000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80630c55699c14604157806317f936fb14605d578063371303c0146079575b600080fd5b60476081565b6040518082815260200191505060405180910390f35b60636087565b6040518082815260200191505060405180910390f35b607f6090565b005b60005481565b60008054905090565b60008081548092919060010191905055507fb0333e0e3a6b99318e4e2e0d7e5e5f93646f9cbf62da1587955a4092bf7df6e733600054604051808373ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a156fea2646970667358221220cad78347f3612ff78c53aa548f2a06dd0a405e38230359614191b0c7081bbd2364736f6c63430007010033a26469706673582212209d23f1a2cbf08731a2aeee1a8e86b1fb1df68ecbb44f19c803a453e5802afd1964736f6c63430007010033";

    public static final String FUNC_GETCALLEEX = "getCalleeX";

    public static final String FUNC_INC = "inc";

    public static final String FUNC_X = "x";

    public static final Event EVENTNAME_EVENT = new Event("EventName", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}, new TypeReference<Uint256>() {}));
    ;

    protected CallerTwo(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected CallerTwo(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
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

    public RemoteCall<BigInteger> getCalleeX() {
        final Function function = new Function(FUNC_GETCALLEEX, 
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

    public RemoteCall<BigInteger> x() {
        final Function function = new Function(FUNC_X, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<CallerTwo> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(CallerTwo.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<CallerTwo> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(CallerTwo.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static CallerTwo load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new CallerTwo(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static CallerTwo load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new CallerTwo(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static class EventNameEventResponse {
        public Log log;

        public String seder;

        public BigInteger x;
    }
}
