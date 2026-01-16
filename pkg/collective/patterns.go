package collective

import (
	"context"
	"sync"
	"time"

	"github.com/squaremind/squaremind/pkg/agent"
	"github.com/squaremind/squaremind/pkg/identity"
)

// SwarmPattern represents a swarm intelligence pattern
type SwarmPattern string

const (
	PatternStigmergy      SwarmPattern = "stigmergy"       // Indirect coordination through environment
	PatternQuorumSensing  SwarmPattern = "quorum_sensing"  // Behavior based on collective size
	PatternChemotaxis     SwarmPattern = "chemotaxis"      // Movement toward/away from signals
	PatternDivisionOfLabor SwarmPattern = "division_of_labor" // Dynamic role specialization
	PatternSwarmOptimization SwarmPattern = "swarm_optimization" // Collective search for optima
)

// PatternExecutor executes swarm patterns
type PatternExecutor struct {
	mu sync.RWMutex

	collective *Collective
	patterns   map[SwarmPattern]PatternHandler
}

// PatternHandler handles a specific swarm pattern
type PatternHandler func(ctx context.Context, collective *Collective, params map[string]interface{}) error

// NewPatternExecutor creates a new pattern executor
func NewPatternExecutor(collective *Collective) *PatternExecutor {
	pe := &PatternExecutor{
		collective: collective,
		patterns:   make(map[SwarmPattern]PatternHandler),
	}

	// Register default patterns
	pe.patterns[PatternStigmergy] = pe.executeStigmergy
	pe.patterns[PatternQuorumSensing] = pe.executeQuorumSensing
	pe.patterns[PatternChemotaxis] = pe.executeChemotaxis
	pe.patterns[PatternDivisionOfLabor] = pe.executeDivisionOfLabor
	pe.patterns[PatternSwarmOptimization] = pe.executeSwarmOptimization

	return pe
}

// Execute runs a swarm pattern
func (pe *PatternExecutor) Execute(ctx context.Context, pattern SwarmPattern, params map[string]interface{}) error {
	pe.mu.RLock()
	handler, ok := pe.patterns[pattern]
	pe.mu.RUnlock()

	if !ok {
		return nil
	}

	return handler(ctx, pe.collective, params)
}

// RegisterPattern registers a custom pattern handler
func (pe *PatternExecutor) RegisterPattern(pattern SwarmPattern, handler PatternHandler) {
	pe.mu.Lock()
	defer pe.mu.Unlock()
	pe.patterns[pattern] = handler
}

// executeStigmergy implements stigmergic coordination
// Agents leave "traces" in shared memory that influence other agents
func (pe *PatternExecutor) executeStigmergy(ctx context.Context, c *Collective, params map[string]interface{}) error {
	signal, _ := params["signal"].(string)
	strength, _ := params["strength"].(float64)
	if strength == 0 {
		strength = 1.0
	}

	// Leave a trace in collective memory
	c.memory.Contribute("stigmergy", signal, map[string]interface{}{
		"strength":  strength,
		"timestamp": time.Now(),
		"pattern":   "stigmergy",
	})

	return nil
}

// executeQuorumSensing implements quorum sensing
// Behavior changes based on collective size/state
func (pe *PatternExecutor) executeQuorumSensing(ctx context.Context, c *Collective, params map[string]interface{}) error {
	threshold, _ := params["threshold"].(int)
	if threshold == 0 {
		threshold = 5
	}

	action, _ := params["action"].(string)

	// Check if quorum is met
	size := c.Size()
	if size >= threshold {
		// Quorum reached - trigger action
		c.memory.Contribute("quorum_sensing", "Quorum reached: "+action, map[string]interface{}{
			"size":      size,
			"threshold": threshold,
			"action":    action,
		})
		return nil
	}

	return nil
}

// executeChemotaxis implements chemotaxis-like behavior
// Agents move toward/away from "signals"
func (pe *PatternExecutor) executeChemotaxis(ctx context.Context, c *Collective, params map[string]interface{}) error {
	signal, _ := params["signal"].(string)
	attract, _ := params["attract"].(bool)

	// In a real implementation, this would influence task prioritization
	// Agents would be drawn to tasks matching the signal
	action := "attracted to"
	if !attract {
		action = "repelled from"
	}

	c.memory.Contribute("chemotaxis", "Collective "+action+" signal: "+signal, map[string]interface{}{
		"signal":  signal,
		"attract": attract,
	})

	return nil
}

