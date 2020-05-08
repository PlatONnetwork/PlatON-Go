
.PHONY: clean-pyc clean-tmp

help:
	@echo "clean-pyc - remove Python file artifacts"
	@echo "lint - check style with autopep8"
	@echo "test - run tests quickly with the default Python"
	@echo "install - install requirements.txt python lib"
	@echo "clean-tmp - remove run tmp"


clean: clean-tmp clean-pyc

clean-pyc:
	find . -name '*.pyc' -exec rm -f {} +
	find . -name '*.pyo' -exec rm -f {} +
	find . -name '*~' -exec rm -f {} +
	find . -name '__pycache__' -exec rm -rf {} +

clean-tmp:
	rm -rf ./report +
	rm -rf ./log +
	rm -rf ./bug_log +
	rm -rf ./allure-report +
	rm -rf ./.pytest_cache +
	rm -rf ./deploy/tmp +

lint:
	find . -name '*.py' -exec autopep8  --max-line-length=120 --in-place --aggressive --ignore=E123,E133,E50 {} +

test:
	py.test --tb native tests

install:
	pip3 install -r requirements.txt

