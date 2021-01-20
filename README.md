# PF Test

![Go](https://github.com/frelon/pftest/workflows/Go/badge.svg)

PF Test is a tool to test [pf](https://man.openbsd.org/pf) rules.
It does this by loading a rule file and a file with tests and then run each test against the rules.

Example rules:

```/etc/pf.conf
block
pass from (self)
```

Example tests:

```/etc/pftest.conf
pass 127.0.0.1 222.111.222.111:80 em0
block 222.111.222.111 127.0.0.1 em0
```