// executeDivisionOfLabor implements dynamic role specialization
func (pe *PatternExecutor) executeDivisionOfLabor(ctx context.Context, c *Collective, params map[string]interface{}) error {
	// Analyze collective capabilities and assign roles dynamically
	agents := c.GetAgents()

	// Group agents by primary capability
	capGroups := make(map[identity.CapabilityType][]*agent.Agent)
	for _, a := range agents {
		caps := a.Capabilities.List()
		if len(caps) > 0 {
			primaryCap := caps[0] // First capability is primary
			capGroups[primaryCap] = append(capGroups[primaryCap], a)
		}
	}

	// Record division of labor
	c.memory.Contribute("division_of_labor", "Labor divided by capability", map[string]interface{}{
		"groups": len(capGroups),
		"agents": len(agents),
	})

	return nil
}

// executeSwarmOptimization implements collective search for optima
func (pe *PatternExecutor) executeSwarmOptimization(ctx context.Context, c *Collective, params map[string]interface{}) error {
	objective, _ := params["objective"].(string)
	iterations, _ := params["iterations"].(int)
	if iterations == 0 {
		iterations = 10
	}

	// Simplified particle swarm optimization concept
	// In real implementation, agents would explore solution space
	c.memory.Contribute("swarm_optimization", "Optimizing: "+objective, map[string]interface{}{
		"objective":  objective,
		"iterations": iterations,
	})

	return nil
}

// EmergentBehavior represents emergent collective behavior
type EmergentBehavior struct {
	Name        string
	Description string
	Trigger     func(*Collective) bool
	Execute     func(context.Context, *Collective) error
}

// EmergentBehaviorMonitor monitors for emergent behaviors
type EmergentBehaviorMonitor struct {
	mu sync.RWMutex

	behaviors []EmergentBehavior
	collective *Collective
}

// NewEmergentBehaviorMonitor creates a new monitor
func NewEmergentBehaviorMonitor(collective *Collective) *EmergentBehaviorMonitor {
	return &EmergentBehaviorMonitor{
		behaviors:  make([]EmergentBehavior, 0),
		collective: collective,
	}
}

// RegisterBehavior registers an emergent behavior
func (m *EmergentBehaviorMonitor) RegisterBehavior(b EmergentBehavior) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.behaviors = append(m.behaviors, b)
}

// Monitor checks for and executes emergent behaviors
func (m *EmergentBehaviorMonitor) Monitor(ctx context.Context) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, behavior := range m.behaviors {
		if behavior.Trigger(m.collective) {
			behavior.Execute(ctx, m.collective)
		}
	}
}

// Start begins continuous monitoring
func (m *EmergentBehaviorMonitor) Start(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.Monitor(ctx)
		}
	}
}

// DefaultEmergentBehaviors returns a set of default emergent behaviors
func DefaultEmergentBehaviors() []EmergentBehavior {
	return []EmergentBehavior{
		{
			Name:        "load_balancing",
			Description: "Redistribute tasks when agents are overloaded",
			Trigger: func(c *Collective) bool {
				// Trigger when any agent has been working too long
				for _, a := range c.GetAgents() {
					if a.GetState() == agent.StateWorking {
						if time.Since(a.LastActive) > 5*time.Minute {
							return true
						}
					}
				}
				return false
			},
			Execute: func(ctx context.Context, c *Collective) error {
				// In real implementation, would reassign tasks
				return nil
			},
		},
		{
			Name:        "collective_learning",
			Description: "Share successful strategies across collective",
			Trigger: func(c *Collective) bool {
				// Trigger periodically based on completed tasks
				return c.Stats().CompletedTasks%10 == 0
			},
			Execute: func(ctx context.Context, c *Collective) error {
				// In real implementation, would aggregate and share learnings
				return nil
			},
		},
	}
}
