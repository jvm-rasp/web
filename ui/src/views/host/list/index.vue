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
          <el-button type="primary" @click="onSubmit">查询</el-button>
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
            <el-button type="text" size="medium">修改</el-button>
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
          <el-select v-model="dialogEditForm.strategyID" placeholder="请选择活动区域">
            <el-option label="策略1" value="1000" />
            <el-option label="策略2" value="2000" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary">确定</el-button>
          <el-button>取消</el-button>
        </el-form-item>
      </el-form>
    </el-dialog>
    <el-dialog title="实例详情" :visible.sync="dialogDetailVisible">
      <div style="margin: 20px;">
        <el-descriptions title="实例信息" :column="3" border>
          <el-descriptions-item label="主机名称">kooriookami</el-descriptions-item>
          <el-descriptions-item label="IP地址">18100000000</el-descriptions-item>
          <el-descriptions-item label="操作系统">苏州市</el-descriptions-item>
          <el-descriptions-item label="系统内存">苏州市</el-descriptions-item>
          <el-descriptions-item label="CPU核数">苏州市</el-descriptions-item>
          <el-descriptions-item label="可用磁盘">苏州市</el-descriptions-item>
          <el-descriptions-item label="接入方式">
            <el-tag size="small">动态</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="在线状态">
            <el-tag size="small">在线</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="配置生效时间">苏州市</el-descriptions-item>
          <el-descriptions-item label="注册时间">苏州市</el-descriptions-item>
          <el-descriptions-item label="备注信息">苏州市</el-descriptions-item>
        </el-descriptions>
      </div>
      <div style="margin: 20px;">
        <el-descriptions title="守护进程信息" :column="3" border>
          <el-descriptions-item label="运行状态">运行中</el-descriptions-item>
          <el-descriptions-item label="版本">运行中</el-descriptions-item>
          <el-descriptions-item label="安装路径">运行中</el-descriptions-item>
          <el-descriptions-item label="MD5">运行中</el-descriptions-item>
          <el-descriptions-item label="编译代码">运行中</el-descriptions-item>
          <el-descriptions-item label="编译时间">运行中</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { deleteHost, getHosts } from '@/api/host/host'

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
      total: 0,
      dialogEditVisible: true,
      dialogDetailVisible: false,
      dialogEditForm: {}
    }
  },
  mounted() {
    this.getTableData()
  },
  methods: {
    onSubmit() {
      console.log('submit!')
      this.getTableData()
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

    // 获取表格数据
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getHosts(this.params)
        console.log(data.list)
        this.tableData = data.list
        this.total = data.total
      } finally {
        this.loading = false
      }
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
