# CyberLab
Some cybersecurity ideas & showcases
Run tests to check the validity:
```shell
go test ./...
```

## Diffie-Hellman man-in-the-middle attack
In [Diffie-Hellman key exchange][diffie-hellman] method there are 4 public parameters.
Two of them form the initial agreement which are the modulus (p) and base (g).
And base (g) should be a [primitive root modulo][primitive-root-n] modulus (p).
Two remaining parameters are public keys computed and shared publicly by parties,
who are going to generate the symmetric secret key individually.
* modulus (p)
* base (g)
* _Alice_'s public key (A)
* _Bob_'s public key (B)

While two parties share neither their private keys (a: _Alice_'s secret, b: _Bob_'s secret) nor the shared secret key
which they compute it individually after exchanging their public keys,
the [GuessDiffieHellmanSharedSecret][guess-diffie-hellman-shared-secret-method] method illustrates
how is it possible to regenerate the shared secret key by only using the publicly exchanged parameters
without any knowledge of the private parameters and equal assertions are made in tests for different private key pairs
in [the test file][guess-diffie-hellman-shared-secret-tests].

### Vulnerability described
Suppose key exchange parties agreed on:
* modulus **p = 11**
* base: **g = 6**
* _Alice_ chooses a secret integer **a = 9** and sends _Bob_ her public key **A = 2** calculated from:
	```
	A = g^a mod p = 6^9 mod 11 = 10077696 mod 11 = 2
	```
* _Bob_ chooses a secret integer **b = 3** and sends _Alice_ his calculated public key **B = 7** calculated from:
	```
	B = g^b mod p = 6^3 mod 11 = 216 mod 11 = 7
	```

Now that public parameters have been exchanged it's time for _Alice_ and _Bob_ to calculate the shared secret key:
* _Alice_:
	```
	sa = B^a mod p = 7^9 mod 11 = 40353607 mod 11 = 8
	```
* _Bob_:
	```
	sb = A^b mod p = 2^3 mod 11 = 8 mod 11 = 8
	```

As it was expected they both calculated the same shared secret key
using counterpart's public key and their own private keys: **sa = sb = 8**. Let's check if it's possible to calculate the same shared secret key without any knowledge of private keys (a and b):
1. Find one cycle of remainders of sequential powers (k) of the base (g) modulo the modulus (p) which results in
one permutation of [reduced residue system][reduced-residue-system] modulo (p):
   * rrs = {r | r = g^k mod p and 0 < k <= φ(p) and k ∈ ℕ}
     * rrs = {6^1 mod 11, 6^2 mod 11, ..., 6^10 mod 11} --> **rrs = {6, 3, 7, 9, 10, 5, 8, 4, 2, 1}**
2. As **g** is a primitive root of **p** and **B** and **A** are remainders of a power of **g** modulo **p**
both **A** and **B** exist in the rrs. Find index of one of **A** or **B** in selected rrs:
	```
	indexB = rrs.indexOf(B=7) = 2
	indexA = rrs.indexOf(A=2) = 8
	```
3. Shared secret key **sh** can be computed using each index individually:
	```
	sh = 1
	for i := 0; i <= indexB; i++ {
	  sh = (sh * A) mod p
	}
	---
	sh = 1
	sh = (1 * 2) mod 11 = 2 // i = 0
	sh = (2 * 2) mod 11 = 4 // i = 1
	sh = (4 * 2) mod 11 = 8 // i = 2
	```
	or
	```
	sh = 1
	for i := 0; i <= indexA; i++ {
	 sh = (sh * B) mod p
	}
	---
	sh = 1
	sh = (1 * 7) mod 11 = 7 // i = 0
	sh = (7 * 7) mod 11 = 5 // i = 1
	sh = (5 * 7) mod 11 = 2 // i = 2
	sh = (2 * 7) mod 11 = 3 // i = 3
	sh = (3 * 7) mod 11 = 10 // i = 4
	sh = (10 * 7) mod 11 = 4 // i = 5
	sh = (4 * 7) mod 11 = 6 // i = 6
	sh = (6 * 7) mod 11 = 9 // i = 7
	sh = (9 * 7) mod 11 = 8 // i = 8
	```
	and surprisingly **sh = sa = sb = 8**

[diffie-hellman]: https://en.wikipedia.org/wiki/Diffie%E2%80%93Hellman_key_exchange
[primitive-root-n]: https://en.wikipedia.org/wiki/Primitive_root_modulo_n
[guess-diffie-hellman-shared-secret-method]: https://github.com/pouyanh/cyberlab/blob/master/diffiehellman.go#L53
[guess-diffie-hellman-shared-secret-tests]: https://github.com/pouyanh/cyberlab/blob/master/diffiehellman_test.go#L51
[reduced-residue-system]: https://en.wikipedia.org/wiki/Reduced_residue_system
