# Squaremind Quick Start Guide

Get up and running with Squaremind in under 5 minutes.

## Prerequisites

- Go 1.21 or later
- (Optional) Anthropic API key for Claude integration

## Installation

### Option 1: Using Go Install

```bash
go install github.com/square-mind/squaremind/cmd/sqm@latest
```

### Option 2: From Source

```bash
git clone https://github.com/square-mind/squaremind.git
cd squaremind
make build
make install
```

### Verify Installation

```bash
sqm --version
# squaremind v0.1.0
```

## Your First Collective

### Step 1: Initialize

Create a new collective with default settings:

```bash
sqm init MyFirstSwarm
```

Output:
```
  Collective 'MyFirstSwarm' initialized

  ID: a1b2c3d4-...
  Max Agents: 100
  Consensus Threshold: 67%
```

### Step 2: Spawn Agents

Create agents with different capabilities:

```bash
# A coding agent
sqm spawn Coder --capabilities code.write,code.review

# A security-focused agent
sqm spawn SecurityBot --capabilities security,code.review

# A documentation agent
sqm spawn DocWriter --capabilities documentation,analysis
```

Each spawn outputs:
```
  Squaremind agent 'Coder' spawned

  SID: e5f6g7h8-...
  Public Key: 3a4b5c6d...
  Capabilities: [code.write code.review]
  Model: claude-sonnet-4-20250514
  Reputation: 50.0
```

### Step 3: Check Status

View your collective:

```bash
sqm status
```

Output:
```
  Collective Status
  ----------------------------------------
  Name: MyFirstSwarm
  Agents: 3
  Tasks Pending: 0
  Tasks Active: 0
  Tasks Completed: 0
  Avg Reputation: 50.0

  Agents:
    - Coder (e5f6g7h8) [idle] rep=50.0
      capabilities: code.write, code.review
    - SecurityBot (i9j0k1l2) [idle] rep=50.0
      capabilities: security, code.review
    - DocWriter (m3n4o5p6) [idle] rep=50.0
      capabilities: documentation, analysis
```

### Step 4: Submit a Task

Send a task to the collective:

```bash
sqm task submit "Write a function to validate email addresses" \
  --requires code.write \
  --complexity medium
```

Output:
```
  Submitting task: Write a function to validate email addresses
  Task ID: q7r8s9t0-...
  Complexity: medium
  Required capabilities: [code.write]

  Task completed!
  Status: completed
  Quality: 0.85
  Duration: 2.3s
  Output: [task output here]
```

### Step 5: Run the Collective

Start the collective for continuous operation:

```bash
sqm run
```

Output:
```
  Starting Squaremind collective...
  Name: MyFirstSwarm
  Agents: 3
  Press Ctrl+C to stop
```

## Using with Claude

To use real LLM capabilities, configure your API key:

```bash
# Set Anthropic API key
sqm config set api-key YOUR_ANTHROPIC_API_KEY

# Or for OpenAI
sqm config set openai-key YOUR_OPENAI_API_KEY
```

Now agents will use the LLM to complete tasks.

## Programmatic Usage

### Go

```go
package main

import (
    "context"
    "fmt"

    "github.com/square-mind/squaremind/pkg/agent"
    "github.com/square-mind/squaremind/pkg/collective"
    "github.com/square-mind/squaremind/pkg/identity"
)

func main() {
    ctx := context.Background()

    // Create collective
    c := collective.NewCollective("GoSwarm", collective.CollectiveConfig{
        MaxAgents:          10,
        ConsensusThreshold: 0.67,
    })

    // Spawn agent
    coder, _ := agent.NewAgent(agent.AgentConfig{
        Name:         "GoCoder",
        Capabilities: []identity.CapabilityType{identity.CapCodeWrite},
    })

    // Join and start
    c.Join(coder)
    c.Start(ctx)

    // Submit task
    task := agent.NewTask(
        "Implement a binary search function",
        []identity.CapabilityType{identity.CapCodeWrite},
    )

    result, err := c.Submit(task)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Output: %s\n", result.Output)
    fmt.Printf("Quality: %.2f\n", result.Quality)

    c.Stop()
}
```

### TypeScript

```typescript
import { Collective, Agent } from '@squaremind/sdk';

async function main() {
  // Create collective
  const collective = new Collective({
    name: 'TSSwarm',
    maxAgents: 10,
  });

  // Spawn agent
  const coder = new Agent({
    name: 'TSCoder',
    capabilities: ['code.write', 'code.review'],
  });

  // Join and start
  collective.join(coder);
  await collective.start();

  // Submit task
  const task = collective.createTask(
    'Implement a debounce function',
    ['code.write']
  );

  const result = await collective.submit(task);
  console.log('Output:', result.output);
  console.log('Quality:', result.quality);

  collective.stop();
}

main();
```

## CLI Reference

| Command | Description |
|---------|-------------|
| `sqm init <name>` | Initialize a new collective |
| `sqm spawn <name>` | Spawn a new agent |
| `sqm run` | Start the collective |
| `sqm status` | Show collective status |
| `sqm task submit <desc>` | Submit a task |
| `sqm agent list` | List all agents |
| `sqm agent stop <sid>` | Stop an agent |
| `sqm config set <key> <val>` | Set configuration |

## Common Options

### sqm init

| Flag | Description | Default |
|------|-------------|---------|
| `--max-agents, -m` | Maximum agents | 100 |
| `--threshold, -t` | Consensus threshold | 0.67 |

### sqm spawn

| Flag | Description | Default |
|------|-------------|---------|
| `--capabilities, -c` | Agent capabilities | code.write |
| `--model, -m` | LLM model | claude-sonnet-4-20250514 |

### sqm task submit

| Flag | Description | Default |
|------|-------------|---------|
| `--requires, -r` | Required capabilities | [] |
| `--complexity, -x` | Task complexity | medium |
| `--reward, -w` | Reputation reward | 10 |
| `--async, -a` | Submit async | false |

## Next Steps

- Read the [Architecture](ARCHITECTURE.md) for technical details
- Explore the [Glossary](GLOSSARY.md) for terminology
- Check out [examples](../examples/) for more use cases
- Read the [API Reference](api.md) for full documentation

## Troubleshooting

### "No suitable agent found"

This means no agent has the required capabilities. Check:
- Agent capabilities match task requirements
- At least one agent is in `idle` state

### "Collective at maximum capacity"

You've reached `maxAgents`. Either:
- Remove idle agents with `sqm agent stop <sid>`
- Increase limit with `--max-agents` during init

### API Key Issues

Ensure your API key is set:
```bash
sqm config set api-key YOUR_KEY
```

Or set environment variable:
```bash
export ANTHROPIC_API_KEY=your_key
```

## Getting Help

- GitHub Issues: [github.com/square-mind/squaremind/issues](https://github.com/square-mind/squaremind/issues)
- X: [@squaremindai](https://x.com/squaremindai)
- Website: [squaremind.xyz](https://squaremind.xyz)
