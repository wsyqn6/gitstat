const API_BASE = '/api'

export async function scanRepositories(path, timeRange) {
  const response = await fetch(`${API_BASE}/scan`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ path, timeRange })
  })
  
  if (!response.ok) {
    throw new Error('Scan failed')
  }
  
  return response.json()
}

export async function getOverviewStats() {
  const response = await fetch(`${API_BASE}/stats/overview`)
  
  if (!response.ok) {
    throw new Error('Failed to fetch stats')
  }
  
  return response.json()
}

export async function getRepositories() {
  const response = await fetch(`${API_BASE}/repositories`)
  
  if (!response.ok) {
    throw new Error('Failed to fetch repositories')
  }
  
  return response.json()
}

export async function exportData() {
  const response = await fetch(`${API_BASE}/export/json`)
  
  if (!response.ok) {
    throw new Error('Export failed')
  }
  
  return response.blob()
}

export async function getDailyStats(email, timeRange = '7', repos = []) {
  const params = new URLSearchParams()
  if (email) params.append('email', email)
  if (timeRange) params.append('range', timeRange)
  if (repos.length > 0) {
    repos.forEach(repo => params.append('repo', repo))
  }
  
  const url = `${API_BASE}/stats/daily?${params.toString()}`
  const response = await fetch(url)
  
  if (!response.ok) {
    throw new Error('Failed to fetch daily stats')
  }
  
  return response.json()
}
