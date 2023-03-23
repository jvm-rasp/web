import request from '@/utils/request'

export function getAttackData() {
  return request({
    url: '/dashboard/attackData',
    method: 'get'
  })
}

export async function getTrendsData() {
  const res = await request({
    url: '/dashboard/attackTrends',
    method: 'get'
  })
  const option = {
    title: {
      text: ''
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
        label: {
          backgroundColor: '#6a7985'
        }
      }
    },
    legend: {
      data: ['拦截请求', '记录日志']
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: [
      {
        type: 'category',
        boundaryGap: false,
        data: ['2022-02-01', '2022-02-02', '2022-02-03', '2022-02-04', '2022-02-05', '2022-02-06', '2022-02-07']
      }
    ],
    yAxis: [
      {
        type: 'value'
      }
    ],
    series: [
      {
        name: '拦截请求',
        type: 'line',
        lineStyle: {
          color: '#FF0000'
        },
        itemStyle: {
          color: '#FF0000'
        },
        areaStyle: {
          color: '#FF0000',
          opacity: 0.3
        },
        data: [60, 132, 101, 134, 90, 230, 210]
      },
      {
        name: '记录日志',
        type: 'line',
        lineStyle: {
          color: '#FFCC00'
        },
        itemStyle: {
          color: '#FFCC00'
        },
        areaStyle: {
          color: '#FFCC00',
          opacity: 0.3
        },
        data: [80, 182, 191, 234, 290, 330, 310]
      }
    ]
  }
  option.xAxis[0].data = res.data.date
  option.series[0].data = res.data.block
  option.series[1].data = res.data.warn
  return option
}

export async function getAttackTypeData() {
  const res = await request({
    url: '/dashboard/attackTypes',
    method: 'get'
  })
  const option = {
    tooltip: {
      trigger: 'item'
    },
    legend: {
      top: '5%',
      left: 'center'
    },
    series: [
      {
        name: '漏洞类型',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 20,
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: [
          { value: 1048, name: 'Search Engine' },
          { value: 735, name: 'Direct' },
          { value: 580, name: 'Email' },
          { value: 484, name: 'Union Ads' },
          { value: 300, name: 'Video Ads' }
        ]
      }
    ]
  }
  option.series[0].data = res.data.list
  return option
}
