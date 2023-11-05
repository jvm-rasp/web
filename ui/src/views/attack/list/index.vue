<template>
  <div id="root">
    <el-card class="container-card" shadow="always">
      <!-- 条件搜索框 -->
      <el-form :size="this.$store.getters.size" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="实例名称">
          <el-input v-model="params.hostName" placeholder="请输入主机名" clearable @clear="onSearch" />
        </el-form-item>
        <el-form-item label="URL关键字">
          <el-input v-model="params.url" placeholder="请输入URL关键字" clearable @clear="onSearch" />
        </el-form-item>
        <el-form-item label="阻断状态">
          <el-select v-model="params.isBlocked" placeholder="请选择" clearable @change="onSearch" @clear="onSearch">
            <el-option label="放行" :value="false" />
            <el-option label="阻断" :value="true" />
            <el-option label="全部" value="" />
          </el-select>
        </el-form-item>
        <el-form-item label="处理状态">
          <el-select v-model="params.handleResult" placeholder="请选择" clearable @change="onSearch" @clear="onSearch">
            <el-option label="未处理" :value="handleStatus['未处理']" />
            <el-option label="已确认" :value="handleStatus['已确认']" />
            <el-option label="误报" :value="handleStatus['误报']" />
            <el-option label="忽略" :value="handleStatus['忽略']" />
            <el-option label="全部" value="" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="el-icon-search" @click="onSearch">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="warning" icon="el-icon-refresh" @click="onClear">重置</el-button>
        </el-form-item>
        <el-form-item>
          <el-button
            :disabled="multipleSelection.length === 0"
            :loading="loading"
            icon="el-icon-delete"
            type="danger"
            @click="batchDelete"
          >批量删除
          </el-button>
        </el-form-item>
      </el-form>
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%"
        :size="this.$store.getters.size"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column label="序号" type="index" width="50" align="center">
          <template slot-scope="scope">
            {{ (params.pageNum - 1) * params.pageSize + scope.$index + 1 }}
          </template>
        </el-table-column>
        <el-table-column prop="attackTime" label="攻击时间" width="180" :formatter="dateFormat" align="center" />
        <!--        <el-table-column prop="hostName" label="实例名称" width="150" align="center" />-->
        <el-table-column prop="hostIp" label="实例IP" width="150" align="center" />
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
        <el-table-column prop="HandleResult" label="处理状态" width="100" align="center">
          <template slot-scope="scope">
            <el-tag size="small" :type="getHandleStatusColor(scope.row.handleResult)" disable-transitions>
              {{ getHandleStatusLabel(scope.row.handleResult) }}
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
        <el-button type="primary" @click="updateStatus(selectRecord.record.ID, 1)">确认</el-button>
        <el-button @click="updateStatus(selectRecord.record.ID, 2)">误报</el-button>
        <el-button @click="updateStatus(selectRecord.record.ID, 3)">忽略</el-button>
      </el-row>
      <el-tabs type="border-card">
        <el-tab-pane label="漏洞详情">
          <el-descriptions title="" size="medium" :column="2" border>
            <el-descriptions-item label="受攻击IP">{{ selectRecord.record.hostIp }}</el-descriptions-item>
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
            <el-descriptions-item label="规则版本">{{ selectRecord.detail.metaInfo }}</el-descriptions-item>
            <el-descriptions-item label="攻击参数">{{ selectRecord.detail.payload }}</el-descriptions-item>
          </el-descriptions>
          <el-descriptions
            title=""
            size="medium"
            direction="vertical"
            :label-style="{'text-align': 'center'}"
            :content-style="{'word-break': 'normal', 'font-size': 'small', 'border-top-width': '0px'}"
            :column="1"
            border
          >
            <el-descriptions-item label="堆栈信息">
              <highlightjs language="java" :code="formatStackTrace(selectRecord.detail.stackTrace)" />
            </el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="请求信息">
          <el-descriptions title="" :column="2" direction="horizontal" border>
            <el-descriptions-item label="请求协议">
              {{ selectRecord.detail.context.protocol }}
            </el-descriptions-item>
            <el-descriptions-item label="请求方法">
              {{ selectRecord.detail.context.method }}
            </el-descriptions-item>
            <el-descriptions-item label="本地IP">
              {{ selectRecord.detail.context.localAddr }}
            </el-descriptions-item>
            <el-descriptions-item label="远程IP">
              {{ selectRecord.detail.context.remoteHost }}
            </el-descriptions-item>
            <el-descriptions-item label="请求类型">
              {{ selectRecord.detail.context.contentType }}
            </el-descriptions-item>
            <el-descriptions-item label="请求长度">
              {{ selectRecord.detail.context.contentLength }}
            </el-descriptions-item>
          </el-descriptions>
          <el-descriptions title="" :column="1" direction="vertical" border>
            <el-descriptions-item label="请求URL">
              {{ selectRecord.detail.context.requestURL }}
            </el-descriptions-item>
            <el-descriptions-item label="查询参数">
              {{ selectRecord.detail.context.queryString }}
            </el-descriptions-item>
          </el-descriptions>
          <el-descriptions title="" :column="2" direction="vertical" :label-style="{'text-align': 'center'}" border>
            <el-descriptions-item label="完整请求头">
              <pre>{{ selectRecord.detail.context.header }}</pre>
            </el-descriptions-item>
          </el-descriptions>
          <el-descriptions title="" :column="2" direction="vertical" :label-style="{'text-align': 'center'}" border>
            <el-descriptions-item label="Form请求参数">
              <pre>{{ selectRecord.detail.context.parameters }}</pre>
            </el-descriptions-item>
          </el-descriptions>
          <el-descriptions title="" :column="2" direction="vertical" :label-style="{'text-align': 'center'}" border>
            <el-descriptions-item label="Body参数">
              <el-button v-if="decodeBtnVisible" type="primary" @click="decodeBody(selectRecord.detail.context.body)">解码数据</el-button>
              <pre>{{ selectRecord.detail.context.body }}</pre>
            </el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>
  </div>
