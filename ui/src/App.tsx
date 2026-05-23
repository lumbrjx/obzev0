import { useState, useEffect, useCallback } from 'react'
import { Experiment, Pod, Report } from './types'
import { ExperimentCard } from './components/ExperimentCard'
import { PodList } from './components/PodList'
import { ReportList } from './components/ReportList'

type Tab = 'experiments' | 'pods' | 'reports'

const POLL_MS = 5000

export default function App() {
  const [tab, setTab] = useState<Tab>('experiments')
  const [experiments, setExperiments] = useState<Experiment[]>([])
  const [pods, setPods] = useState<Pod[]>([])
  const [reports, setReports] = useState<Report[]>([])
  const [error, setError] = useState<string | null>(null)

  const fetchAll = useCallback(async () => {
    try {
      const [exps, ps, reps] = await Promise.all([
        fetch('/api/experiments').then(r => r.json()),
        fetch('/api/pods').then(r => r.json()),
        fetch('/api/reports').then(r => r.json()),
      ])
      setExperiments(exps)
      setPods(ps)
      setReports(reps)
      setError(null)
    } catch (e) {
      setError(String(e))
    }
  }, [])

  useEffect(() => {
    fetchAll()
    const id = setInterval(fetchAll, POLL_MS)
    return () => clearInterval(id)
  }, [fetchAll])

  const deleteExperiment = async (ns: string, name: string) => {
    await fetch(`/api/experiments/${ns}/${name}`, { method: 'DELETE' })
    fetchAll()
  }

  return (
    <div style={styles.root}>
      <header style={styles.header}>
        <span style={styles.logo}>obzev0</span>
        <span style={styles.subtitle}>chaos engineering dashboard</span>
      </header>

      {error && <div style={styles.error}>{error}</div>}

      <nav style={styles.nav}>
        {(['experiments', 'pods', 'reports'] as Tab[]).map(t => (
          <button
            key={t}
            style={{ ...styles.tab, ...(tab === t ? styles.activeTab : {}) }}
            onClick={() => setTab(t)}
          >
            {t}
            {t === 'experiments' && experiments.length > 0 && (
              <span style={styles.badge}>{experiments.length}</span>
            )}
            {t === 'pods' && pods.length > 0 && (
              <span style={styles.badge}>{pods.length}</span>
            )}
          </button>
        ))}
      </nav>

      <main style={styles.main}>
        {tab === 'experiments' && (
          experiments.length === 0
            ? <p style={styles.empty}>No active experiments.</p>
            : experiments.map(e => (
                <ExperimentCard key={`${e.namespace}/${e.name}`} exp={e} onDelete={deleteExperiment} />
              ))
        )}
        {tab === 'pods' && <PodList pods={pods} />}
        {tab === 'reports' && <ReportList reports={reports} />}
      </main>
    </div>
  )
}

const styles: Record<string, React.CSSProperties> = {
  root: { minHeight: '100vh', background: '#181825', color: '#cdd6f4', fontFamily: 'system-ui, sans-serif' },
  header: { padding: '16px 24px', borderBottom: '1px solid #313244', display: 'flex', alignItems: 'baseline', gap: 12 },
  logo: { fontSize: 22, fontWeight: 800, color: '#89b4fa' },
  subtitle: { fontSize: 13, color: '#6c7086' },
  error: { background: '#f38ba820', color: '#f38ba8', padding: '8px 24px', fontSize: 13 },
  nav: { display: 'flex', gap: 4, padding: '12px 24px 0', borderBottom: '1px solid #313244' },
  tab: {
    background: 'none', border: 'none', color: '#6c7086', cursor: 'pointer',
    padding: '8px 14px', fontSize: 14, borderBottom: '2px solid transparent',
  },
  activeTab: { color: '#89b4fa', borderBottomColor: '#89b4fa' },
  badge: {
    background: '#313244', color: '#a6adc8', borderRadius: 10,
    padding: '1px 6px', fontSize: 11, marginLeft: 6,
  },
  main: { padding: 24 },
  empty: { color: '#6c7086', fontSize: 13 },
}
