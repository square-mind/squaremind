/**
 * Squaremind SDK Type Definitions
 */

// Capability types
export type CapabilityType =
  | 'code.write'
  | 'code.review'
  | 'code.refactor'
  | 'research'
  | 'analysis'
  | 'security'
  | 'documentation'
  | 'testing'
  | 'architecture';

// Agent state
export type AgentState =
  | 'initializing'
  | 'idle'
  | 'working'
  | 'paused'
  | 'terminated';

// Task status
export type TaskStatus =
  | 'pending'
  | 'assigned'
  | 'running'
  | 'completed'
  | 'failed';

/**
 * Capability with proficiency
 */
export interface Capability {
  type: CapabilityType;
  proficiency: number; // 0.0 - 1.0
  metadata?: Record<string, unknown>;
}

/**
 * Squaremind Identity
 */
export interface SquaremindIdentity {
  sid: string;
  name: string;
  publicKey: string;
  createdAt: Date;
  parentSid?: string;
  generation: number;
}

/**
 * Agent reputation
 */
export interface Reputation {
  overall: number;
  reliability: number;
  quality: number;
  cooperation: number;
  honesty: number;
  tasksCompleted: number;
  tasksFailed: number;
  lastActive: Date;
}

/**
 * Task definition
 */
export interface Task {
  id: string;
  description: string;
  requirements?: string;
  complexity: 'low' | 'medium' | 'high';
  required: CapabilityType[];
  deadline?: Date;
  reward: number;
  status: TaskStatus;
  assignedTo?: string;
  createdAt: Date;
}

/**
 * Task result
 */
export interface TaskResult {
  taskId: string;
  agentSid: string;
  status: TaskStatus;
  output: string;
  error?: string;
  quality: number;
  duration: number; // milliseconds
  timestamp: Date;
}

/**
 * Agent configuration
 */
export interface AgentConfig {
  name: string;
  capabilities: CapabilityType[];
  model?: string;
  parentSid?: string;
}

/**
 * Collective configuration
 */
export interface CollectiveConfig {
  name: string;
  minAgents?: number;
  maxAgents?: number;
  consensusThreshold?: number;
  reputationDecay?: number;
}

/**
 * Collective statistics
 */
export interface CollectiveStats {
  name: string;
  agentCount: number;
  activeTasks: number;
  completedTasks: number;
  pendingTasks: number;
  avgReputation: number;
}

/**
 * LLM Provider configuration
 */
export interface ProviderConfig {
  apiKey: string;
  baseUrl?: string;
  model?: string;
}

/**
 * Message for gossip protocol
 */
export interface Message {
  id: string;
  type: string;
  from: string;
  payload: unknown;
  timestamp: Date;
  ttl: number;
}

/**
 * Bid on a task
 */
export interface Bid {
  agentSid: string;
  taskId: string;
  capabilityScore: number;
  reputationStake: number;
  estimatedTime: number; // milliseconds
  timestamp: Date;
}

/**
 * Event types for event emitter
 */
export interface SquaremindEvents {
  'agent:spawned': (agent: { sid: string; name: string }) => void;
  'agent:terminated': (sid: string) => void;
  'task:submitted': (task: Task) => void;
  'task:completed': (result: TaskResult) => void;
  'task:failed': (result: TaskResult) => void;
  'collective:started': () => void;
  'collective:stopped': () => void;
}
