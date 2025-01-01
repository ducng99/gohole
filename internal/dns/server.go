package dns

import (
	"os"
	"os/signal"
	"syscall"

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

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Printf(logger.LogError, "[DNS] Server error: %v\n", err)
		}
	}()

	logger.Printf(logger.LogSuccess, "DNS server listening on port 53\n")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	if err := server.Shutdown(); err != nil {
		logger.Fatalf("DNS server shutdown error: %v\n", err)
	}

	logger.Printf(logger.LogNormal, "DNS server graceful shutdown.\n")
	return nil
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
