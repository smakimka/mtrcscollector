package interceptors

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type SubnetInterseptor struct {
	TrustedSubnet *net.IPNet
}

func NewSubnetInterseptor(subnet *net.IPNet) *SubnetInterseptor {
	return &SubnetInterseptor{TrustedSubnet: subnet}
}

func (i *SubnetInterseptor) AllowTrusted(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var err error

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		realIP := md.Get("X-Real-IP")

		if len(realIP) != 1 {
			return nil, status.Error(codes.Unauthenticated, "unauthorized")
		}

		if realIP[0] == "" {
			return nil, status.Error(codes.Unauthenticated, "unauthorized")
		}

		ip := net.ParseIP(realIP[0])

		if !i.TrustedSubnet.Contains(ip) {
			return nil, status.Error(codes.Unauthenticated, "unauthorized")
		}

		return handler(ctx, req)
	}
	return err, status.Error(codes.Unauthenticated, "unauthorized")
}
