const BASE = import.meta.env.VITE_API_BASE || '/api'

function buildParams(opt = {}) {
  const { startDate, endDate, repos, timeRange, path, email, offset, limit, prevStartDate, prevEndDate } = opt
  const p = new URLSearchParams()
  if (startDate) p.append('startDate', startDate)
  if (endDate) p.append('endDate', endDate)
  if (prevStartDate) p.append('prevStartDate', prevStartDate)
  if (prevEndDate) p.append('prevEndDate', prevEndDate)
  if (timeRange) p.append('range', timeRange)
  if (email) p.append('email', email)
  if (path) p.append('path', path)
  if (offset != null) p.append('offset', offset)
  if (limit != null) p.append('limit', limit)
  if (repos?.length && !repos.includes('all'))
    repos.forEach(r => p.append('repo', r))
  return p
}

async function request(url, opts = {}) {
  const { method = 'GET', headers, body, errorMsg, blob, signal } = opts
  const reqBody = body && typeof body === 'object' ? JSON.stringify(body) : body
  const reqHeaders = body && typeof body === 'object'
    ? { 'Content-Type': 'application/json', ...headers }
    : headers
  const res = await fetch(url, { method, headers: reqHeaders, body: reqBody, signal })
  if (!res.ok) throw new Error(errorMsg || `HTTP ${res.status}`)
  return blob ? res.blob() : res.json()
}

function makeAborter() {
  const ctrl = new AbortController()
  return { signal: ctrl.signal, abort: () => ctrl.abort() }
}

export { BASE, buildParams, request, makeAborter }
