<img src="static/images/gopher.png" alt="tipsy gopher" width="300" height="300"/>

# tipsy 🍻

Tipsy is a Go library that provides the building blocks for typo-tolerant authentication systems

## Future work

- [ ] check if user already has requested token before creating new one in PasswordReset before storeToken
- [ ] use per-user salt
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
* seclists https://github.com/danielmiessler/SecLists/tree/master/Passwords/Leaked-Databases
* OPAQUE implementations
    * https://github.com/schollz/pake
    * https://github.com/cretz/gopaque
    * https://github.com/frekui/opaque