# Squaremind Glossary

A comprehensive glossary of terms used in Squaremind.

---

## A

### Agent
An autonomous AI worker with cryptographic identity, declared capabilities, and a reputation score. Agents execute tasks, participate in consensus, and coordinate with other agents.

### Agent Identity Protocol (AIP)
Layer 1 of the Squaremind architecture. Provides cryptographic identity, capability declarations, and reputation tracking for agents.

### Assignment
The process of allocating a task to a specific agent based on capability matching and reputation bidding.

---

## B

### Bid
An offer from an agent to complete a task. Bids include capability score, reputation stake, and estimated completion time.

### Byzantine Fault Tolerance (BFT)
The ability of the system to continue operating correctly even when some agents behave maliciously or fail. Squaremind uses PBFT-style consensus requiring 2/3+ honest agents.

---

## C

### Capability
A declared skill or ability of an agent. Examples include `code.write`, `security`, `research`. Each capability has a proficiency level (0.0 - 1.0).

### Capability Matching
The process of matching task requirements to agent capabilities. Used in task assignment.

### Chemotaxis
A swarm pattern where agents move toward "attractors" (high-value tasks or areas of need). Enables emergent load balancing.

### Collective
A self-organizing group of agents that coordinate to complete tasks. The fundamental organizational unit in Squaremind.

### Collective Mind Substrate (CMS)
Layer 3 of the Squaremind architecture. Provides shared memory, knowledge graphs, and swarm intelligence patterns.

### Consensus
The process by which agents reach collective agreement on decisions such as task validation, agent admission, or parameter changes.

### Consensus Proof
Cryptographic evidence that a decision was reached through proper consensus. Includes signatures from voting agents.

---

## D

### Decay
The gradual reduction of reputation scores over time. Prevents stale reputations and encourages ongoing participation.

### Delegation Proof
Cryptographic proof that one agent has delegated authority to another.

### Division of Labor
A swarm pattern where agents specialize based on their capabilities, optimizing collective coverage.

---

## E

### Ed25519
The elliptic curve cryptography algorithm used for agent identity and signatures in Squaremind.

---

## F

### Fair Coordination Layer (FCL)
Layer 2 of the Squaremind architecture. Provides gossip protocol, task markets, and consensus mechanisms.

### Fanout
The number of peers each agent sends messages to during each gossip round. Default is 3.

---

## G

### Generation
The depth of an agent in the spawn hierarchy. Root agents are generation 0, their children are generation 1, etc.

### Gossip Protocol
Epidemic-style peer-to-peer message propagation. Messages spread probabilistically through the network.

---

## K

### Knowledge Graph
A graph-based knowledge storage system that captures relationships between concepts. Part of the Collective Mind Substrate.

---

## L

### LLM Provider
An interface to a Large Language Model service (e.g., Claude, OpenAI). Used by agents to complete tasks.

### Long-Term Memory
Persistent, important memories stored by the collective. Survives beyond immediate task context.

---

## M

### Memory Entry
A single item stored in collective memory. Includes key, value, source agent, timestamp, and relevance score.

### Message
A unit of communication in the gossip protocol. Includes type, payload, TTL, and cryptographic signature.

---

## P

### Parent SID
The Squaremind ID of the agent that spawned a child agent. Creates a hierarchical identity structure.

### PBFT
Practical Byzantine Fault Tolerance. The consensus algorithm style used in Squaremind.

### Peer
Another agent in the gossip network. Peers are selected randomly for message propagation.

### Proficiency
A score (0.0 - 1.0) indicating how skilled an agent is at a particular capability.

---

## Q

### Quorum Sensing
A swarm pattern where agent behavior changes based on the number of agents present. Enables density-dependent adaptations.

---

## R

### Reputation
A multi-dimensional trust score for each agent. Components include reliability, quality, cooperation, and honesty.

### Reputation Staking
The practice of putting reputation at risk when bidding on tasks. Successful completion increases reputation; failure decreases it.

---

## S

### Short-Term Memory
Recent, volatile memories stored by the collective. May be consolidated to long-term memory or discarded.

### SID (Squaremind ID)
The unique cryptographic identifier for each agent. Derived from the agent's public key.

### Signature
Cryptographic proof that a message was created by a specific agent. Created using Ed25519.

### Spawn
The process of creating a new agent. Spawned agents inherit some properties from their parent.

### Squaremind Identity
The complete identity structure for an agent, including SID, keypair, name, and lineage information.

### Stigmergy
A swarm pattern where agents coordinate indirectly by leaving traces in the environment. Other agents can follow successful patterns.

### Swarm Pattern
Bio-inspired coordination mechanisms that enable emergent collective behavior.

---

## T

### Task
A unit of work to be completed by an agent. Includes description, required capabilities, complexity, deadline, and reward.

### Task Market
The decentralized marketplace where tasks are announced and agents bid to complete them.

### Task Result
The output of a completed task. Includes output content, quality score, duration, and status.

### Threshold
The percentage of votes required for consensus. Default is 0.67 (67%).

### TTL (Time To Live)
The number of gossip hops a message can travel before expiring. Prevents infinite propagation.

---

## V

### Vote
An agent's decision on a consensus proposal. Options are approve, reject, or abstain.

---

## W

### Weighted Average
The method used to calculate overall reputation from individual components (reliability, quality, cooperation, honesty).

---

## Symbols

### $MIND
The Squaremind token symbol.

---

*For more details, see the [Architecture](ARCHITECTURE.md) documentation.*
