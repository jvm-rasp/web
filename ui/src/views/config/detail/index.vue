<template>
  <div>
    <el-card class="container-card" shadow="always">
      <!-- 条件搜索框 -->
      <el-row>
        <el-form :size="this.$store.getters.size" :inline="true" :model="params" class="demo-form-inline">
          <el-form-item label="!!!!配置名称">
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
            <el-upload
              class="upload"
              action="./config/import"
              :headers="{Authorization: 'Bearer ' + getUserToken()}"
              :show-file-list="false"
              :before-upload="onBeforeUpload"
              :on-preview="handlePreview"
              :on-remove="handleRemove"
              :before-remove="beforeRemove"
              :on-success="onUploadSuccess"
              multiple
              name="files"
              :limit="1"
              :on-exceed="handleExceed"
              :file-list="fileList"
            >
              <el-button :loading="loading" icon="el-icon-upload2" type="success">导入</el-button>
            </el-upload>
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
        <el-table-column show-overflow-tooltip prop="name" label="名称" align="center" />
        <el-table-column show-overflow-tooltip prop="version" label="版本" align="center">
          <template slot-scope="scope">
            <el-tag :size="$store.getters.size" disable-transitions>
              {{ 'V' + scope.row.version }}
            </el-tag>
            <el-tag v-if="scope.row.isModified" :size="$store.getters.size" style="margin-left: 5px;" type="danger" disable-transitions>
              {{ '配置已更新未推送' }}
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
        <el-table-column show-overflow-tooltip sortable prop="isModified" label="配置状态" align="center">
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
            <el-button type="text" size="medium" @click="handle()">操作</el-button>
            <el-button type="text" size="medium" @click="handlePush(scope.row)">推送</el-button>
            <el-dropdown trigger="click" @command="handleCommand">
              <span class="el-dropdown-link">
                更多<i class="el-icon-arrow-down el-icon--right" />
              </span>
              <el-dropdown-menu slot="dropdown">
                <el-dropdown-item :command="{cmd: 'edit', args: scope.row}">修 改</el-dropdown-item>
                <el-dropdown-item :command="{cmd: 'copy', args: scope.row}">复 制</el-dropdown-item>
                <el-dropdown-item :command="{cmd: 'export', args: scope.row}">导 出</el-dropdown-item>
                <el-dropdown-item :command="{cmd: 'sync', args: scope.row}">同 步</el-dropdown-item>
                <el-dropdown-item :command="{cmd: 'delete', args: scope.row}">删 除</el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
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
              <el-form-item label="配置版本" prop="version">
                <el-select v-model="bindConfigData.version" placeholder="请选择">
                  <el-option
                    v-for="item in configHistoryData"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  />
                </el-select>
                <el-button type="primary" icon="el-icon-search" style="margin-left: 10px;">版本信息</el-button>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item label="防护模式" prop="agentMode">
                <el-radio-group v-model="bindConfigData.agentMode">
                  <el-radio :label="1">开启防护</el-radio>
                  <!--                  <el-radio :label="2">静态</el-radio>-->
                  <el-radio :label="0">关闭防护</el-radio>
                </el-radio-group>
              </el-form-item>
            </el-col>
          </el-row>
          <div v-if="bindConfigData.agentMode === 1">
            <el-row>
              <el-col :span="24">
                <el-form-item label="模块列表" prop="moduleConfigs">
                  <el-checkbox v-model="checkAll" :indeterminate="isIndeterminate" @change="handleCheckAllChange">全选</el-checkbox>
                  <div style="margin: 15px 0;" />
                  <el-checkbox-group v-model="selectedModuleName" @change="handleCheckedModulesChange">
                    <el-col v-for="module in moduleList" :key="module.moduleName" :span="4">
                      <el-checkbox :label="module.moduleName">{{ module.moduleName }}</el-checkbox>
                    </el-col>
                  </el-checkbox-group>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row>
              <el-col :span="24">
                <el-form-item label="拦截选项">
                  <el-switch
                    v-model="showBlockParameters"
                    :active-value="true"
                    :inactive-value="false"
                  />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row v-if="showBlockParameters">
              <el-form-item v-for="(module, index) in bindConfigData.moduleConfigs" :key="module.moduleName" :label="module.moduleName">
                <el-tooltip class="item" effect="dark" content="刷新参数并重置为默认" placement="top-start">
                  <i class="el-icon-refresh" @click="refreshModuleParameter(index, module.moduleName)" />
                </el-tooltip>
                <el-col v-for="name in Object.keys(module.parameters.action)" :key="name" :span="24">
                  <el-radio-group v-model="module.parameters.action[name]" size="mini" fill="#409eff">
                    <el-radio-button border :label="1">拦截攻击</el-radio-button>
                    <el-radio-button border :label="0">记录日志</el-radio-button>
                  </el-radio-group>
                  <span style="margin-left: 10px;">{{ convertParameters2cn(module.parameters.cn_map, name) }}</span>
                </el-col>
              </el-form-item>
            </el-row>
            <el-row>
              <el-form-item label="高级配置">
                <el-switch
                  v-model="showAdvanced"
                  :active-value="true"
                  :inactive-value="false"
                />
              </el-form-item>
            </el-row>
            <el-row v-if="showAdvanced">
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
            <el-row v-if="showAdvanced">
              <el-col :span="24">
                <el-form-item label="json阻断文本" label-width="120px" prop="agentConfigs.json_block_content">
                  <el-input v-model.trim="bindConfigData.agentConfigs.json_block_content" placeholder="json格式的阻断文本" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row v-if="showAdvanced">
              <el-col :span="24">
                <el-form-item label="xml阻断文本" label-width="120px" prop="agentConfigs.xml_block_content">
                  <el-input v-model.trim="bindConfigData.agentConfigs.xml_block_content" placeholder="xml格式的阻断文本" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row v-if="showAdvanced">
              <el-col :span="24">
                <el-form-item label="html阻断文本" label-width="120px" prop="agentConfigs.html_block_content">
                  <el-input v-model.trim="bindConfigData.agentConfigs.html_block_content" placeholder="html格式的阻断文本" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row v-if="showAdvanced">
              <el-col :span="24">
                <el-form-item label="日志文件路径" label-width="120px" prop="logPath">
                  <el-input v-model.trim="bindConfigData.logPath" placeholder="日志文件路径" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row v-if="showAdvanced">
              <el-col :span="12">
                <el-form-item label="RASP Bin更新包" label-width="120px" prop="binFileUrl">
                  <el-input v-model.trim="bindConfigData.raspBinInfo.downLoadUrl" placeholder="RASP Bin更新包下载链接">
                    <i slot="suffix" class="el-input__icon el-icon-more" @click="openUploadForm('raspBinInfo')" />
                  </el-input>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="文件hash" label-width="120px" prop="binFileHash">
                  <el-input v-model.trim="bindConfigData.raspBinInfo.md5" placeholder="RASP Bin更新包hash" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row v-if="showAdvanced">
              <el-col :span="12">
                <el-form-item label="RASP Lib更新包" label-width="120px" prop="binFileUrl">
                  <el-input v-model.trim="bindConfigData.raspLibInfo.downLoadUrl" placeholder="RASP Lib更新包下载链接">
                    <i slot="suffix" class="el-input__icon el-icon-more" @click="openUploadForm('raspLibInfo')" />
                  </el-input>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="文件hash" label-width="120px" prop="binFileHash">
                  <el-input v-model.trim="bindConfigData.raspLibInfo.md5" placeholder="RASP Lib更新包hash" />
                </el-form-item>
              </el-col>
            </el-row>
          </div>
          <el-row>
            <el-col :span="24">
              <el-form-item label="配置描述" prop="desc">
                <el-input v-model.trim="bindConfigData.desc" type="textarea" placeholder="配置描述" />
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="closeCreateConfig()">关 闭</el-button>
          <el-button size="mini" :loading="submitCreateConfigLoading" type="primary" @click="submitNewConfigForm()">保 存</el-button>
        </div>
      </el-dialog>
      <!-- 新建配置对话框 -->
      <el-dialog title="修改配置" :visible.sync="editConfigVisible" width="50%">
        <el-form size="small" :model="bindConfigData" label-width="100px">
          <el-row>
            <el-col :span="12">
              <el-form-item label="配置名称" prop="name">
                <el-input v-model.trim="bindConfigData.name" placeholder="配置名称" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="配置版本" prop="version">
                <el-select v-model="bindConfigData.version" placeholder="请选择" @change="configHistoryChanged">
                  <el-option
                    v-for="item in configHistoryData"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  />
                </el-select>
                <el-button type="primary" icon="el-icon-search" style="margin-left: 10px;" @click="showConfigHistoryDesc(selectConfigHistory)">版本信息</el-button>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item label="修改配置">
                <el-switch
                  v-model="isLocked"
                  active-text="解锁"
                  inactive-text="锁定"
                  :active-value="false"
                  :inactive-value="true"
                />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item label="防护模式" prop="agentMode">
                <el-radio-group v-model="bindConfigData.agentMode" :disabled="isLocked">
                  <el-radio :label="1">开启防护</el-radio>
                  <!--                  <el-radio :label="2">静态</el-radio>-->
                  <el-radio :label="0">关闭防护</el-radio>
                </el-radio-group>
              </el-form-item>
            </el-col>
          </el-row>
          <div v-if="bindConfigData.agentMode === 1">
            <el-row>
              <el-col :span="24">
                <el-form-item label="模块列表" prop="moduleConfigs">
                  <el-checkbox v-model="checkAll" :disabled="isLocked" :indeterminate="isIndeterminate" @change="handleCheckAllChange">全选</el-checkbox>
                  <div style="margin: 15px 0;" />
                  <el-checkbox-group v-model="selectedModuleName" :disabled="isLocked" @change="handleCheckedModulesChange">
                    <el-col v-for="module in moduleList" :key="module.moduleName" :span="6">
                      <el-checkbox :label="module.moduleName">{{ module.moduleName }}</el-checkbox>
                    </el-col>
                  </el-checkbox-group>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row>
              <el-col :span="24">
                <el-form-item label="拦截选项">
                  <el-switch
                    v-model="showBlockParameters"
                    :disabled="isLocked"
                    :active-value="true"
                    :inactive-value="false"
                  />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row v-if="showBlockParameters">
              <el-form-item v-for="(module, index) in bindConfigData.moduleConfigs" :key="module.moduleName" :label="module.moduleName">
                <el-tooltip class="item" effect="dark" content="刷新参数并重置为默认" placement="top-start">
                  <i class="el-icon-refresh" @click="refreshModuleParameter(index, module.moduleName)" />
                </el-tooltip>
                <el-col v-for="name in Object.keys(module.parameters.action)" :key="name" :span="24">
                  <el-radio-group v-model="module.parameters.action[name]" :disabled="isLocked" size="mini" fill="#409eff">
                    <el-radio-button border :label="1">拦截攻击</el-radio-button>
                    <el-radio-button border :label="0">记录日志</el-radio-button>
                  </el-radio-group>
                  <span style="margin-left: 10px;">{{ convertParameters2cn(module.parameters.cn_map, name) }}</span>
                </el-col>
              </el-form-item>
            </el-row>
            <el-row>
              <el-form-item label="高级配置">
                <el-switch
                  v-model="showAdvanced"
                  :disabled="isLocked"
                  :active-value="true"
                  :inactive-value="false"
                />
              </el-form-item>
            </el-row>
            <el-row v-if="showAdvanced">
              <el-col :span="12">
                <el-form-item label="阻断反馈链接" label-width="120px" prop="agentConfigs.redirect_url">
                  <el-input v-model.trim="bindConfigData.agentConfigs.redirect_url" :disabled="isLocked" placeholder="配置名称" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="阻断状态码" label-width="120px" prop="agentConfigs.block_status_code">
                  <el-input v-model.number="bindConfigData.agentConfigs.block_status_code" :disabled="isLocked" placeholder="阻断状态码" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row v-if="showAdvanced">
              <el-col :span="24">
                <el-form-item label="json阻断文本" label-width="120px" prop="agentConfigs.json_block_content">
                  <el-input v-model.trim="bindConfigData.agentConfigs.json_block_content" :disabled="isLocked" placeholder="json格式的阻断文本" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row v-if="showAdvanced">
              <el-col :span="24">
                <el-form-item label="xml阻断文本" label-width="120px" prop="agentConfigs.xml_block_content">
                  <el-input v-model.trim="bindConfigData.agentConfigs.xml_block_content" :disabled="isLocked" placeholder="xml格式的阻断文本" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row v-if="showAdvanced">
              <el-col :span="24">
                <el-form-item label="html阻断文本" label-width="120px" prop="agentConfigs.html_block_content">
                  <el-input v-model.trim="bindConfigData.agentConfigs.html_block_content" :disabled="isLocked" placeholder="html格式的阻断文本" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row v-if="showAdvanced">
              <el-col :span="24">
                <el-form-item label="日志文件路径" label-width="120px" prop="logPath">
                  <el-input v-model.trim="bindConfigData.logPath" :disabled="isLocked" placeholder="日志文件路径" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row v-if="showAdvanced">
              <el-col :span="12">
                <el-form-item label="RASP Bin更新包" label-width="120px" prop="binFileUrl">
                  <el-input v-model.trim="bindConfigData.raspBinInfo.downLoadUrl" :disabled="isLocked" placeholder="RASP Bin更新包下载链接">
                    <i slot="suffix" class="el-input__icon el-icon-more" @click="openUploadForm('raspBinInfo')" />
                  </el-input>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="文件hash" label-width="120px" prop="binFileHash">
                  <el-input v-model.trim="bindConfigData.raspBinInfo.md5" :disabled="isLocked" placeholder="RASP Bin更新包hash" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row v-if="showAdvanced">
              <el-col :span="12">
                <el-form-item label="RASP Lib更新包" label-width="120px" prop="binFileUrl">
                  <el-input v-model.trim="bindConfigData.raspLibInfo.downLoadUrl" :disabled="isLocked" placeholder="RASP Lib更新包下载链接">
                    <i slot="suffix" class="el-input__icon el-icon-more" @click="openUploadForm('raspLibInfo')" />
                  </el-input>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="文件hash" label-width="120px" prop="binFileHash">
                  <el-input v-model.trim="bindConfigData.raspLibInfo.md5" :disabled="isLocked" placeholder="RASP Lib更新包hash" />
                </el-form-item>
              </el-col>
            </el-row>
          </div>
          <el-row>
            <el-col :span="24">
              <el-form-item label="配置描述" prop="desc">
                <el-input v-model.trim="bindConfigData.desc" type="textarea" placeholder="配置描述" />
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="editConfigVisible = false">关 闭</el-button>
          <el-button size="mini" :loading="submitCreateConfigLoading" type="primary" @click="updateConfigForm()">保 存</el-button>
        </div>
      </el-dialog>
      <!-- 选择更新模块-->
      <el-dialog title="选择模块" :visible.sync="selectUploadFileVisible" width="60%">
        <div style="margin-bottom: 30px;">
          <!-- 条件搜索框 -->
          <el-row>
            <el-form :size="this.$store.getters.size" :inline="true" :model="uploadFilesParams" class="demo-form-inline">
              <el-form-item label="文件名称">
                <el-input v-model.trim="uploadFilesParams.fileName" clearable placeholder="文件名称" @clear="filesSearch" />
              </el-form-item>
              <el-form-item label="创建人">
                <el-input v-model.trim="uploadFilesParams.creator" clearable placeholder="创建人" @clear="filesSearch" />
              </el-form-item>
              <el-form-item>
                <el-button :loading="loading" icon="el-icon-search" type="primary" :size="this.$store.getters.size" @click="filesSearch">查询</el-button>
              </el-form-item>
            </el-form>
          </el-row>
          <el-table
            v-loading="loading"
            :data="uploadFileTableData"
            border
            stripe
            style="width: 100%;"
            highlight-current-row
            @current-change="updateCurrentUploadFile"
          >
            <!-- 单选方法,返回选中的表格行 -->
            <el-table-column width="35" :resizable="false" align="center">
              <template slot-scope="scope">
                <el-radio v-model="selectedRadio" :label="scope.row.ID" style="color: #fff;" @change.native="setCurrentUploadFile(scope.row)" />
              </template>
            </el-table-column>
            <el-table-column show-overflow-tooltip sortable prop="fileName" label="文件名称" align="center" />
            <el-table-column show-overflow-tooltip sortable prop="fileHash" label="文件hash" align="center" width="300" />
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
          </el-table>
          <!--分页配置-->
          <el-pagination
            :current-page="uploadFilesParams.pageNum"
            :page-size="uploadFilesParams.pageSize"
            :total="uploadFileTotal"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, prev, pager, next, sizes"
            background
            style="margin-top: 10px;float:right;margin-bottom: 10px;"
            @size-change="handleSizeChange2"
            @current-change="handleCurrentChange2"
          />
        </div>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="closeSelectUploadFileDialog">关 闭</el-button>
          <el-button size="mini" type="primary" @click="addSelectedFile">确 定</el-button>
        </div>
      </el-dialog>
      <!-- 同步模块-->
      <el-dialog title="同步策略" :visible.sync="syncConfigVisible" width="35%">
        <el-form ref="syncConfigForm" :size="this.$store.getters.size">
          <el-row>
            <el-col :span="12">
              <el-form-item label="配置名称" prop="name">
                <el-select v-model="bindSyncConfigData.configId" placeholder="请选择" @change="handleSyncConfigChange">
                  <el-option
                    v-for="item in bindSyncConfigList"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="配置版本" prop="version">
                <el-select v-model="bindSyncConfigData.configVersion" placeholder="请选择">
                  <el-option
                    v-for="item in bindSyncConfigVersionList"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="同步选项" prop="name">
                <el-radio-group v-model="bindSyncConfigData.syncOptions">
                  <el-tooltip class="item" effect="dark" content="不改变您启用/禁用的模块和拦截选项" placement="top-start">
                    <el-radio :label="1">同步配置</el-radio>
                  </el-tooltip>
                  <el-tooltip class="item" effect="dark" content="会改变您启用/禁用的模块和拦截选项（慎用）" placement="top-start">
                    <el-radio :label="2">同步所有</el-radio>
                  </el-tooltip>
                </el-radio-group>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="closeSyncConfig()">关 闭</el-button>
          <el-button size="mini" :loading="loading" type="primary" @click="submitSyncConfig()">同 步</el-button>
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
  updateDefaultById,
  pushConfigById,
  copyConfigById,
  exportConfigById,
  getModules,
  getConfigHistoryByGuid, syncConfigById
} from '@/api/config/config'
import vueJsonEditor from 'vue-json-editor'
import { getUploadFiles } from '@/api/module/module'
import moment from 'moment/moment'
import { getToken } from '@/utils/auth'

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
      },
      // 上传附件查询参数
      uploadFilesParams: {
        fileName: '',
        creator: '',
        fileHash: '',
        mimeType: '',
        pageNum: 1,
        pageSize: 10
      },
      selectedRadio: '',
      selectUploadData: {
        ID: '',
        fileName: '',
        fileHash: '',
        moduleVersion: '',
        creator: '',
        CreatedAt: '',
        UpdatedAt: '',
        downLoadUrl: ''
      },
      // 表格数据
      tableData: [],
      // 上传附件表格数据
      uploadFileTableData: [],
      total: 0,
      uploadFileTotal: 0,
      loading: false,

      // 可用模块
      moduleList: [
      ],

      // 已选择的模块
      selectedModuleName: [],

      // 配置数据绑定
      bindConfigData: {
        ID: '',
        name: '',
        version: '',
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
        raspBinInfo: {
          fileName: '',
          downLoadUrl: '',
          md5: ''
        },
        raspLibInfo: {
          fileName: '',
          downLoadUrl: '',
          md5: ''
        }
      },
      selectKeyName: null,
      isIndeterminate: true,
      checkAll: false,

      // 编辑配置
      editConfigVisible: false,

      // 创建配置
      createConfigVisible: false,
      submitCreateConfigLoading: false,
      selectUploadFileVisible: false,
      // 同步配置
      syncConfigVisible: false,

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
      },
      showAdvanced: false,
      showBlockParameters: false,
      fileList: [],
      configHistoryData: [
        { label: '无', value: 0 }
      ],
      // 选定策略的历史记录列表
      bindConfigHistoryData: [],
      // 当前选定记录数据
      selectConfigHistory: {},
      isCreateNewVersion: false,
      isLocked: true,
      // 绑定同步策略数据
      bindSyncDstData: {},
      bindSyncConfigList: [],
      bindSyncConfigVersionList: [],
      bindSyncConfigData: {
        configId: '',
        configVersion: '',
        syncOptions: 1
      }
    }
  },
  watch: {
    isLocked: {
      handler: function(newValue, oldValue) {
        if (this.editConfigVisible && newValue === false) {
          console.log('create new version')
          this.isCreateNewVersion = true
        }
      }
    }
  },
  created() {
    this.getConfigTableData()
    this.getModuleListData()
  },
  methods: {
    prev() {
      this.$router.go(-1)
    },
    // 查询
    search() {
      this.params.pageNum = 1
      this.getConfigTableData()
    },
    // 上传模块查询
    filesSearch() {
      this.uploadFilesParams.pageNum = 1
      this.getUploadTableData()
    },

    handleCommand(command) {
      this.isCreateNewVersion = false
      switch (command.cmd) {
        case 'delete':
          this.handleDelete(command.args)
          break
        case 'edit':
          this.handleEdit(command.args)
          break
        case 'copy':
          this.handleCopy(command.args)
          break
        case 'export':
          this.handleExport(command.args)
          break
        case 'sync':
          this.handleSync(command.args)
          break
      }
    },

    handleEdit(record) {
      this.loading = true
      this.isLocked = true
      this.showBlockParameters = false
      this.showAdvanced = false
      this.bindConfigData = JSON.parse(JSON.stringify(record))
      this.getConfigHistoryList(record)
      this.editConfigVisible = true
      this.loading = false
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

    handleCopy(record) {
      this.$confirm('此操作将复制当前配置, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        let msg = ''
        try {
          const { message } = await copyConfigById({ id: record.ID })
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
          message: '已取消复制'
        })
      })
    },
    handle() {
      // 进入到详情页面
      //  路径/home对应我在router目录下index.js中定义的path属性值
      this.$router.push({ path: '/config/detail' })
    },
    handlePush(record) {
      this.$confirm('确定推送此配置至所有服务器?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        let msg = ''
        try {
          const { message } = await pushConfigById({ id: record.ID })
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
      })
    },

    handleExport(record) {
      this.$confirm('确定导出此配置?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const fileName = 'export_' + record.name + '.zip'
        try {
          const data = await exportConfigById({ id: record.ID })
          console.log(data)
          if (!data || data.size === 0) {
            this.$message.warning('文件导出失败')
            return
          }
          if (typeof window.navigator.msSaveBlob !== 'undefined') {
            window.navigator.msSaveBlob(new Blob([data]), fileName)
          } else {
            const url = window.URL.createObjectURL(new Blob([data]))
            const link = document.createElement('a')
            link.style.display = 'none'
            link.href = url
            link.download = fileName
            // link.setAttribute('download', fileName)
            document.body.appendChild(link)
            link.click()
            document.body.removeChild(link) // 下载完成移除元素
            window.URL.revokeObjectURL(url) // 释放掉blob对象
          }
        } finally {
          this.loading = false
        }
        await this.getConfigTableData()
      })
    },

    handleSync(record) {
      this.bindSyncConfigData.syncOptions = 1
      this.bindSyncDstData = record
      this.bindSyncConfigList = []
      this.tableData.forEach((item) => {
        if (item.ID === record.ID) {
          return
        }
        this.bindSyncConfigList.push({ label: item.name, value: item.ID })
      })
      this.syncConfigVisible = true
    },

    async handleSyncConfigChange(value) {
      const record = this.tableData.find((item) => {
        return item.ID === value
      })
      const { data } = await getConfigHistoryByGuid({ configGuid: record.rowGuid })
      this.bindSyncConfigVersionList = []
      data.list.forEach((item) => {
        this.bindSyncConfigVersionList.push({ label: 'V' + item.version, value: item.version })
      })
      this.bindSyncConfigData.configVersion = this.bindSyncConfigVersionList[0].value
    },

    handleCheckAllChange(value) {
      if (value) {
        this.checkAll = true
        this.moduleList.forEach((item) => {
          this.selectedModuleName.push(item.moduleName)
        })
      } else {
        this.checkAll = false
        this.selectedModuleName = []
      }
      this.isIndeterminate = false
      this.handleCheckedModulesChange(this.selectedModuleName)
    },

    handleCheckedModulesChange(value) {
      const checkedCount = value.length
      this.checkAll = checkedCount === this.moduleList.length
      this.isIndeterminate = checkedCount > 0 && checkedCount < this.moduleList.length
      this.bindConfigData.moduleConfigs = []
      value.forEach((checkedItem) => {
        const matches = this.moduleList.filter((moduleItem) => {
          return moduleItem.moduleName === checkedItem
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
      this.loading = true
      try {
        const { data } = await getModules(this.moduleQueryParams)
        this.moduleList = data.list === null ? [] : data.list
      } finally {
        this.loading = false
      }
    },

    async getConfigHistoryList(record) {
      this.loading = true
      try {
        const { data } = await getConfigHistoryByGuid({ configGuid: record.rowGuid })
        this.bindConfigHistoryData = data.list
        // 历史版本信息数据
        this.configHistoryData = []
        this.bindConfigHistoryData.forEach((item) => {
          this.configHistoryData.push({ label: 'V' + item.version, value: item.version })
          if (item.version === record.version) {
            this.selectConfigHistory = item
            this.bindConfigData.agentMode = item.agentMode
            this.bindConfigData.moduleConfigs = item.moduleConfigs
            this.bindConfigData.logPath = item.logPath
            this.bindConfigData.agentConfigs = item.agentConfigs
            this.bindConfigData.raspBinInfo = item.raspBinInfo
            this.bindConfigData.raspLibInfo = item.raspLibInfo
          }
        })
        // 当前版本勾选模块数据
        this.selectedModuleName = []
        this.selectConfigHistory.moduleConfigs.forEach((item) => {
          this.selectedModuleName.push(item.moduleName)
          this.moduleList.forEach((item2, index) => {
            if (item.moduleName === item2.moduleName && item.moduleVersion === item2.moduleVersion) {
              this.moduleList[index].parameters = item.parameters
            }
          })
        })
        const checkedCount = this.selectedModuleName.length
        this.checkAll = checkedCount === this.moduleList.length
        this.isIndeterminate = checkedCount > 0 && checkedCount < this.moduleList.length
      } finally {
        this.loading = false
      }
    },

    // 配置创建
    create() {
      this.isLocked = true
      this.showBlockParameters = false
      this.showAdvanced = false
      this.selectedModuleName = []
      this.configHistoryData = [
        { label: '无', value: 0 }
      ]
      this.bindConfigData = {
        id: '',
        name: '',
        version: 0,
        desc: '',
        status: true,
        creator: '',
        operator: '',
        CreatedAt: '',
        UpdatedAt: '',
        agentMode: 1,
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
        raspBinInfo: {
          fileName: '',
          downLoadUrl: '',
          md5: ''
        },
        raspLibInfo: {
          fileName: '',
          downLoadUrl: '',
          md5: ''
        }
      }
      const checkedCount = this.selectedModuleName.length
      this.checkAll = checkedCount === this.moduleList.length
      this.isIndeterminate = checkedCount > 0 && checkedCount < this.moduleList.length
      this.createConfigVisible = true
    },
    // 关闭配置创建
    closeCreateConfig() {
      this.createConfigVisible = false
    },
    // 关闭上传文件选择窗口
    closeSelectUploadFileDialog() {
      this.selectUploadFileVisible = false
      this.selectKeyName = null
    },
    // 关闭配置同步窗口
    closeSyncConfig() {
      this.syncConfigVisible = false
    },
    // 提交表单
    submitNewConfigForm() {
      // 填写版本信息
      this.$prompt('请输入版本更新内容', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      }).then(({ value }) => {
        this.bindConfigData.historyDesc = value
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
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '取消输入'
        })
      })
    },

    // 提交同步配置表单
    async submitSyncConfig() {
      this.loading = true
      let msg = ''
      try {
        const { message } = await syncConfigById(
          {
            srcConfigId: this.bindSyncConfigData.configId,
            srcConfigVersion: this.bindSyncConfigData.configVersion,
            dstConfigId: this.bindSyncDstData.ID,
            syncOptions: this.bindSyncConfigData.syncOptions
          })
        msg = message
      } finally {
        this.loading = false
      }
      await this.getConfigTableData()
      this.syncConfigVisible = false
      this.$message({
        showClose: true,
        message: msg,
        type: 'success'
      })
    },

    // 更新表单
    async updateConfigForm() {
      // 填写版本信息
      if (!this.isCreateNewVersion) {
        this.submitCreateConfigLoading = true
        let msg = ''
        try {
          this.bindConfigData.IsNewVersion = false
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
        this.$prompt('请输入版本更新内容', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消'
        }).then(async({ value }) => {
          this.bindConfigData.historyDesc = value
          this.submitCreateConfigLoading = true
          let msg = ''
          try {
            this.bindConfigData.IsNewVersion = true
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
        }).catch(() => {
          this.$message({
            type: 'info',
            message: '取消输入'
          })
        })
      }
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
    },
    // 上传模块
    openUploadForm(selectKeyName) {
      this.selectKeyName = selectKeyName
      this.selectUploadFileVisible = true
      this.getUploadTableData()
    },
    // 获取上传文件表格数据
    async getUploadTableData() {
      this.loading = true
      try {
        const { data } = await getUploadFiles(this.uploadFilesParams)
        this.uploadFileTableData = data.list
        this.uploadFileTotal = data.total
      } finally {
        this.loading = false
      }
    },
    // 点击选中的行也可以选中单选按钮
    updateCurrentUploadFile(row) {
      if (!row) return
      this.selectedRadio = row.ID
      this.selectUploadData = row
    },
    handleSizeChange2(val) {
      this.uploadFilesParams.pageSize = val
      this.getUploadTableData()
    },
    handleCurrentChange2(val) {
      this.uploadFilesParams.pageNum = val
      this.getUploadTableData()
    },
    // 单选选中
    setCurrentUploadFile(row) {
      this.selectUploadData = row
    },
    async addSelectedFile() {
      if (this.selectedRadio) {
        this.bindConfigData[this.selectKeyName].fileName = this.selectUploadData.fileName
        this.bindConfigData[this.selectKeyName].downLoadUrl = this.selectUploadData.downLoadUrl
        this.bindConfigData[this.selectKeyName].md5 = this.selectUploadData.fileHash
      }
      this.selectUploadFileVisible = false
      this.selectKeyName = null
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
    onBeforeUpload(file) {
      this.loading = true
    },
    onUploadSuccess(response, file, fileList) {
      this.loading = false
      this.fileList = fileList
      const { message, code } = response
      const type = code === 200 ? 'success' : 'error'
      this.$message({
        showClose: true,
        message: message,
        type: type
      })
      this.getConfigTableData()
      this.fileList = []
    },
    getUserToken() {
      return getToken()
    },
    convertParameters2cn(cn_map, en_name) {
      if (cn_map[en_name]) {
        return cn_map[en_name]
      } else {
        return en_name
      }
    },
    showConfigHistoryDesc(record) {
      const time = moment(record.CreatedAt).format('YYYY-MM-DD HH:mm:ss')
      const source = record.source
      const desc = record.desc
      const message = '时间: ' + time + '\n' + '来源: ' + source + '\n' + '描述: ' + desc
      this.$alert('<pre>' + message + '</pre>', '版本信息', {
        confirmButtonText: '确定',
        dangerouslyUseHTMLString: true
      })
    },
    configHistoryChanged(value) {
      const index = this.bindConfigHistoryData.findIndex((item) => {
        if (item.version === value) {
          return true
        }
      })
      const item = this.bindConfigHistoryData[index]
      this.selectConfigHistory = item
      // 当前版本勾选模块数据
      this.selectedModuleName = []
      this.selectConfigHistory.moduleConfigs.forEach((item) => {
        this.selectedModuleName.push(item.moduleName)
        this.moduleList.forEach((item2, index) => {
          if (item.moduleName === item2.moduleName && item.moduleVersion === item2.moduleVersion) {
            this.moduleList[index].parameters = item.parameters
          }
        })
      })
      const checkedCount = this.selectedModuleName.length
      this.checkAll = checkedCount === this.moduleList.length
      this.isIndeterminate = checkedCount > 0 && checkedCount < this.moduleList.length
      // 更新配置参数记录
      this.bindConfigData.agentMode = item.agentMode
      this.bindConfigData.moduleConfigs = item.moduleConfigs
      this.bindConfigData.logPath = item.logPath
      this.bindConfigData.agentConfigs = item.agentConfigs
      this.bindConfigData.raspBinInfo = item.raspBinInfo
      this.bindConfigData.raspLibInfo = item.raspLibInfo
    },
    async refreshModuleParameter(index, moduleName) {
      await this.getModuleListData()
      this.moduleList.forEach((item) => {
        if (item.moduleName === moduleName) {
          this.bindConfigData.moduleConfigs[index].parameters = item.parameters
        }
      })
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

.el-dropdown-link {
  cursor: pointer;
  color: #1890ff;
  margin-left: 5px;
}
.el-icon-arrow-down {
  font-size: 14px;
}

.el-dropdown-menu__item {
  font-size: 14px;
}
</style>
