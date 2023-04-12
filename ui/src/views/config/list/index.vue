<template>
  <div>
    <el-card class="container-card" shadow="always">
      <!-- 条件搜索框 -->
      <el-row>
        <el-form :size="this.$store.getters.size" :inline="true" :model="params" class="demo-form-inline">
          <el-form-item label="配置名称">
            <el-input v-model.trim="params.name" clearable placeholder="名称" @clear="search" />
          </el-form-item>
          <el-form-item label="配置状态">
            <el-select v-model.trim="params.status" clearable placeholder="状态" @change="search" @clear="search">
              <el-option label="可用" value="true" />
              <el-option label="禁用" value="false" />
              <el-option label="全部" value="" />
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
        <el-table-column show-overflow-tooltip sortable prop="name" label="配置名称" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="desc" label="配置描述" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="agentMode" label="接入模式" align="center">
          <template slot-scope="scope">
            <el-tag size="small" :type="getAgentModeColor(scope.row.agentMode)" disable-transitions>
              {{ getAgentModeLabel(scope.row.agentMode) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="isDefault" label="设为默认" align="center">
          <template slot-scope="scope">
            <el-switch
              v-model="scope.row.isDefault"
              active-color="#13ce66"
              inactive-color="grey"
              :active-value="true"
              :inactive-value="false"
              @change="switchDefaultConfig($event,scope.row,scope.$index)"
            />
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="status" label="配置状态" align="center">
          <template slot-scope="scope" label="配置状态" align="center">
            <el-switch
              v-model="scope.row.status"
              active-color="#13ce66"
              inactive-color="grey"
              :active-value="true"
              :inactive-value="false"
              @change="switchChange($event,scope.row,scope.$index)"
            />
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="creator" label="创建人" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="operator" label="操作人" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="UpdatedAt" label="更新时间" align="center" :formatter="dateFormat" />
        <el-table-column show-overflow-tooltip sortable prop="CreatedAt" label="创建时间" align="center" :formatter="dateFormat" />
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

      <!-- 新建配置对话框 -->
      <el-dialog title="创建配置" :visible.sync="createConfigVisible" width="50%">
        <el-form ref="createConfigForm" size="small" :model="bindConfigData" :rules="createConfigFormRules" label-width="100px">
          <el-row>
            <el-col :span="12">
              <el-form-item label="配置名称" prop="name">
                <el-input v-model.trim="bindConfigData.name" placeholder="配置名称" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="接入模式" prop="agentMode">
                <el-radio-group v-model="bindConfigData.agentMode">
                  <el-radio :label="1">动态</el-radio>
                  <el-radio :label="2">静态</el-radio>
                  <el-radio :label="0">关闭</el-radio>
                </el-radio-group>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="模块列表" prop="moduleConfigs">
                <el-checkbox v-model="checkAll" :indeterminate="isIndeterminate" @change="handleCheckAllChange">全选</el-checkbox>
                <div style="margin: 15px 0;" />
                <el-checkbox-group v-model="selectedModuleId" @change="handleCheckedModulesChange">
                  <el-checkbox v-for="module in moduleList" :key="module.ID" :label="module.ID">{{ module.moduleName }}</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item label="阻断反馈链接" label-width="120px" prop="agentConfigs.redirect_url">
                <el-input v-model.trim="bindConfigData.agentConfigs.redirect_url" placeholder="配置名称" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="阻断状态码" label-width="120px" prop="agentConfigs.block_status_code">
                <el-input v-model.trim="bindConfigData.agentConfigs.block_status_code" placeholder="阻断状态码" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="json阻断文本" label-width="120px" prop="agentConfigs.json_block_content">
                <el-input v-model.trim="bindConfigData.agentConfigs.json_block_content" placeholder="json格式的阻断文本" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="xml阻断文本" label-width="120px" prop="agentConfigs.xml_block_content">
                <el-input v-model.trim="bindConfigData.agentConfigs.xml_block_content" placeholder="xml格式的阻断文本" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="html阻断文本" label-width="120px" prop="agentConfigs.html_block_content">
                <el-input v-model.trim="bindConfigData.agentConfigs.html_block_content" placeholder="html格式的阻断文本" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="日志文件路径" label-width="120px" prop="logPath">
                <el-input v-model.trim="bindConfigData.logPath" placeholder="日志文件路径" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="守护进程URL" label-width="120px" prop="binFileUrl">
                <el-input v-model.trim="bindConfigData.binFileUrl" placeholder="守护进程下载链接" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="守护进程hash" label-width="120px" prop="binFileHash">
                <el-input v-model.trim="bindConfigData.binFileHash" placeholder="守护进程文件hash" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="配置描述" label-width="120px" prop="desc">
                <el-input v-model.trim="bindConfigData.desc" placeholder="配置描述" />
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="closeCreateConfig()">关 闭</el-button>
          <el-button size="mini" :loading="submitCreateConfigLoading" type="primary" @click="submitNewConfigForm()">确 定</el-button>
        </div>
      </el-dialog>
      <!-- 新建配置对话框 -->
      <el-dialog title="修改配置" :visible.sync="editConfigVisible" width="50%">
        <el-form ref="editConfigForm" size="small" :model="bindConfigData" :rules="createConfigFormRules" label-width="100px">
          <el-row>
            <el-col :span="12">
              <el-form-item label="配置名称" prop="name">
                <el-input v-model.trim="bindConfigData.name" placeholder="配置名称" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="接入模式" prop="agentMode">
                <el-radio-group v-model="bindConfigData.agentMode">
                  <el-radio :label="1">动态</el-radio>
                  <el-radio :label="2">静态</el-radio>
                  <el-radio :label="0">关闭</el-radio>
                </el-radio-group>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="模块列表" prop="moduleConfigs">
                <el-checkbox v-model="checkAll" :indeterminate="isIndeterminate" @change="handleCheckAllChange">全选</el-checkbox>
                <div style="margin: 15px 0;" />
                <el-checkbox-group v-model="selectedModuleId" @change="handleCheckedModulesChange">
                  <el-checkbox v-for="module in moduleList" :key="module.ID" :label="module.ID">{{ module.moduleName }}</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item label="阻断反馈链接" label-width="120px" prop="agentConfigs.redirect_url">
                <el-input v-model.trim="bindConfigData.agentConfigs.redirect_url" placeholder="配置名称" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="阻断状态码" label-width="120px" prop="agentConfigs.block_status_code">
                <el-input v-model.number="bindConfigData.agentConfigs.block_status_code" placeholder="阻断状态码" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="json阻断文本" label-width="120px" prop="agentConfigs.json_block_content">
                <el-input v-model.trim="bindConfigData.agentConfigs.json_block_content" placeholder="json格式的阻断文本" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="xml阻断文本" label-width="120px" prop="agentConfigs.xml_block_content">
                <el-input v-model.trim="bindConfigData.agentConfigs.xml_block_content" placeholder="xml格式的阻断文本" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="html阻断文本" label-width="120px" prop="agentConfigs.html_block_content">
                <el-input v-model.trim="bindConfigData.agentConfigs.html_block_content" placeholder="html格式的阻断文本" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="日志文件路径" label-width="120px" prop="logPath">
                <el-input v-model.trim="bindConfigData.logPath" placeholder="日志文件路径" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="守护进程URL" label-width="120px" prop="binFileUrl">
                <el-input v-model.trim="bindConfigData.binFileUrl" placeholder="守护进程下载链接" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="守护进程hash" label-width="120px" prop="binFileHash">
                <el-input v-model.trim="bindConfigData.binFileHash" placeholder="守护进程文件hash" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="配置描述" label-width="120px" prop="desc">
                <el-input v-model.trim="bindConfigData.desc" placeholder="配置描述" />
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="editConfigVisible = false">关 闭</el-button>
          <el-button size="mini" :loading="submitCreateConfigLoading" type="primary" @click="editConfigForm()">更 新</el-button>
        </div>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>
import {
  getConfigs,
  createConfig,
  batchDeleteConfigByIds,
  updateConfig,
  updateStatusById,
  updateDefaultById
} from '@/api/config/config'
import vueJsonEditor from 'vue-json-editor'
import { getModules } from '@/api/module/module'
import moment from 'moment/moment'

export default {
  name: 'Config',
  // 注册组件
  // eslint-disable-next-line vue/no-unused-components
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

      // 模块查询参数
      moduleQueryParams: {
        name: '',
        pageNum: 1,
        pageSize: 1000
      },
      // 表格数据
      tableData: [],
      total: 0,
      loading: false,

      // 可用模块
      moduleList: [
      ],

      // 已选择的模块
      selectedModuleId: [],

      // 配置数据绑定
      bindConfigData: {
        ID: '',
        name: '',
        desc: '',
        status: true,
        creator: '',
        operator: '',
        CreatedAt: '',
        UpdatedAt: '',
        agentMode: '2',
        moduleConfigs: [],
        logPath: '',
        agentConfigs: {
          check_disable: false,
          redirect_url: '',
          block_status_code: '',
          json_block_content: '',
          xml_block_content: '',
          html_block_content: ''
        },
        binFileUrl: '',
        binFileHash: ''
      },

      isIndeterminate: true,
      checkAll: false,

      // 编辑配置
      editConfigVisible: false,

      // 创建配置
      createConfigVisible: false,
      submitCreateConfigLoading: false,

      // 配置项约束
      createConfigFormRules: {
        name: [
          { required: true, message: '请输入配置名称', trigger: 'blur' },
          { min: 2, max: 200, message: '长度在 2 到 20 个字符', trigger: 'blur' }
        ],
        desc: [
          { required: true, message: '请输入配置描述', trigger: 'blur' },
          { min: 6, max: 300, message: '长度在 6 到 30 个字符', trigger: 'blur' }
        ]
      },

      hasJsonFlag: true, // json是否验证通过
      // 删除按钮弹出框
      popoverVisible: false,
      // 表格多选
      multipleSelection: [],
      // agent模式
      agentMode: {
        '关闭': 0,
        '动态': 1,
        '静态': 2
      }
    }
  },
  created() {
    this.getConfigTableData()
    this.getModuleListData()
  },
  methods: {
    // 查询
    search() {
      this.params.pageNum = 1
      this.getConfigTableData()
    },

    handleEdit(record) {
      this.bindConfigData = JSON.parse(JSON.stringify(record))
      this.selectedModuleId = []
      record.moduleConfigs.forEach((item) => {
        this.selectedModuleId.push(item.ID)
      })
      const checkedCount = this.selectedModuleId.length
      this.checkAll = checkedCount === this.moduleList.length
      this.isIndeterminate = checkedCount > 0 && checkedCount < this.moduleList.length
      this.editConfigVisible = true
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
          const { message } = await batchDeleteConfigByIds({ ids: [record.ID] })
          msg = message
        } finally {
          this.loading = false
        }

        await this.getConfigTableData()
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

    handleCheckAllChange(value) {
      if (value) {
        this.checkAll = true
        this.moduleList.forEach((item) => {
          this.selectedModuleId.push(item.ID)
        })
      } else {
        this.checkAll = false
        this.selectedModuleId = []
      }
      this.isIndeterminate = false
      this.handleCheckedModulesChange(this.selectedModuleId)
    },

    handleCheckedModulesChange(value) {
      const checkedCount = value.length
      this.checkAll = checkedCount === this.moduleList.length
      this.isIndeterminate = checkedCount > 0 && checkedCount < this.moduleList.length
      this.bindConfigData.moduleConfigs = []
      value.forEach((checkedItem) => {
        const matches = this.moduleList.filter((moduleItem) => {
          return moduleItem.ID === checkedItem
        })
        this.bindConfigData.moduleConfigs = this.bindConfigData.moduleConfigs.concat(matches)
      })
    },

    // 获取表格数据
    async getConfigTableData() {
      this.loading = true
      try {
        const { data } = await getConfigs(this.params)
        this.tableData = data.list
        this.total = data.total
      } finally {
        this.loading = false
      }
    },

    async getModuleListData() {
      const { data } = await getModules(this.moduleQueryParams)
      this.moduleList = data.list
    },

    // 配置创建
    create() {
      this.selectedModuleId = []
      this.bindConfigData = {
        id: '',
        name: '',
        desc: '',
        status: true,
        creator: '',
        operator: '',
        CreatedAt: '',
        UpdatedAt: '',
        agentMode: '',
        moduleConfigs: [],
        logPath: '../logs',
        agentConfigs: {
          check_disable: false,
          redirect_url: 'https://www.jrasp.com/block.html',
          block_status_code: 500,
          json_block_content: '{"controls":[],"status":{"code":"503","text":"您的请求可能存在攻击行为, 已被EpointRASP拦截</br>攻击类型: 【%attack_name%】</br>攻击时间: 【%attack_time%】","url":""}}',
          xml_block_content: '<?xml version="1.0"?><doc><error>true</error><reason>Request blocked by EpointRASP</reason><attack_name>%attack_name%</attack_name><attack_time>%attack_time%</attack_time></doc>',
          html_block_content: '<html class="no-js" style="background-color: transparent"><head><title>安全拦截信息</title><style>.blockquote, body, button, code, dd, div, dl, dt, fieldset, form, h1, h2, h3, h4, h5, h6, input, legend, li, ol, p, pre, td, textarea, th, ul{margin: 0;padding: 0;}body{font-size: 14px;font-family: \'Microsoft YaHei\';}.sys-panel-cover{position: absolute;top: 0;left: 0;width: 100%;height: 100%;background: #000;opacity: 0.6;filter: alpha(opacity=60);}.sys-panel.in-center{ position: absolute;top: 50%;left: 50%;margin-left: -240px;margin-top: -240px;}.sys-panel{border-radius: 5px;border: 1px solid #cdcdcd;background-color: #fff;box-shadow: 0 0 8px rgba(0,0,0,0.8);width: 550px;box-sizing: border-box;padding: 0 30px;padding-bottom:20px;}.sys-panel .panel-title.danger{color: #a94442;}.sys-panel .panel-hd{border-bottom: 1px solid #dcdcdc;}.sys-panel .panel-title{font-size: 18px;line-height: 2.5;}.sys-panel .panel-ft{border-top: 1px solid #dcdcdc;padding: 10px 0;}.sys-panel .alert{padding: 10px;border-radius: 3px;line-height: 1.8;border: 1px solid transparent;height:100px;}.sys-panel .alert-warn{color: #8a6d3b;background-color: #fcf8e3;border-color: #faebcc;}</style></head><body><form id="form1"><div class="sys-panel-cover"></div><div class="sys-panel in-center" id="err-panel"><div class="panel-hd mb20"><h4 class="panel-title danger">安全提示</h4></div><div class="panel-bd mb20"><p class="alert alert-warn">安全提示：您的请求可能存在攻击行为, 已被EpointRASP拦截</br>攻击类型: 【%attack_name%】</br>攻击时间: 【%attack_time%】</p></div></div></form></body></html>'
        },
        binFileUrl: '',
        binFileHash: ''
      }
      const checkedCount = this.selectedModuleId.length
      this.checkAll = checkedCount === this.moduleList.length
      this.isIndeterminate = checkedCount > 0 && checkedCount < this.moduleList.length
      this.createConfigVisible = true
    },
    // 关闭配置创建
    closeCreateConfig() {
      this.createConfigVisible = false
    },

    // 提交表单
    submitNewConfigForm() {
      this.$refs['createConfigForm'].validate(async valid => {
        if (valid) {
          this.submitCreateConfigLoading = true
          let msg = ''
          try {
            const { message } = await createConfig(this.bindConfigData)
            msg = message
          } finally {
            this.submitCreateConfigLoading = false
          }
          this.resetForm()
          await this.getConfigTableData()
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

    // 更新表单
    editConfigForm() {
      this.$refs['editConfigForm'].validate(async valid => {
        if (valid) {
          this.submitCreateConfigLoading = true
          let msg = ''
          try {
            const { message } = await updateConfig(this.bindConfigData)
            msg = message
          } finally {
            this.submitCreateConfigLoading = false
          }
          await this.getConfigTableData()
          this.editConfigVisible = false
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
      this.createConfigVisible = false
      this.$refs['createConfigForm'].resetFields()
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
          const { message } = await batchDeleteConfigByIds({ ids: configIds })
          msg = message
        } finally {
          this.loading = false
        }

        await this.getConfigTableData()
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

    switchChange(e, row, index) {
      updateStatusById({ id: row.ID })
    },

    async switchDefaultConfig(e, row, index) {
      const { message, code } = await updateDefaultById({ id: row.ID, isDefault: e })
      const type = code === 200 ? 'success' : 'error'
      this.$message({
        showClose: true,
        message: message,
        type: type
      })
      await this.getConfigTableData()
    },

    // 分页
    handleSizeChange(val) {
      this.params.pageSize = val
      this.getConfigTableData()
    },
    handleCurrentChange(val) {
      this.params.pageNum = val
      this.getConfigTableData()
    },
    dateFormat(row, column) {
      const date = row[column.property]
      if (date === undefined) {
        return ''
      }
      return moment(date).format('YYYY-MM-DD HH:mm:ss')
    },
    getAgentModeColor(status) {
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
    getAgentModeLabel(status) {
      for (const k in this.agentMode) {
        const val = this.agentMode[k]
        if (val === status) {
          return k
        }
      }
      return '未知'
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
