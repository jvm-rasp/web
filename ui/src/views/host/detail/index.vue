<template>
  <div>
    <el-card class="container-card" shadow="always">
      <el-descriptions class="margin-top" title="主机详情" :column="2" size="medium" border>
        <template slot="extra">
          <el-button type="primary" size="small">操作</el-button>
        </template>
        <el-descriptions-item label="主机名称">{{ hostInfo.hostName }}</el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ hostInfo.ip }}</el-descriptions-item>
        <el-descriptions-item label="操作系统">{{ hostInfo.osType }}</el-descriptions-item>
        <el-descriptions-item label="系统内存(GB)">{{ hostInfo.totalMem }}</el-descriptions-item>
        <el-descriptions-item label="CPU核数">{{ hostInfo.cpuCounts }}</el-descriptions-item>
        <el-descriptions-item label="可用磁盘(GB)">{{ hostInfo.freeDisk }}</el-descriptions-item>
        <el-descriptions-item label="RASP版本">{{ hostInfo.version }}</el-descriptions-item>
        <el-descriptions-item label="接入方式">
          <template>
            <el-tag size="medium" disable-transitions>
              {{ getAgentMode(hostInfo.agentMode).value }}
            </el-tag>
          </template>
        </el-descriptions-item>
        <el-descriptions-item label="编译时间">{{ hostInfo.buildDateTime }}</el-descriptions-item>
        <el-descriptions-item label="编译分支">{{ hostInfo.buildGitBranch }}</el-descriptions-item>
        <el-descriptions-item label="代码commit">{{ hostInfo.buildGitCommit }}</el-descriptions-item>
        <el-descriptions-item label="心跳时间">{{ hostInfo.heartbeatTime }}</el-descriptions-item>
        <el-descriptions-item label="注册时间">{{ hostInfo.CreatedAt }}</el-descriptions-item>
        <el-descriptions-item label="安装目录">{{ hostInfo.installDir }}</el-descriptions-item>
        <el-descriptions-item label="可执行文件hash">{{ hostInfo.exeFileHash }}</el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script>
import { getDetail } from '@/api/host/host'

export default {
  data() {
    return {
      hostInfo: {
        // system
        hostName: '',
        ip: '',
        osType: '',
        totalMem: '',
        cpuCounts: '',
        freeDisk: '',
        version: '',
        agentMode: '',
        buildDateTime: '',
        buildGitBranch: '',
        buildGitCommit: '',
        heartbeatTime: '',
        CreatedAt: '',
        installDir: '',
        exeFileHash: ''
      }
    }
  },
  created() {
    this.hostName = this.$route.query.hostName
    this.getHostDetail(this.hostName)
  },
  methods: {
    // 获取表格数据
    async getHostDetail(hostName) {
      const { data } = await getDetail(hostName)
      this.hostInfo.hostName = data.data.hostName
      this.hostInfo.osType = data.data.osType
      this.hostInfo.ip = data.data.ip
      this.hostInfo.totalMem = data.data.totalMem
      this.hostInfo.cpuCounts = data.data.cpuCounts
      this.hostInfo.freeDisk = data.data.freeDisk
      this.hostInfo.version = data.data.version
      this.hostInfo.agentMode = data.data.agentMode
      this.hostInfo.buildDateTime = data.data.buildDateTime
      this.hostInfo.buildGitBranch = data.data.buildGitBranch
      this.hostInfo.buildGitCommit = data.data.buildGitCommit
      this.hostInfo.heartbeatTime = data.data.heartbeatTime
      // 时间格式转换
      this.hostInfo.CreatedAt = this.$dayjs(data.data.CreatedAt).format('YYYY-MM-DD HH:mm:ss.SSS')
      this.hostInfo.installDir = data.data.installDir
      this.hostInfo.exeFileHash = data.data.exeFileHash
    },
    getAgentMode(mode) {
      switch (mode) {
        case 'static':
          return { value: '静态', color: '' }
        case 'dynamic':
          return { value: '动态', color: 'success' }
        case 'disable':
          return { value: '禁用', color: 'info' }
        case '':
          return { value: '全部', color: '' }
        default:
          return { value: '未知', color: 'danger' }
      }
    }
  }
}
</script>

<style scoped>

</style>
