export interface BlastRadius {
  namespaces?: string[]
  labelSelector?: Record<string, string>
  percentage?: number
}

export interface Experiment {
  name: string
  namespace: string
  schedule?: string
  stepCount?: number
  hasRollback: boolean
  blastRadius?: BlastRadius
}

export interface Pod {
  name: string
  namespace: string
  node: string
  phase: string
  ip: string
}

export interface Report {
  name: string
  namespace: string
  report: string
}

export interface ParsedReport {
  experimentName: string
  namespace: string
  startTime: string
  endTime: string
  durationSeconds: number
  affectedPods: string[]
  stopReason: string
  stepCount?: number
}
