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

var root_ips = map[string]string{"a.root-servers.net": "198.41.0.4"}

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

		recur := string(parsed_inputs[0][2:]) // should return t or f
		query := string(parsed_inputs[1])

		if recur == "t" {
			fmt.Println("performing recursive resolver")
			ans := resolver.Recursive_resolve(query)

			if ans != nil {
				print(ans.String())
			} else {
				print("No answer found for this query")
			}
			
		} else if recur == "f" {

			//send initial query to root server
			firstResponse, err := resolver.Send_query(root_ips["a.root-servers.net"], query, false)

			if err != nil {
				fmt.Println("Error in asking root: ", err)
				return
			}

			//call Iterative resolve on the first response...
			if (len(firstResponse.Answer) == 1) {
				ans := firstResponse.Answer[0]  //if root server immediately returns answer
				print(ans)
			} else if (len(firstResponse.Extra) >= 1) { //otherwise, need to call the iterative resolver on the first set of responses
				ans := resolver.Iterative_resolve(query, firstResponse.Extra)
				print(ans)
			}

			
			
		}

	}
}
