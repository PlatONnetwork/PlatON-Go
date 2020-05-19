from dataclasses import dataclass
from dacite import from_dict
from typing import Optional
import json


@dataclass
class Cbft:
    initialNodes: list
    amount: int
    validatorMode: str
    period: int


@dataclass
class Config:
    chainId: int
    eip155Block: int
    cbft: Cbft
    genesisVersion: int


@dataclass
class Common:
    maxEpochMinutes: int
    maxConsensusVals: int
    additionalCycleTime: int


@dataclass
class Staking:
    stakeThreshold: int
    operatingThreshold: int
    maxValidators: int
    unStakeFreezeDuration: int
    rewardPerMaxChangeRange: int
    rewardPerChangeInterval: int


@dataclass
class Slashing:
    slashFractionDuplicateSign: int
    duplicateSignReportReward: int
    slashBlocksReward: int
    maxEvidenceAge: int
    zeroProduceCumulativeTime: int
    zeroProduceNumberThreshold: int


@dataclass
class Gov:
    versionProposalVoteDurationSeconds: int
    versionProposalSupportRate: int
    textProposalVoteDurationSeconds: int
    textProposalVoteRate: int
    textProposalSupportRate: int
    cancelProposalVoteRate: int
    cancelProposalSupportRate: int
    paramProposalVoteDurationSeconds: int
    paramProposalVoteRate: int
    paramProposalSupportRate: int


@dataclass
class Reward:
    newBlockRate: int
    platONFoundationYear: int
    increaseIssuanceRatio: int


@dataclass
class InnerAcc:
    platonFundAccount: str
    platonFundBalance: int
    cdfAccount: str
    cdfBalance: int


@dataclass
class EconomicModel:
    common: Common
    staking: Staking
    slashing: Slashing
    gov: Gov
    reward: Reward
    innerAcc: InnerAcc


@dataclass
class Genesis:
    config: Config
    economicModel: EconomicModel
    nonce: str
    timestamp: str
    extraData: str
    gasLimit: str
    alloc: dict
    number: str
    gasUsed: str
    parentHash: str

    def to_dict(self):
        from copy import copy
        data = copy(self.__dict__)
        data["config"] = copy(self.config.__dict__)
        data["config"]["cbft"] = copy(self.config.cbft.__dict__)
        data["economicModel"] = copy(self.economicModel.__dict__)
        data["economicModel"]["common"] = copy(self.economicModel.common.__dict__)
        data["economicModel"]["staking"] = copy(self.economicModel.staking.__dict__)
        data["economicModel"]["slashing"] = copy(self.economicModel.slashing.__dict__)
        data["economicModel"]["gov"] = copy(self.economicModel.gov.__dict__)
        data["economicModel"]["reward"] = copy(self.economicModel.reward.__dict__)
        data["economicModel"]["innerAcc"] = copy(self.economicModel.innerAcc.__dict__)
        return data

    def to_file(self, file):
        data = self.to_dict()
        with open(file, "w") as f:
            f.write(json.dumps(data, indent=4))


def to_genesis(genesis_conf) -> Genesis:
    return from_dict(Genesis, genesis_conf)


if __name__ == "__main__":
    # from common.load_file import LoadFile
    # from conf.settings import GENESIS_FILE
    # import json
    # genesis_data = LoadFile(GENESIS_FILE).get_data()
    # genesis = from_dict(data_class=Genesis, data=genesis_data)
    # genesis.config.chainId = 1
    # # print(genesis_to_dict(genesis))
    # print(genesis.to_dict())

    pass
