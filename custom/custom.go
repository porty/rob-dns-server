package custom

import (
	"net"
	"sync"

	"github.com/miekg/dns"
)

type Custom struct {
	entriesMutex *sync.Mutex
	entries      map[string]net.IP
}

func New() Custom {
	return Custom{
		entriesMutex: new(sync.Mutex),
		entries:      make(map[string]net.IP),
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

func (c *Custom) Find(req *dns.Msg) *dns.Msg {
	q := req.Question[0]
	if q.Qclass == dns.ClassINET && q.Qtype == dns.TypeA {
		ip, ok := c.FindA(q.Name)
		if ok {
			return buildAReply(req, q, ip)
		}
	}
	return nil
}

func (c *Custom) FindA(address string) (net.IP, bool) {
	c.entriesMutex.Lock()
	defer c.entriesMutex.Unlock()

	ip, ok := c.entries[address]
	return ip, ok
}

func (c *Custom) AddA(address string, ip net.IP) {
	c.entriesMutex.Lock()
	defer c.entriesMutex.Unlock()

	c.entries[address] = ip
}
