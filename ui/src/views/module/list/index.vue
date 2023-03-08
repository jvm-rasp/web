<template>
  <div>
    <el-card class="container-card" shadow="always">

      <!-- 条件搜索框 -->
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="名称">
          <el-input v-model.trim="params.name" clearable placeholder="名称" @clear="search" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model.trim="params.status" clearable placeholder="状态" @change="search" @clear="search">
            <el-option label="正常" value="1" />
            <el-option label="禁用" value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="创建人">
          <el-input v-model.trim="params.creator" clearable placeholder="创建人" @clear="search" />
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-search" type="primary" @click="search">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-plus" type="warning" @click="create">新增</el-button>
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

      <!-- 配置列表 -->
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="ID" with="55" label="模块ID" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="name" with="120" label="模块名称" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="desc" with="120" label="模块描述" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="status" with="60" label="模块状态" align="center">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.status === 1 ? 'success':'danger'" disable-transitions>
              {{ scope.row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="status" with="60" label="配置状态" align="center">
          <template slot="status" slot-scope="scope">
            <el-switch
              v-model="scope.row.status"
              on-value="true"
              off-value="false"
              @change="switchChange(scope.row)"
            />
          </template>
        </el-table-column>

        <el-table-column show-overflow-tooltip sortable prop="creator" label="创建人" align="center" />
        <el-table-column fixed="right" label="操作" align="center" width="240">
          <template slot-scope="scope">
            <el-tooltip content="查看模块" effect="dark" placement="top">
              <el-button size="mini" circle type="primary" icon="el-icon-view" @click="viewModuleDetail(scope.row)" />
            </el-tooltip>
            <el-tooltip class="delete-popover" content="删除模块" effect="dark" placement="top">
              <el-popconfirm title="确定删除吗？" @onConfirm="singleDelete(scope.row.ID)">
                <el-button slot="reference" size="mini" icon="el-icon-delete" circle type="danger" />
              </el-popconfirm>
            </el-tooltip>
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

      <!-- 模块内容详情对话框 -->
      <el-dialog title="配置详情" :visible.sync="moduleDetailVisible" width="60%">
        <el-form size="small" :model="moduleDetailData" label-width="100px">
          <json-viewer v-model="moduleDetailData.parameters" copyable boxed />
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="closeModuleDetail()">关 闭</el-button>
        </div>
      </el-dialog>

      <el-dialog title="新增防护模块" :visible.sync="createModuleVisible" width="60%">
        <el-form
          ref="createModuleForm"
          size="small"
          :model="createModuleData"
          label-width="100px"
        >
          <el-form-item label="模块名称" prop="name">
            <el-input v-model.trim="createModuleData.name" placeholder="模块名称" />
          </el-form-item>
          <el-form-item label="模块版本" prop="version">
            <el-input v-model.number="createModuleData.version" />
          </el-form-item>
          <el-form-item label="下载链接" prop="url">
            <el-input v-model.number="createModuleData.url" />
          </el-form-item>
          <el-form-item label="文件hash" prop="hash">
            <el-input v-model.number="createModuleData.hash" />
          </el-form-item>
          <el-form-item label="模块描述" prop="desc">
            <el-input v-model.trim="createModuleData.desc" placeholder="模块描述" />
          </el-form-item>
          <el-form-item label="配置参数" prop="content">
            <vue-json-editor
              v-model="createModuleData.parameters"
              :show-btns="false"
              :mode="'code'"
              lang="zh"
              @json-change="onJsonChange"
              @json-save="onJsonSave"
              @has-error="onError"
            />
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="closeCreateModule()">关 闭</el-button>
          <el-button size="mini" :loading="submitCreateModuleLoading" type="primary" @click="submitNewModuleForm()">
            确 定
          </el-button>
        </div>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>

import vueJsonEditor from 'vue-json-editor'
import { batchDeleteModuleByIds, createModule, getModules } from '@/api/module/module'
import { batchDeleteConfigByIds } from '@/api/config/config'

export default {
  name: 'Module',
  components: { vueJsonEditor },
  data() {
    return {
      // 查询参数
      params: {
        name: '',
        status: '',
        pageNum: 1,
        pageSize: 10
      },

      // 表格数据
      tableData: [],
      total: 0,
      loading: false,

      // 模块详情
      moduleDetailVisible: false,
      moduleDetailData: {
        id: '',
        name: '',
        desc: '',
        status: '',
        parameters: {}
      },

      // 创建模块
      createModuleVisible: false,
      submitCreateModuleLoading: false,
      createModuleData: {
        name: '',
        version: '',
        url: '',
        hash: '',
        desc: '',
        status: '',
        parameters: {}
      },
      // 配置项约束
      createModuleFormRules: {
        name: [
          { required: true, message: '请输入配置名称', trigger: 'blur' },
          { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
        ],
        desc: [
          { required: true, message: '请输入配置描述', trigger: 'blur' },
          { min: 6, max: 30, message: '长度在 6 到 30 个字符', trigger: 'blur' }
        ]
      },

      hasJsonFlag: true, // json是否验证通过
      // 删除按钮弹出框
      popoverVisible: false,
      // 表格多选
      multipleSelection: []
    }
  },
  created() {
    this.getModuleTableData()
  },
  methods: {
    // 查询
    search() {
      this.params.pageNum = 1
      this.getModuleTableData()
    },

    // 获取表格数据
    async getModuleTableData() {
      this.loading = true
      try {
        const { data } = await getModules(this.params)
        this.tableData = data.data
        this.total = data.total
      } finally {
        this.loading = false
      }
    },

    // 查看配置详情
    viewModuleDetail(row) {
      this.configDetailVisible = true
      this.configDetailData.configId = row.ID
      this.configDetailData.configName = row.name
      this.configDetailData.configDesc = row.desc
      this.configDetailData.configStatus = row.status
      this.configDetailData.configContent = JSON.parse(row.content) // 字符串转json对象
    },
    closeCreateModule() {
      this.createModuleVisible = false
    },

    closeModuleDetail() {
      this.moduleDetailVisible = false
    },

    create() {
      this.createModuleVisible = true
    },

    // 提交表单
    submitNewModuleForm() {
      this.$refs['createModuleForm'].validate(async valid => {
        if (valid) {
          this.submitCreateConfigLoading = true
          let msg = ''
          try {
            const { message } = await createModule(this.createConfigData)
            msg = message
          } finally {
            this.submitCreateModuleLoading = false
          }
          this.resetForm()
          this.getModuleTableData()
          this.$message({
            showClose: true,
            message: msg,
            type: 'success'
          })
        } else {
          this.$message({
            showClose: true,
            message: '表单校验失败',
            type: 'error'
          })
          return false
        }
      })
    },

    resetForm() {
      this.createModuleVisible = false
      this.$refs['createModuleForm'].resetFields()
      this.dialogFormData = {
        id: '',
        name: '',
        desc: '',
        status: '',
        content: {}
      }
    },

    onJsonChange(value) {
      this.onJsonSave(value)
    },

    onJsonSave(value) {
      this.resultInfo = value
      this.hasJsonFlag = true
    },
    onError(value) {
      this.hasJsonFlag = false
    },

    // 检查json
    checkJson() {
      if (this.hasJsonFlag === false) {
        alert('json验证失败')
        return false
      } else {
        alert('json验证成功')
        return true
      }
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
          const { message } = await batchDeleteModuleByIds({ ids: configIds })
          msg = message
        } finally {
          this.loading = false
        }

        this.getModuleTableData()
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

    // 单个删除
    async singleDelete(Id) {
      this.loading = true
      let msg = ''
      try {
        const { message } = await batchDeleteConfigByIds({ ids: [Id] })
        msg = message
      } finally {
        this.loading = false
      }

      this.getModuleTableData()
      this.$message({
        showClose: true,
        message: msg,
        type: 'success'
      })
    },

    // 分页
    handleSizeChange(val) {
      this.params.pageSize = val
      this.getModuleTableData()
    },
    handleCurrentChange(val) {
      this.params.pageNum = val
      this.getModuleTableData()
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