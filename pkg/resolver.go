package pkg

import (
	"net" 
	"fmt"
	dns "github.com/miekg/dns" 
)


func Iterative_resolve(query string, msg *dns.Msg) (*dns.Msg) {
	resourceRecords := msg.Extra
	// fmt.Println(msg.String())
	// fmt.Println("first this in additional section: ", resourceRecords[0])

	for _, rr := range resourceRecords {
		//only consider A types

		//if it is type A
		if aRecord, ok := rr.(*dns.A); ok { //only do this for ipv4 addersses
			data := aRecord.A.String()  //get the ip address
		
			//rdata is just the IP record
			dnsResponse, err := Send_query(data, query, false) //will return the entire message
			if err != nil {
				continue // couldn't establish connection, go to the next server 
			} 

			// fmt.Println("response: ", dnsResponse.String())
			// fmt.Println("answer len: ", len(dnsResponse.Answer))

			if (len(dnsResponse.Answer) >= 1) {
				return dnsResponse  //will return the entire rr
			} else if (len(dnsResponse.Extra) >= 1) {
				return Iterative_resolve(query, dnsResponse)  //returns a list of all the authority servers' IP addresses
			} else {
				//return
				return nil
			}
		}
	}
	return nil //this is the case if we do not find an answer?
}

func Recursive_resolve(query string) (*dns.Msg) {
	dnsResponse, err := Send_query("8.8.8.8", query, true) //can instead make the IP custom
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	if len(dnsResponse.Answer) > 0 { //we have found an answer
		return dnsResponse
		// TODO: Should we return all possible answers?
	} else { //no response
		return nil
	}
 }

//recursive:
//let the DNS server query OTHER DNS servers
//can only recieve the answer


//creates a new socket
//sends a query to the ip address
//receives on that new socket?
func Send_query(server_ip_addr string, query string, recur bool) (*dns.Msg, error){

	
	conn, err := create_socket(server_ip_addr)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	//serialize the query
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(query), dns.TypeA)
	msg.RecursionDesired = recur


	msgBytes, err := msg.Pack()
	if err != nil {
		return nil, fmt.Errorf("failed to pack query: %w", err)
	}

	// Send the query to the server
	_, err = conn.Write(msgBytes)

	if err != nil {
		return nil, fmt.Errorf("failed to write message bytes: %w", err)
	}

	//wait to recieve in this time
	
	//receiving from the socket
	buffer := make([]byte, 1024)

    n, err := conn.Read(buffer)

    if err != nil {
        fmt.Println("Error:", err)
        return nil, err
    }

	m := new(dns.Msg)

	err = m.Unpack(buffer[:n]) 	
    if err != nil {
        fmt.Println("Error unpacking message:", err)
        return nil, err
    }

	//return the whole message, not just the answer or ns section
	return m, nil

	// TODO: close socket once you get a response
}
//create socket each time you send a query
//serialize into DNS packet
//send through socket to server 

// read msg from socket
// store in cache? 
 
//if recursive = 8.8.8.8
//if iterative = root_ips["a.root-servers.net"]
func create_socket(server_ip string) (net.Conn, error) { //this will always be the root server
	conn, err := net.Dial("udp4", server_ip + ":53")
	//53 is default port?

	//missing port in address
    if err != nil {
        fmt.Println("Error:", err)
        return nil, err
    }

	return conn, nil
}
