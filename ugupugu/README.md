# Ugupugu

![Ugupugu](ugupugu.jpeg "ugupugu")

A problem to solve!

## Desc

A cipher containing flag is given. but its encrypted using RSA algo with (e, N), (d, N) key pair.
Ugupugu is a service that will decrypt ciphers using same priv key (d, N). but its not designed to decrypt flags!
so fool developer (not me: its a fairy tale), will check deciphered msg and rejects them if it contains:

```
xeroctf{ugupugu?use_padd1ng-for-rsa!}
```

## Solution

Eve here should multiply the ciphered flag by arbitrary cipher that is encrypted with same pub key.

Therefore:

```
flag^e * arbitrary^e = (flag * arbitrary) ^ e
```

and decrypt that using Ugupugu. Therefore:

```
((flag * arbitrary) ^ e) ^ d = flag * arbitrary
```

Hence Eve has the flag ezly by multiplying `flag * arbitrary` by `arbitrary^-1`.

For more information see:
https://en.wikipedia.org/wiki/RSA_(cryptosystem)#Attacks_against_plain_RSA
