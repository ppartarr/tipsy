# TypTop

### Registration with password w: Reg(w)
user submits password w
<!-- (𝑝𝑘, 𝑠𝑘) ←$ 𝒦 -->
generate public/private key pair for user
<!-- T[0] ←$ E𝑤 (𝑠𝑘) -->
encrypt the private key using the password and add it to the typo cache T.add(E(sk, w))
<!-- For 𝑖 = 1, . . . , 𝑡 do T[𝑖] ←$ 𝒞E -->
warm up the typo cache T.add(E(sk, w1))where w1 is a likely typo of password w
<!-- For 𝑗 = 1, . . . , 𝜔 do W[𝑗] ←$ ℰ𝑝𝑘 (𝜀) -->
warm up the wait list {epsilon}pk where epsilon is the empty string
<!-- (S0, 𝒰0) ←$ CacheInit(𝑤) -->
init cache returns an initial state for the caching scheme, and a set of typo/index pairs
    S0 initial state is simply (password, frequency)
    U0 typo/index pair is null
<!-- 𝑐 ←$ ℰ𝑝𝑘 (S0) -->
encrypt the initial state of the caching scheme: c = {state}pk
<!-- For (𝑤, 𝑖 ˜ ) ∈ 𝒰0 do -->
iterate over the (typo, index) pairs
    <!-- T[𝑖] ←$ E𝑤˜(𝑠𝑘) -->
    add E(sk, typo) to the typo cache T
<!-- 𝛾 ←$ Z𝜔 -->
gamma = pick a random number between 0 & t, Z is the set of integers [1]
<!-- 𝑠 ← (𝑝𝑘, 𝑐, T, W, 𝛾) -->
set the serverState = (pk, {initialState}pk, T, W, gamma)
<!-- Return s -->

[1] https://en.wikipedia.org/wiki/Integer

### Cache initialisation CacheInit(w)
<!-- For 𝑖 = 1, . . . , 𝑡 do F[𝑖] ← 0 -->
set the frequency of passwords in typo cache to 0
<!-- S ← (𝑤, F) -->
set the state of the cache as (w, F)
<!-- 𝒰 ← 𝜑 -->
create the set of typo / index pairs such that it is the empty set...
    U = (w0, i) where w0 is the typo & i is the index at which to store it in the initial cache
<!-- Return (S, 𝒰) -->
return the state & the set of typo / index pairs

