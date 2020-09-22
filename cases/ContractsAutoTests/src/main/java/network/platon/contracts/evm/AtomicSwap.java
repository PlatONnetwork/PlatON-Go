package network.platon.contracts.evm;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.EventEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.Event;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Bytes32;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.abi.datatypes.generated.Uint8;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.DefaultBlockParameter;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.request.PlatonFilter;
import org.web3j.protocol.core.methods.response.Log;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple8;
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
public class AtomicSwap extends Contract {
    private static final String BINARY = "6080604052336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506114eb806100536000396000f3fe608060405260043610610078576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680632b68b9c61461007d578063446bffba1461009457806348e558da146100c25780637249fbb614610124578063b31597ad1461015f578063eb84e7f2146101a4575b600080fd5b34801561008957600080fd5b5061009261028a565b005b6100c0600480360360208110156100aa57600080fd5b8101908080359060200190929190505050610417565b005b610122600480360360808110156100d857600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919080359060200190929190505050610648565b005b34801561013057600080fd5b5061015d6004803603602081101561014757600080fd5b8101908080359060200190929190505050610a30565b005b34801561016b57600080fd5b506101a26004803603604081101561018257600080fd5b810190808035906020019092919080359060200190929190505050610d5c565b005b3480156101b057600080fd5b506101dd600480360360208110156101c757600080fd5b810190808035906020019092919050505061131d565b604051808981526020018881526020018773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200185815260200184815260200183815260200182600381111561026f57fe5b60ff1681526020019850505050505050505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561034e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600a8152602001807f6f6e6c79206f776e65720000000000000000000000000000000000000000000081525060200191505060405180910390fd5b60003073ffffffffffffffffffffffffffffffffffffffff16311415156103dd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f62616c616e6365206973206e6f74207a65726f0000000000000000000000000081525060200191505060405180910390fd5b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16ff5b806001600381111561042557fe5b6001600083815260200190815260200160002060070160009054906101000a900460ff16600381111561045457fe5b1415156104ef576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602c8152602001807f7377617020666f722074686973206861736820697320656d707479206f72206181526020017f6c7265616479207370656e74000000000000000000000000000000000000000081525060400191505060405180910390fd5b816001600082815260200190815260200160002060040154421115151561057e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f726566756e6454696d652068617320616c726561647920636f6d65000000000081525060200191505060405180910390fd5b6105a73460016000868152602001908152602001600020600501546113b290919063ffffffff16565b6001600085815260200190815260200160002060050181905550827fd760a88b05be4d78a2815eb20f72049b7c89e1dca4fc467139fe3f2224a37423336001600087815260200190815260200160002060050154604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a2505050565b83826000600381111561065757fe5b6001600084815260200190815260200160002060070160009054906101000a900460ff16600381111561068657fe5b141515610721576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260278152602001807f7377617020666f722074686973206861736820697320616c726561647920696e81526020017f697469617465640000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b42811115156107be576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001807f726566756e6454696d657374616d702068617320616c7265616479207061737381526020017f656400000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b6107d1833461143a90919063ffffffff16565b6001600088815260200190815260200160002060050181905550856001600088815260200190815260200160002060000181905550336001600088815260200190815260200160002060020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550846001600088815260200190815260200160002060030160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550836001600088815260200190815260200160002060040181905550826001600088815260200190815260200160002060060181905550600180600088815260200190815260200160002060070160006101000a81548160ff0219169083600381111561091857fe5b02179055506001600087815260200190815260200160002060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16867f5e919055312829285818d366d1cfe50a1ba27ce2c752b655cb2faa0179e1422733600160008b815260200190815260200160002060040154600160008c815260200190815260200160002060050154600160008d815260200190815260200160002060060154604051808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200184815260200183815260200182815260200194505050505060405180910390a3505050505050565b8060016003811115610a3e57fe5b6001600083815260200190815260200160002060070160009054906101000a900460ff166003811115610a6d57fe5b141515610b08576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602c8152602001807f7377617020666f722074686973206861736820697320656d707479206f72206181526020017f6c7265616479207370656e74000000000000000000000000000000000000000081525060400191505060405180910390fd5b8160016000828152602001908152602001600020600401544210151515610b97576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f726566756e6454696d657374616d7020686173206e6f7420706173736564000081525060200191505060405180910390fd5b60036001600085815260200190815260200160002060070160006101000a81548160ff02191690836003811115610bca57fe5b0217905550827ffe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf060405160405180910390a26001600084815260200190815260200160002060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc610c8b600160008781526020019081526020016000206006015460016000888152602001908152602001600020600501546113b290919063ffffffff16565b9081150290604051600060405180830381858888f19350505050158015610cb6573d6000803e3d6000fd5b506001600084815260200190815260200160002060008082016000905560018201600090556002820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690556003820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690556004820160009055600582016000905560068201600090556007820160006101000a81549060ff02191690555050505050565b8160016003811115610d6a57fe5b6001600083815260200190815260200160002060070160009054906101000a900460ff166003811115610d9957fe5b141515610e34576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602c8152602001807f7377617020666f722074686973206861736820697320656d707479206f72206181526020017f6c7265616479207370656e74000000000000000000000000000000000000000081525060400191505060405180910390fd5b8282600160008381526020019081526020016000206004015442101515610ee9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001807f726566756e6454696d657374616d702068617320616c7265616479207061737381526020017f656400000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b8160028083604051602001808281526020019150506040516020818303038152906040526040518082805190602001908083835b602083101515610f425780518252602082019150602081019050602083039250610f1d565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa158015610f84573d6000803e3d6000fd5b5050506040513d6020811015610f9957600080fd5b8101908080519060200190929190505050604051602001808281526020019150506040516020818303038152906040526040518082805190602001908083835b602083101515610ffe5780518252602082019150602081019050602083039250610fd9565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa158015611040573d6000803e3d6000fd5b5050506040513d602081101561105557600080fd5b81019080805190602001909291905050501415156110db576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f736563726574206973206e6f7420636f7272656374000000000000000000000081525060200191505060405180910390fd5b83600160008781526020019081526020016000206001018190555060026001600087815260200190815260200160002060070160006101000a81548160ff0219169083600381111561112957fe5b0217905550847f489e9ee921192823d1aa1ef800c9ffc642993538b1e7e43a4d46a91965e894ab856040518082815260200191505060405180910390a26001600086815260200190815260200160002060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc60016000888152602001908152602001600020600501549081150290604051600060405180830381858888f193505050501580156111f8573d6000803e3d6000fd5b50600060016000878152602001908152602001600020600601541115611276573373ffffffffffffffffffffffffffffffffffffffff166108fc60016000888152602001908152602001600020600601549081150290604051600060405180830381858888f19350505050158015611274573d6000803e3d6000fd5b505b6001600086815260200190815260200160002060008082016000905560018201600090556002820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690556003820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690556004820160009055600582016000905560068201600090556007820160006101000a81549060ff021916905550505050505050565b60016020528060005260406000206000915090508060000154908060010154908060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060040154908060050154908060060154908060070160009054906101000a900460ff16905088565b60008183019050828110151515611431576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f536166654d617468206164642077726f6e672076616c7565000000000000000081525060200191505060405180910390fd5b80905092915050565b60008282111515156114b4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f536166654d617468207375622077726f6e672076616c7565000000000000000081525060200191505060405180910390fd5b81830390509291505056fea165627a7a72305820262d7838a7f9c176cd2535a9b6696914bdc215ab8f7e0415663609ca73c3ebb50029";

