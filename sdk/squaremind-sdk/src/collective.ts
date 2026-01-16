import { EventEmitter } from 'events';
import { Agent } from './agent';
import {
  CollectiveConfig,
  CollectiveStats,
  CapabilityType,
  Task,
  TaskResult,
  TaskStatus,
} from './types';

/**
 * Generate a UUID v4
 */
function generateUUID(): string {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    const v = c === 'x' ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

/**
 * Collective - A group of squaremind agents
 */
export class Collective extends EventEmitter {
  public readonly id: string;
  public readonly name: string;

  private agents: Map<string, Agent> = new Map();
  private config: Required<CollectiveConfig>;
  private pendingTasks: Task[] = [];
  private activeTasks: Map<string, Task> = new Map();
  private completedTasks: TaskResult[] = [];
  private running: boolean = false;

  constructor(config: CollectiveConfig) {
    super();

    this.id = generateUUID();
    this.name = config.name;
    this.config = {
      name: config.name,
      minAgents: config.minAgents ?? 2,
      maxAgents: config.maxAgents ?? 100,
      consensusThreshold: config.consensusThreshold ?? 0.67,
      reputationDecay: config.reputationDecay ?? 0.01,
    };
  }

  /**
   * Add an agent to the collective
   */
  join(agent: Agent): void {
    if (this.agents.size >= this.config.maxAgents) {
      throw new Error('Collective at maximum capacity');
    }

    this.agents.set(agent.identity.sid, agent);

    // Forward agent events
    agent.on('task:completed', (result: TaskResult) => {
      this.onTaskCompleted(result);
    });

    agent.on('task:failed', (result: TaskResult) => {
      this.onTaskFailed(result);
    });

    this.emit('agent:joined', agent);
  }

  /**
   * Remove an agent from the collective
   */
  leave(sid: string): void {
    const agent = this.agents.get(sid);
    if (!agent) {
      throw new Error('Agent not found');
    }

    agent.stop();
    this.agents.delete(sid);
    this.emit('agent:left', sid);
  }

  /**
   * Get an agent by SID
   */
  getAgent(sid: string): Agent | undefined {
    return this.agents.get(sid);
  }

  /**
   * Get all agents
   */
  getAgents(): Agent[] {
    return Array.from(this.agents.values());
  }

  /**
   * Get collective size
   */
  size(): number {
    return this.agents.size;
  }

  /**
   * Start the collective
   */
  async start(): Promise<void> {
    this.running = true;

    // Start all agents
    for (const agent of this.agents.values()) {
      await agent.start();
    }

    this.emit('started');
  }

  /**
   * Stop the collective
   */
  stop(): void {
    this.running = false;

    // Stop all agents
    for (const agent of this.agents.values()) {
      agent.stop();
    }

    this.emit('stopped');
  }

  /**
   * Submit a task to the collective
   */
  async submit(task: Task): Promise<TaskResult> {
    this.pendingTasks.push(task);
    this.emit('task:submitted', task);

    // Find best agent for the task
    const assignment = this.assignTask(task);
    if (!assignment) {
      throw new Error('No suitable agent found');
    }

    const agent = this.agents.get(assignment.agentSid);
    if (!agent) {
      throw new Error('Assigned agent not found');
    }

    // Move to active
    this.activeTasks.set(task.id, task);
    task.status = 'assigned';
    task.assignedTo = assignment.agentSid;

    // Submit to agent and wait for result
    return new Promise((resolve, reject) => {
      const timeout = setTimeout(() => {
        reject(new Error('Task timeout'));
      }, 60000); // 1 minute timeout

      agent.once('task:completed', (result: TaskResult) => {
        if (result.taskId === task.id) {
          clearTimeout(timeout);
          resolve(result);
        }
      });

      agent.once('task:failed', (result: TaskResult) => {
        if (result.taskId === task.id) {
          clearTimeout(timeout);
          reject(new Error(result.error || 'Task failed'));
        }
      });

      agent.submitTask(task);
    });
  }

  /**
   * Submit a task asynchronously (fire and forget)
   */
  submitAsync(task: Task): string {
    this.submit(task).catch((err) => {
      this.emit('task:error', { taskId: task.id, error: err });
    });
    return task.id;
  }

  /**
   * Create a new task
   */
  createTask(
    description: string,
    required: CapabilityType[] = []
  ): Task {
    return {
      id: generateUUID(),
      description,
      complexity: 'medium',
      required,
      reward: 10,
      status: 'pending',
      createdAt: new Date(),
    };
  }

  /**
   * Get collective statistics
   */
  stats(): CollectiveStats {
    let totalReputation = 0;
    for (const agent of this.agents.values()) {
      totalReputation += agent.reputation.overall;
    }

    return {
      name: this.name,
      agentCount: this.agents.size,
      activeTasks: this.activeTasks.size,
      completedTasks: this.completedTasks.length,
      pendingTasks: this.pendingTasks.length,
      avgReputation:
        this.agents.size > 0 ? totalReputation / this.agents.size : 0,
    };
  }

  private assignTask(task: Task): { agentSid: string; score: number } | null {
    let bestAgent: string | null = null;
    let bestScore = 0;

    for (const agent of this.agents.values()) {
      if (agent.state !== 'idle') continue;

      const capScore = agent.getCapabilityScore(task.required);
      if (capScore > 0.5) {
        // Combine capability and reputation
        const totalScore =
          capScore * 0.6 + (agent.reputation.overall / 100) * 0.4;

        if (totalScore > bestScore) {
          bestScore = totalScore;
          bestAgent = agent.identity.sid;
        }
      }
    }

    return bestAgent ? { agentSid: bestAgent, score: bestScore } : null;
  }

  private onTaskCompleted(result: TaskResult): void {
    this.activeTasks.delete(result.taskId);
    this.completedTasks.push(result);
    this.emit('task:completed', result);
  }

  private onTaskFailed(result: TaskResult): void {
    this.activeTasks.delete(result.taskId);
    this.completedTasks.push(result);
    this.emit('task:failed', result);
  }
}
