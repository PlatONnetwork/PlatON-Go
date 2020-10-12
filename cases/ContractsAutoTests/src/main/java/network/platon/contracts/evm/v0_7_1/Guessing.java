package network.platon.contracts.evm.v0_7_1;

import com.alaya.abi.solidity.EventEncoder;
import com.alaya.abi.solidity.FunctionEncoder;
import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Address;
import com.alaya.abi.solidity.datatypes.Bool;
import com.alaya.abi.solidity.datatypes.DynamicArray;
import com.alaya.abi.solidity.datatypes.Event;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Bytes32;
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
 * <p>Generated with web3j version 0.13.2.1.
 */
public class Guessing extends Contract {
    private static final String BINARY = "60806040526000600260006101000a81548160ff021916908315150217905550674563918244f400006003556001600a556000600d5534801561004157600080fd5b5060405161122b38038061122b8339818101604052602081101561006457600080fd5b810190808051906020019092919050505033600c60006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550806000819055505061115e806100cd6000396000f3fe6080604052600436106101445760003560e01c806382d333a0116100b6578063b69ef8a81161006f578063b69ef8a814610602578063c2930f911461062d578063c29bde8c14610658578063c7f40b6414610683578063ecbde5e6146106c4578063ff37be16146106ef57610155565b806382d333a01461041d5780638941f2f01461044857806394696a9214610497578063a46c363714610503578063ae4479411461052e578063b03e00771461059357610155565b806313eaca431161010857806313eaca43146102af5780631ef9c56f146102dc578063220bc55e1461034857806327ebd9ab14610352578063629374ab1461037d5780636e5ab671146103b857610155565b8063045f9c9714610164578063062d6a98146101c9578063083c6323146101f4578063094cc1ab1461021f5780630b8b85021461024a57610155565b3661015557610153343361071a565b005b34801561016157600080fd5b50005b34801561017057600080fd5b5061019d6004803603602081101561018757600080fd5b81019080803590602001909291905050506108c2565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156101d557600080fd5b506101de6108fe565b6040518082815260200191505060405180910390f35b34801561020057600080fd5b50610209610904565b6040518082815260200191505060405180910390f35b34801561022b57600080fd5b5061023461090a565b6040518082815260200191505060405180910390f35b34801561025657600080fd5b506102836004803603602081101561026d57600080fd5b8101908080359060200190929190505050610910565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156102bb57600080fd5b506102c4610943565b60405180821515815260200191505060405180910390f35b3480156102e857600080fd5b506102f1610956565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b83811015610334578082015181840152602081019050610319565b505050509050019250505060405180910390f35b6103506109e4565b005b34801561035e57600080fd5b506103676109f0565b6040518082815260200191505060405180910390f35b34801561038957600080fd5b506103b6600480360360208110156103a057600080fd5b8101908080359060200190929190505050610a1a565b005b3480156103c457600080fd5b50610407600480360360208110156103db57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610ddc565b6040518082815260200191505060405180910390f35b34801561042957600080fd5b50610432610df4565b6040518082815260200191505060405180910390f35b34801561045457600080fd5b506104816004803603602081101561046b57600080fd5b8101908080359060200190929190505050610dfa565b6040518082815260200191505060405180910390f35b3480156104a357600080fd5b506104ac610e23565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b838110156104ef5780820151818401526020810190506104d4565b505050509050019250505060405180910390f35b34801561050f57600080fd5b50610518610eb8565b6040518082815260200191505060405180910390f35b34801561053a57600080fd5b5061057d6004803603602081101561055157600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610ec5565b6040518082815260200191505060405180910390f35b34801561059f57600080fd5b506105ec600480360360408110156105b657600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610edd565b6040518082815260200191505060405180910390f35b34801561060e57600080fd5b50610617610f0b565b6040518082815260200191505060405180910390f35b34801561063957600080fd5b50610642610f11565b6040518082815260200191505060405180910390f35b34801561066457600080fd5b5061066d610f17565b6040518082815260200191505060405180910390f35b34801561068f57600080fd5b50610698610f1d565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156106d057600080fd5b506106d9610f43565b6040518082815260200191505060405180910390f35b3480156106fb57600080fd5b50610704610f4b565b6040518082815260200191505060405180910390f35b43600054106108be5760035482101561073257600080fd5b600860008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600a5490806001815401808255809150506001900390600052602060002001600090919091909150558060076000600a54815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600a6000815480929190600101919050555081600660008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282540192505081905550816004600082825401925050819055507fe842aea7a5f1b01049d752008c53c52890b1a6daf660cf39e8eec506112bbdf681836001604051808473ffffffffffffffffffffffffffffffffffffffff1681526020018381526020018215158152602001935050505060405180910390a15b5050565b600b81815481106108cf57fe5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600a5481565b60005481565b600d5481565b60076020528060005260406000206000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600260009054906101000a900460ff1681565b6060600b8054806020026020016040519081016040528092919081815260200182805480156109da57602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311610990575b5050505050905090565b6109ee343361071a565b565b600060011515600260009054906101000a900460ff16151514610a1257600080fd5b600d54905090565b600054431115610dd957600260009054906101000a900460ff16158015610a8e57503373ffffffffffffffffffffffffffffffffffffffff16600c60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16145b8015610a9c57506001600a54115b15610dd857600081604051602001808281526020019150506040516020818303038152906040528051906020012060001c90506000600a548281610adc57fe5b0690506064600a541015610afa57610af581600a610f55565b610b24565b612710600a541015610b1657610b11816064610f55565b610b23565b610b22816103e8610f55565b5b5b600b8054905060045481610b3457fe5b04600581905550600080600090505b600b80549050811015610db157600160096000600b8481548110610b6357fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020540160096000600b8481548110610bdb57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550600073ffffffffffffffffffffffffffffffffffffffff16600b8281548110610c6957fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16141580610d1b5750600073ffffffffffffffffffffffffffffffffffffffff16600b8281548110610cd757fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614155b15610da457600b8181548110610d2d57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1691508173ffffffffffffffffffffffffffffffffffffffff166108fc6005549081150290604051600060405180830381858888f19350505050158015610da2573d6000803e3d6000fd5b505b8080600101915050610b43565b506001600260006101000a81548160ff021916908315150217905550836001819055505050505b5b50565b60096020528060005260406000206000915090505481565b60055481565b6000610101430382118015610e1157506001430382105b610e1a57600080fd5b81409050919050565b6060600860003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020805480602002602001604051908101604052809291908181526020018280548015610eae57602002820191906000526020600020905b815481526020019060010190808311610e9a575b5050505050905090565b6000600b80549050905090565b60066020528060005260406000206000915090505481565b60086020528160005260406000208181548110610ef657fe5b90600052602060002001600091509150505481565b60045481565b60035481565b60015481565b600c60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600047905090565b6000600154905090565b808281610f5e57fe5b06600d819055506000600d541415611041576000600190505b600a5481101561103b57600082600d54830381610f9057fe5b06141561102e57600b6007600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505b8080600101915050610f77565b50611124565b6000600190505b600a5481101561112257600082828161105d57fe5b06141580156110795750600082600d5483038161107657fe5b06145b1561111557600b6007600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505b8080600101915050611048565b505b505056fea264697066735822122019788075e7b7a29c1813cbaa81c1e403b514ffac5068d38ac73d45e70e55e8fe64736f6c63430007010033";

