package middleware

import (
	"net"
	"net/http"
)

type SubnetMiddleware struct {
	TrustedSubnet *net.IPNet
}

func NewSubnetMiddleware(subnet *net.IPNet) *SubnetMiddleware {
	return &SubnetMiddleware{TrustedSubnet: subnet}
}

func (m *SubnetMiddleware) AllowTrusted(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		realIP := r.Header.Get("X-Real-IP")
		if realIP == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ip := net.ParseIP(realIP)

		if !m.TrustedSubnet.Contains(ip) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)

	})
}
