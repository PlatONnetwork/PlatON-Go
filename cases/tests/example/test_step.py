import pytest


@pytest.fixture(scope="class")
def staking_a(global_test_env):
    global_test_env.deploy_all()
    print("I'm here a")
    yield global_test_env
    global_test_env.shutdown()


@pytest.fixture(scope="class")
def staking_b(global_test_env):
    print("I'm here b")
    return global_test_env


class TestStep:
    def test_staking_a(self, staking_a):
        staking_a.account.generate_account(staking_a.get_rand_node().web3)
        print("failed")
        assert False

    def test_staking_b(self, staking_b):
        assert True

    def test_staking_a_b(self, staking_a, staking_b):
        assert True
