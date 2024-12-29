package dns

import (
	"github.com/ducng99/gohole/internal/logger"
	"github.com/miekg/dns"
)

type dnsHandler struct{}

func StartDNSServer() error {
	StartCacheCleaner()

	handler := new(dnsHandler)
	server := &dns.Server{
		Addr:      ":53",
		Net:       "udp",
		Handler:   handler,
		UDPSize:   65535,
		ReusePort: true,
	}

	logger.Printf(logger.LogNormal, "Starting DNS server on port 53\n")

	return server.ListenAndServe()
}

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true

	for _, question := range r.Question {
		answers := resolver(question.Name, question.Qtype)
		msg.Answer = append(msg.Answer, answers...)
	}

	w.WriteMsg(msg)
}
