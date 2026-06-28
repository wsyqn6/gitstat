export function pad(n) {
  return String(n).padStart(2, '0')
}

export function eachDay(start, end) {
  const days = []
  const cur = new Date(start)
  const endDate = new Date(end)
  while (cur <= endDate) {
    days.push(cur.toISOString().slice(0, 10))
    cur.setDate(cur.getDate() + 1)
  }
  return days
}

export function eachMonth(start, end) {
  const months = []
  const cur = new Date(start)
  const endDate = new Date(end)
  while (cur <= endDate) {
    months.push(cur.toISOString().slice(0, 7))
    cur.setMonth(cur.getMonth() + 1)
  }
  return months
}
