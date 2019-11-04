import pytest


@pytest.fixture(scope="class")
def staking_a(global_test_env):
    print("I'm here a")
    return global_test_env


@pytest.fixture(scope="class")
def staking_b(global_test_env):
    print("I'm here b")
    return global_test_env


class TestStep:
    def test_staking_a(self, staking_a):
        print("failed")
        assert False

    def test_staking_b(self, staking_b):
        assert True

    def test_staking_a_b(self, staking_a, staking_b):
        assert True
