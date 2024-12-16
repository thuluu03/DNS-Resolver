package pkg

import (
	"net" 
	"fmt"
	dns "github.com/miekg/dns" 
)





//how to do root IPs as resource record?


func Iterative_resolve(query string, resourceRecords []dns.RR) (dns.RR) {
	for _, rr := range resourceRecords {
		data := rr.String()[:rr.Header().Rdlength] //Rdlength is length of data after header
		//rdata is just the IP record
		dnsResponse, err := Send_query(data, query, false) //will return the entire message
		if err != nil {
			continue // couldn't establish connection, go to the next server 
		} 
		if (len(dnsResponse.Answer) == 1) {
			return dnsResponse.Answer[0]  //will return the entire rr
		} else if (len(dnsResponse.Extra) >= 1) {
			Iterative_resolve(query, dnsResponse.Extra)  //returns a list of all the authority servers' IP addresses
		}
	}
	return nil //this is the case if we do not find an answer?
}

func Recursive_resolve(query string) (dns.RR) {
	dnsResponse, err := Send_query("8.8.8.8", query, true)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return dnsResponse.Answer[0]
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

	//serialize the query
	msg := new(dns.Msg)
	msg.SetQuestion(query, dns.TypeA)
	msg.RecursionDesired = recur
	
	//receiving from the socket
	buffer := make([]byte, 512)
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
	conn, err := net.Dial("udp4", server_ip)
    if err != nil {
        fmt.Println("Error:", err)
        return nil, err
    }

	return conn, nil
}
