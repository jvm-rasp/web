<template>
  <div id="root" class="dashboard-container">
    <el-row :gutter="20">
      <el-col :span="24">
        <el-card class="box-card">
          <el-descriptions title="攻击数据" :extra="currentDatetime" />
          <el-col :span="8">
            <el-statistic group-separator="," :value="attackData.total" title="总攻击次数" />
          </el-col>
          <el-col :span="8">
            <el-statistic group-separator="," :value="attackData.high" title="高危漏洞" />
          </el-col>
          <el-col :span="8">
            <el-statistic group-separator="," :value="attackData.block" title="拦截漏洞" />
          </el-col>
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span>近7天攻击趋势</span>
          </div>
          <div id="attack_trends" style="width: auto;height:400px;" />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span>漏洞类型</span>
          </div>
          <div id="attack_type" style="width: auto;height:400px;" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { getAttackData, getAttackTypeData, getTrendsData } from '@/api/dashboard/dashboard'
import * as echarts from 'echarts'

export default {
  data() {
    return {
      currentDatetime: '',
      attackData: {
        total: 0,
        high: 0,
        block: 0
      }
    }
  },
  mounted() {
    this.refreshDate()
    this.getDashboardInfo()
    this.getEChartsInfo()
  },
  methods: {
    refreshDate() {
      setInterval(() => {
        // 获取当前时间
        const date = new Date()
        // 格式化为本地时间格式
        this.currentDatetime = date.toLocaleString()
      }, 1000)
    },

    async getDashboardInfo() {
      const res = await getAttackData()
      this.attackData = res.data
    },

    async getEChartsInfo() {
      const trendsDom = document.getElementById('attack_trends')
      const trendsChart = echarts.init(trendsDom)
      const trendsOption = await getTrendsData()
      trendsOption && trendsChart.setOption(trendsOption)

      const typeDom = document.getElementById('attack_type')
      const typeChart = echarts.init(typeDom)
      const typeOption = await getAttackTypeData()
      typeOption && typeChart.setOption(typeOption)
    }
  }
}
</script>

<style scoped>
.el-row {
  margin: 0px;
}
.el-col {
  padding: 20px;
}
</style>
