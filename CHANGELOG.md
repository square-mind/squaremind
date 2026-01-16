# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Persistent agent storage
- Web dashboard for collective monitoring
- Plugin system for custom capabilities
- Multi-collective networking

## [0.1.0] - 2026-01-16

### Added

#### Core Protocol
- **Agent Identity Protocol (AIP)**: Ed25519 cryptographic identity for all agents
- **Fair Coordination Layer (FCL)**: Gossip protocol, task markets, BFT consensus
- **Collective Mind Substrate (CMS)**: Shared memory and knowledge graphs

#### Identity & Security
- Ed25519 keypair generation for agent identity (SID)
- Hierarchical capability system with proficiency levels
- Signed actions and delegation proofs
- Consensus proofs for verifiable collective decisions

#### Agent Runtime
- Agent lifecycle management (spawn, start, pause, resume, terminate)
- Task execution with quality scoring
- Multi-dimensional reputation tracking (reliability, quality, cooperation, honesty)
- Agent memory with short-term and long-term storage

#### Coordination
- Epidemic gossip protocol with configurable fanout and TTL
- Task market with capability-based matching
- Reputation-weighted bidding system
- PBFT consensus engine with configurable threshold

#### Collective Intelligence
- Collective management with min/max agent constraints
- Distributed shared memory
- Knowledge graph for collective learning
- Swarm patterns: stigmergy, quorum sensing, chemotaxis, division of labor

#### LLM Integration
- Provider interface for pluggable LLM backends
- Claude (Anthropic) provider
- OpenAI provider
- Simulated provider for testing

#### CLI (`sqm`)
- `sqm init` - Initialize a new collective
- `sqm spawn` - Spawn agents with capabilities
- `sqm run` - Start the collective
- `sqm status` - Display collective status
- `sqm task submit` - Submit tasks to the collective
- `sqm agent list` - List all agents
- `sqm agent stop` - Stop an agent
- `sqm config set` - Configure API keys

#### SDK
- TypeScript SDK (`@squaremind/sdk`)
- Full type definitions
- Event-driven architecture
- Promise-based async API

#### Documentation
- API reference documentation
- Architecture documentation
- Quick start guide
- Example applications

### Architecture

```
┌─────────────────────────────────────────┐
│       Collective Mind Substrate         │
│    (Shared Memory, Knowledge Graph)     │
├─────────────────────────────────────────┤
│       Fair Coordination Layer           │
│   (Gossip, Task Market, Consensus)      │
├─────────────────────────────────────────┤
│       Agent Identity Protocol           │
│  (Ed25519 Identity, Capabilities, Rep)  │
└─────────────────────────────────────────┘
```

### Technical Details
- Written in Go 1.21+
- Zero external runtime dependencies
- Concurrent agent execution with goroutines
- Context-based cancellation for clean shutdown

---

[Unreleased]: https://github.com/square-mind/squaremind/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/square-mind/squaremind/releases/tag/v0.1.0
