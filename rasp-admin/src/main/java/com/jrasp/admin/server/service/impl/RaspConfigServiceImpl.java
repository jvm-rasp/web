package com.jrasp.admin.server.service.impl;

import cn.hutool.core.util.StrUtil;
import com.alibaba.fastjson.JSONObject;
import com.alibaba.nacos.api.annotation.NacosInjected;
import com.alibaba.nacos.api.config.ConfigService;
import com.alibaba.nacos.api.exception.NacosException;
import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.mapper.RaspConfigMapper;
import com.jrasp.admin.server.pojo.RaspConfig;
import com.jrasp.admin.server.pojo.RaspModule;
import com.jrasp.admin.server.pojo.RbacUser;
import com.jrasp.admin.server.service.RaspConfigService;
import com.jrasp.admin.server.service.RaspHostService;
import com.jrasp.admin.server.service.RaspModuleService;
import com.jrasp.admin.server.vo.CreateConfigVo;
import com.jrasp.admin.server.vo.QueryVo;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Service
@Slf4j
public class RaspConfigServiceImpl extends ServiceImpl<RaspConfigMapper, RaspConfig> implements RaspConfigService {

    public static final String GROUP_ID = "DEFAULT_GROUP";

    @NacosInjected
    private ConfigService nacosConfigService;

    @Autowired
    private RaspConfigMapper raspConfigMapper;

    @Autowired
    private RaspHostService raspHostService;

    @Autowired
    private RaspModuleService raspModuleService;

    @Override
    public boolean publishConfig(RaspConfig raspConfig) throws NacosException {
        List<String> dataIds = raspConfig.getDataId();
        String groupId = raspConfig.getGroupId();
        if (StringUtils.isEmpty(groupId)) {
            groupId = GROUP_ID;
        }

        for (int i = 0; i < dataIds.size(); i++) {
            boolean isPublishOK = nacosConfigService.publishConfig(dataIds.get(i), groupId, raspConfig.getConfigContent());
            if (!isPublishOK) {
                return false;
            } else {
                raspHostService.updateHostConfigId(dataIds.get(i), raspConfig.getId());
            }
        }
        return true;
    }

    @Override
    public RaspConfig getConfigByName(String configName) {
        return null;
    }

    @Override
    public RaspConfig getConfigById(Integer id) {
        return null;
    }

    @Override
    public void updateConfigByName(RaspConfig raspConfig) {

    }

    @Value("${nacos.config.server-addr}")
    private String nacosAddr;

    @Override
    public void createConfig(CreateConfigVo createConfigVo, RbacUser rbacUser) {
        List<Integer> modules = createConfigVo.getModules();
        List<RaspModule> raspModules = raspModuleService.getRaspModules(modules);
        JSONObject json = new JSONObject();

        // nacos配置写入到config中
        if (nacosAddr != null && nacosAddr.trim().length() > 0) {
            String[] nacosAddrs = nacosAddr.split(",");
            String[] nacosIps = new String[nacosAddrs.length];
            for (int i = 0; i < nacosAddrs.length; i++) {
                String[] ipAndPort = nacosAddrs[i].split(":", 2);
                nacosIps[i] = ipAndPort[0];
            }
            json.put("nacosIps", nacosIps);
        } else {
            throw new RuntimeException("nacos addr is null in application.yml ");
        }
        json.put("agentMode", createConfigVo.getAgentMode());
        json.put("binFileUrl", createConfigVo.getBinFileUrl());
        json.put("binFileHash", createConfigVo.getBinFileHash());

        String logPath = createConfigVo.getLogPath();
        if (logPath != null && !"".equals(logPath.trim())) {
            json.put("logPath", logPath);
        }

        // moduleConfigs
        List<Map<String, Object>> moduleConfigMap = new ArrayList<>();
        json.put("moduleConfigs", moduleConfigMap);
        for (RaspModule raspModule : raspModules) {
            Map<String, Object> tmp = new HashMap();
            tmp.put("moduleName", raspModule.getName());
            tmp.put("downLoadURL", raspModule.getUrl());
            tmp.put("md5", raspModule.getHash());
            List<RaspModule.ParameterItem> parameters = raspModule.getParameters();
            HashMap<String, String> parametersItemMap = new HashMap<>();
            for (RaspModule.ParameterItem item : parameters) {
                parametersItemMap.put(item.getKey(), item.getValue());
            }
            tmp.put("parameters", parametersItemMap);
            moduleConfigMap.add(tmp);
        }

        // agent 全局配置 agentConfigs
        HashMap<String, Object> agentConfigsMap = new HashMap<>();
        json.put("agentConfigs", agentConfigsMap);
        Boolean checkDisable = createConfigVo.getCheckDisable();
        if (checkDisable != null) {
            agentConfigsMap.put("check_disable", checkDisable);
        }
        Integer blockStatusCode = createConfigVo.getBlockStatusCode();
        if (blockStatusCode != null) {
            agentConfigsMap.put("block_status_code", blockStatusCode);
        }
        String redirectUrl = createConfigVo.getRedirectUrl();
        if (redirectUrl != null && !"".equals(redirectUrl.trim())) {
            agentConfigsMap.put("redirect_url", redirectUrl.trim());
        }
        String htmlBlockContent = createConfigVo.getHtmlBlockContent();
        if (htmlBlockContent != null && !"".equals(htmlBlockContent.trim())) {
            agentConfigsMap.put("html_block_content", htmlBlockContent.trim());
        }
        String jsonBlockContent = createConfigVo.getJsonBlockContent();
        if (jsonBlockContent != null && !"".equals(jsonBlockContent.trim())) {
            agentConfigsMap.put("json_block_content", jsonBlockContent.trim());
        }
        String xmlBlockContent = createConfigVo.getXmlBlockContent();
        if (xmlBlockContent != null && !"".equals(xmlBlockContent.trim())) {
            agentConfigsMap.put("xml_block_content", xmlBlockContent);
        }

        RaspConfig raspConfig = new RaspConfig();
        raspConfig.setConfigName(createConfigVo.getConfigName());
        raspConfig.setConfigContent(JSONObject.toJSONString(json));
        raspConfig.setMessage(createConfigVo.getMessage());
        raspConfig.setUsername(rbacUser.getUsername());
        save(raspConfig);
    }

