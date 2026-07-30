export function cellLevel(commits, maxC) {
  if (!commits || !maxC) return 0
  const ratio = commits / maxC
  if (ratio <= 0.25) return 1
  if (ratio <= 0.5) return 2
  if (ratio <= 0.75) return 3
  return 4
}
