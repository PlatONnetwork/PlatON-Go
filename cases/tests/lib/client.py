from environment.env import TestEnvironment
from environment.node import Node
from .config import StakingConfig
from .economic import Economic
from .restricting import Restricting
from .delegate import Delegate
from .staking import Staking
from .pip import Pip


class Client:
    """
    Test client, the interface call method required for the collection test
    example:
    >>>cfg = StakingConfig("1111","test-node","http://test-node.com","I'm tester")
    >>>client = Client(env, node, cfg)
    >>>client.staking.create_staking(...)
    """
    def __init__(self, env: TestEnvironment, node: Node, cfg: StakingConfig):
        self.node = node
        self.economic = Economic(env)
        self.staking = Staking(env, node, cfg)
        self.restricting = Restricting(env, node)
        self.delegate = Delegate(env, node)
        self.pip = Pip(env, node)

    @property
    def ppos(self):
        """
        use sdk ppos object
        :return:
        """
        return self.node.ppos