    @Override
    public PageResult<RaspConfig> getIndex(String configName, Integer status, Long pageNum, Long pageSize) {
        Page<RaspConfig> page = new Page<>(pageNum, pageSize);
        Page<RaspConfig> pageResult = lambdaQuery()
                .eq(StrUtil.isNotBlank(configName), RaspConfig::getConfigName, configName)
                .eq(status != null && status < Integer.MAX_VALUE, RaspConfig::getStatus, status)
                .orderByDesc(RaspConfig::getId)
                .page(page);
        return PageResult.page(pageResult.getRecords(), page);
    }

    @Override
    public void deleteConfig(Integer id) {
        removeById(id);
    }

    @Override
    public List<QueryVo> listConfig() {
        QueryWrapper<RaspConfig> query = new QueryWrapper<>();
        query.select(" DISTINCT config_name,id,message ").lambda()
                .eq(RaspConfig::getStatus, 1).orderByDesc(RaspConfig::getId);
        List<RaspConfig> raspConfigs = raspConfigMapper.selectList(query);
        List<QueryVo> result = new ArrayList<>();
        for (RaspConfig config : raspConfigs) {
            String message = config.getMessage();
            QueryVo moduleType = new QueryVo(config.getConfigName() + " （" + message + "）", config.getId());
            result.add(moduleType);
        }
        return result;
    }

    @Override
    public boolean updateConfig(String hostName, Integer id) throws Exception {
        RaspConfig config = raspConfigMapper.selectById(id);
        boolean isPublishOk = nacosConfigService.publishConfig(hostName, GROUP_ID, config.getConfigContent());
        if (!isPublishOk) {
            return false;
        } else {
            raspHostService.updateHostConfigId(hostName, config.getId());
        }
        return true;
    }

    @Override
    public Map<String, Boolean> batchUpdateConfig(String[] hostNames, Integer configId) throws Exception {
        Map<String, Boolean> result = new HashMap<>();
        RaspConfig config = raspConfigMapper.selectById(configId);
        for (String hostName : hostNames) {
            // 手动输入时，防止空格出现
            hostName = hostName.trim();
            boolean isPublishOk = nacosConfigService.publishConfig(hostName, GROUP_ID, config.getConfigContent());
            if (!isPublishOk) {
                result.put(hostName, false);
            } else {
                raspHostService.updateHostConfigId(hostName, config.getId());
                result.put(hostName, true);
            }
        }
        return result;
    }

    @Override
    public synchronized boolean copyConfig(int id) {
        RaspConfig raspConfig = raspConfigMapper.selectById(id);
        String configName = raspConfig.getConfigName() + "_copy";
        raspConfig.setId(null);
        raspConfig.setConfigName(configName);
        return raspConfigMapper.insert(raspConfig) >= 1;
    }

}
