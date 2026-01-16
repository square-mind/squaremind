/**
 * Squaremind SDK
 *
 * TypeScript SDK for building autonomous AI collectives
 *
 * Many Agents. One Mind.
 *
 * @example
 * ```typescript
 * import { Collective, Agent } from '@squaremind/sdk';
 *
 * // Create a collective
 * const collective = new Collective({ name: 'DevSwarm' });
 *
 * // Spawn agents
 * const coder = new Agent({
 *   name: 'Coder',
 *   capabilities: ['code.write', 'code.review'],
 * });
 *
 * // Join collective
 * collective.join(coder);
 *
 * // Start
 * await collective.start();
 *
 * // Submit tasks
 * const result = await collective.submit(
 *   collective.createTask('Implement user authentication')
 * );
 * ```
 */

export { Agent } from './agent';
export { Collective } from './collective';
export * from './types';

// Version
export const VERSION = '0.1.0';

// Default export for convenience
import { Collective } from './collective';
import { Agent } from './agent';

export default {
  Collective,
  Agent,
  VERSION: '0.1.0',
};
