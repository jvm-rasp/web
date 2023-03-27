<template>
  <div id="root">
    <el-card class="container-card" shadow="always">
      <el-form size="medium" :inline="true" :model="params" class="demo-form-inline">
        <el-col :span="6">
          <el-form-item label="实例名称">
            <el-input v-model="params.name" placeholder="请输入主机名" />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="阻断状态">
            <el-select v-model="params.block" placeholder="请选择">
              <el-option label="放行" :value="0" />
              <el-option label="阻断" :value="1" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="6">
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
        </el-col>
        <el-col :span="6">
          <el-form-item>
            <el-button type="primary" @click="onSearch">查询</el-button>
            <el-button type="primary" @click="onClear">重置</el-button>
          </el-form-item>
        </el-col>
      </el-form>
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column prop="attackTime" label="攻击时间" width="180" :formatter="dateFormat" align="center" />
        <el-table-column prop="hostName" label="实例名称" width="150" align="center" />
        <el-table-column prop="attackType" label="攻击类型" width="180" align="center">
          <template slot-scope="scope">
            <el-tag size="small" color="#cd201f" style="color: #ffffff; font-weight: bold; border-color: #cd201f" effect="dark">{{ scope.row.attackType }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="requestUri" label="URL" header-align="center" />
        <el-table-column prop="remoteIp" label="攻击IP" width="150" align="center" />
        <el-table-column prop="level" label="风险等级" width="100" align="center">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.level >= 90 ? 'danger' : ''" disable-transitions>
              {{ scope.row.level >= 90 ? '高危' : '中危' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="isBlocked" label="阻断状态" width="100" align="center">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.isBlocked ? 'danger' : 'success'" disable-transitions>
              {{ scope.row.isBlocked ? '阻断' : '放行' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="handleStatus" label="处理状态" width="100" align="center">
          <template slot-scope="scope">
            <el-tag size="small" :type="getHandleStatusColor(scope.row.handleStatus)" disable-transitions>
              {{ getHandleStatusLabel(scope.row.handleStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" width="150" align="center">
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
        :page-sizes="[10, 20, 50, 100]"
        layout="total, prev, pager, next, sizes"
        background
        style="margin-top: 10px;float:right;margin-bottom: 10px;"
        @size-change="handlePageSizeChanged"
        @current-change="handleCurrentPageChanged"
      />
    </el-card>

    <el-dialog title="告警详情" :visible.sync="dialogDetailVisible" width="1000px">
      <el-row style="text-align: right; margin-bottom: 20px;">
        <el-button type="primary">确认</el-button>
        <el-button>误报</el-button>
        <el-button>忽略</el-button>
      </el-row>
      <el-tabs type="border-card">
        <el-tab-pane label="漏洞详情">
          <el-descriptions title="" size="medium" :column="2" border>
            <el-descriptions-item label="主机名称">{{ selectRecord.record.hostName }}</el-descriptions-item>
            <el-descriptions-item label="受攻击IP">{{ selectRecord.record.localIp }}</el-descriptions-item>
            <el-descriptions-item label="攻击IP">{{ selectRecord.record.remoteIp }}</el-descriptions-item>
            <el-descriptions-item label="攻击时间">{{ selectRecord.record.attackTime }}</el-descriptions-item>
            <el-descriptions-item label="危险等级">
              <el-tag size="small" :type="selectRecord.record.level >= 90 ? 'danger' : ''" disable-transitions>
                {{ selectRecord.record.level >= 90 ? '高危' : '中危' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="阻断状态">
              <el-tag size="small" :type="selectRecord.record.isBlocked ? 'danger' : 'success'" disable-transitions>
                {{ selectRecord.record.isBlocked ? '阻断' : '放行' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="攻击类型">{{ selectRecord.record.attackType }}</el-descriptions-item>
            <el-descriptions-item label="检测算法">{{ selectRecord.detail.algorithm }}</el-descriptions-item>
            <el-descriptions-item label="告警详情">{{ selectRecord.detail.extend }}</el-descriptions-item>
          </el-descriptions>
          <el-descriptions
            title=""
            size="medium"
            direction="vertical"
            :label-style="{'text-align': 'center'}"
            :content-style="{'word-break': 'normal'}"
            :column="1"
            border
          >
            <el-descriptions-item label="堆栈信息">
              <highlightjs class="stack" language="java" :code="formatStackTrace(selectRecord.detail.stackTrace)" />
            </el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="请求信息">
          <el-descriptions title="" :column="2" border>
            <el-descriptions-item label="请求协议">{{ selectRecord.requestProtocol }}</el-descriptions-item>
            <el-descriptions-item label="请求方法">{{ selectRecord.httpMethod }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>
  </div>
</template>

<script>
import { deleteAttackLog, getAttackDetail, getAttackLogs } from '@/api/log/attackLog'
import moment from 'moment'

export default {
  data() {
    return {
      params: {
        name: '',
        block: '',
        status: '',
        pageNum: 1,
        pageSize: 10
      },
      tableData: [],
      total: 0,
      loading: false,
      handleStatus: {
        '未处理': 0,
        '处理中': 1,
        '已处理': 2,
        '误报': 3,
        '忽略': 4
      },
      dialogDetailVisible: false,
      selectRecord: {
        record: {},
        detail: {}
      },
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
      this.params.pageNum = 1
      this.getTableData()
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

    // 表格多选
    handleSelectionChange(val) {
      this.multipleSelection = val
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
      this.selectRecord.record = record
      const { data } = await getAttackDetail({ id: record.rowGuid })
      this.selectRecord.detail = data.detail
      this.dialogDetailVisible = true
    },

    // 获取表格数据
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getAttackLogs(this.params)
        this.tableData = data.list
        this.total = data.total
      } finally {
        this.loading = false
      }
    },
    dateFormat(row, column) {
      const date = row[column.property]
      if (date === undefined) {
        return ''
      }
      return moment(date).format('YYYY-MM-DD HH:mm:ss')
    },
    formatStackTrace(text) {
      if (text) {
        const splitText = text.split(',')
        let str = ''
        splitText.forEach((item, index) => {
          if (item !== '') {
            str = str + '\n' + item
          }
        })
        return str.trim()
      } else {
        return ''
      }
    }
  }
}
</script>

<style scoped>
.stack {
  font-size: small;
}
</style>
