# gnc-api-d

A read-only REST server for GnuCash file

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/vinymeuh/gnc-api-d)](https://goreportcard.com/report/github.com/vinymeuh/gnc-api-d)
[![Build Status](https://travis-ci.org/vinymeuh/gnc-api-d.svg?branch=master)](https://travis-ci.org/vinymeuh/gnc-api-d)
[![codecov](https://codecov.io/gh/vinymeuh/gnc-api-d/branch/master/graph/badge.svg)](https://codecov.io/gh/vinymeuh/gnc-api-d)

## How to use

```
~> export GNUCASH_FILE_PATH=models/testdata/empty.gnucash
~> ./gnc-api-d
```

The root URL list all available commands.

```
~> curl localhost:8000/
/accounts
/accounts/{id}
/accountypes
/balance
```

### Retrieve accounts

An account is uniquely identified by its ID.

```
~> curl -v localhost:8000/accounts/4c7a43144b99496ea74b135d65da4f10
{"id":"4c7a43144b99496ea74b135d65da4f10","name":"Education","type":"EXPENSE"}
```

But accounts can also be search by **name** or **type**:

```
~> curl -v localhost:8000/accounts?type=ROOT
[{"id":"121045e62ce042faa249f1f997afd5a0","name":"Root Account","type":"ROOT"}]
```

```
~> curl -v localhost:8000/accounts?name=Education
[{"id":"4c7a43144b99496ea74b135d65da4f10","name":"Education","type":"EXPENSE"}]
```

Finally it is possible to retrieve the breakdown of accounts by type.

```
~> curl -v localhost:8000/accounttypes
{"ASSET":2,"BANK":2,"CASH":1,"CREDIT":1,"EQUITY":2,"EXPENSE":45,"INCOME":9,"LIABILITY":1,"ROOT":1}
```

### Accounts balance