    public static final String FUNC_AVERAGEAMOUNT = "averageAmount";

    public static final String FUNC_BALANCE = "balance";

    public static final String FUNC_BASEUNIT = "baseUnit";

    public static final String FUNC_BLOCK_HASH = "block_hash";

    public static final String FUNC_CREATEADDRESS = "createAddress";

    public static final String FUNC_DRAW = "draw";

    public static final String FUNC_ENDBLOCK = "endBlock";

    public static final String FUNC_GENERATEBLOCKHASH = "generateBlockHash";

    public static final String FUNC_GETBALANCEOF = "getBalanceOf";

    public static final String FUNC_GETENDBLOCKHASH = "getEndBlockHash";

    public static final String FUNC_GETMYGUESSCODES = "getMyGuessCodes";

    public static final String FUNC_GETPOSTFIX = "getPostfix";

    public static final String FUNC_GETWINNERADDRESSES = "getWinnerAddresses";

    public static final String FUNC_GETWINNERCOUNT = "getWinnerCount";

    public static final String FUNC_GUESSINGCLOSED = "guessingClosed";

    public static final String FUNC_GUESSINGWITHLAT = "guessingWithLat";

    public static final String FUNC_GUSSINGERCODES = "gussingerCodes";

    public static final String FUNC_GUSSINGERLAT = "gussingerLat";

