package resolver

import (
	"log"
	"net"

	"github.com/miekg/dns"
)

type Resolver struct {
}

func New() Resolver {
	return Resolver{}
}

func (r Resolver) ResolveA(addr string) (net.IP, bool) {
	// log.Println("Looking up A for " + addr)

	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{addr, dns.TypeA, dns.ClassINET}

	c := new(dns.Client)
	in, _, err := c.Exchange(m1, "10.1.3.254:53")
	if err != nil {
		log.Println("Failed on c.Exchange call: " + err.Error())
		return nil, false
	}

	if a, ok := in.Answer[0].(*dns.A); ok {
		// log.Println("I found an answer for " + addr)
		return a.A, true
	}
	// log.Println("Nothing found for " + addr)
	return nil, false
}
