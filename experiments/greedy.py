import string
import numpy as np

allowed_chars = string.printable[:-5]

def apply_edits(w):
    yield w.capitalize()
    yield w[0].upper() + w[1:]
    yield w[0].lower() + w[1:]
    for c in allowed_chars:
        for i in range(len(w)):
            yield w[:i] + c + w[i:]
            yield w[:i] + c + w[i+1:]
        yield w + c
    for i in range(len(w)):
        yield w[:i] + w[i+1:]


def getball(w):
    return np.array(filter(
        lambda x: x>=0,
        [ x for x in apply_edits(w) if len(x)>=6 ]
    ))

def greedy_maxcoverage_heap(filename, rateLimit, **kwargs):
    global passwordModel
    passwordModel = Passwords(filename)
    subsetHeap = priority_dict()
    guessList = []
    ballsize = 2000 # I don't care any bigger ball
    done = set()
    passwordFrequencies = np.copy(passwordModel.values()) # deep copy of the frequencies
    length = 1
    startTime = time.time()
    # pool of workers for multiprocessing
    pool = multiprocessing.Pool(5)
    for index, (passwordID, frequency) in enumerate(passwordModel):
        # set the registered password
        registeredPassword = passwordModel.id2pw(passwordID)
        if len(registeredPassword) < 6:
            continue
        password = passwordModel.id2pw(passwordID)
        probability = passwordModel.prob(password)
        # apply edits to password
        neighbors = set(apply_edits(password.encode('ascii', errors='ignore'))) - done

        # pops password with smallest frequency
        for submittedPassword, submittedPasswordFrequency in subsetHeap.sorted_iter():
            # wtf ?
            submittedPasswordFrequency = -submittedPasswordFrequency
            ball = getball(submittedPassword)
            # ballFrequencySum is the sum of frequencies in the ball of of submittedPasswordFrequency
            ballFrequencySum = passwordFrequencies[ball].sum()
            # why is this here?
            if submittedPasswordFrequency == ballFrequencySum:
                if submittedPasswordFrequency >= frequency * ballsize:
                    print("Guess({}/{}): {} weight: {}".format(
                        len(guessList),
                        rateLimit,
                        submittedPassword,
                        submittedPasswordFrequency/passwordModel.totalf()))
                    done.add(submittedPassword)
                    guessList.append(submittedPassword)
                    passwordFrequencies[ball] = 0
                    if len(guessList) >= rateLimit:
                        break
                else:  # The ball weight is still small
                    subsetHeap[submittedPassword] = -ballFrequencySum
                    break
            else:
                subsetHeap[submittedPassword] = -ballFrequencySum
        ballMax = 0

        # HERE
        for submittedPassword, ball in zip(neighbors, pool.map(getball, iter(neighbors))):
            subsetHeap[submittedPassword] = -passwordFrequencies[ball].sum()
            ballMax = max(ballMax, ball.shape[0])
        ballsize = ballsize*0.9 + ballMax*0.1

        if len(subsetHeap) > length:
            print(">< ({}) : Heap size: {} ballsize: {}".format(
                time.time()-startTime,
                len(subsetHeap),
                ballsize))
            length = len(subsetHeap) * 2
        if index % 10 == 0:
            print("({}) : {}: {} ({})".format(
                time.time()-startTime,
                index,
                registeredPassword,
                frequency))
        if len(guessList)>=rateLimit:
            break
    normalSuccess = passwordModel.sumvalues(rateLimit=rateLimit)/passwordModel.totalf()
    guessedPasswords = np.unique(np.concatenate(pool.map(getball, guessList)))
    fuzzyQ = passwordModel.values()[guessedPasswords].sum()/passwordModel.totalf()
    print("normal succ: {}, fuzzy succ: {}".format(normalSuccess, fuzzyQ))
    with open('guess_{}.json'.format(rateLimit), 'submittedPasswordFrequency') as frequency:
        json.dump(guessList, frequency)
    return guessList