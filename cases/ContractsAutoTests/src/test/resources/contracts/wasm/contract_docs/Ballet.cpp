#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * @author liweic
 * 投票合约
 * */

struct Voter{
    public:
        uint64_t weight;          // 256bit 的非负整数投票权重
        bool voted;               // 用户是否已经投票
        Address delegate;         // 被委托人账户
        uint64_t vote;            // 投票提案编号
        Voter(){}
        Voter(uint64_t &weight, bool &voted, Address &delegate, uint64_t &vote):weight(weight),voted(voted),delegate(delegate),vote(vote){}
        PLATON_SERIALIZE(Voter,(weight)(voted)(delegate)(vote))
}voter;

struct Proposal{
    public:
        uint64_t voteCount;          // 得票数
        Proposal(){}
        Proposal(uint64_t &voteCount):voteCount(voteCount){}
        PLATON_SERIALIZE(Proposal,(voteCount))
};

Proposal proposals;

CONTRACT Ballot : public platon::Contract{

    private:
        platon::StorageType<"chairperson"_n,Address> chairperson;
        platon::StorageType<"provector"_n,std::vector<std::string>> proposal_vector;
        platon::StorageType<"mapperson"_n,std::map<Address,Voter>> storage_address_voter;
        platon::StorageType<"mapperson"_n,std::map<std::vector<std::string>,Proposal>> storage_vector_proposals;

    public:
        ACTION void init(uint8_t _numProposals)
        {
            chairperson.self() = platon::platon_address();
            storage_address_voter.self()[chairperson.self()].weight = 1;
            proposals.voteCount = _numProposals;
        }

        ACTION void giveRightToVote(Address toVoter)
        {
            if(platon::platon_address() != chairperson.self() || storage_address_voter.self()[toVoter].voted != false)
                return;
            storage_address_voter.self()[toVoter].weight = 1;
        }

        ACTION void delegate(Address to)
        {
            voter = storage_address_voter.self()[platon::platon_address()];
            if (voter.voted) return;
            while (storage_address_voter.self()[to].delegate != Address(0))
            {
                to = storage_address_voter.self()[to].delegate;
                // 受委托人不能又将自己的票委托给委托人，形成循环
                assert(to != platon::platon_address());
            }
            if(platon::platon_address() == to) return;
            voter.voted = true;
            voter.delegate = to;
            Voter delegateTo = storage_address_voter.self()[to];
            proposals = storage_vector_proposals.self()[proposal_vector.self()];
            if (delegateTo.voted)
                proposals.voteCount += voter.weight;
            else
                delegateTo.weight += voter.weight;
        }

        ACTION void vote(uint8_t toProposal)
        {
            voter = storage_address_voter.self()[platon::platon_address()];
            if (voter.voted || toProposal >= proposal_vector.self().size()) return;
            voter.voted = true;
            voter.vote = toProposal;
            proposals.voteCount += voter.weight;
        }

        CONST uint8_t winningProposal(uint8_t _winningProposal)
        {
            u128 winningVoteCount = 0;
            for (uint8_t prop = 0; prop < proposal_vector.self().size(); prop++)
                if (proposals.voteCount > winningVoteCount)
                {
                    winningVoteCount = proposals.voteCount;
                    _winningProposal = prop;
                }
            return _winningProposal;
        }

};

PLATON_DISPATCH(Ballot, (init)(giveRightToVote)(delegate)(vote)(winningProposal))