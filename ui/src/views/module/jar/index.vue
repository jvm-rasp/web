<template>
  <div>
    <el-card class="container-card" shadow="always">
      <!-- 条件搜索框 -->
      <el-row>
        <el-form size="medium" :inline="true" :model="params" class="demo-form-inline">
          <el-col :span="6">
            <el-form-item label="模块名称">
              <el-input v-model.trim="params.moduleName" clearable placeholder="模块名称" @clear="search" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="模块hash">
              <el-input v-model.trim="params.fileHash" clearable placeholder="模块hash" @clear="search" />
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
              <el-button :loading="loading" icon="el-icon-plus" type="warning" @click="upload">批量上传</el-button>
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
        <el-table-column show-overflow-tooltip sortable prop="fileHash" label="模块hash" align="center" width="300" />
        <el-table-column show-overflow-tooltip sortable prop="moduleVersion" label="版本" align="center" width="100">
          <template slot-scope="scope">
            <el-tag size="small" disable-transitions>
              {{ scope.row.moduleVersion }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="creator" label="创建人" align="center" />
        <el-table-column
          show-overflow-tooltip
          sortable
          prop="CreatedAt"
          label="创建时间"
          align="center"
          :formatter="dateFormat"
        />
        <el-table-column
          show-overflow-tooltip
          sortable
          prop="UpdatedAt"
          label="更新时间"
          align="center"
          :formatter="dateFormat"
        />
        <el-table-column fixed="right" label="操作" align="center">
          <template slot-scope="scope">
            <el-button type="text" size="medium" @click="handleDelete(scope.row)">删除</el-button>
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
      <el-dialog title="上传防护模块" :visible.sync="uploadModuleVisible" :before-close="onUploadDialogClosed" width="60%">
        <el-upload
          class="upload"
          action="./file/upload"
          :headers="{Authorization: 'Bearer ' + getUserToken()}"
          :show-file-list="true"
          :on-preview="handlePreview"
          :on-remove="handleRemove"
          :before-remove="beforeRemove"
          :on-success="onUploadSuccess"
          multiple
          accept=".jar"
          name="files"
          :limit="10"
          :on-exceed="handleExceed"
          :file-list="fileList"
        >
          <el-button size="medium" type="primary">点击上传</el-button>
          <div slot="tip" class="el-upload__tip">只能上传jar文件</div>
        </el-upload>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>
import {
  batchDeleteFilesByIds,
  getUploadFiles
} from '@/api/module/module'
import moment from 'moment'
import { getToken } from '@/utils/auth'

export default {
  name: 'Module',
  data() {
    return {
      // 查询参数
      params: {
        moduleName: '',
        fileHash: '',
        creator: '',
        pageNum: 1,
        pageSize: 10
      },

      // 表格数据
      tableData: [],
      total: 0,
      loading: false,

      // 上传模块
      uploadModuleVisible: false,
      submitCreateModuleLoading: false,
      bindModuleData: {
        ID: '',
        moduleName: '',
        moduleVersion: '',
        downLoadURL: '',
        md5: '',
        moduleType: 1,
        tag: '',
        desc: '',
        status: '',
        parameters: {}
      },
      fileList: [],
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
        const { data } = await getUploadFiles(this.params)
        this.tableData = data.list
        this.total = data.total
      } finally {
        this.loading = false
      }
    },

    onUploadDialogClosed() {
      this.getModuleTableData()
      this.uploadModuleVisible = false
    },

    upload() {
      this.fileList = []
      this.uploadModuleVisible = true
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
          const { message } = await batchDeleteFilesByIds({ ids: [record.ID] })
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
          const { message } = await batchDeleteFilesByIds({ ids: configIds })
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
    dateFormat(row, column) {
      const date = row[column.property]
      if (date === undefined) {
        return ''
      }
      return moment(date).format('YYYY-MM-DD HH:mm:ss')
    },
    handleRemove(file, fileList) {
      console.log(file, fileList)
    },
    handlePreview(file) {
      console.log(file)
    },
    handleExceed(files, fileList) {
      this.$message.warning(`当前限制选择 10 个文件，本次选择了 ${files.length} 个文件，共选择了 ${files.length + fileList.length} 个文件`)
    },
    beforeRemove(file, fileList) {
      this.$alert('文件已上传，请在管理页面删除')
      return false
    },
    onUploadSuccess(response, file, fileList) {
      this.fileList = fileList
    },
    getUserToken() {
      return getToken()
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
