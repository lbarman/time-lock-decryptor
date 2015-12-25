package main

import (
	"os"
	"bufio"
	"fmt"
	"math/rand"
	"time"
	"strings"
	"math/big"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
)

func newPrime(pow int) *big.Int {
	//we pick the first prime
	base      := int64(2)
	power     := int64(pow)
	init      := ExpByPowOfTwo(big.NewInt(base), big.NewInt(power))
	precision := 50
	one       := big.NewInt(1)
	inc       := 0
	for !init.ProbablyPrime(precision) {
		init = big.NewInt(0).Add(init, one)
		inc += 1
	}
	prime := init
	fmt.Println("Parameters are :")
	fmt.Println("Prime number is", prime)
	fmt.Println("Corresponding to", base, "^", power, "+", inc)
	fmt.Println("Its hash is", sha256AndHex(prime.Bytes()))

	return prime
}

func WriteBig(p *big.Int, file string) {
	data := p.Bytes()
	ioutil.WriteFile(file, data[:], 0644)
}

func ReadBig(file string) *big.Int {
	data, err := ioutil.ReadFile(file)

	if err != nil {
		fmt.Println("Can't read n.dat, please regenerate the primes.")
		panic(err.Error())
	}

	p := big.NewInt(0).SetBytes(data)
	return p
}

func main() {
	fmt.Println("Time-Lock Puzzle")
	fmt.Println("****************")

	/*
	//we generate the parameters
	p := newPrime(1234)
	q := newPrime(1010)
	n := big.NewInt(0).Mul(p, q)

	WriteBig(n, "n.dat")

	fmt.Println("Parameters are :")
	fmt.Println("N is", n)
	fmt.Println("Its hash is", sha256AndHex(n.Bytes()))
	*/

	n := ReadBig("n.dat")

	fmt.Println("Parameters are :")
	fmt.Println("N is", n)
	fmt.Println("Its hash is", sha256AndHex(n.Bytes()))

	//generate benchmark, a 1-second exponentiation
	repeat := 10
	times := make([]time.Duration, repeat)
	expOneSecond :=ExpByPowOfTwo(big.NewInt(2), big.NewInt(5000))

	for i:=0; i<repeat; i++ {

		rndSource := rand.NewSource(time.Now().UnixNano())
		rnd       := rand.New(rndSource)
		rand      := big.NewInt(0).Rand(rnd, n)
		start     := time.Now()

		ExpByPowOfTwoModular(rand, expOneSecond, n)

		diff := time.Now().Sub(start)
		times[i] = diff
		fmt.Println("Benchmark", i, "took", diff)
	}

	fmt.Println("****************")
	fmt.Println("All clear !")

	reader := bufio.NewReader(os.Stdin)

	proceed := false
	x := big.NewInt(0)
	y := big.NewInt(0)

	for !proceed{
		 fmt.Print("Enter base : ")
	    base, _ := reader.ReadString('\n')
	    base = strings.Replace(base, "\n", "", 1)

	    fmt.Print("Enter number of iterations : ")
	    iter, _ := reader.ReadString('\n')
	    iter = strings.Replace(iter, "\n", "", 1)

	    x.SetString(base, 10)
	    y.SetString(iter, 10)

	    fmt.Println("Base is", x, "Iter is", y)
	    duration := time.Duration(int(y.Int64())) * time.Second
	    fmt.Println("Expected duration is", duration)
	    fmt.Print("Proceed [y/n]")
	    ans, _ := reader.ReadString('\n')
	    ans = strings.Replace(ans, "\n", "", 1)

	    if ans == "y" {
	    	proceed = true
	    }
	}

	curr := big.NewInt(0)
	one := big.NewInt(1)
	zero := big.NewInt(0)
	hundred := big.NewInt(100)

	for curr.Cmp(y) < 0 {
		if mod(curr, hundred).Cmp(zero) == 0 {
			diff := big.NewInt(0).Sub(y, curr)
			rem := time.Duration(int(diff.Int64())) * time.Second
			fmt.Println("Iteration", curr.String()+"/"+y.String(), "remaning time is", rem)
		}
		
		x = ExpByPowOfTwoModular(x, expOneSecond, n)

		curr.Add(curr,one)

	}

	fmt.Println("################")
	fmt.Println("Time-Lock Puzzle finished !")
	fmt.Println("")
	fmt.Println(x)
	fmt.Println("")
	fmt.Println(sha256AndHex(x.Bytes()))
	fmt.Println("")
	fmt.Println("################")
	fmt.Println("****************")



}

func sha256AndHex(x []byte) string {
	hash := sha256.Sum256(x)
	return hex.EncodeToString(hash[:])
}

func ExpByPowOfTwoModular(base, power, modulus *big.Int) *big.Int {

	//if power == 0, return 1
	if power.Cmp(big.NewInt(0)) == 0 {
		return mod(big.NewInt(1), modulus)
	}
	//if power == 1,  base
	if power.Cmp(big.NewInt(1)) == 0 {
		return mod(base, modulus)
	}

    result := big.NewInt(1)
    one := big.NewInt(1)
    for power.Cmp(one) > 0 {
        if modBy2(power).Cmp(one) == 0 {
            result = mod(multiply(result, base), modulus)
        	power = big.NewInt(0).Sub(power, one)
        }
        power = divideBy2(power)
        base = mod(multiply(base, base), modulus)
    }
    return mod(big.NewInt(0).Mul(result, base), modulus)
}


func ExpByPowOfTwo(base, power *big.Int) *big.Int {

	//if power == 0, return 1
	if power.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(1)
	}
	//if power == 1,  base
	if power.Cmp(big.NewInt(1)) == 0 {
		return base
	}

    result := big.NewInt(1)
    one := big.NewInt(1)
    for power.Cmp(one) > 0 {
        if modBy2(power).Cmp(one) == 0 {
            result = multiply(result, base)
        	power = big.NewInt(0).Sub(power, one)
        }
        power = divideBy2(power)
        base = multiply(base, base)
    }
    return big.NewInt(0).Mul(result, base)
}

func modBy2(x *big.Int) *big.Int {
    return big.NewInt(0).Mod(x, big.NewInt(2))
}

func divideBy2(x *big.Int) *big.Int {
    return big.NewInt(0).Div(x, big.NewInt(2))
}

func multiply(x, y *big.Int) *big.Int {
    return big.NewInt(0).Mul(x, y)
}

func mod(x, m *big.Int) *big.Int {
	return big.NewInt(0).Mod(x, m)
}