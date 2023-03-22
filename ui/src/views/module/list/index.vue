<template>
  <div>
    <el-card class="container-card" shadow="always">
      <!-- 条件搜索框 -->
      <el-row>
        <el-form size="medium" :inline="true" :model="params" class="demo-form-inline">
          <el-col :span="6">
            <el-form-item label="名称">
              <el-input v-model.trim="params.name" clearable placeholder="名称" @clear="search" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="状态">
              <el-select v-model.trim="params.status" clearable placeholder="状态" @change="search" @clear="search">
                <el-option label="启用" value="true" />
                <el-option label="禁用" value="false" />
                <el-option label="全部" value="" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="创建人">
              <el-input v-model.trim="params.creator" clearable placeholder="创建人" @clear="search" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
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
          </el-col>
        </el-form>
      </el-row>
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
        <el-table-column show-overflow-tooltip sortable prop="moduleName" label="模块名称" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="moduleType" label="模块类型" align="center">
          <template slot-scope="scope">
            {{ getModuleType(scope.row.moduleType) }}
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="desc" label="模块描述" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="moduleVersion" label="版本" align="center">
          <template slot-scope="scope">
            <el-tag size="small" disable-transitions>
              {{ scope.row.moduleVersion }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="status" label="模块状态" align="center">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.status ? 'success':'danger'" disable-transitions>
              {{ scope.row.status ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="operator" label="操作人" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="CreatedAt" label="创建时间" align="center" :formatter="dateFormat" />
        <el-table-column show-overflow-tooltip sortable prop="UpdatedAt" label="更新时间" align="center" :formatter="dateFormat" />
        <el-table-column fixed="right" label="操作" align="center">
          <template slot-scope="scope">
            <el-button type="text" size="medium" @click="handleDelete(scope.row)">删除</el-button>
            <el-button type="text" size="medium" @click="handleEdit(scope.row)">修改</el-button>
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
      <!-- 新增防护模块对话框-->
      <el-dialog title="新增防护模块" :visible.sync="createModuleVisible" width="60%">
        <el-form
          ref="createModuleForm"
          size="small"
          :model="bindModuleData"
          label-width="100px"
        >
          <el-row>
            <el-col :span="7">
              <el-form-item label="模块名称" prop="moduleName">
                <el-input v-model="bindModuleData.moduleName" placeholder="模块名称" />
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="模块版本" prop="moduleVersion">
                <el-input v-model="bindModuleData.moduleVersion" />
              </el-form-item>
            </el-col>
            <el-col :span="11">
              <el-form-item label="模块类型" prop="moduleType">
                <el-radio-group v-model="bindModuleData.moduleType">
                  <el-radio :label="1">hook模块</el-radio>
                  <el-radio :label="2">algorithm模块</el-radio>
                  <el-radio :label="3">其他模块</el-radio>
                </el-radio-group>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="7">
              <el-form-item label="目标中间件" prop="middlewareName">
                <el-input v-model="bindModuleData.middlewareName" placeholder="模块名称" />
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="目标版本号" prop="middlewareVersion">
                <el-input v-model="bindModuleData.middlewareVersion" />
              </el-form-item>
            </el-col>
            <el-col :span="11">
              <el-form-item label="模块状态" prop="status">
                <el-select v-model="bindModuleData.status" clearable placeholder="模块状态">
                  <el-option label="启用" :value="true" />
                  <el-option label="禁用" :value="false" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item label="下载链接" prop="downLoadURL">
                <el-input v-model="bindModuleData.downLoadURL" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="文件hash" prop="md5">
                <el-input v-model="bindModuleData.md5" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="模块描述" prop="desc">
            <el-input v-model="bindModuleData.desc" placeholder="模块描述" />
          </el-form-item>
          <el-form-item label="配置参数" prop="parameters">
            <vue-json-editor
              v-model="bindModuleData.parameters"
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
      <!-- 修改防护模块对话框-->
      <el-dialog title="修改防护模块" :visible.sync="editModuleVisible" width="60%">
        <el-form
          ref="editModuleForm"
          size="small"
          :model="bindModuleData"
          label-width="100px"
        >
          <el-row>
            <el-col :span="7">
              <el-form-item label="模块名称" prop="moduleName">
                <el-input v-model="bindModuleData.moduleName" placeholder="模块名称" />
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="模块版本" prop="moduleVersion">
                <el-input v-model="bindModuleData.moduleVersion" />
              </el-form-item>
            </el-col>
            <el-col :span="11">
              <el-form-item label="模块类型" prop="moduleType">
                <el-radio-group v-model="bindModuleData.moduleType">
                  <el-radio :label="1">hook模块</el-radio>
                  <el-radio :label="2">algorithm模块</el-radio>
                  <el-radio :label="3">其他模块</el-radio>
                </el-radio-group>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="7">
              <el-form-item label="目标中间件" prop="middlewareName">
                <el-input v-model="bindModuleData.middlewareName" placeholder="模块名称" />
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="目标版本号" prop="middlewareVersion">
                <el-input v-model="bindModuleData.middlewareVersion" />
              </el-form-item>
            </el-col>
            <el-col :span="11">
              <el-form-item label="模块状态" prop="status">
                <el-select v-model="bindModuleData.status" clearable placeholder="模块状态">
                  <el-option label="启用" :value="true" />
                  <el-option label="禁用" :value="false" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item label="下载链接" prop="downLoadURL">
                <el-input v-model="bindModuleData.downLoadURL" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="文件hash" prop="md5">
                <el-input v-model="bindModuleData.md5" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="模块描述" prop="desc">
            <el-input v-model="bindModuleData.desc" placeholder="模块描述" />
          </el-form-item>
          <el-form-item label="配置参数" prop="parameters">
            <vue-json-editor
              v-model="bindModuleData.parameters"
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
          <el-button size="mini" @click="closeEditModule()">关 闭</el-button>
          <el-button size="mini" :loading="submitCreateModuleLoading" type="primary" @click="updateModuleForm()">更 新</el-button>
        </div>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>

import vueJsonEditor from 'vue-json-editor'
import { batchDeleteModuleByIds, createModule, deleteModule, getModules, updateModule } from '@/api/module/module'
import moment from 'moment'

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

      // 创建模块
      createModuleVisible: false,
      submitCreateModuleLoading: false,
      bindModuleData: {
        ID: '',
        moduleName: '',
        moduleVersion: '',
        downLoadURL: '',
        md5: '',
        moduleType: 1,
        middlewareName: '',
        middlewareVersion: '',
        tag: '',
        desc: '',
        status: true,
        parameters: {}
      },
      // 编辑模块
      editModuleVisible: false,

      selectedModuleData: null,

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
        this.tableData = data.list
        this.total = data.total
      } finally {
        this.loading = false
      }
    },

    closeCreateModule() {
      this.createModuleVisible = false
    },

    closeEditModule() {
      this.editModuleVisible = false
    },

    create() {
      this.bindModuleData = {
        moduleName: '',
        moduleVersion: '',
        downLoadURL: '',
        md5: '',
        moduleType: 1,
        middlewareName: '',
        middlewareVersion: '',
        desc: '',
        status: true,
        parameters: {}
      }
      this.createModuleVisible = true
    },

    handleEdit(record) {
      this.selectedModuleData = record
      this.bindModuleData = record
      this.editModuleVisible = true
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
          const { message } = await deleteModule({ id: record.ID })
          msg = message
        } finally {
          this.loading = false
        }

        await this.getModuleTableData()
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

    // 提交表单
    submitNewModuleForm() {
      this.$refs['createModuleForm'].validate(async valid => {
        if (valid) {
          this.submitCreateConfigLoading = true
          let msg = ''
          try {
            const { message } = await createModule(this.bindModuleData)
            msg = message
          } finally {
            this.submitCreateModuleLoading = false
          }
          this.resetForm()
          await this.getModuleTableData()
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

    updateModuleForm() {
      this.$refs['editModuleForm'].validate(async valid => {
        if (valid) {
          this.submitCreateConfigLoading = true
          let msg = ''
          try {
            const { message } = await updateModule(this.bindModuleData)
            msg = message
          } finally {
            this.submitCreateModuleLoading = false
          }
          await this.getModuleTableData()
          this.editModuleVisible = false
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

        await this.getModuleTableData()
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
      this.getModuleTableData()
    },
    handleCurrentChange(val) {
      this.params.pageNum = val
      this.getModuleTableData()
    },
    getModuleType(type) {
      if (type === 1) {
        return 'hook模块'
      } else if (type === 2) {
        return 'algorithm模块'
      } else if (type === 3) {
        return '其他模块'
      }
    },
    dateFormat(row, column) {
      const date = row[column.property]
      if (date === undefined) {
        return ''
      }
      return moment(date).format('YYYY-MM-DD HH:mm:ss')
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
