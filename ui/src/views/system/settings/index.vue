<template>
  <div>
    <el-card class="container-card" shadow="always">
      <el-tabs tab-position="left" style="height: 100%;">
        <el-tab-pane label="更新管理">
          <el-form ref="autoUpdateForm" label-width="100px" :model="bindConfigData" :rules="autoUpdateRules">
            <el-form-item label="自动更新" prop="autoUpdate">
              <el-tooltip class="item" effect="dark" content="开启后系统将自动更新通用安全策略并上报拦截告警日志" placement="top-start">
                <el-switch
                  v-model="bindConfigData.autoUpdate"
                  active-color="#13ce66"
                  inactive-color="#C0CCDA"
                  :active-value="true"
                  :inactive-value="false"
                  active-text="开启"
                  inactive-text="关闭"
                  :validate-event="true"
                  @change="onAutoUpdateChanged"
                />
              </el-tooltip>
            </el-form-item>
            <el-form-item label="上报地址" prop="reportUrl">
              <el-row>
                <el-col :span="12">
                  <el-input v-model.trim="bindConfigData.reportUrl" clearable placeholder="上报地址" />
                </el-col>
                <el-col :span="12">
                  <el-button type="primary" style="margin-left: 10px;" :size="this.$store.getters.size" :loading="loading" @click="saveSetting('reportUrl', bindConfigData.reportUrl)">保 存</el-button>
                </el-col>
              </el-row>
            </el-form-item>
            <el-form-item label="项目guid" prop="projectGuid">
              <el-col :span="12">
                <el-input v-model.trim="bindConfigData.projectGuid" readonly clearable placeholder="项目guid" />
              </el-col>
            </el-form-item>
            <el-form-item label="项目类型" prop="projectType">
              <el-col :span="12">
                <el-select v-model="bindConfigData.projectType" placeholder="请选择项目类型" @change="saveSetting('projectType', bindConfigData.projectType)">
                  <el-option
                    v-for="item in projectTypeOptions"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  />
                </el-select>
              </el-col>
            </el-form-item>
            <el-form-item label="更新服务器" prop="updateUrl">
              <el-col :span="12">
                <el-input v-model.trim="bindConfigData.updateUrl" clearable placeholder="更新服务器" />
              </el-col>
              <el-col :span="12">
                <el-button type="primary" style="margin-left: 10px;" :size="this.$store.getters.size" :loading="loading" @click="saveSetting('updateUrl', bindConfigData.updateUrl)">保 存</el-button>
              </el-col>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        <el-tab-pane label="配置管理">配置管理</el-tab-pane>
        <el-tab-pane label="角色管理">角色管理</el-tab-pane>
        <el-tab-pane label="定时任务补偿">定时任务补偿</el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script>
import { getProjectGuid, getSettings, updateSetting } from '@/api/system/settings'

export default {
  name: 'Index',
  data() {
    return {
      bindConfigData: {
        autoUpdate: false,
        reportUrl: '',
        projectGuid: '',
        updateUrl: '',
        projectType: ''
      },
      projectTypeOptions: [
        {
          value: 'internal',
          label: '内部OA系统组'
        },
        {
          value: 'self-employed',
          label: '自营系统组'
        },
        {
          value: 'live-project',
          label: '现场项目组'
        }
      ],
      loading: false,
      autoUpdateRules: {
        reportUrl: [
          { required: true, message: '请输入运维监控report地址', trigger: 'blur' },
          { min: 2, max: 200, message: '长度在 2 到 20 个字符', trigger: 'blur' }
        ],
        updateUrl: [
          { required: true, message: '请输入服务器更新地址', trigger: 'blur' },
          { min: 2, max: 200, message: '长度在 2 到 20 个字符', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.getSettingData()
  },
  methods: {
    // 获取表格数据
    async getSettingData() {
      this.loading = true
      try {
        const { data } = await getSettings()
        this.bindConfigData = data.list
      } finally {
        this.loading = false
      }
    },
    onAutoUpdateChanged() {
      if (this.bindConfigData.autoUpdate) {
        this.$refs['autoUpdateForm'].validate(async valid => {
          if (valid) {
            try {
              const { code, message } = await getProjectGuid({ reportUrl: this.bindConfigData.reportUrl })
              if (code === 200) {
                await this.saveSetting('autoUpdate', this.bindConfigData.autoUpdate)
              } else {
                this.$message({
                  showClose: true,
                  message: message,
                  type: 'error'
                })
              }
            } finally {
              await this.getSettingData()
            }
          } else {
            await this.getSettingData()
            this.$message({
              showClose: true,
              message: '表单校验失败',
              type: 'error'
            })
            return false
          }
        })
      } else {
        this.saveSetting('autoUpdate', this.bindConfigData.autoUpdate)
      }
    },
    async saveSetting(name, value) {
      this.loading = true
      try {
        const { code, message } = await updateSetting(name, value)
        const type = code === 200 ? 'success' : 'error'
        this.$message({
          showClose: true,
          message: message,
          type: type
        })
      } finally {
        this.loading = false
        await this.getSettingData()
      }
    }
  }
}
</script>

<style scoped>

</style>
