package main

import (
	"os"
	"flag"
	"strconv"
	"bufio"
	"fmt"
	"math/rand"
	"time"
	"strings"
	"math/big"
)
func main() {

	gen   := flag.Bool("gen", false, "Start generating a new puzzle")
	solve := flag.Bool("solve", false, "Start solving a puzzle")

	fmt.Println(*gen)
	fmt.Println(*solve)
	
	fmt.Println("")
	fmt.Println("Time-Lock Puzzle")
	fmt.Println("****************")

	fmt.Println(`
A time-lock puzzle is a mathematical puzzle that takes a bounded time
to solve, but can be created much more efficiently. It was proposed by
Rivest, Shamir, and Wagner in [1]. 
In effect, solving a time-lock puzzle yields a solution only after fi-
xed number of operations which are inherently non-parallelizable; it is
used to force the solver to use a certain quantity of resources (time, 
and/or computational power) to discover the solution.

Concretely, Alice creates a puzzle with some public parameters and a 
number of cycles C, each cycle taking around 1 sec for a common laptop 
to compute. Alice is able to compute the solution immediately, using RSA
properties and Fermat's little theorem. Alice then sends the puzzle to 
Bob, which cannot find the solution without computing all the C cycles, 
which takes C * 1 seconds.

Alice can use this system to make sure some content is not available be-
fore some time (by encrypting with the puzzle's solution as a key, and
by setting C high enough).

[1] https://people.csail.mit.edu/rivest/pubs/RSW96.pdf

******

The software is meant as a fun exercice, not something secure under all 
conditions. All usage is at your own risk!

******
`)

	if *gen {
		preparePuzzle()
	} else if *solve {
		solvePuzzle()
	} else {
		fmt.Println("Please add the flag -gen or -solve in the command line for the desired action.")
	}

	fmt.Println("****************")
}

func preparePuzzle() {

	//randomness
	rndSource := rand.NewSource(time.Now().UnixNano())
	rnd       := rand.New(rndSource)

	expOfExp     := 50000
	expOneSecond := ExpByPowOfTwo(big.NewInt(2), big.NewInt(int64(expOfExp)))

	//we generate the parameters
	p := newRandomPrime(1024)
	q := newRandomPrime(1024)
	one := big.NewInt(1)
	p_1 := big.NewInt(0).Sub(p, one)
	q_1 := big.NewInt(0).Sub(q, one)
	phi := big.NewInt(0).Mul(p_1, q_1)
	n := big.NewInt(0).Mul(p, q)

	fmt.Println("Generation of parameters done...")
	fmt.Println("Benchmarking, each step should take around 1 sec.")

	//generate benchmark, a 1-second exponentiation
	repeat := 10
	sum := int64(0)

	for i:=0; i<repeat; i++ {
		rand      := big.NewInt(0).Rand(rnd, n)
		start     := time.Now()

		ExpByPowOfTwoModular(rand, expOneSecond, n)

		diff := time.Now().Sub(start)
		sum += diff.Nanoseconds()
		fmt.Println("Benchmark", i, "of", repeat, "...")
	}

	mean := time.Duration(sum / 10) * time.Nanosecond
	fmt.Println("On this machine, one cycle takes ", mean)
	fmt.Println("Gonna generate the puzzle...")

	proceed := false
	x := big.NewInt(0)
	y := big.NewInt(0)

	for !proceed{
		x = big.NewInt(0).Rand(rnd, n)

		y = ReadBigFromConsole("Enter number of cycles needed to solve the puzzle (keep in mind that on other machines, one cycle might go faster or slower than the expected 1 sec) : ")

	    duration := time.Duration(int(y.Int64())) * time.Second
	    durationThisMachine := time.Duration(int(y.Int64())) * mean
	    fmt.Println("You entered ", y, ", expected duration is", duration, " (", durationThisMachine, "on this machine). Proceed [y/n] ?")
 	    reader := bufio.NewReader(os.Stdin)
	    ans, _ := reader.ReadString('\n')
	    ans = strings.Replace(ans, "\n", "", 1)

	    if ans == "y" || ans == "Y" {
	    	proceed = true
	    }
	}

	curr    := big.NewInt(0)
	base    := x
	iter    := y

	fmt.Println("Generating...")

	exponent := one
	for curr.Cmp(y) < 0 {
		exponent = mod(big.NewInt(0).Mul(exponent, expOneSecond), phi)
		curr.Add(curr,one)

	}

	x = ExpByPowOfTwoModular(x, exponent, n)
	timeToUnlock := time.Duration(int(y.Int64()))*time.Second
	sha := sha256AndHex(x.Bytes())
	shaOfSha := sha256AndHex([]byte(sha))

	fmt.Println("Time-Lock Puzzle finished ! the solution is : ")
	fmt.Println("")
	fmt.Println(sha)
	fmt.Println("")
	fmt.Println("Please save the following JSON to recompute the solution (expected time :", timeToUnlock, ")")

	filledPuzzle := &Puzzle{n.String(), timeToUnlock.String(), iter.String(), base.String(), strconv.Itoa(expOfExp), shaOfSha}

	fmt.Println("")
	fmt.Println(JsonToString(filledPuzzle))
	fmt.Println("")
}


func solvePuzzle() {

	fmt.Println("Please paste the JSON puzzle :")

	reader := bufio.NewReader(os.Stdin)
	ans, _ := reader.ReadString('\n')
	ans    = strings.Replace(ans, "\n", "", 1)

	puzzle := StringToJson(ans)

	//we generate the parameters
	n, _        := big.NewInt(0).SetString(puzzle.N, 10)
	base, _     := big.NewInt(0).SetString(puzzle.Base, 10)
	iter, _     := big.NewInt(0).SetString(puzzle.NCycles, 10)
	expOfExp, _ := big.NewInt(0).SetString(puzzle.ExponentOneSecond, 10)

	exp         := ExpByPowOfTwo(big.NewInt(2), expOfExp)

	fmt.Println("All right ! Going to solve...")
	fmt.Println("Expected time to unlock : ", puzzle.TimeToUnlock)

	x       := base
	curr    := big.NewInt(0)
	one     := big.NewInt(1)
	zero    := big.NewInt(0)
	five := big.NewInt(5)
	start   := time.Now()

	for curr.Cmp(iter) < 0 {
		if mod(curr, five).Cmp(zero) == 0 {
			diff := big.NewInt(0).Sub(iter, curr)
			rem := time.Duration(int(diff.Int64())) * time.Second
			fmt.Println("Iteration", curr.String()+"/"+iter.String(), ", remaining time is", rem)
		}
		
		x = ExpByPowOfTwoModular(x, exp, n)

		curr.Add(curr,one)
	}

	fmt.Println("Puzzle solved, verifying...")
	sha := sha256AndHex(x.Bytes())
	shaOfSha := sha256AndHex([]byte(sha))

	if shaOfSha != puzzle.Sha256OfSha256 {
		fmt.Println("An error occured. The solution isn't the expected one.")
	} else {
		timeNeeded := time.Now().Sub(start)
		fmt.Println("################")
		fmt.Println("Time-Lock Puzzle solved in", timeNeeded, "!")
		fmt.Println("")
		fmt.Println(sha)
		fmt.Println("")
		fmt.Println("################")
	}
}