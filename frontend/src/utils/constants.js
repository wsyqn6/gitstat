export const THEMES = {
  neon: {
    chartColors: ['#00f5ff', '#ff00ff', '#ffd700', '#00ff88', '#ff6b6b', '#a78bfa'],
    chartColors8: ['#00d4ff', '#ff6b9d', '#00ff88', '#ffd700', '#a78bfa', '#f472b6', '#34d399', '#fb923c'],
    primary: '#00d4ff',
    primaryRgb: '0, 212, 255',
    accent: '#00f5ff',
    accentRgb: '0, 245, 255',
    titleColor: '#64748b',
    tooltipText: '#fff',
    axisLabel: '#94a3b8',
    axisLine: '#334155',
    splitLine: '#1e293b',
    tooltipBg: 'rgba(10, 14, 39, 0.95)',
    calLevels: [
      'rgba(0, 212, 255, 0.04)',
      'rgba(0, 212, 255, 0.18)',
      'rgba(0, 212, 255, 0.38)',
      'rgba(0, 140, 255, 0.6)',
      'rgba(0, 100, 255, 0.85)'
    ],
    heatmapColors: ['rgba(0, 245, 255, 0.1)', 'rgba(0, 245, 255, 0.4)', 'rgba(0, 245, 255, 0.7)', '#00f5ff']
  },
  ios26: {
    chartColors: ['#007AFF', '#AF52DE', '#FF9500', '#34C759', '#FF2D55', '#FFCC00'],
    chartColors8: ['#007AFF', '#FF2D55', '#34C759', '#FF9500', '#AF52DE', '#FFCC00', '#5AC8FA', '#FF3B30'],
    primary: '#007AFF',
    primaryRgb: '0, 122, 255',
    accent: '#AF52DE',
    accentRgb: '175, 82, 222',
    titleColor: '#1C1C1E',
    tooltipText: '#1C1C1E',
    axisLabel: '#636366',
    axisLine: '#C7C7CC',
    splitLine: '#E5E5EA',
    tooltipBg: 'rgba(255, 255, 255, 0.95)',
    calLevels: [
      'rgba(0, 122, 255, 0.04)',
      'rgba(0, 122, 255, 0.15)',
      'rgba(0, 122, 255, 0.3)',
      'rgba(0, 100, 220, 0.5)',
      'rgba(0, 80, 200, 0.75)'
    ],
    heatmapColors: ['rgba(0, 122, 255, 0.08)', 'rgba(0, 122, 255, 0.3)', 'rgba(0, 122, 255, 0.6)', '#007AFF']
  }
}

export function getChartConfig(theme) {
  return THEMES[theme] || THEMES.neon
}

export const CHART_COLORS = THEMES.neon.chartColors8
