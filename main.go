package main

import (
	"log"
	"flag"
)

func main() {
	// default secret = 1
	g := flag.Int("g", 11, "g")
	p := flag.Int("p", 34, "p")
	A := flag.Int("A", 33, "A")
	B := flag.Int("B", 19, "B")
	flag.Parse()
	prms := primitiveRootsModulo(*p)
	log.Printf("g: %d, p: %d, A: %d, B: %d, prms: %v", *g, *p, *A, *B, prms)

	if ms, ok := prms[*g]; ok {
		period := len(ms)

		log.Printf("Modulos: %v(%d)", ms, period);
		for k, v := range ms {
			if v == *B {
				c := k + 1
				s := 1
				for i := 1; i <= c; i++ {
					s = (s * *A) % *p
				}

				log.Printf("c: %d, secret: %d", c, s);
			}
		}
	}
}

func primitiveRootsModulo(n int) map[int][]int {
	prms := make(map[int][]int)
	cprs := coprimes(n)

	for _, cpr := range cprs {
		if 1 == cpr {
			if 1 == len(cprs) {
				prms = map[int][]int{cpr: []int{cpr}}
				break
			} else {
				continue
			}
		}

		var modulo []int
		f := 1
		for i := 1; i < n; i++ {
			m := (f * cpr) % n
			if len(modulo) > 0 && modulo[0] == m {
				break
			} else {
				modulo = append(modulo, m)
			}
			f = m
		}

		if len(modulo) == len(cprs) {
			prms[cpr] = modulo
		}
	}

	return prms
}

func coprimes(n int) []int {
	var cprs []int
	prs := primes(n)

	for i := 1; i < n; i++ {
		suc := true
		for _, p := range prs {
			if 0 == i % p {
				suc = false
				break
			}
		}

		if suc {
			cprs = append(cprs, i)
		}
	}

	return cprs
}

func primes(n int) []int {
	var prs []int

	d := 0
	for {
		p := nextPrime(d)
		suc := false
		for {
			if n != 1 && 0 == n % p {
				suc = true
				n = n / p
			} else {
				break
			}
		}

		if suc {
			prs = append(prs, p)
		}

		if 1 == n {
			break
		}
		d = p
	}

	return prs
}

func nextPrime(n int) int {
	for {
		n++
		if isPrime(n) {
			return n
		}
	}
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}

	for i := 2; i < n; i++ {
		if 0 == n % i {
			return false
		}
	}

	return true
}