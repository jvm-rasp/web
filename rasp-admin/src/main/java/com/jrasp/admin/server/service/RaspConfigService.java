package com.jrasp.admin.server.service;

import com.alibaba.nacos.api.exception.NacosException;
import com.baomidou.mybatisplus.extension.service.IService;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.pojo.RaspConfig;
import com.jrasp.admin.server.pojo.RbacUser;
import com.jrasp.admin.server.vo.CreateConfigVo;
import com.jrasp.admin.server.vo.QueryVo;
import com.jrasp.admin.server.vo.UpdateConfigVo;

import java.util.List;
import java.util.Map;

public interface RaspConfigService extends IService<RaspConfig> {
    /**
     * 发布配置，往nacos写配置
     *
     * @param raspConfig
     * @return
     * @throws NacosException
     */
    boolean publishConfig(RaspConfig raspConfig) throws NacosException;

    /**
     * 获取配置
     *
     * @param configName 配置名称
     * @return
     */
    RaspConfig getConfigByName(String configName);

    /**
     * 获取配置
     *
     * @param id 配置id
     * @return
     */
    RaspConfig getConfigById(Integer id);

    /**
     * 更新配置
     *
     * @param raspConfig
     */
    void updateConfigByName(RaspConfig raspConfig);

    /**
     * 新增配置
     *
     * @param raspConfig
     * @param rbacUser
     */
    void createConfig(CreateConfigVo raspConfig, RbacUser rbacUser);

    /**
     * 分页获取配置
     *
     * @param configName
     * @param status
     * @param pageNum
     * @param pageSize
     * @return
     */
    PageResult<RaspConfig> getIndex(String configName, Integer status, Long pageNum, Long pageSize);

    void deleteConfig(Integer id);

    List<QueryVo> listConfig();

    boolean updateConfig(String hostName, Integer configId) throws Exception;

    Map<String,Boolean> batchUpdateConfig(String[] hostNames, Integer configId) throws Exception;

    boolean copyConfig(int id);

}
