# PlatON-Tests
This is some automated test cases

### Run test:
windows Environmental variables
add both './utils/ethkey' and './utils/pubkey' to your 'PATH" variable.

Run the following command in your project directory
```js
python run.py --node='./deploy/node/4_node.yml' --case='all'

```
help:
```js
python run.py -h
```


### Dir introduce:
- [case](docs/case_example.md)
- common：common utils
- conf：Global variables
- data：Some necessary test dependencies, or data to drive use cases
- [deploy](docs/deploy.md)
- docs：Some instructions
- utils：Basic library