package main

import (
	"log"
	"net"
	"time"

	"./cache"
	"./resolver"

	"github.com/miekg/dns"
)

var lookupCache = cache.New()
var r = resolver.New()

func main() {
	server := &dns.Server{Addr: ":5533", Net: "udp"}

	dns.HandleFunc(".", handleRequest)

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func resolve(address string, c chan net.IP) bool {
	// time.Sleep(1 * time.Second)
	// ip := net.IPv4(1, 2, 3, 4)
	// c <- ip
	// //addToCache(address, ip)
	// lookupCache.AddA(address, ip)
	// return true

	ip, ok := r.ResolveA(address)
	if ok {
		c <- ip
		lookupCache.AddA(address, ip)
		return true
	}
	return false
}

func handleRequest(w dns.ResponseWriter, req *dns.Msg) {
	q := req.Question[0]
	if q.Qclass != dns.ClassINET || q.Qtype != dns.TypeA {
		log.Println("Unknown/uninteresting query type")

		m := new(dns.Msg)
		m.SetReply(req)
		w.WriteMsg(m)
		return
	}
	log.Println("I have a query for " + q.Name)

	rrHeader := dns.RR_Header{
		Name:   q.Name,
		Rrtype: dns.TypeA,
		Class:  dns.ClassINET,
		Ttl:    3, // lol I dunno
	}

	m := new(dns.Msg)
	m.SetReply(req)
	c := make(chan net.IP, 1)
	var ip net.IP
	resolved := false
	go func(addr string, c chan net.IP) {
		//if !searchCache(addr, c) {
		//	resolve(addr, c)
		//}
		ip, ok := lookupCache.FindA(addr)
		if ok {
			c <- ip
			log.Println("  Answer found in cache")
			return
		}
		if resolve(addr, c) {
			log.Println("  Answer found through query")
		}
	}(q.Name, c)
	select {
	case ip = <-c:
		resolved = true
		break
	case <-time.After(2 * time.Second):
		break
	}
	if resolved {
		a := &dns.A{
			Hdr: rrHeader,
			A:   ip, //net.IPv4(1, 2, 3, 4),
		}
		m.Answer = append(m.Answer, a)
	}
	w.WriteMsg(m)
}
