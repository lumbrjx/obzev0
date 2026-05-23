import { Experiment } from '../types'

interface Props {
  exp: Experiment
  onDelete: (ns: string, name: string) => void
}

export function ExperimentCard({ exp, onDelete }: Props) {
  return (
    <div style={styles.card}>
      <div style={styles.header}>
        <span style={styles.name}>{exp.name}</span>
        <span style={styles.ns}>{exp.namespace}</span>
        <button style={styles.stopBtn} onClick={() => onDelete(exp.namespace, exp.name)}>
          Stop
        </button>
      </div>
      <div style={styles.meta}>
        {exp.schedule && <Tag label="cron" value={exp.schedule} />}
        {exp.stepCount ? <Tag label="steps" value={String(exp.stepCount)} /> : null}
        <Tag label="rollback" value={exp.hasRollback ? 'yes' : 'no'} />
        {exp.blastRadius?.percentage ? (
          <Tag label="blast%" value={`${exp.blastRadius.percentage}%`} />
        ) : null}
      </div>
    </div>
  )
}

function Tag({ label, value }: { label: string; value: string }) {
  return (
    <span style={styles.tag}>
      <span style={styles.tagLabel}>{label}</span>
      <span style={styles.tagValue}>{value}</span>
    </span>
  )
}

const styles: Record<string, React.CSSProperties> = {
  card: {
    background: '#1e1e2e',
    border: '1px solid #313244',
    borderRadius: 8,
    padding: '12px 16px',
    marginBottom: 10,
  },
  header: { display: 'flex', alignItems: 'center', gap: 12, marginBottom: 8 },
  name: { fontWeight: 700, fontSize: 15, color: '#cdd6f4' },
  ns: { fontSize: 12, color: '#6c7086', flex: 1 },
  stopBtn: {
    background: '#f38ba8',
    color: '#1e1e2e',
    border: 'none',
    borderRadius: 4,
    padding: '3px 10px',
    cursor: 'pointer',
    fontWeight: 600,
    fontSize: 12,
  },
  meta: { display: 'flex', flexWrap: 'wrap', gap: 6 },
  tag: {
    display: 'inline-flex',
    border: '1px solid #45475a',
    borderRadius: 4,
    overflow: 'hidden',
    fontSize: 11,
  },
  tagLabel: { background: '#313244', color: '#a6adc8', padding: '2px 6px' },
  tagValue: { color: '#cdd6f4', padding: '2px 6px' },
}
