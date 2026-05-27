package sim

import "strings"

type Router interface {
	Route(node *Node, packet *Packet) *Node
}

type RoundRobinRouter struct{}

func (this *RoundRobinRouter) Route(node *Node, packet *Packet) *Node {
	if len(node.Outbound) == 0 {
		return nil
	}

	outbound := node.Outbound
	if packet != nil && packet.RouteTarget != "" {
		outbound = matchingOutbound(node.Outbound, packet.RouteTarget)
		if len(outbound) == 0 {
			return nil
		}
	}

	target := outbound[node.roundRobinIdx%len(outbound)]
	node.roundRobinIdx = (node.roundRobinIdx + 1) % len(outbound)
	return target
}

func matchingOutbound(outbound []*Node, routeTarget string) []*Node {
	var matches []*Node
	for _, candidate := range outbound {
		if matchesRouteTarget(candidate, routeTarget) {
			matches = append(matches, candidate)
		}
	}
	return matches
}

func matchesRouteTarget(node *Node, routeTarget string) bool {
	if node == nil {
		return false
	}

	switch strings.ToLower(routeTarget) {
	case "load_balancer", "loadbalancer":
		return node.Type == TypeLoadBalancer
	case "server_pool", "server", "servers":
		return node.Type == TypeServer
	case "cache_tier", "cache":
		return node.Type == TypeCache
	case "database", "data_base", "db":
		return node.Type == TypeDatabase
	case "message_queue", "queue":
		return node.Type == TypeMessageQueue
	default:
		return node.ID == routeTarget
	}
}
