package sim

type Processor interface {
	Process(node *Node, packet *Packet) ([]*Packet, error)
}

type Action string

const (
	ActionTransform Action = "TRANSFORM"
	ActionDelay     Action = "DELAY"
	ActionFetch     Action = "FETCH_STATE"
	ActionStore     Action = "STORE_STATE"
	ActionForward   Action = "FORWARD"
	ActionReturn    Action = "RETURN"
	ActionBlock     Action = "BLOCK"
	ActionResume    Action = "RESUME"
)

type WorkflowStep struct {
	Action string      `json:"action"` //"FETCH", "STORE", "TRANSFORM
	Target string      `json:"target"` // eg: cache_1, db_read
	Value  interface{} `json:"value"`  // Parameters for action
}

type WorkflowProcessor struct {
	Workflows map[string][]WorkflowStep
}

func (this *WorkflowProcessor) Process(node *Node, packet *Packet) ([]*Packet, error) {
	steps, ok := this.Workflows[packet.Type]
	if !ok {
		return []*Packet{packet}, nil
	}

	var outbound []*Packet

	// Execution loop
	for _, step := range steps {
		// Add process steps
		switch step.Action {
		case string(ActionTransform):
			packet.Type = step.Target
		case string(ActionStore):
			node.State[step.Target] = packet.Payload
		case string(ActionFetch):
			if val, exists := node.State[step.Target]; exists {
				if packet.Payload == nil {
					packet.Payload = make(map[string]any)
				}
				packet.Payload["fetched_data"] = val
			}
		case string(ActionBlock):
			node.WaitingRoom[packet.TraceId] = packet
			return nil, nil
		case string(ActionResume):
			if original, exists := node.WaitingRoom[packet.TraceId]; exists {
				delete(node.WaitingRoom, packet.TraceId)
				outbound = append(outbound, original)
			}
		case string(ActionForward):
			outbound = append(outbound, packet)
		}
	}

	if len(outbound) == 0 {
		outbound = append(outbound, packet)
	}

	return outbound, nil
}
