package com.jrasp.admin.server.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.dto.StatusDto;
import com.jrasp.admin.server.pojo.RaspModule;
import com.jrasp.admin.server.pojo.RbacUser;
import com.jrasp.admin.server.vo.QueryVo;

import java.util.List;

public interface RaspModuleService extends IService<RaspModule> {

    /**
     * @param name     配置名称
     * @param status   配置状态
     * @param pageNum
     * @param pageSize
     * @return
     */
    PageResult<RaspModule> getIndex(String name, Integer status, Long pageNum, Long pageSize);

    /**
     * 删除模块信息
     *
     * @param id
     */
    void deleteModule(Integer id);

    /**
     * 插件模块
     *
     * @param createModuleVo
     * @param user
     */
    void createModule(RaspModule createModuleVo, RbacUser user);


    /**
     * 获取模块信息
     *
     * @param id
     */
    RaspModule getModule(Integer id);

    void updateModule(RaspModule raspModule);

    void setModuleStatus(StatusDto status);

    List<QueryVo> getModuleType();

    List<RaspModule> getRaspModules(List<Integer> ids);

}
