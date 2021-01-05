Explain the project:
  * Typo-tolerant authentication system. Rather than denying a login attempt with an incorrect password, the system would try to correct for common typographical mistakes
    * only works is authentiction scheme has rate-limiting
    * attacker should not know the distribution of password in the DB e.g. password, password1, password! are all very commonly used passwords
  * Example a couple of years ago, mobile keyboards would always capitalise the first letter => many typos can be fixed by simply changing the capitalisation of the first letter of the password input by the user.
  * During authentication, take the string provided by the user, add the class of corrections, hash both strings and compare against the one in the DB
  * we don't calculate levenhstein distance between two strings because it's computationally expensive and makes the system less secure!

Analyse the most common classes of (easily) correctable typos (caps lock, extra char etc)
PoC of authentication system that corrects typoes

I had also contacted Romke because he of his talk at the NFI and his experience with password cracking & hashcat!

Add hashcat

Have you supervised a project at SNE before? I saw your name on the list of projects this year

requirements of supervisor:
* 30+ min meeting a week to discuss progresss
* note that I am working on this project alone (usually RP1 is in pairs) so I might ask more questions that usual

Schedule like in January?

Monday 2pm @ reccurent

