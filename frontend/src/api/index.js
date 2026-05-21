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

export async function setScanPath(path) {
  const response = await fetch(`${API_BASE}/scan/path`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ path })
  })
  
  if (!response.ok) {
    throw new Error('Set path failed')
  }
  
  return response.json()
}

export async function getOverviewStats(startDate = null, endDate = null) {
  const params = new URLSearchParams()
  if (startDate) params.append('startDate', startDate)
  if (endDate) params.append('endDate', endDate)
  
  const url = `${API_BASE}/stats/overview?${params.toString()}`
  const response = await fetch(url)
  
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

export async function getDailyStats(email, timeRange = 'week', repos = [], startDate = null, endDate = null) {
  const params = new URLSearchParams()
  if (email) params.append('email', email)
  if (timeRange) params.append('range', timeRange)
  if (startDate) params.append('startDate', startDate)
  if (endDate) params.append('endDate', endDate)
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

export async function getVersion() {
  const response = await fetch(`${API_BASE}/version`)
  
  if (!response.ok) {
    throw new Error('Failed to fetch version')
  }
  
  const data = await response.json()
  return data.version
}
