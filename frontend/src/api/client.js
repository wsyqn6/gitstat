const BASE = import.meta.env.VITE_API_BASE || '/api'

function buildParams(opt = {}) {
  const { startDate, endDate, repos, timeRange, path, email } = opt
  const p = new URLSearchParams()
  if (startDate) p.append('startDate', startDate)
  if (endDate) p.append('endDate', endDate)
  if (timeRange) p.append('range', timeRange)
  if (email) p.append('email', email)
  if (path) p.append('path', path)
  if (repos?.length && !repos.includes('all'))
    repos.forEach(r => p.append('repo', r))
  return p
}

async function request(url, opts = {}) {
  const { method = 'GET', headers, body, errorMsg, blob } = opts
  const res = await fetch(url, { method, headers, body })
  if (!res.ok) throw new Error(errorMsg || `HTTP ${res.status}`)
  return blob ? res.blob() : res.json()
}

export { BASE, buildParams, request }
