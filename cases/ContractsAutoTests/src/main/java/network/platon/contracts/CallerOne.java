package network.platon.contracts;

import org.web3j.abi.TypeReference;
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

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.0.
 */
public class CallerOne extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610625806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80630c55699c146100465780637b8ed01814610064578063f88bef8c14610082575b600080fd5b61004e61008c565b6040518082815260200191505060405180910390f35b61006c610092565b6040518082815260200191505060405180910390f35b61008a61009b565b005b60005481565b60008054905090565b60006040516100a9906101d2565b604051809103906000f0801580156100c5573d6000803e3d6000fd5b5090508073ffffffffffffffffffffffffffffffffffffffff1660405180807f696e63282900000000000000000000000000000000000000000000000000000081525060050190506040518091039020604051602001808281526020019150506040516020818303038152906040526040518082805190602001908083835b602083106101675780518252602082019150602081019050602083039250610144565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381855af49150503d80600081146101c7576040519150601f19603f3d011682016040523d82523d6000602084013e6101cc565b606091505b50505050565b610411806101e08339019056fe608060405234801561001057600080fd5b506103f1806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80630c55699c14610046578063371303c0146100645780635a3617561461006e575b600080fd5b61004e61008c565b6040518082815260200191505060405180910390f35b61006c610092565b005b610076610236565b6040518082815260200191505060405180910390f35b60005481565b60006040516100a09061023f565b604051809103906000f0801580156100bc573d6000803e3d6000fd5b5090508073ffffffffffffffffffffffffffffffffffffffff1660405180807f696e63282900000000000000000000000000000000000000000000000000000081525060050190506040518091039020604051602001808281526020019150506040516020818303038152906040526040518082805190602001908083835b6020831061015e578051825260208201915060208101905060208303925061013b565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381855af49150503d80600081146101be576040519150601f19603f3d011682016040523d82523d6000602084013e6101c3565b606091505b5050507fb0333e0e3a6b99318e4e2e0d7e5e5f93646f9cbf62da1587955a4092bf7df6e733600054604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a150565b60008054905090565b6101708061024d8339019056fe608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80630c55699c1461004657806317f936fb14610064578063371303c014610082575b600080fd5b61004e61008c565b6040518082815260200191505060405180910390f35b61006c610092565b6040518082815260200191505060405180910390f35b61008a61009b565b005b60005481565b60008054905090565b60008081548092919060010191905055507fb0333e0e3a6b99318e4e2e0d7e5e5f93646f9cbf62da1587955a4092bf7df6e733600054604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a156fea265627a7a72315820c1d1db72903ec8f14966fb38c712bf97738afdc3f4570ea953ef919f4461353364736f6c634300050d0032a265627a7a72315820dd0d97f1ef97f31be1ea3b2f9dbfa2f290af2f15617808ae5c06cb7792bc9a1e64736f6c634300050d0032a265627a7a7231582092cd8694670ff40d78c81bc0c055da365690865099d22512888daaa697e6770b64736f6c634300050d0032";

    public static final String FUNC_GETCALLERX = "getCallerX";

    public static final String FUNC_INC_DELEGATECALL = "inc_delegatecall";

    public static final String FUNC_X = "x";

    @Deprecated
    protected CallerOne(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected CallerOne(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected CallerOne(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected CallerOne(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<BigInteger> getCallerX() {
        final Function function = new Function(FUNC_GETCALLERX,
                Arrays.<Type>asList(),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> inc_delegatecall() {
        final Function function = new Function(
                FUNC_INC_DELEGATECALL, 
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

    public static RemoteCall<CallerOne> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(CallerOne.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<CallerOne> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(CallerOne.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<CallerOne> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(CallerOne.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<CallerOne> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(CallerOne.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static CallerOne load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new CallerOne(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static CallerOne load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new CallerOne(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static CallerOne load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new CallerOne(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static CallerOne load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new CallerOne(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
