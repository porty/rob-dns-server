package cache

import (
	"net"
	"sync"

	"github.com/miekg/dns"
)

type Cache struct {
	m        *sync.Mutex
	aEntries map[string]net.IP
}

func New() Cache {
	return Cache{
		m:        new(sync.Mutex),
		aEntries: make(map[string]net.IP),
	}
}

func buildEmptyReply(req *dns.Msg, q dns.Question) *dns.Msg {
	reply := new(dns.Msg)
	reply.SetReply(req)
	return reply
}

func buildAReply(req *dns.Msg, q dns.Question, ip net.IP) *dns.Msg {
	reply := buildEmptyReply(req, q)
	rrHeader := dns.RR_Header{
		Name:   q.Name,
		Rrtype: dns.TypeA,
		Class:  dns.ClassINET,
		Ttl:    3, // lol I dunno
	}
	a := &dns.A{
		Hdr: rrHeader,
		A:   ip,
	}
	reply.Answer = append(reply.Answer, a)
	return reply
}

func (c *Cache) Find(req *dns.Msg) *dns.Msg {
	q := req.Question[0]
	if q.Qclass == dns.ClassINET && q.Qtype == dns.TypeA {
		ip, ok := c.FindA(q.Name)
		if ok {
			return buildAReply(req, q, ip)
		}
	}
	return nil
}

func (c *Cache) FindA(address string) (net.IP, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	ip, ok := c.aEntries[address]
	return ip, ok
}

func (c *Cache) AddA(address string, ip net.IP) {
	c.m.Lock()
	defer c.m.Unlock()

	c.aEntries[address] = ip
}
