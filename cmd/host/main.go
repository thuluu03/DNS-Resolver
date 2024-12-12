package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	resolver "github.com/thuluu03/DNS-Resolver/pkg"
)

/*
reads in repl input
*/

func main() {
	log.SetOutput(os.Stdout)

	keyboardChan := make(chan string, 1)

	// make a thread for taking in keyboard commands
	go func() {
		scanner := bufio.NewScanner(os.Stdin)

		// Wait for a line from stdin, convert it to an int
		for scanner.Scan() {
			line := scanner.Text()
			// Send an integer to the channel
			keyboardChan <- line
		}
	}()

	for {
		fmt.Print("> ")
		input := <-keyboardChan

		parsed_inputs := strings.Split(input, " ")

		if len(parsed_inputs) != 2 {
			fmt.Println("-[r] [domain-name] \n \t r=t or r=f")
			continue
		}

		recur := string(parsed_inputs[1][1:]) // should return t or f
		domain_name := string(parsed_inputs[2])

		if recur == "t" {
			ans := resolver.Recursive_resolve(domain_name)
			print(ans)
		} else if recur == "f" {
			ans := resolver.Iterative_resolve(domain_name)
			print(ans)
		}

	}
}
