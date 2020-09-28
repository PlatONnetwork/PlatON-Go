package network.platon.contracts.evm;

import com.alaya.abi.solidity.EventEncoder;
import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Address;
import com.alaya.abi.solidity.datatypes.Event;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Bytes32;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.abi.solidity.datatypes.generated.Uint8;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.DefaultBlockParameter;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.protocol.core.methods.request.PlatonFilter;
import com.alaya.protocol.core.methods.response.Log;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import com.alaya.tuples.generated.Tuple8;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Callable;
import rx.Observable;
import rx.functions.Func1;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.0.
 */
public class AtomicSwap extends Contract {
    private static final String BINARY = "6080604052336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555061137f806100536000396000f3fe6080604052600436106100555760003560e01c80632b68b9c61461005a578063446bffba1461007157806348e558da1461009f5780637249fbb614610101578063b31597ad1461013c578063eb84e7f214610181575b600080fd5b34801561006657600080fd5b5061006f610267565b005b61009d6004803603602081101561008757600080fd5b81019080803590602001909291905050506103d9565b005b6100ff600480360360808110156100b557600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190803590602001909291905050506105c3565b005b34801561010d57600080fd5b5061013a6004803603602081101561012457600080fd5b8101908080359060200190929190505050610921565b005b34801561014857600080fd5b5061017f6004803603604081101561015f57600080fd5b810190808035906020019092919080359060200190929190505050610c06565b005b34801561018d57600080fd5b506101ba600480360360208110156101a457600080fd5b8101908080359060200190929190505050611137565b604051808981526020018881526020018773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200185815260200184815260200183815260200182600381111561024c57fe5b60ff1681526020019850505050505050505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610329576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600a8152602001807f6f6e6c79206f776e65720000000000000000000000000000000000000000000081525060200191505060405180910390fd5b6000471461039f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f62616c616e6365206973206e6f74207a65726f0000000000000000000000000081525060200191505060405180910390fd5b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16ff5b80600160038111156103e757fe5b6001600083815260200190815260200160002060070160009054906101000a900460ff16600381111561041657fe5b1461046c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602c81526020018061131f602c913960400191505060405180910390fd5b8160016000828152602001908152602001600020600401544211156104f9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f726566756e6454696d652068617320616c726561647920636f6d65000000000081525060200191505060405180910390fd5b6105223460016000868152602001908152602001600020600501546111cc90919063ffffffff16565b6001600085815260200190815260200160002060050181905550827fd760a88b05be4d78a2815eb20f72049b7c89e1dca4fc467139fe3f2224a37423336001600087815260200190815260200160002060050154604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a2505050565b8382600060038111156105d257fe5b6001600084815260200190815260200160002060070160009054906101000a900460ff16600381111561060157fe5b14610657576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260278152602001806112f86027913960400191505060405180910390fd5b4281116106af576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001806112d66022913960400191505060405180910390fd5b6106c2833461125290919063ffffffff16565b6001600088815260200190815260200160002060050181905550856001600088815260200190815260200160002060000181905550336001600088815260200190815260200160002060020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550846001600088815260200190815260200160002060030160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550836001600088815260200190815260200160002060040181905550826001600088815260200190815260200160002060060181905550600180600088815260200190815260200160002060070160006101000a81548160ff0219169083600381111561080957fe5b02179055506001600087815260200190815260200160002060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16867f5e919055312829285818d366d1cfe50a1ba27ce2c752b655cb2faa0179e1422733600160008b815260200190815260200160002060040154600160008c815260200190815260200160002060050154600160008d815260200190815260200160002060060154604051808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200184815260200183815260200182815260200194505050505060405180910390a3505050505050565b806001600381111561092f57fe5b6001600083815260200190815260200160002060070160009054906101000a900460ff16600381111561095e57fe5b146109b4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602c81526020018061131f602c913960400191505060405180910390fd5b816001600082815260200190815260200160002060040154421015610a41576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f726566756e6454696d657374616d7020686173206e6f7420706173736564000081525060200191505060405180910390fd5b60036001600085815260200190815260200160002060070160006101000a81548160ff02191690836003811115610a7457fe5b0217905550827ffe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf060405160405180910390a26001600084815260200190815260200160002060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc610b35600160008781526020019081526020016000206006015460016000888152602001908152602001600020600501546111cc90919063ffffffff16565b9081150290604051600060405180830381858888f19350505050158015610b60573d6000803e3d6000fd5b506001600084815260200190815260200160002060008082016000905560018201600090556002820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690556003820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690556004820160009055600582016000905560068201600090556007820160006101000a81549060ff02191690555050505050565b8160016003811115610c1457fe5b6001600083815260200190815260200160002060070160009054906101000a900460ff166003811115610c4357fe5b14610c99576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602c81526020018061131f602c913960400191505060405180910390fd5b828260016000838152602001908152602001600020600401544210610d09576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001806112d66022913960400191505060405180910390fd5b8160028083604051602001808281526020019150506040516020818303038152906040526040518082805190602001908083835b60208310610d605780518252602082019150602081019050602083039250610d3d565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa158015610da2573d6000803e3d6000fd5b5050506040513d6020811015610db757600080fd5b8101908080519060200190929190505050604051602001808281526020019150506040516020818303038152906040526040518082805190602001908083835b60208310610e1a5780518252602082019150602081019050602083039250610df7565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa158015610e5c573d6000803e3d6000fd5b5050506040513d6020811015610e7157600080fd5b810190808051906020019092919050505014610ef5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f736563726574206973206e6f7420636f7272656374000000000000000000000081525060200191505060405180910390fd5b83600160008781526020019081526020016000206001018190555060026001600087815260200190815260200160002060070160006101000a81548160ff02191690836003811115610f4357fe5b0217905550847f489e9ee921192823d1aa1ef800c9ffc642993538b1e7e43a4d46a91965e894ab856040518082815260200191505060405180910390a26001600086815260200190815260200160002060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc60016000888152602001908152602001600020600501549081150290604051600060405180830381858888f19350505050158015611012573d6000803e3d6000fd5b50600060016000878152602001908152602001600020600601541115611090573373ffffffffffffffffffffffffffffffffffffffff166108fc60016000888152602001908152602001600020600601549081150290604051600060405180830381858888f1935050505015801561108e573d6000803e3d6000fd5b505b6001600086815260200190815260200160002060008082016000905560018201600090556002820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690556003820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690556004820160009055600582016000905560068201600090556007820160006101000a81549060ff021916905550505050505050565b60016020528060005260406000206000915090508060000154908060010154908060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060040154908060050154908060060154908060070160009054906101000a900460ff16905088565b6000818301905082811015611249576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f536166654d617468206164642077726f6e672076616c7565000000000000000081525060200191505060405180910390fd5b80905092915050565b6000828211156112ca576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f536166654d617468207375622077726f6e672076616c7565000000000000000081525060200191505060405180910390fd5b81830390509291505056fe726566756e6454696d657374616d702068617320616c7265616479207061737365647377617020666f722074686973206861736820697320616c726561647920696e697469617465647377617020666f722074686973206861736820697320656d707479206f7220616c7265616479207370656e74a265627a7a72315820ae65fa24fe5d20dd26bc780ba2f25425e73859dacbc0fc6c4cc6adbe5f49c81f64736f6c63430005110032";