### Check password
<!-- Parse 𝑠 as (𝑝𝑘, 𝑐, 𝛾, T, W) -->
parse the server state
<!-- 𝑏 ← false -->
set b to false where b = typo matches typo in typo cache
<!-- For 𝑖 = 0, . . . , 𝑡 do -->
iterate t time
    <!-- 𝑠𝑘 ← D𝑤(T[𝑖]) -->
    try derive the secret key by D(c, sk)
    <!-- If 𝑠𝑘 , ⊥ then -->
        wtf is this?
        <!-- 𝑏 ← true; 𝜋 ←$ Perm(𝑡); S ← 𝒟𝑠𝑘 (𝑐) -->
        b = true
        pi <- sample from permutation(t)
        state = [c]sk
        <!-- For 𝑗 = 1, . . . , 𝜔  -->
        iterate over the wait list
            <!-- do 𝑤˜ 𝑗 ← 𝒟𝑠𝑘 (W[𝑗]) -->
            D(waitList.at(j), sk)
        <!-- (S′, 𝒰) ← CacheUpdt(𝜋, S, (𝑤, 𝑖 ˜ ), 𝑤˜ 1, . . . , 𝑤˜ 𝜔) -->
        cache update
        <!-- 𝑐′ ←$ ℰ𝑝𝑘 (S′) -->
        newCipheredState = {S'}pk
        <!-- For (𝑤˜′, 𝑗) ∈ 𝒰 -->
        iterate over every (typo, index) pair
            <!-- do T[𝑗] ←$ E𝑤˜′ (𝑠𝑘) -->
            update typo cache with T[index] = E(sk, 𝑤˜′)
        <!-- For 𝑗 = 1, . . . , 𝑡 -->
        iterate over the typo cache
            <!-- do T′[𝜋[𝑗]] ← T[𝑗]𝑤˜′ -->
            randomise the order of the elements in the typo cache
        <!-- For 𝑗 = 1, . . . , 𝜔 -->
        iterate over the wait list
            <!-- do W[𝑗] ←$ ℰ𝑝𝑘 (𝜀) -->
            add {epsilon}pk to the wait list 
        <!-- 𝑠 ← (𝑝𝑘, 𝑐′, 𝛾, T′, W) -->
        update the server state
<!-- If 𝑏 = false then -->
if b == false
    <!-- W[𝛾] ←$ ℰ𝑝𝑘 (𝑤˜ ); 𝛾′ ← 𝛾 + 1 mod 𝜔 -->
    add {𝑤˜}pk to wait list using an LFU eviction policy
    <!-- 𝑠 ← (𝑝𝑘, 𝑐, T, W, 𝛾′) -->
    set the serverState = (pk, {initialState}pk, T, W, gammaPrime)
<!-- Return (𝑏, 𝑠) -->
return (b, s)

### CacheUpdate (𝜋, S, (𝑤, 𝑖 ˜ ), 𝑤˜ 1, . . . , 𝑤˜ 𝜔)
<!-- Parse S as (𝑤, F) -->
self explanatory
<!-- If 𝑖 > 0 then F[𝑖] ← F[𝑖] + 1 -->
if i > 0 then increment the frequency in F
<!-- For 𝑗 = 1, . . . , 𝜔 do -->
iterate over the wait list
    <!-- If valid(𝑤, 𝑤˜ 𝑗 ) = true then -->
    valid() checks 3 conditions (see next section)
        <!-- ℳ[𝑤˜ 𝑗 ] ← ℳ[𝑤˜ 𝑗 ] + 1 -->
        increment frequency in M
<!-- Sort ℳ in decreasing order of values -->
self explanatory
<!-- For each 𝑤˜′ such that ℳ[𝑤˜′] > 0 do -->
for every password in the wait list with a frequency > 0 
    <!-- 𝑘 ← argmin𝑗 F[𝑗] -->
    k is password with the index of the lowest frequency in the typo cache
    <!-- 𝜈 ← ℳ[𝑤˜′]/(F[𝑘] + ℳ[𝑤˜′]) -->
    nu is the (frequency of the password in wait list) / (frequency of least used password in typo cache) + (frequency of the password in wait list)
    <!-- 𝑑 ← 𝜈 {0, 1} -->
    if nu < 0.5 d = 0 else d = 1
    <!-- If 𝑑 = 1 then -->
    if d == 1
        <!-- F[𝑘] ← F[𝑘] + ℳ[𝑤˜′] -->
        set the frequency of the new typo in the typo cache to be the (frequency of least used password in typo cache) + (frequency of password in wait list)
        <!-- 𝒰 ← 𝒰 ∪ {(𝑤˜′, 𝑘)} -->
        add the pair (typo, index) to the list
<!-- For 𝑗 = 1, . . . , 𝑡 do -->
iterate over frequency list
    <!-- F′[𝜋(𝑗)] ← F[𝑗] -->
    apply random permutation to order of frequencies
<!-- S′ ← (𝑤, F′) -->
create the new server state s prime
<!-- Return (S′, 𝒰) -->
return

#### CacheUpdate valid() conditions
1. damerau-levenshtein distance is < 2
2. strength estimation of typo password >= 10   (ensures easily guessed passwords are never cached using zxcvbn)
3. strength estimation of typo password >= strength estimation of password - 3  (prevent caching of typos significantly more guessable than the real password)

## Definitionss

This motherfucker is using omega � & w...


## public key notation PKE = (K, ℰ, D)
sign message M with Alice's private key   => [M]Alice
encrypt message M with Alice's public key => {M}Alice

## symmetrict key notation
encrypt with symmetric key K   => C = E(P, K)
decrypt C with symmetric key K => P = D(C, K)

## symbols
U is a pair made up of a (typo, index) where index is the index in the cache

c is the ciphered server state S

gamma is a random number from Zt

pi is a random sample from the Perm(t)

epsilon is the empty string

Perm(t) is the set of all permutations for on Zt:
    given t = 3, we have Zt = 0, 1, 2
    Perm(t) = {
        [0, 1, 2],
        [0, 2, 1],
        [1, 0, 2],
        [1, 2, 0],
        [2, 0, 1],
        [2, 1, 0]
    }


## personalised typo correction
    * implement encrypted wait list: W
        * public-key pk encryption of recent incorrect password submissions that are not the registered password or one of the typos in the typo-cache e.g. for registered password p and incorrect submissions w0, w1, ..., wn such that w != w0 and for typo cache T, w ∉ T
            * { {w0}pk, {w1}pk, ..., {wn}pk }
                * ensure that p != w, for any w in {w0, w1, ..., wn}
                * given a typo cache T, w ∉ T, for any {w0, w1, ..., wn}
                * has a configurable length that should match the rateLimit in practice
            * uses least recently used (LRU) eviction policy
        * private-key sk is encrypted using the registered password:
            * C = E(sk, p)
            * p is the user's registered password
            * sk is the private-key (secret key)
    * implement typo cache: T
        * set of strings that the user is allowed to authenticate with (registed passwords + accepted typos of password) e.g. for password w and accepted submissions wa, wb, ..., wn
            * { E(sk, wa), E(sk, wb), ..., E(sk, wn) }
            * has configurable length
        * cache initialisation
        * cache update policy (configurable policy e.g. LFU, PLFU etc)
        * after every change in typo cache, randomly permute order of cached typos
    * write isTypo() based on the Damerau-Levenshtein distance
    * state = (publicKey, cachingSchemeState, typoCache, waitList, nextWaitListEntry)
    * admissable typos must respect 3 (configurable) conditions:
        * Damerau-Levenshtein < 1
        * strength estimation of typo password >= 10   (ensures easily guessed passwords are never cached using zxcvbn)
        * strength estimation of typo password >= strength estimation of password - 3  (prevent caching of typos significantly more guessable than the real password)