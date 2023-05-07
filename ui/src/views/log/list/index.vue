<template>
  <div id="root">
    <el-card class="container-card" shadow="always">
      <!-- 条件搜索框 -->
      <el-row>
        <el-form :size="this.$store.getters.size" :inline="true" :model="params" class="demo-form-inline">
          <el-form-item label="日志类型">
            <el-select v-model.trim="params.topic" clearable placeholder="日志类型" @change="search" @clear="search">
              <el-option label="Agent" value="jrasp-agent" />
              <el-option label="Daemon" value="jrasp-daemon" />
              <el-option label="Module" value="jrasp-module" />
              <el-option label="全部" value="" />
            </el-select>
          </el-form-item>
          <el-form-item label="日志等级">
            <el-select v-model.trim="params.level" clearable placeholder="日志等级" @change="search" @clear="search">
              <el-option label="ERROR" value="ERROR" />
              <el-option label="WARN" value="WARN" />
              <el-option label="全部" value="" />
            </el-select>
          </el-form-item>
          <el-form-item label="时间范围">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              align="right"
              unlink-panels
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              :picker-options="pickerOptions"
              value-format="yyyy-MM-dd"
              @change="handleDatePickerChange"
              @clear="search"
            />
          </el-form-item>
          <el-form-item>
            <el-button :loading="loading" icon="el-icon-search" type="primary" @click="search">查询</el-button>
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
      </el-row>
      <!-- 配置列表 -->
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
        <el-table-column show-overflow-tooltip sortable prop="topic" label="日志类型" align="center" width="150" />
        <el-table-column show-overflow-tooltip sortable prop="time" label="时间" align="center" :formatter="dateFormat" width="200" />
        <el-table-column show-overflow-tooltip sortable prop="level" label="日志等级" align="center" width="150" />
        <el-table-column show-overflow-tooltip sortable prop="message" label="日志消息" align="center" />
        <el-table-column fixed="right" label="操作" align="center" width="150">
          <template slot-scope="scope">
            <el-button type="text" size="medium" @click="handleDelete(scope.row)">删除</el-button>
            <el-button type="text" size="medium" @click="handleDetail(scope.row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
      <!--分页配置-->
      <el-pagination
        :current-page="params.pageNum"
        :page-size="params.pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, prev, pager, next, sizes"
        background
        style="margin-top: 10px;float:right;margin-bottom: 10px;"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
      <!-- 新建配置对话框 -->
      <el-dialog title="日志详情" :visible.sync="logDetailVisible" width="70%">
        <el-descriptions title="" size="medium" :column="3" border>
          <el-descriptions-item label="日志类型">{{ bindLogData.topic }}</el-descriptions-item>
          <el-descriptions-item label="时间">{{ bindLogData.time }}</el-descriptions-item>
          <el-descriptions-item label="日志等级">{{ bindLogData.level }}</el-descriptions-item>
        </el-descriptions>
        <el-descriptions
          title=""
          size="medium"
          direction="vertical"
          :label-style="{'text-align': 'center', 'word-break': 'keep-all'}"
          :content-style="{'word-break': 'break-all', 'font-size': 'small', 'border-top-width': '0px'}"
          :column="1"
          border
        >
          <el-descriptions-item label="日志消息">
            <div v-if="bindLogData.topic==='jrasp-daemon'">
              <vue-json-editor
                v-model="bindLogData.message"
                :show-btns="false"
                :mode="'code'"
                lang="zh"
              />
            </div>
            <div v-else>
              <pre>{{ bindLogData.message }}</pre>
            </div>
          </el-descriptions-item>
        </el-descriptions>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="closeDetail()">关 闭</el-button>
          <el-button size="mini" :loading="loading" type="primary" @click="closeDetail()">
            确 定
          </el-button>
        </div>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>
import { batchDeleteLogsByIds, getLogs } from '@/api/log/errorLog'
import moment from 'moment'
import vueJsonEditor from 'vue-json-editor'

export default {
  name: 'ErrorLog',
  components: { vueJsonEditor },
  data() {
    return {
      // 查询参数
      params: {
        topic: '',
        startDate: '',
        endDate: '',
        level: '',
        pageNum: 1,
        pageSize: 10
      },
      // 表格数据
      tableData: [],
      total: 0,
      loading: false,
      // 表格多选
      multipleSelection: [],
      logDetailVisible: false,
      bindLogData: {
        ID: '',
        topic: '',
        time: '',
        level: '',
        message: ''
      },
      pickerOptions: {
        shortcuts: [
          {
            text: '最近一周',
            onClick(picker) {
              const end = new Date()
              const start = new Date()
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
              picker.$emit('pick', [start, end])
            }
          },
          {
            text: '最近一个月',
            onClick(picker) {
              const end = new Date()
              const start = new Date()
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
              picker.$emit('pick', [start, end])
            }
          },
          {
            text: '最近三个月',
            onClick(picker) {
              const end = new Date()
              const start = new Date()
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
              picker.$emit('pick', [start, end])
            }
          }
        ]
      },
      dateRange: []
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    // 查询
    search() {
      this.params.pageNum = 1
      this.getTableData()
    },
    // 获取表格数据
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getLogs(this.params)
        this.tableData = data.list
        this.total = data.total
      } finally {
        this.loading = false
      }
    },
    // 表格多选
    handleSelectionChange(val) {
      this.multipleSelection = val
    },
    // 批量删除
    batchDelete() {
      this.$confirm('此操作将永久删除, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const configIds = []
        this.multipleSelection.forEach(x => {
          configIds.push(x.ID)
        })
        let msg = ''
        try {
          const { message } = await batchDeleteLogsByIds({ ids: configIds })
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
    handleDelete(record) {
      this.$confirm('此操作将永久删除, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        let msg = ''
        try {
          const { message } = await batchDeleteLogsByIds({ ids: [record.ID] })
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
    // 分页
    handleSizeChange(val) {
      this.params.pageSize = val
      this.getTableData()
    },
    handleCurrentChange(val) {
      this.params.pageNum = val
      this.getTableData()
    },
    dateFormat(row, column) {
      const date = row[column.property]
      if (date === undefined) {
        return ''
      }
      return moment(date).format('YYYY-MM-DD HH:mm:ss')
    },
    handleDetail(record) {
      this.bindLogData = JSON.parse(JSON.stringify(record))
      if (record.topic === 'jrasp-daemon') {
        this.bindLogData.message = JSON.parse(this.bindLogData.message)
      }
      this.logDetailVisible = true
    },
    closeDetail() {
      this.logDetailVisible = false
    },
    handleDatePickerChange(val) {
      console.log(val)
      if (val) {
        this.params.startDate = val[0]
        this.params.endDate = val[1]
      } else {
        this.params.startDate = ''
        this.params.endDate = ''
      }
      this.search()
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
