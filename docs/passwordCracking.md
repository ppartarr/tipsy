# Password cracking

1. use pcfg cracker https://github.com/lakiw/pcfg_cracker or faster version https://github.com/lakiw/compiled-pcfg!
```
python3 pcfg_guesser -r $ruleset -s $session-name | hashcat
```
2. send a `kill -SIGUSR1 $hashcat_pid` if hashcat blocs
Common rules to target:
    l33t sp34k
    case mangling with hashcat pipe for capitalisation

list of cracked passwords doesn't tell us much about the list of non-cracked password. Can we infer something about the distribution?

# Links
Defcon talk: https://www.youtube.com/watch?v=rri9vBYR60M&feature=emb_title