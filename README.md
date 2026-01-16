<p align="center">
  <img src="assets/logo-wire-512.svg" width="140" alt="Squaremind">
</p>

<h1 align="center">Squaremind</h1>

<p align="center">
  <strong>Many Agents. One Mind.</strong><br>
  <sub>Multi-agent orchestration protocol for autonomous AI collectives</sub>
</p>

<p align="center">
  <a href="#quick-start">Quick Start</a> •
  <a href="#features">Features</a> •
  <a href="#architecture">Architecture</a> •
  <a href="#documentation">Docs</a> •
  <a href="https://squaremind.xyz">Website</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/go-1.21+-00ADD8?style=flat&logo=go" alt="Go 1.21+">
  <img src="https://img.shields.io/badge/license-MIT-22C55E?style=flat" alt="MIT License">
  <img src="https://img.shields.io/github/stars/square-mind/squaremind?style=flat&color=22C55E" alt="Stars">
  <img src="https://img.shields.io/github/actions/workflow/status/square-mind/squaremind/build.yml?branch=main&style=flat" alt="Build">
</p>

---

## What is Squaremind?

Squaremind is a **multi-agent orchestration protocol** that enables truly autonomous AI collectives. Unlike traditional agent frameworks with rigid hierarchies, Squaremind creates self-organizing swarms that coordinate through:

- **Cryptographic Identity** — Ed25519 keypairs for every agent
- **Reputation Economy** — Trust built through performance
- **Fair Markets** — Transparent task allocation
- **BFT Consensus** — Decentralized decision making

The result: emergent collective intelligence. Many agents, one mind.

## Quick Start

### Install

```bash
# Clone and build
git clone https://github.com/square-mind/squaremind.git
cd squaremind
make build

# Or use Go directly
go install github.com/square-mind/squaremind/cmd/sqm@latest
```

### Configure API Key (BYOK)

Squaremind uses **Bring Your Own Key** — you provide your own LLM API key.

**Claude is the primary/default LLM** (recommended). OpenAI is supported as a fallback.

```bash
# Option 1: Environment variable (recommended)
export ANTHROPIC_API_KEY=your-key-here

# Option 2: Config file (~/.squaremind/config.yaml)
mkdir -p ~/.squaremind
echo "anthropic_api_key: your-key-here" > ~/.squaremind/config.yaml

# Option 3: CLI flag (per command)
sqm demo --api-key your-key-here
```

**Verify your setup:**

```bash
sqm demo
```

This runs an interactive demo that:
- Creates a collective of AI agents
- Assigns cryptographic identities
- Demonstrates fair task markets
- Shows real-time coordination

### Swarm Intelligence

Try the swarm command for complex multi-agent tasks:

```bash
sqm swarm "Design a microservices architecture for an e-commerce platform"
```

This spawns 5 specialized agents (Researcher, Architect, Implementer, Critic, Synthesizer) that work together through parallel coordination, demonstrating emergent collective intelligence.

### Create a Collective

```bash
# Initialize
sqm init MySwarm --max-agents 10

# Spawn agents
sqm spawn Coder --capabilities code.write,code.review
sqm spawn Reviewer --capabilities code.review,security

# Submit a task
sqm task submit "Implement user authentication" --requires code.write

# Run
sqm run
```

### Use in Go

```go
package main

import (
    "github.com/square-mind/squaremind/pkg/agent"
    "github.com/square-mind/squaremind/pkg/collective"
    "github.com/square-mind/squaremind/pkg/identity"
)

func main() {
    // Create collective
    c := collective.NewCollective("DevSwarm", collective.DefaultCollectiveConfig())

    // Spawn agent
    coder, _ := agent.NewAgent(agent.AgentConfig{
        Name:         "Coder",
        Capabilities: []identity.CapabilityType{identity.CapCodeWrite},
    })

    // Join and run
    c.Join(coder)
    c.Start(ctx)

    // Submit work
    task := agent.NewTask("Build REST API", []identity.CapabilityType{identity.CapCodeWrite})
    c.Submit(task)
}
```

### Use in TypeScript

```typescript
import { Collective, Agent } from '@squaremind/sdk';

const collective = new Collective({ name: 'DevSwarm' });

const coder = new Agent({
  name: 'Coder',
  capabilities: ['code.write', 'code.review'],
});

collective.join(coder);
await collective.start();

const result = await collective.submit({
  description: 'Build REST API',
  requires: ['code.write'],
});
```

## Features

| Feature | Description |
|---------|-------------|
| **Agent Identity** | Ed25519 cryptographic identity for every agent |
| **Capabilities** | Declarative skills matching (`code.write`, `security`, etc.) |
| **Reputation** | Multi-dimensional trust scores staked on tasks |
| **Task Markets** | Fair bidding with transparent matching |
| **Gossip Protocol** | Epidemic message propagation |
| **BFT Consensus** | Byzantine fault-tolerant collective decisions |
| **Shared Memory** | Distributed knowledge across the swarm |

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                 COLLECTIVE MIND SUBSTRATE                   │
│   Shared Memory  •  Knowledge Graph  •  Swarm Patterns      │
├─────────────────────────────────────────────────────────────┤
│                 FAIR COORDINATION LAYER                     │
│   Gossip Protocol  •  Task Market  •  PBFT Consensus        │
├─────────────────────────────────────────────────────────────┤
│                 AGENT IDENTITY PROTOCOL                     │
│   Ed25519 Identity  •  Capabilities  •  Reputation          │
└─────────────────────────────────────────────────────────────┘
```

## Project Structure

```
squaremind/
├── cmd/sqm/           # CLI
├── pkg/
│   ├── identity/      # Cryptographic identity
│   ├── agent/         # Agent runtime
│   ├── coordination/  # Gossip, market, consensus
│   ├── collective/    # Collective management
│   └── llm/           # LLM integrations
├── sdk/               # TypeScript SDK
└── docs/              # Documentation
```

## Documentation

| Doc | Description |
|-----|-------------|
| [Architecture](docs/ARCHITECTURE.md) | Technical deep-dive |
| [Quick Start](docs/QUICKSTART.md) | Getting started |
| [Glossary](docs/GLOSSARY.md) | Terms & concepts |
| [API Reference](docs/api.md) | Go & TypeScript APIs |

## $MIND Token

$MIND is the native token of the Squaremind protocol on Solana.

```
CA: 7X6QxLvNddxsSTt6Ai9c5sMrukEqVpJq2P9J6ujvBAGS
```

| | |
|---|---|
| **Buy** | [BAGS.FM](https://bags.fm/token/7X6QxLvNddxsSTt6Ai9c5sMrukEqVpJq2P9J6ujvBAGS) |
| **Chart** | [DexScreener](https://dexscreener.com/solana/7X6QxLvNddxsSTt6Ai9c5sMrukEqVpJq2P9J6ujvBAGS) |
| **Explorer** | [Solscan](https://solscan.io/token/7X6QxLvNddxsSTt6Ai9c5sMrukEqVpJq2P9J6ujvBAGS) |

## Community

- **Website**: [squaremind.xyz](https://squaremind.xyz)
- **X**: [@squaremindai](https://x.com/squaremindai)
- **GitHub**: [square-mind/squaremind](https://github.com/square-mind/squaremind)

## Contributing

```bash
# Fork, clone, branch
git checkout -b feature/your-feature

# Make changes, test
make test

# Submit PR
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

MIT — see [LICENSE](LICENSE)

---

<p align="center">
  <strong>Many Agents. One Mind.</strong>
</p>
