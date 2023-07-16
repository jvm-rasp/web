<template>
  <div>
    <el-card class="container-card" shadow="always">
      <el-descriptions class="margin-top" title="主机详情" :column="2" size="medium" border>
        <template slot="extra">
          <el-button type="primary" size="small" @click="backup()">返回</el-button>
        </template>
        <el-descriptions-item label="主机名称">{{ hostInfo.hostName }}</el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ hostInfo.ip }}</el-descriptions-item>
        <el-descriptions-item label="操作系统">{{ hostInfo.osType }}</el-descriptions-item>
        <el-descriptions-item label="系统内存">{{ hostInfo.totalMem }} GB</el-descriptions-item>
        <el-descriptions-item label="CPU核数">{{ hostInfo.cpuCounts }}</el-descriptions-item>
        <el-descriptions-item label="可用磁盘">{{ hostInfo.freeDisk }} GB</el-descriptions-item>
        <el-descriptions-item label="RASP版本">{{ hostInfo.version }}</el-descriptions-item>
        <el-descriptions-item label="接入方式">
          <template>
            <el-tag size="medium" disable-transitions>
              {{ getAgentMode(hostInfo.agentMode).value }}
            </el-tag>
          </template>
        </el-descriptions-item>
        <el-descriptions-item label="编译信息">{{ hostInfo.buildGitBranch }}:{{
          hostInfo.buildGitCommit
        }}:{{ hostInfo.buildDateTime }}
        </el-descriptions-item>
        <el-descriptions-item label="心跳时间">{{ hostInfo.heartbeatTime }}</el-descriptions-item>
        <el-descriptions-item label="注册时间">{{ hostInfo.CreatedAt }}</el-descriptions-item>
        <el-descriptions-item label="安装目录">{{ hostInfo.installDir }}</el-descriptions-item>
        <el-descriptions-item label="可执行文件hash">{{ hostInfo.exeFileHash }}</el-descriptions-item>
      </el-descriptions>
      <br>
      <el-table
        v-loading="processLoading"
        :data="processTableData"
        border
        stripe
        style="width: 100%"
      >
        <el-table-column show-overflow-tooltip prop="pid" label="PID" align="center" />
        <el-table-column show-overflow-tooltip prop="startTime" label="启动时间" align="center" />
        <el-table-column show-overflow-tooltip prop="status" label="防护状态" align="center">
          <template slot-scope="scope">
            <el-tag size="small" :type="getAgentStatusColor(scope.row.status)" disable-transitions>
              {{ getAgentStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip prop="cmdlineInfo" label="命令行信息" align="center" />
      </el-table>
      <el-pagination
        :current-page="processParams.pageNum"
        :page-size="processParams.pageSize"
        :total="processTotal"
        :page-sizes="[1, 2, 5, 10]"
        layout="total, prev, pager, next, sizes"
        background
        style="margin-top: 10px;float:right;margin-bottom: 10px;"
        @size-change="processHandleSizeChange"
        @current-change="processHandleCurrentChange"
      />

    </el-card>
  </div>
</template>

<script>
import { getDetail, getProcesss } from '@/api/host/host'

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
      },

      // 进程详情查询参数
      processParams: {
        hostName: '',
        status: '',
        pageNum: 1,
        pageSize: 10
      },
      processTableData: [],
      processTotal: 0
    }
  },
  created() {
    this.hostName = this.$route.query.hostName
    this.processParams.hostName = this.hostName
    this.getHostDetail(this.hostName)
    this.getProcessTable()
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
    },

    // 分页
    processHandleSizeChange(val) {
      this.processParams.pageSize = val
      this.getProcessTable()
    },

    processHandleCurrentChange(val) {
      this.processParams.pageNum = val
      this.getProcessTable()
    },

    async getProcessTable() {
      this.processLoading = true
      try {
        const { data } = await getProcesss(this.processParams)
        this.processTableData = data.data
        this.processTotal = data.total
      } finally {
        this.processLoading = false
      }
    },
    getAgentStatusColor(status) {
      if (status === 1) {
        return 'success'
      } else if (status === 0) {
        return 'primary'
      } else if (status === 2) {
        return 'danger'
      }
    },

    getAgentStatusText(status) {
      if (status === 1) {
        return '防护中'
      } else if (status === 0) {
        return '未安装'
      } else if (status === 2) {
        return '安装失败'
      }
    },
    backup() {
      this.$router.go(-1)
    }
  }
}
</script>

<style scoped>

</style>