    public static final String FUNC_DESTRUCT = "destruct";

    public static final String FUNC_ADD = "add";

    public static final String FUNC_INITIATE = "initiate";

    public static final String FUNC_REFUND = "refund";

    public static final String FUNC_REDEEM = "redeem";

    public static final String FUNC_SWAPS = "swaps";

    public static final Event INITIATED_EVENT = new Event("Initiated", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>(true) {}, new TypeReference<Address>(true) {}, new TypeReference<Address>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}));
    ;

    public static final Event ADDED_EVENT = new Event("Added", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>(true) {}, new TypeReference<Address>() {}, new TypeReference<Uint256>() {}));
    ;

    public static final Event REDEEMED_EVENT = new Event("Redeemed", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>(true) {}, new TypeReference<Bytes32>() {}));
    ;

    public static final Event REFUNDED_EVENT = new Event("Refunded", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>(true) {}));
    ;

    protected AtomicSwap(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected AtomicSwap(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> destruct() {
        final Function function = new Function(
                FUNC_DESTRUCT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> add(byte[] _hashedSecret, BigInteger vonValue) {
        final Function function = new Function(
                FUNC_ADD, 
                Arrays.<Type>asList(new Bytes32(_hashedSecret)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public RemoteCall<TransactionReceipt> initiate(byte[] _hashedSecret, String _participant, BigInteger _refundTimestamp, BigInteger _payoff, BigInteger vonValue) {
        final Function function = new Function(
                FUNC_INITIATE, 
                Arrays.<Type>asList(new Bytes32(_hashedSecret),
                new Address(_participant),
                new Uint256(_refundTimestamp),
                new Uint256(_payoff)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public RemoteCall<TransactionReceipt> refund(byte[] _hashedSecret) {
        final Function function = new Function(
                FUNC_REFUND, 
                Arrays.<Type>asList(new Bytes32(_hashedSecret)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> redeem(byte[] _hashedSecret, byte[] _secret) {
        final Function function = new Function(
                FUNC_REDEEM, 
                Arrays.<Type>asList(new Bytes32(_hashedSecret),
                new Bytes32(_secret)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<Tuple8<byte[], byte[], String, String, BigInteger, BigInteger, BigInteger, BigInteger>> swaps(byte[] param0) {
        final Function function = new Function(FUNC_SWAPS, 
                Arrays.<Type>asList(new Bytes32(param0)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}, new TypeReference<Bytes32>() {}, new TypeReference<Address>() {}, new TypeReference<Address>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint8>() {}));
        return new RemoteCall<Tuple8<byte[], byte[], String, String, BigInteger, BigInteger, BigInteger, BigInteger>>(
                new Callable<Tuple8<byte[], byte[], String, String, BigInteger, BigInteger, BigInteger, BigInteger>>() {
                    @Override
                    public Tuple8<byte[], byte[], String, String, BigInteger, BigInteger, BigInteger, BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple8<byte[], byte[], String, String, BigInteger, BigInteger, BigInteger, BigInteger>(
                                (byte[]) results.get(0).getValue(), 
                                (byte[]) results.get(1).getValue(), 
                                (String) results.get(2).getValue(), 
                                (String) results.get(3).getValue(), 
                                (BigInteger) results.get(4).getValue(), 
                                (BigInteger) results.get(5).getValue(), 
                                (BigInteger) results.get(6).getValue(), 
                                (BigInteger) results.get(7).getValue());
                    }
                });
    }

    public List<InitiatedEventResponse> getInitiatedEvents(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(INITIATED_EVENT, transactionReceipt);
        ArrayList<InitiatedEventResponse> responses = new ArrayList<InitiatedEventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
            InitiatedEventResponse typedResponse = new InitiatedEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse._hashedSecret = (byte[]) eventValues.getIndexedValues().get(0).getValue();
            typedResponse._participant = (String) eventValues.getIndexedValues().get(1).getValue();
            typedResponse._initiator = (String) eventValues.getNonIndexedValues().get(0).getValue();
            typedResponse._refundTimestamp = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
            typedResponse._value = (BigInteger) eventValues.getNonIndexedValues().get(2).getValue();
            typedResponse._payoff = (BigInteger) eventValues.getNonIndexedValues().get(3).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<InitiatedEventResponse> initiatedEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, InitiatedEventResponse>() {
            @Override
            public InitiatedEventResponse call(Log log) {
                EventValuesWithLog eventValues = extractEventParametersWithLog(INITIATED_EVENT, log);
                InitiatedEventResponse typedResponse = new InitiatedEventResponse();
                typedResponse.log = log;
                typedResponse._hashedSecret = (byte[]) eventValues.getIndexedValues().get(0).getValue();
                typedResponse._participant = (String) eventValues.getIndexedValues().get(1).getValue();
                typedResponse._initiator = (String) eventValues.getNonIndexedValues().get(0).getValue();
                typedResponse._refundTimestamp = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
                typedResponse._value = (BigInteger) eventValues.getNonIndexedValues().get(2).getValue();
                typedResponse._payoff = (BigInteger) eventValues.getNonIndexedValues().get(3).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<InitiatedEventResponse> initiatedEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(INITIATED_EVENT));
        return initiatedEventObservable(filter);
    }

    public List<AddedEventResponse> getAddedEvents(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(ADDED_EVENT, transactionReceipt);
        ArrayList<AddedEventResponse> responses = new ArrayList<AddedEventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
            AddedEventResponse typedResponse = new AddedEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse._hashedSecret = (byte[]) eventValues.getIndexedValues().get(0).getValue();
            typedResponse._sender = (String) eventValues.getNonIndexedValues().get(0).getValue();
            typedResponse._value = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<AddedEventResponse> addedEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, AddedEventResponse>() {
            @Override
            public AddedEventResponse call(Log log) {
                EventValuesWithLog eventValues = extractEventParametersWithLog(ADDED_EVENT, log);
                AddedEventResponse typedResponse = new AddedEventResponse();
                typedResponse.log = log;
                typedResponse._hashedSecret = (byte[]) eventValues.getIndexedValues().get(0).getValue();
                typedResponse._sender = (String) eventValues.getNonIndexedValues().get(0).getValue();
                typedResponse._value = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<AddedEventResponse> addedEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(ADDED_EVENT));
        return addedEventObservable(filter);
    }

    public List<RedeemedEventResponse> getRedeemedEvents(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(REDEEMED_EVENT, transactionReceipt);
        ArrayList<RedeemedEventResponse> responses = new ArrayList<RedeemedEventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
            RedeemedEventResponse typedResponse = new RedeemedEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse._hashedSecret = (byte[]) eventValues.getIndexedValues().get(0).getValue();
            typedResponse._secret = (byte[]) eventValues.getNonIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<RedeemedEventResponse> redeemedEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, RedeemedEventResponse>() {
            @Override
            public RedeemedEventResponse call(Log log) {
                EventValuesWithLog eventValues = extractEventParametersWithLog(REDEEMED_EVENT, log);
                RedeemedEventResponse typedResponse = new RedeemedEventResponse();
                typedResponse.log = log;
                typedResponse._hashedSecret = (byte[]) eventValues.getIndexedValues().get(0).getValue();
                typedResponse._secret = (byte[]) eventValues.getNonIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<RedeemedEventResponse> redeemedEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(REDEEMED_EVENT));
        return redeemedEventObservable(filter);
    }

    public List<RefundedEventResponse> getRefundedEvents(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(REFUNDED_EVENT, transactionReceipt);
        ArrayList<RefundedEventResponse> responses = new ArrayList<RefundedEventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
            RefundedEventResponse typedResponse = new RefundedEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse._hashedSecret = (byte[]) eventValues.getIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<RefundedEventResponse> refundedEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, RefundedEventResponse>() {
            @Override
            public RefundedEventResponse call(Log log) {
                EventValuesWithLog eventValues = extractEventParametersWithLog(REFUNDED_EVENT, log);
                RefundedEventResponse typedResponse = new RefundedEventResponse();
                typedResponse.log = log;
                typedResponse._hashedSecret = (byte[]) eventValues.getIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<RefundedEventResponse> refundedEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(REFUNDED_EVENT));
        return refundedEventObservable(filter);
    }

    public static RemoteCall<AtomicSwap> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AtomicSwap.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<AtomicSwap> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AtomicSwap.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static AtomicSwap load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new AtomicSwap(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static AtomicSwap load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new AtomicSwap(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static class InitiatedEventResponse {
        public Log log;

        public byte[] _hashedSecret;

        public String _participant;

        public String _initiator;

        public BigInteger _refundTimestamp;

        public BigInteger _value;

        public BigInteger _payoff;
    }

    public static class AddedEventResponse {
        public Log log;

        public byte[] _hashedSecret;

        public String _sender;

        public BigInteger _value;
    }

    public static class RedeemedEventResponse {
        public Log log;

        public byte[] _hashedSecret;

        public byte[] _secret;
    }

    public static class RefundedEventResponse {
        public Log log;

        public byte[] _hashedSecret;
    }
}
