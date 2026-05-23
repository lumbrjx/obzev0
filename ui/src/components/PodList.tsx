import { Pod } from '../types'

interface Props {
  pods: Pod[]
}

const phaseColor: Record<string, string> = {
  Running: '#a6e3a1',
  Pending: '#f9e2af',
  Failed: '#f38ba8',
  Succeeded: '#89dceb',
}

export function PodList({ pods }: Props) {
  if (pods.length === 0) {
    return <p style={{ color: '#6c7086', fontSize: 13 }}>No daemon pods found.</p>
  }
  return (
    <table style={styles.table}>
      <thead>
        <tr>
          {['Name', 'Namespace', 'Node', 'Phase', 'IP'].map(h => (
            <th key={h} style={styles.th}>{h}</th>
          ))}
        </tr>
      </thead>
      <tbody>
        {pods.map(pod => (
          <tr key={pod.name} style={styles.tr}>
            <td style={styles.td}>{pod.name}</td>
            <td style={styles.td}>{pod.namespace}</td>
            <td style={styles.td}>{pod.node}</td>
            <td style={styles.td}>
              <span style={{ color: phaseColor[pod.phase] ?? '#cdd6f4' }}>{pod.phase}</span>
            </td>
            <td style={{ ...styles.td, fontFamily: 'monospace' }}>{pod.ip}</td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}

const styles: Record<string, React.CSSProperties> = {
  table: { width: '100%', borderCollapse: 'collapse', fontSize: 13 },
  th: { textAlign: 'left', color: '#6c7086', padding: '6px 10px', borderBottom: '1px solid #313244' },
  tr: { borderBottom: '1px solid #1e1e2e' },
  td: { color: '#cdd6f4', padding: '6px 10px' },
}
