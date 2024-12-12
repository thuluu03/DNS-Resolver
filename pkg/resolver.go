package pkg

import (
	"net" 
	"fmt"
	dns "github.com/miekg/dns" 
)




const (
	root_ips = map[string]string{"a.root-servers.net": "198.41.0.4"}
)

func Iterative_resolve(query string, resourceRecords []dns.RR) (dns.RR) {
	for _, rr := range resourceRecords {
		data := rr.String()[rr.Header().Rdlength:]

		ans, err := send_query(data, query, false)
		if err != nil {
			continue // couldn't establish connection, go to the next server 
		} 
		if (len(ans.Answer) == 1) {
			return ans.Answer[0]  //will return the entire rr
		} else if (len(ans.Extra) >= 1) {
			iterative_resolve(ans.Extra)  //returns a list of all the authority servers
		}
	}
		
		//iterative => can receive either the answer or any type of response
		
}

func Recursive_resolve(query string) []dns.RR {
	ans, err := send_query("8.8.8.8", query, true)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return ans
}

//recursive:
//let the DNS server query OTHER DNS servers
//can only recieve the answer


func send_query(server_ip_addr string, query string, recur bool) ([]dns.RR, error){

	conn, err := create_socket(server_ip_addr)
	if err != nil {
		return nil, err
	}

	//serialize the query
	msg := new(dns.Msg)
	msg.SetQuestion(query, dns.TypeA)
	msg.RecursionDesired = recur
	
	//receiving from the socket
	var buf bytes.Buffer //dynamic buffer size
    dnsPacket, err := conn.Read(buf)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

	m := new(dns.Msg)
	m.unpack(dnsPacket) 	

	
	return m.Answer
}
//create socket each time you send a query
//serialize into DNS packet
//send through socket to server 

// read msg from socket
// store in cache? 

// close the socket once we receive a response



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
