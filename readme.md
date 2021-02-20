<img src="static/images/gopher.png" alt="tipsy gopher" width="300" height="300"/>

# tipsy 🍻
![go test](https://github.com/ppartarr/tipsy/actions/workflows/go.yml/badge.svg)

Tipsy is a Go library that provides the building blocks for typo-tolerant authentication systems

## Future work

- [ ] JS to collect data about users with password managers & how many typos could be corrected
- [ ] check if user already has requested token before creating new one in PasswordReset before storeToken
- [x] use per-user salt in typtop
- [ ] make PBE & PKE configurable
- [ ] block attempts by IP
- [ ] add support for whitespace passwords
- [ ] make correctors work with different keyboard layouts

## Running the experiments
```bash
# all tests
go test -v -timeout=0

# tests for a single checker
go test -v -timeout=0 -run TestSecLossAlways
go test -v -timeout=0 -run TestSecLossBlacklist
go test -v -timeout=0 -run TestSecLossOptimal

# read results for q = 10, 100, 10000
go test -v -timeout=0 -run TestSecLoss
```

### Generating plots
```bash
# comparing checkers using a given dataset
go test -v -run TestPlotDataset

# comparing dataset using a checker
go test -v -run TestPlotChecker

# generate a plot for a single checker & dataset
go test -v -run TestPlot
```


## Links
* password lists
    * https://github.com/danielmiessler/SecLists/tree/master/Passwords/Leaked-Databases
    * https://github.com/berzerk0/Probable-Wordlists/tree/master/Real-Passwords
* OPAQUE implementations
    * https://github.com/schollz/pake
    * https://github.com/cretz/gopaque
    * https://github.com/frekui/opaque
* Research
    * [pASSWORD tYPOS and How to Correct Them Securely, Chatterjee et al., 2016](https://ieeexplore.ieee.org/document/7546536)
    * [The TypTop System: Personalized Typo-Tolerant Password Checking, Chatterjee et al., 2017](https://eprint.iacr.org/2017/810.pdf)
    * [Tipsy: How to Correct Password Typos Safely, Partarrieu, 2020](https://delaat.net/rp/2020-2021/p67/report.pdf)