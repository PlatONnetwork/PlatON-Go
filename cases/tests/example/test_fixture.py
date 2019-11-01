import pytest
import time


@pytest.fixture
def step_1(global_test_env):
    return 1


@pytest.fixture
def step_2(step_1):
    return step_1 + 1


@pytest.fixture
def step(step_2):
    return step_2 + 1


def test_case(step):
    assert step == 3


def test_case_2(step):
    assert step == 3
