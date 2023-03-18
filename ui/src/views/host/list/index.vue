<template>
  <div>
    <el-card class="container-card" shadow="always">
      <!-- 条件搜索框 -->
      <el-row>
        <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
          <el-form-item label="用户名">
            <el-input v-model.trim="params.hostName" clearable placeholder="名称" @clear="search" />
          </el-form-item>
          <el-form-item label="IP">
            <el-input v-model.trim="params.ip" clearable placeholder="IP" @clear="search" />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model.trim="params.agentMode" clearable placeholder="接入模式" @change="search" @clear="search">
              <el-option label="动态" value="dynamic" />
              <el-option label="静态" value="static" />
              <el-option label="禁用" value="disable" />
            </el-select>
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
          <el-form-item>
            <el-button
              :disabled="multipleSelection.length === 0"
              :loading="loading"
              icon="el-icon-delete"
              type="warning"
              @click="batchUpdate"
            >批量更新(TODO)
            </el-button>
          </el-form-item>
        </el-form>
      </el-row>
      <!-- 主机列表 -->
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="ID" label="ID" align="center" width="55" />
        <el-table-column show-overflow-tooltip sortable prop="hostName" label="实例名称" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="ip" label="IP" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="agentMode" label="接入方式" align="center">
          <template slot-scope="scope">
            <el-tag size="small" disable-transitions>
              {{ scope.row.agentMode }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip prop="version" label="RASP版本" align="center" width="100" />
        <el-table-column show-overflow-tooltip sortable prop="configId" label="配置ID" align="center" width="100" />
        <el-table-column show-overflow-tooltip sortable prop="heatbeatTime" label="心跳时间" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="CreatedAt" label="注册时间" align="center" />
        <el-table-column fixed="right" label="操作" align="center">
          <template slot-scope="scope">
            <el-button type="text" size="medium" @click="handleDelete(scope.row)">删除</el-button>
            <el-button type="text" size="medium" @click="showHostAndProcessDetail(scope.row)">详情</el-button>
            <el-button type="text" size="medium" @click="updateConfig(scope.row)">更新策略</el-button>
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

      <!-- 新建进程详情对话框 -->
      <el-dialog title="进程详情" :visible.sync="showProcessInfoVisible">
        <!-- 进程列表 -->
        <el-table
          v-loading="processLoading"
          :data="processTableData"
          border
          stripe
          style="width: 100%"
        >
          <el-table-column show-overflow-tooltip sortable prop="pid" label="PID" align="center" />
          <el-table-column show-overflow-tooltip sortable prop="startTime" label="启动时间" align="center" />
          <el-table-column show-overflow-tooltip sortable prop="status" label="保护状态" align="center" />
          <el-table-column show-overflow-tooltip sortable prop="status" label="状态" align="center">
            <template slot-scope="scope">
              <el-tag size="small" :type="getAgentStatusColor(scope.row.status)" disable-transitions>{{ getAgentStatusText(scope.row.status) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column show-overflow-tooltip sortable prop="cmdlineInfo" label="命令行信息" align="center" />
        </el-table>
        <el-pagination
          :current-page="processParams.pageNum"
          :page-size="processParams.pageSize"
          :total="processTotal"
          :page-sizes="[1, 2, 5, 10]"
          layout="total, prev, pager, next, sizes"
          background
          style="margin-top: 10px;float:left;margin-bottom: 10px;"
          @size-change="processHandleSizeChange"
          @current-change="processHandleCurrentChange"
        />
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="closeProcessInfoVisible()">关 闭</el-button>
        </div>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>

import { batchDeleteHost, getHosts, getProcesss } from '@/api/host/host'

export default {
  name: 'Host',
  data() {
    return {
      // 查询参数
      params: {
        hostName: '',
        ip: '',
        pageNum: 1,
        pageSize: 10
      },

      // 表格数据
      tableData: [],
      total: 0,
      loading: false,
      // 表格多选
      multipleSelection: [],

      // 进程详情弹窗
      // 查询参数
      processParams: {
        hostName: '',
        status: '',
        pageNum: 1,
        pageSize: 2
      },
      currentHostName: '',
      showProcessInfoVisible: false,
      processTableData: [],
      processTotal: 0,
      processLoading: false
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
        const { data } = await getHosts(this.params)
        this.tableData = data.data
        this.total = data.total
      } finally {
        this.loading = false
      }
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

    handleDelete(record) {
      this.$confirm('此操作将永久删除, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        let msg = ''
        try {
          const { message } = await batchDeleteHost({ ids: [record.ID] })
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
          const { message } = await batchDeleteHost({ ids: configIds })
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
    // 表格多选
    handleSelectionChange(val) {
      this.multipleSelection = val
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

    // process table
    showHostAndProcessDetail(row) {
      this.showProcessInfoVisible = true
      this.currentHostName = row.hostName
      this.getProcessTable(row)
    },
    // 分页
    processHandleSizeChange(val) {
      this.processParams.pageSize = val
      this.getProcessTable(this.currentHostName)
    },
    processHandleCurrentChange(val) {
      this.processParams.pageNum = val
      this.getProcessTable(this.currentHostName)
    },
    // 关闭配置创建
    closeProcessInfoVisible() {
      this.showProcessInfoVisible = false
      this.currentHostName = ''
    },
    getAgentStatusColor(type) {
      if (type === 'success inject' || type === 'success degrade') {
        return 'success'
      } else if (type === 'not inject') {
        return 'primary'
      } else if (type === 'failed inject' || type === 'failed uninstall agent' || type === 'failed degrade') {
        return 'danger'
      }
    },
    getAgentStatusText(type) {
      if (type === 'success inject' || type === 'success degrade') {
        return '防护中'
      } else if (type === 'not inject') {
        return '未安装'
      } else if (type === 'failed inject' || type === 'failed uninstall agent' || type === 'failed degrade') {
        return '安装失败'
      }
    }
  }
}
</script>

<style scoped>
.container-card {
  margin: 10px;
}

.delete-popover {
  margin-left: 10px;
}

.show-pwd {
  position: absolute;
  right: 10px;
  top: 3px;
  font-size: 16px;
  color: #889aa4;
  cursor: pointer;
  user-select: none;
}
</style>
