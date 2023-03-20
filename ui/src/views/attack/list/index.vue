<template>
  <div id="root">
    <el-card class="app-container" :body-style="{ padding: '10px' }">
      <el-form :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="实例名称">
          <el-input v-model="params.name" placeholder="请输入主机名" />
        </el-form-item>
        <el-form-item label="阻断状态">
          <el-select v-model="params.block" placeholder="请选择">
            <el-option label="放行" :value="0" />
            <el-option label="阻断" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item label="处理状态">
          <el-select v-model="params.status" placeholder="请选择">
            <el-option label="未处理" :value="handleStatus['未处理']" />
            <el-option label="处理中" :value="handleStatus['处理中']" />
            <el-option label="已处理" :value="handleStatus['已处理']" />
            <el-option label="误报" :value="handleStatus['误报']" />
            <el-option label="忽略" :value="handleStatus['忽略']" />
            <el-option label="全部" value="" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSearch">查询</el-button>
          <el-button type="primary" @click="onClear">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    <el-card class="app-container" :body-style="{ padding: '10px' }">
      <el-table :data="tableData" border stripe style="width: 100%">
        <el-table-column prop="attackTime" label="攻击时间" />
        <el-table-column prop="hostName" label="实例名称" />
        <el-table-column prop="attackType" label="攻击类型">
          <template slot-scope="scope">
            <el-tag size="small" type="" effect="plain" disable-transitions>{{ scope.row.attackType }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remoteIp" label="攻击IP" />
        <el-table-column prop="level" label="风险等级">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.level >= 90 ? 'danger' : ''" disable-transitions>
              {{ scope.row.level >= 90 ? '高危' : '中危' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="isBlocked" label="阻断状态">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.isBlocked ? 'danger' : 'success'" disable-transitions>
              {{ scope.row.isBlocked ? '阻断' : '放行' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="handleStatus" label="处理状态">
          <template slot-scope="scope">
            <el-tag size="small" :type="getHandleStatusColor(scope.row.handleStatus)" disable-transitions>
              {{ getHandleStatusLabel(scope.row.handleStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" align="left">
          <template slot-scope="scope">
            <el-button type="text" size="medium" @click="handleDelete(scope.row)">删除</el-button>
            <el-button type="text" size="medium" @click="handleDetail(scope.row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        :current-page="params.pageNum"
        :page-size="params.pageSize"
        :total="total"
        :page-sizes="[1, 5, 10, 30]"
        layout="total, prev, pager, next, sizes"
        background
        style="margin-top: 10px;float:right;margin-bottom: 10px;"
        @size-change="handlePageSizeChanged"
        @current-change="handleCurrentPageChanged"
      />
    </el-card>

    <el-dialog title="攻击详情" :visible.sync="dialogDetailVisible" width="1000px">
      <div style="margin: 20px;">
        <el-descriptions title="攻击详情" :column="2" border>
          <el-descriptions-item label="主机名称">{{ attackDetail.hostName }}</el-descriptions-item>
          <el-descriptions-item label="受攻击IP">{{ attackDetail.localIp }}</el-descriptions-item>
          <el-descriptions-item label="攻击IP">{{ attackDetail.remoteIp }}</el-descriptions-item>
          <el-descriptions-item label="攻击时间">{{ attackDetail.attackTime }}</el-descriptions-item>
          <el-descriptions-item label="危险等级">{{ attackDetail.level }}</el-descriptions-item>
          <el-descriptions-item label="阻断状态">{{ attackDetail.isBlocked }}</el-descriptions-item>
          <el-descriptions-item label="请求协议">{{ attackDetail.requestProtocol }}</el-descriptions-item>
          <el-descriptions-item label="请求方法">{{ attackDetail.httpMethod }}</el-descriptions-item>
        </el-descriptions>
      </div>
      <!--      <div style="margin: 20px;">-->
      <!--        <el-descriptions title="请求信息" direction="vertical" :column="1">-->
      <!--          <el-descriptions-item label="请求uri">{{ attackDetail.requestUri }}</el-descriptions-item>-->
      <!--          <el-descriptions-item label="请求参数">{{ attackDetail.requestParameters }}</el-descriptions-item>-->
      <!--          <el-descriptions-item label="攻击参数">{{ attackDetail.attackParameters }}</el-descriptions-item>-->
      <!--          <el-descriptions-item label="检测结果">{{ attackDetail.checkType }}</el-descriptions-item>-->
      <!--          <el-descriptions-item label="堆栈信息">-->
      <!--            <json-viewer v-model="attackDetail.stackTrace" copyable :boxed="false"/>-->
      <!--          </el-descriptions-item>-->
      <!--        </el-descriptions>-->
      <!--      </div>-->
    </el-dialog>
  </div>
</template>

<script>
import { deleteAttackLog, getAttackLogs } from '@/api/log/attackLog'

export default {
  data() {
    return {
      params: {
        name: '',
        block: '',
        status: '',
        pageNum: 1,
        pageSize: 20
      },
      tableData: [],
      total: 0,
      handleStatus: {
        '未处理': 0,
        '处理中': 1,
        '已处理': 2,
        '误报': 3,
        '忽略': 4
      },
      dialogDetailVisible: false,
      attackDetail: {
        id: null,
        hostName: null,
        remoteIp: null,
        localIp: null,
        attackType: null,
        checkType: null,
        isBlocked: null,
        level: null,
        tag: null,
        handleStatus: null,
        handleResult: null,
        stackTrace: null,
        httpMethod: null,
        requestProtocol: null,
        requestUri: null,
        requestParameters: null,
        attackParameters: null,
        cookies: null,
        header: null,
        body: null,
        attackTime: null,
        createTime: null,
        updateTime: null
      }
    }
  },
  mounted() {
    this.getTableData()
  },
  methods: {
    onSearch() {
    },
    onClear() {
      this.params.name = ''
      this.params.ip = ''
      this.params.method = ''
    },

    getHandleStatusLabel(status) {
      for (const k in this.handleStatus) {
        const val = this.handleStatus[k]
        if (val === status) {
          return k
        }
      }
      return '未知'
    },

    getHandleStatusColor(status) {
      switch (status) {
        case 0:
          return 'danger'
        case 1:
          return ''
        case 2:
          return 'success'
        case 3:
          return 'info'
        case 4:
          return 'info'
        default:
          return 'danger'
      }
    },

    handlePageSizeChanged(val) {
      this.params.pageSize = val
      this.getTableData()
    },

    handleCurrentPageChanged(val) {
      this.params.pageNum = val
      this.getTableData()
    },

    handleDelete(record) {
      this.$confirm('您确定要删除该日志吗？删除后不能恢复！', '删除日志', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async() => {
        const data = await deleteAttackLog({ ids: [record.ID] })
        await this.getTableData()
        if (data.code === 200) {
          this.$message({
            type: 'success',
            message: data.message
          })
        } else {
          this.$message({
            type: 'error',
            message: data.message
          })
        }
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        })
      })
    },

    async handleDetail(record) {
      this.attackDetail.attackParameters = record.attackParameters
      this.attackDetail.attackTime = record.attackTime
      this.attackDetail.hostName = record.hostName
      this.attackDetail.remoteIp = record.remoteIp
      this.attackDetail.localIp = record.localIp
      this.attackDetail.checkType = record.checkType
      this.attackDetail.isBlocked = record.isBlocked
      this.attackDetail.level = record.level
      this.attackDetail.handleStatus = record.handleStatus
      this.attackDetail.handleResult = record.handleResult
      this.attackDetail.stackTrace = JSON.parse(this.attackDetail.stackTrace)
      this.attackDetail.httpMethod = record.httpMethod
      this.attackDetail.requestProtocol = record.requestProtocol
      this.attackDetail.requestUri = record.requestUri
      this.attackDetail.requestParameters = record.requestParameters
      this.attackDetail.attackParameters = record.attackParameters
      this.attackDetail.cookies = record.cookies
      this.attackDetail.header = record.header
      this.attackDetail.body = record.body
      this.dialogDetailVisible = true
    },

    // 获取表格数据
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getAttackLogs(this.params)
        this.tableData = data.data
        this.total = data.total
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>

</style>
