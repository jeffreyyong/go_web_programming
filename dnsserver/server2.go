package main

import (
	"log"
	"net"
	"strconv"

	"github.com/miekg/dns"
)

func main() {
	dns.HandleFunc(".", handleDnsRequest)

	// start server
	port := 5002
	server := &dns.Server{Addr: "127.0.0.1:" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()

	if err != nil {
		log.Fatalf("Failed to start server: %s\n", err.Error())
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	var dnsAnswerIP net.IP
	var dnsAnswerTTL uint32
	var dnsAnswerTarget string

	dnsAnswerIP = net.ParseIP("1.1.1.1")
	// dnsAnswerTTL = 3600
	dnsAnswerTTL = 0
	dnsAnswerTarget = "cdomain-name-test."

	m := new(dns.Msg)
	m.SetReply(r)
	if dnsAnswerIP != nil {
		m.Answer = []dns.RR{
			&dns.A{
				Hdr: dns.RR_Header{
					Name:   m.Question[0].Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    dnsAnswerTTL,
				},
				A: dnsAnswerIP,
			},
			&dns.CNAME{
				Hdr: dns.RR_Header{
					Name:   m.Question[0].Name,
					Rrtype: dns.TypeCNAME,
					Class:  dns.ClassINET,
					Ttl:    dnsAnswerTTL,
				},
				Target: dnsAnswerTarget,
			},
		}
	} else {
		m.Answer = []dns.RR{}
	}
	w.WriteMsg(m)
}