</template>

<script>
import { deleteAttackLog, getAttackDetail, getAttackLogs, updateStatus } from '@/api/log/attackLog'
import moment from 'moment'

export default {
  data() {
    return {
      params: {
        hostName: '',
        url: '',
        isBlocked: '',
        handleResult: '',
        pageNum: 1,
        pageSize: 10
      },
      tableData: [],
      total: 0,
      loading: false,
      handleStatus: {
        '未处理': 0,
        '已确认': 1,
        '误报': 2,
        '忽略': 3
      },
      dialogDetailVisible: false,
      selectRecord: {
        record: {
          hostName: '',
          hostIp: '',
          remoteIp: '',
          attackTime: '',
          attackType: '',
          algorithm: '',
          extend: '',
          metaInfo: '',
          level: '',
          isBlocked: ''
        },
        detail: {
          context: {}
        }
      },
      multipleSelection: [],
      decodeBtnVisible: true
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
      this.params.hostName = ''
      this.params.isBlocked = ''
      this.params.handleResult = ''
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
        const data = await deleteAttackLog({ guids: [record.rowGuid] })
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

    // 批量删除
    batchDelete() {
      this.$confirm('此操作将永久删除, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const guids = []
        this.multipleSelection.forEach(x => {
          guids.push(x.rowGuid)
        })
        let msg = ''
        try {
          const { message } = await deleteAttackLog({ guids: guids })
          msg = message
        } finally {
          this.loading = false
        }
        await this.getTableData()
        this.$message({
          showClose: true,
          message: msg,
          type: 'success'
        })
      }).catch(() => {
        this.$message({
          showClose: true,
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
      this.decodeBtnVisible = true
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
    // 更新状态
    async updateStatus(id, result) {
      this.$confirm('您确定要更新状态吗？', '更新处理状态', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async() => {
        const data = await updateStatus({ id: id, result: result })
        await this.getTableData()
        this.dialogDetailVisible = false
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
          message: '已取消操作'
        })
      })
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
    },
    decodeBody(body) {
      const asciiArray = body.split(',')
      let str = ''
      for (let i = 0; i < asciiArray.length; i++) {
        str += String.fromCharCode(asciiArray[i])
      }
      this.selectRecord.detail.context.body = str
      this.decodeBtnVisible = false
    }
  }
}
</script>

<style scoped>
pre {
  word-wrap: break-word;
  white-space: pre-wrap;
}
</style>
