package com.jrasp.admin.server.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.pojo.RaspHost;

import java.text.ParseException;

public interface RaspHostService extends IService<RaspHost> {
    /**
     * @param hostName
     * @param ip
     * @param pageNum
     * @param agentMode
     * @param pageSize
     * @return
     */
    PageResult<RaspHost> getIndex(String hostName, String ip, String agentMode, Long pageNum, Long pageSize) ;

    /**
     * @param id
     */
    void deleteHost(Integer id);

    /**
     * @param raspHost
     */
    void updateOrInsert(RaspHost raspHost);

    /**
     * 更新主机id对应的配置id
     *
     * @param hostId
     * @param configId
     */
    void updateHostConfigId(Long hostId, Integer configId);

    /**
     * 更新主机对应的配置id
     *
     * @param hostName
     * @param configId
     */
    void updateHostConfigId(String hostName, Integer configId);

    /**
     * 更新主机对应的配置nacos信息
     *
     * @param hostName
     * @param nacosInfo
     */
    void updateHostNacosInfo(String hostName, String nacosInfo);

    /**
     * 更新主机上Javaagent的配置更新时间
     *
     * @param hostName
     * @param agentConfigUpdateTime
     */
    void updateAgentConfigUpdateTime(String hostName, String agentConfigUpdateTime);


}