    public static final String FUNC_ADD = "add";

    public static final String FUNC_DESTRUCT = "destruct";

    public static final String FUNC_INITIATE = "initiate";

    public static final String FUNC_REDEEM = "redeem";

    public static final String FUNC_REFUND = "refund";

    public static final String FUNC_SWAPS = "swaps";

    public static final Event ADDED_EVENT = new Event("Added", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>(true) {}, new TypeReference<Address>() {}, new TypeReference<Uint256>() {}));
    ;

    public static final Event INITIATED_EVENT = new Event("Initiated", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>(true) {}, new TypeReference<Address>(true) {}, new TypeReference<Address>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}));
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

    public RemoteCall<TransactionReceipt> add(byte[] _hashedSecret, BigInteger vonValue) {
        final Function function = new Function(
                FUNC_ADD, 
                Arrays.<Type>asList(new Bytes32(_hashedSecret)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public RemoteCall<TransactionReceipt> destruct() {
        final Function function = new Function(
                FUNC_DESTRUCT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
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

    public RemoteCall<TransactionReceipt> redeem(byte[] _hashedSecret, byte[] _secret) {
        final Function function = new Function(
                FUNC_REDEEM, 
                Arrays.<Type>asList(new Bytes32(_hashedSecret),
                new Bytes32(_secret)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> refund(byte[] _hashedSecret) {
        final Function function = new Function(
                FUNC_REFUND, 
                Arrays.<Type>asList(new Bytes32(_hashedSecret)),
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

    public static class AddedEventResponse {
        public Log log;

        public byte[] _hashedSecret;

        public String _sender;

        public BigInteger _value;
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
