<template>
  <div id="root">
    <el-card class="app-container" :body-style="{ padding: '10px' }">
      <el-form :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="实例名称">
          <el-input v-model="params.name" placeholder="请输入主机名" />
        </el-form-item>
        <el-form-item label="IP">
          <el-input v-model="params.ip" placeholder="请输入IP地址" />
        </el-form-item>
        <el-form-item label="接入方式">
          <el-select v-model="params.method" placeholder="请选择">
            <el-option label="动态" value="动态" />
            <el-option label="静态" value="静态" />
            <el-option label="关闭" value="关闭" />
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
        <el-table-column prop="serverName" label="实例名称" />
        <el-table-column prop="serverIP" label="IP地址" />
        <el-table-column prop="agentVersion" label="RASP版本" width="90" />
        <el-table-column prop="agentType" label="接入方式" width="90">
          <template slot-scope="scope">
            <el-tag size="small" :type="getAgentTypeColor(scope.row.agentType)" disable-transitions>{{ scope.row.agentType }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="strategyID" label="策略ID" width="90" />
        <el-table-column prop="online" label="在线状态" width="90">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.online === 1 ? 'success':'danger'" disable-transitions>{{ scope.row.online === 1 ? '在线':'离线' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="configEffectTime" label="配置生效时间" />
        <el-table-column prop="agentRegisterTime" label="注册时间" />
        <el-table-column fixed="right" label="操作" align="center">
          <template slot-scope="scope">
            <el-button type="text" size="medium" @click="handleDelete(scope.row.serverID)">删除</el-button>
            <el-button type="text" size="medium" @click="handleEdit(scope.row)">修改</el-button>
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
    <el-dialog title="实例修改" :visible.sync="dialogEditVisible">
      <el-form ref="form" :model="dialogEditForm" label-width="80px">
        <el-form-item label="选择策略">
          <el-select v-model="dialogEditForm.strategyID" style="width: 100%" placeholder="请选择策略">
            <el-option
              v-for="item in configData"
              :key="item.ID"
              :label="item.name"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onUpdate">确定</el-button>
          <el-button @click="dialogEditVisible=false">取消</el-button>
        </el-form-item>
      </el-form>
    </el-dialog>
    <el-dialog title="实例详情" :visible.sync="dialogDetailVisible">
      <div style="margin: 20px;">
        <el-descriptions title="实例信息" :column="3" border>
          <el-descriptions-item label="主机名称">{{ recordDetail.host.name }}</el-descriptions-item>
          <el-descriptions-item label="IP地址">{{ recordDetail.host.ip }}</el-descriptions-item>
          <el-descriptions-item label="操作系统">{{ recordDetail.host.os }}</el-descriptions-item>
          <el-descriptions-item label="系统内存">{{ recordDetail.host.memory }}</el-descriptions-item>
          <el-descriptions-item label="CPU核数">{{ recordDetail.host.cpu }}</el-descriptions-item>
          <el-descriptions-item label="可用磁盘">{{ recordDetail.host.disk }}</el-descriptions-item>
        </el-descriptions>
      </div>
      <div style="margin: 20px;">
        <el-descriptions title="守护进程信息" :column="2" border>
          <el-descriptions-item label="运行状态">
            <el-tag size="small" :type="recordDetail.daemon.state === 1 ? 'success':'danger'" disable-transitions>{{ recordDetail.daemon.state === 1 ? '运行中':'离线中' }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="版本">{{ recordDetail.daemon.version }}</el-descriptions-item>
          <el-descriptions-item label="安装路径">{{ recordDetail.daemon.path }}</el-descriptions-item>
          <el-descriptions-item label="MD5">{{ recordDetail.daemon.md5 }}</el-descriptions-item>
          <el-descriptions-item label="编译代码">{{ recordDetail.daemon.buildCode }}</el-descriptions-item>
          <el-descriptions-item label="编译时间">{{ recordDetail.daemon.buildTime }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { deleteHost, getDetail, getHosts, updateHost } from '@/api/host/host'
import { getConfigs } from '@/api/config/config'

export default {
  data() {
    return {
      params: {
        name: '',
        ip: '',
        method: '',
        pageNum: 1,
        pageSize: 10
      },

      tableData: [],
      configData: [],
      total: 0,
      dialogEditVisible: false,
      dialogDetailVisible: false,
      dialogEditForm: {},
      selectRecord: {},
      recordDetail: {
        host: {
          name: null,
          ip: null,
          os: null,
          memory: null,
          cpu: null,
          disk: null
        },
        daemon: {
          state: null,
          version: null,
          path: null,
          md5: null,
          buildCode: null,
          buildTime: null
        }
      }
    }
  },
  mounted() {
    this.getTableData()
    this.getConfigData()
  },
  methods: {
    onSearch() {
      this.getTableData()
    },

    async onUpdate() {
      const { code, message } = await updateHost(this.dialogEditForm)
      console.log(code)
      if (code === 200) {
        this.$message({
          type: 'success',
          message: message
        })
        this.dialogEditVisible = false
      } else {
        this.$message({
          type: 'error',
          message: message
        })
      }
    },

    onClear() {
      this.params.name = ''
      this.params.ip = ''
      this.params.method = ''
    },

    handlePageSizeChanged(val) {
      this.params.pageSize = val
      this.getTableData()
    },

    handleCurrentPageChanged(val) {
      this.params.pageNum = val
      this.getTableData()
    },

    handleDelete(serverID) {
      this.$confirm('您确定要删除该实例吗？删除后不能恢复！', '删除实例', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async() => {
        const data = await deleteHost({ id: serverID })
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

    handleEdit(record) {
      this.dialogEditForm = record
      this.dialogEditVisible = true
    },

    async handleDetail(record) {
      this.selectRecord = record
      const { data } = await getDetail(this.params)
      this.recordDetail = data
      this.dialogDetailVisible = true
    },
    // 获取表格数据
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getHosts(this.params)
        this.tableData = data.list
        this.total = data.total
      } finally {
        this.loading = false
      }
    },

    async getConfigData() {
      const { data } = await getConfigs({})
      this.configData = data.data
    },

    getAgentTypeColor(type) {
      if (type === '动态') {
        return 'success'
      } else if (type === '静态') {
        return ''
      } else if (type === '关闭') {
        return 'danger'
      }
    }
  }
}
</script>

<style scoped>

</style>
