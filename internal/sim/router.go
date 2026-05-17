package sim

type Router interface {
	Route(node *Node, packet *Packet) *Node
}

type RoundRobinRouter struct{}

func (this *RoundRobinRouter) Route(node *Node, packet *Packet) *Node {
	if len(node.Outbound) == 0 {
		return nil
	}
	target := node.Outbound[node.roundRobinIdx]
	node.roundRobinIdx = (node.roundRobinIdx + 1) % len(node.Outbound)
	return target
}
