import { EventEmitter } from 'events';
import {
  AgentConfig,
  AgentState,
  Capability,
  CapabilityType,
  Reputation,
  SquaremindIdentity,
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
 * Squaremind Agent - An autonomous AI agent
 */
export class Agent extends EventEmitter {
  public readonly identity: SquaremindIdentity;
  public readonly capabilities: Map<CapabilityType, Capability>;
  public reputation: Reputation;
  public state: AgentState;
  public currentTask: Task | null = null;

  private model: string;
  private taskQueue: Task[] = [];

  constructor(config: AgentConfig) {
    super();

    // Create identity
    this.identity = {
      sid: generateUUID(),
      name: config.name,
      publicKey: this.generatePublicKey(),
      createdAt: new Date(),
      parentSid: config.parentSid,
      generation: config.parentSid ? 1 : 0,
    };

    // Initialize capabilities
    this.capabilities = new Map();
    for (const capType of config.capabilities) {
      this.capabilities.set(capType, {
        type: capType,
        proficiency: 0.5, // Start at 50%
      });
    }

    // Initialize reputation
    this.reputation = {
      overall: 50,
      reliability: 50,
      quality: 50,
      cooperation: 50,
      honesty: 50,
      tasksCompleted: 0,
      tasksFailed: 0,
      lastActive: new Date(),
    };

    this.state = 'initializing';
    this.model = config.model || 'claude-sonnet-4-20250514';
  }

  private generatePublicKey(): string {
    // Simplified - in production, use proper Ed25519
    const bytes = new Uint8Array(32);
    if (typeof crypto !== 'undefined' && crypto.getRandomValues) {
      crypto.getRandomValues(bytes);
    } else {
      for (let i = 0; i < 32; i++) {
        bytes[i] = Math.floor(Math.random() * 256);
      }
    }
    return Array.from(bytes)
      .map((b) => b.toString(16).padStart(2, '0'))
      .join('');
  }

  /**
   * Start the agent
   */
  async start(): Promise<void> {
    this.state = 'idle';
    this.emit('started');
    this.processQueue();
  }

  /**
   * Stop the agent
   */
  stop(): void {
    this.state = 'terminated';
    this.emit('stopped');
  }

  /**
   * Pause the agent
   */
  pause(): void {
    if (this.state === 'idle') {
      this.state = 'paused';
      this.emit('paused');
    }
  }

  /**
   * Resume the agent
   */
  resume(): void {
    if (this.state === 'paused') {
      this.state = 'idle';
      this.emit('resumed');
      this.processQueue();
    }
  }

  /**
   * Submit a task to the agent
   */
  submitTask(task: Task): void {
    this.taskQueue.push(task);
    this.emit('task:received', task);

    if (this.state === 'idle') {
      this.processQueue();
    }
  }

  /**
   * Get capability match score
   */
  getCapabilityScore(required: CapabilityType[]): number {
    if (required.length === 0) return 1.0;

    let totalScore = 0;
    let matched = 0;

    for (const req of required) {
      const cap = this.capabilities.get(req);
      if (cap) {
        totalScore += cap.proficiency;
        matched++;
      }
    }

    if (matched === 0) return 0;

    const coverage = matched / required.length;
    const avgProficiency = totalScore / matched;

    return coverage * avgProficiency;
  }

  /**
   * Check if agent has a capability
   */
  hasCapability(capType: CapabilityType): boolean {
    return this.capabilities.has(capType);
  }

  private async processQueue(): Promise<void> {
    while (this.taskQueue.length > 0 && this.state === 'idle') {
      const task = this.taskQueue.shift();
      if (task) {
        await this.executeTask(task);
      }
    }
  }

  private async executeTask(task: Task): Promise<TaskResult> {
    this.state = 'working';
    this.currentTask = task;
    this.reputation.lastActive = new Date();

    const startTime = Date.now();

    try {
      // Simulate task execution
      // In real implementation, this would call the LLM
      const output = await this.performTask(task);

      const result: TaskResult = {
        taskId: task.id,
        agentSid: this.identity.sid,
        status: 'completed',
        output,
        quality: 0.8,
        duration: Date.now() - startTime,
        timestamp: new Date(),
      };

      this.recordSuccess(result.quality);
      this.emit('task:completed', result);

      this.state = 'idle';
      this.currentTask = null;

      return result;
    } catch (error) {
      const result: TaskResult = {
        taskId: task.id,
        agentSid: this.identity.sid,
        status: 'failed',
        output: '',
        error: error instanceof Error ? error.message : 'Unknown error',
        quality: 0,
        duration: Date.now() - startTime,
        timestamp: new Date(),
      };

      this.recordFailure();
      this.emit('task:failed', result);

      this.state = 'idle';
      this.currentTask = null;

      return result;
    }
  }

  private async performTask(task: Task): Promise<string> {
    // Simulated task execution
    // In real implementation, this would use the LLM provider
    return `[Simulated] Completed: ${task.description}`;
  }

  private recordSuccess(quality: number): void {
    this.reputation.tasksCompleted++;
    this.reputation.quality =
      this.reputation.quality * 0.9 + quality * 100 * 0.1;
    this.reputation.reliability =
      this.reputation.reliability * 0.95 + 100 * 0.05;
    this.recalculateOverall();
  }

  private recordFailure(): void {
    this.reputation.tasksFailed++;
    this.reputation.reliability *= 0.9;
    this.recalculateOverall();
  }

  private recalculateOverall(): void {
    this.reputation.overall =
      (this.reputation.reliability +
        this.reputation.quality +
        this.reputation.cooperation +
        this.reputation.honesty) /
      4;
  }
}
