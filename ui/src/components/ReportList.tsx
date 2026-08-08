import { Report, ParsedReport } from '../types'

interface Props {
  reports: Report[]
}

export function ReportList({ reports }: Props) {
  if (reports.length === 0) {
    return <p style={{ color: '#6c7086', fontSize: 13 }}>No reports yet.</p>
  }
  return (
    <div>
      {reports.map(r => {
        let parsed: ParsedReport | null = null
        try { parsed = JSON.parse(r.report) } catch { /* malformed */ }
        return (
          <div key={r.name} style={styles.card}>
            <div style={styles.title}>{r.name}</div>
            {parsed ? (
              <div style={styles.grid}>
                <Kv k="experiment" v={parsed.experimentName} />
                <Kv k="stop reason" v={parsed.stopReason} />
                <Kv k="duration" v={`${parsed.durationSeconds.toFixed(1)}s`} />
                <Kv k="pods" v={parsed.affectedPods?.join(', ') ?? '—'} />
                {parsed.stepCount ? <Kv k="steps" v={String(parsed.stepCount)} /> : null}
              </div>
            ) : (
              <pre style={styles.raw}>{r.report}</pre>
            )}
          </div>
        )
      })}
    </div>
  )
}

function Kv({ k, v }: { k: string; v: string }) {
  return (
    <div style={{ fontSize: 12 }}>
      <span style={{ color: '#6c7086' }}>{k}: </span>
      <span style={{ color: '#cdd6f4' }}>{v}</span>
    </div>
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
  title: { fontWeight: 600, color: '#89b4fa', marginBottom: 8, fontSize: 13 },
  grid: { display: 'flex', flexDirection: 'column', gap: 3 },
  raw: { color: '#a6adc8', fontSize: 11, overflow: 'auto', margin: 0 },
}
