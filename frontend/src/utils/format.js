export function formatNumber(n, decimals = 1) {
  if (!n) return '0'
  if (n >= 1000) return (n / 1000).toFixed(decimals).replace(/\.0$/, '') + 'k'
  return String(n)
}
