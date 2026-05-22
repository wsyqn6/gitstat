import * as echarts from 'echarts/core'
import { LineChart, BarChart, HeatmapChart, RadarChart } from 'echarts/charts'
import {
  TooltipComponent,
  LegendComponent,
  GridComponent,
  RadarComponent,
  VisualMapComponent
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

echarts.use([
  LineChart,
  BarChart,
  HeatmapChart,
  RadarChart,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  RadarComponent,
  VisualMapComponent,
  CanvasRenderer
])

export default echarts
