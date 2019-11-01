from environment.env import TestEnvironment
from environment.node import Node
from .economic import Economic


class Restricting:
    """
    Used to initiate a Restricting transaction,
    if you need to use the call method, please call ppos
    example:
    >>>res=Restricting(env, node)
    >>>res.ppos.getRestrictingInfo(...)
    """

    def __init__(self, env: TestEnvironment, node: Node):
        self.node = node
        self.economic = Economic(env)

    def createRestrictingPlan(self, account, plan, from_address, transaction_cfg=None):
        """
        Create a lockout plan
        :param account: Locked account release account
        :param plan:
        An is a list of RestrictingPlan types (array), and RestrictingPlan is defined as follows:
        type RestrictingPlan struct {
            Epoch uint64
            Amount *big.Int
            }
         where Epoch: represents a multiple of the billing period.
         The product of the number of blocks per billing cycle indicates that the locked fund
         s are released at the target block height. Epoch * The number of blocks per cycle is
         at least greater than the maximum irreversible block height.
         Amount: indicates the amount to be released on the target block.
        :param from_address: address for transaction
        :param transaction_cfg: Transaction basic configuration
              type: dict
              example:cfg = {
                  "gas":100000000,
                  "gasPrice":2000000000000,
                  "nonce":1,
              }
        :return: if is need analyze return transaction result dict
                if is not need analyze return transaction hash
        """
        pri_key = self.economic.account.find_pri_key(from_address)
        return self.ppos.createRestrictingPlan(account, plan, pri_key, transaction_cfg)

    @property
    def ppos(self):
        return self.node.ppos
