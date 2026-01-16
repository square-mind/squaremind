# @squaremind/sdk

TypeScript SDK for **Squaremind** - the multi-agent orchestration protocol for autonomous AI collectives.

Many Agents. One Mind. Fair coordination. Transparent markets.

## Installation

```bash
npm install @squaremind/sdk
```

## Quick Start

```typescript
import { Collective, Agent } from '@squaremind/sdk';

// Create a collective
const collective = new Collective({
  name: 'DevSwarm',
  maxAgents: 100,
  consensusThreshold: 0.67,
});

// Spawn agents with capabilities
const coder = new Agent({
  name: 'Coder',
  capabilities: ['code.write', 'code.review'],
  model: 'claude-sonnet-4-20250514',
});

const reviewer = new Agent({
  name: 'Reviewer',
  capabilities: ['code.review', 'security'],
  model: 'claude-sonnet-4-20250514',
});

// Join the collective
collective.join(coder);
collective.join(reviewer);

// Start the collective
await collective.start();

// Submit a task
const task = collective.createTask(
  'Implement user authentication with OAuth2',
  ['code.write']
);

const result = await collective.submit(task);
console.log('Result:', result.output);
console.log('Quality:', result.quality);
```

## API Reference

### Collective

```typescript
const collective = new Collective({
  name: string;
  minAgents?: number;      // Default: 2
  maxAgents?: number;      // Default: 100
  consensusThreshold?: number; // Default: 0.67
  reputationDecay?: number;    // Default: 0.01
});

// Methods
collective.join(agent);          // Add agent to collective
collective.leave(sid);           // Remove agent
collective.getAgent(sid);        // Get agent by SID
collective.getAgents();          // Get all agents
collective.size();               // Number of agents
collective.start();              // Start collective
collective.stop();               // Stop collective
collective.submit(task);         // Submit task and wait
collective.submitAsync(task);    // Submit without waiting
collective.createTask(desc, caps); // Create a task
collective.stats();              // Get statistics
```

### Agent

```typescript
const agent = new Agent({
  name: string;
  capabilities: CapabilityType[];
  model?: string;
  parentSid?: string;
});

// Properties
agent.identity;      // SquaremindIdentity
agent.capabilities;  // Map<CapabilityType, Capability>
agent.reputation;    // Reputation
agent.state;         // AgentState

// Methods
agent.start();                    // Start agent
agent.stop();                     // Stop agent
agent.pause();                    // Pause agent
agent.resume();                   // Resume agent
agent.submitTask(task);           // Submit a task
agent.hasCapability(type);        // Check capability
agent.getCapabilityScore(types);  // Get match score
```

### Capability Types

```typescript
type CapabilityType =
  | 'code.write'
  | 'code.review'
  | 'code.refactor'
  | 'research'
  | 'analysis'
  | 'security'
  | 'documentation'
  | 'testing'
  | 'architecture';
```

### Events

```typescript
// Collective events
collective.on('agent:joined', (agent) => {});
collective.on('agent:left', (sid) => {});
collective.on('task:submitted', (task) => {});
collective.on('task:completed', (result) => {});
collective.on('task:failed', (result) => {});
collective.on('started', () => {});
collective.on('stopped', () => {});

// Agent events
agent.on('started', () => {});
agent.on('stopped', () => {});
agent.on('task:received', (task) => {});
agent.on('task:completed', (result) => {});
agent.on('task:failed', (result) => {});
```

## Advanced Usage

### Custom LLM Integration

```typescript
// Coming soon: Custom provider support
```

### Consensus Decisions

```typescript
// Coming soon: Consensus API
```

## Learn More

- [Whitepaper](https://github.com/squaremind/squaremind/blob/main/SQUAREMIND_WHITEPAPER.md)
- [Documentation](https://squaremind.xyz/docs)
- [GitHub](https://github.com/squaremind/squaremind)

## License

MIT
