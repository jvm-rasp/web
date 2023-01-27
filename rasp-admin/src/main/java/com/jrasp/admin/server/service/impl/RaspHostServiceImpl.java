package com.jrasp.admin.server.service.impl;

import cn.hutool.core.util.StrUtil;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.mapper.RaspHostMapper;
import com.jrasp.admin.server.pojo.RaspConfig;
import com.jrasp.admin.server.pojo.RaspHost;
import com.jrasp.admin.server.service.RaspHostService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.text.DateFormat;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.List;
import java.util.Locale;

@Service
@Slf4j
public class RaspHostServiceImpl extends ServiceImpl<RaspHostMapper, RaspHost> implements RaspHostService {

    @Override
    public PageResult<RaspHost> getIndex(String hostName, String ip, String agentMode, Long pageNum, Long pageSize) {
        Page<RaspHost> page = new Page<>(pageNum, pageSize);
        Page<RaspHost> pageResult = lambdaQuery().like(StrUtil.isNotBlank(hostName), RaspHost::getHostName, hostName)
                .like(StrUtil.isNotBlank(ip), RaspHost::getIp, ip)
                .eq(StrUtil.isNotBlank(agentMode), RaspHost::getAgentMode, agentMode)
                .orderByDesc(RaspHost::getId)
                .page(page);
        // 是否在线以服务端时间为准
        SimpleDateFormat dateFormat = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss", Locale.getDefault());
        List<RaspHost> records = pageResult.getRecords();
        for (RaspHost raspHost : records) {
            String heartbeatTime = raspHost.getHeartbeatTime();
            try {
                long heatbeatTime = dateFormat.parse(heartbeatTime).getTime();
                raspHost.setOnline((System.currentTimeMillis() - heatbeatTime) <= 6 * 60 * 1000);
            } catch (Exception e) {
                raspHost.setOnline(true);
                log.error("parse heatbeat time error, heatbeatTime: " + heartbeatTime);
            }

        }
        return PageResult.page(pageResult.getRecords(), page);
    }

    @Override
    public void deleteHost(Integer id) {
        removeById(id);
    }

    @Override
    public void updateOrInsert(RaspHost raspHost) {
        RaspHost raspHostDb = lambdaQuery().eq(RaspHost::getHostName, raspHost.getHostName())
                .select(RaspHost::getId)
                .last("limit 1")
                .one();
        if (raspHostDb != null) {
            raspHost.setId(raspHostDb.getId());
            updateById(raspHost);
        } else {
            save(raspHost);
        }
    }

    @Override
    public void updateHostConfigId(Long hostId, Integer configId) {
        RaspHost raspHost = new RaspHost(hostId, configId);
        updateById(raspHost);
    }

    @Override
    public void updateHostConfigId(String hostName, Integer configId) {
        lambdaUpdate().eq(RaspHost::getHostName, hostName).set(RaspHost::getConfigId, configId).update();
    }

    @Override
    public void updateHostNacosInfo(String hostName, String nacosInfo) {
        lambdaUpdate().eq(RaspHost::getHostName, hostName).set(RaspHost::getNacosInfo, nacosInfo).update();
    }

    @Override
    public void updateAgentConfigUpdateTime(String hostName, String agentConfigUpdateTime) {
        lambdaUpdate().eq(RaspHost::getHostName, hostName).set(RaspHost::getAgentConfigUpdateTime, agentConfigUpdateTime).update();
    }

}