    public static final String FUNC_INDEXKEY = "indexKey";

    public static final String FUNC_INDEXOFGUSSINGER = "indexOfgussinger";

    public static final String FUNC_POSTFIX = "postfix";

    public static final String FUNC_WINNERADDRESSES = "winnerAddresses";

    public static final String FUNC_WINNERMAP = "winnerMap";

    public static final Event CURRENTBALANCE_EVENT = new Event("CurrentBalance", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}, new TypeReference<Uint256>() {}));
    ;

    public static final Event FUNDTRANSFER_EVENT = new Event("FundTransfer", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}, new TypeReference<Uint256>() {}, new TypeReference<Bool>() {}));
    ;

    protected Guessing(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected Guessing(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<Guessing> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId, BigInteger _endBlock) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new Uint256(_endBlock)));
        return deployRemoteCall(Guessing.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static RemoteCall<Guessing> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId, BigInteger _endBlock) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new Uint256(_endBlock)));
        return deployRemoteCall(Guessing.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public List<CurrentBalanceEventResponse> getCurrentBalanceEvents(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(CURRENTBALANCE_EVENT, transactionReceipt);
        ArrayList<CurrentBalanceEventResponse> responses = new ArrayList<CurrentBalanceEventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
            CurrentBalanceEventResponse typedResponse = new CurrentBalanceEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse._msgSenderAddress = (String) eventValues.getNonIndexedValues().get(0).getValue();
            typedResponse._balance = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<CurrentBalanceEventResponse> currentBalanceEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, CurrentBalanceEventResponse>() {
            @Override
            public CurrentBalanceEventResponse call(Log log) {
                EventValuesWithLog eventValues = extractEventParametersWithLog(CURRENTBALANCE_EVENT, log);
                CurrentBalanceEventResponse typedResponse = new CurrentBalanceEventResponse();
                typedResponse.log = log;
                typedResponse._msgSenderAddress = (String) eventValues.getNonIndexedValues().get(0).getValue();
                typedResponse._balance = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<CurrentBalanceEventResponse> currentBalanceEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(CURRENTBALANCE_EVENT));
        return currentBalanceEventObservable(filter);
    }

    public List<FundTransferEventResponse> getFundTransferEvents(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(FUNDTRANSFER_EVENT, transactionReceipt);
        ArrayList<FundTransferEventResponse> responses = new ArrayList<FundTransferEventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
            FundTransferEventResponse typedResponse = new FundTransferEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse._backer = (String) eventValues.getNonIndexedValues().get(0).getValue();
            typedResponse._amount = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
            typedResponse._isSuccess = (Boolean) eventValues.getNonIndexedValues().get(2).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<FundTransferEventResponse> fundTransferEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, FundTransferEventResponse>() {
            @Override
            public FundTransferEventResponse call(Log log) {
                EventValuesWithLog eventValues = extractEventParametersWithLog(FUNDTRANSFER_EVENT, log);
                FundTransferEventResponse typedResponse = new FundTransferEventResponse();
                typedResponse.log = log;
                typedResponse._backer = (String) eventValues.getNonIndexedValues().get(0).getValue();
                typedResponse._amount = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
                typedResponse._isSuccess = (Boolean) eventValues.getNonIndexedValues().get(2).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<FundTransferEventResponse> fundTransferEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(FUNDTRANSFER_EVENT));
        return fundTransferEventObservable(filter);
    }

    public RemoteCall<BigInteger> averageAmount() {
        final Function function = new Function(FUNC_AVERAGEAMOUNT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> balance() {
        final Function function = new Function(FUNC_BALANCE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> baseUnit() {
        final Function function = new Function(FUNC_BASEUNIT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<byte[]> block_hash() {
        final Function function = new Function(FUNC_BLOCK_HASH, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<String> createAddress() {
        final Function function = new Function(FUNC_CREATEADDRESS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> draw(byte[] _block_hash) {
        final Function function = new Function(
                FUNC_DRAW, 
                Arrays.<Type>asList(new Bytes32(_block_hash)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> endBlock() {
        final Function function = new Function(FUNC_ENDBLOCK, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<byte[]> generateBlockHash(BigInteger _blocknumber) {
        final Function function = new Function(FUNC_GENERATEBLOCKHASH, 
                Arrays.<Type>asList(new Uint256(_blocknumber)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<BigInteger> getBalanceOf() {
        final Function function = new Function(FUNC_GETBALANCEOF, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<byte[]> getEndBlockHash() {
        final Function function = new Function(FUNC_GETENDBLOCKHASH, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<List> getMyGuessCodes() {
        final Function function = new Function(FUNC_GETMYGUESSCODES, 
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

    public RemoteCall<BigInteger> getPostfix() {
        final Function function = new Function(FUNC_GETPOSTFIX, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<List> getWinnerAddresses() {
        final Function function = new Function(FUNC_GETWINNERADDRESSES, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicArray<Address>>() {}));
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

    public RemoteCall<BigInteger> getWinnerCount() {
        final Function function = new Function(FUNC_GETWINNERCOUNT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<Boolean> guessingClosed() {
        final Function function = new Function(FUNC_GUESSINGCLOSED, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<TransactionReceipt> guessingWithLat(BigInteger vonValue) {
        final Function function = new Function(
                FUNC_GUESSINGWITHLAT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public RemoteCall<BigInteger> gussingerCodes(String param0, BigInteger param1) {
        final Function function = new Function(FUNC_GUSSINGERCODES, 
                Arrays.<Type>asList(new Address(param0),
                new Uint256(param1)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> gussingerLat(String param0) {
        final Function function = new Function(FUNC_GUSSINGERLAT, 
                Arrays.<Type>asList(new Address(param0)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> indexKey() {
        final Function function = new Function(FUNC_INDEXKEY, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> indexOfgussinger(BigInteger param0) {
        final Function function = new Function(FUNC_INDEXOFGUSSINGER, 
                Arrays.<Type>asList(new Uint256(param0)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> postfix() {
        final Function function = new Function(FUNC_POSTFIX, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> winnerAddresses(BigInteger param0) {
        final Function function = new Function(FUNC_WINNERADDRESSES, 
                Arrays.<Type>asList(new Uint256(param0)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> winnerMap(String param0) {
        final Function function = new Function(FUNC_WINNERMAP, 
                Arrays.<Type>asList(new Address(param0)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static Guessing load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new Guessing(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static Guessing load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new Guessing(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static class CurrentBalanceEventResponse {
        public Log log;

        public String _msgSenderAddress;

        public BigInteger _balance;
    }

    public static class FundTransferEventResponse {
        public Log log;

        public String _backer;

        public BigInteger _amount;

        public Boolean _isSuccess;
    }
}
