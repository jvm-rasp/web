package com.jrasp.admin.server.service.impl;

import cn.hutool.core.util.StrUtil;
import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.dto.StatusDto;
import com.jrasp.admin.server.mapper.RaspModuleMapper;
import com.jrasp.admin.server.pojo.RaspModule;
import com.jrasp.admin.server.pojo.RbacUser;
import com.jrasp.admin.server.service.RaspModuleService;
import com.jrasp.admin.server.vo.QueryVo;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;

@Service
@Slf4j
public class RaspModuleServiceImpl extends ServiceImpl<RaspModuleMapper, RaspModule> implements RaspModuleService {

    @Autowired
    private RaspModuleMapper raspModuleMapper;

    @Override
    public PageResult<RaspModule> getIndex(String name, Integer status, Long pageNum, Long pageSize) {
        Page<RaspModule> page = new Page<>(pageNum, pageSize);
        Page<RaspModule> pageResult = lambdaQuery().eq(StrUtil.isNotBlank(name), RaspModule::getName, name)
                .eq(status != null && status != 0, RaspModule::getStatus, status)
                .orderByDesc(RaspModule::getId)
                .page(page);
        return PageResult.page(pageResult.getRecords(), page);
    }

    @Override
    public void deleteModule(Integer id) {
        removeById(id);
    }

    @Override
    public void createModule(RaspModule raspModule, RbacUser user) {
        raspModule.setUsername(user.getUsername());
        save(raspModule);
    }

    @Override
    public RaspModule getModule(Integer id) {
        return getById(id);
    }

    @Override
    public void updateModule(RaspModule raspModule) {
        updateById(raspModule);
    }

    @Override
    public void setModuleStatus(StatusDto status) {
        lambdaUpdate().eq(RaspModule::getId, status.getId())
                .set(RaspModule::getStatus, status.getStatus())
                .set(RaspModule::getUpdateTime, LocalDateTime.now())
                .update();
    }

    @Override
    public List<QueryVo> getModuleType() {
        QueryWrapper<RaspModule> query = new QueryWrapper<>();
        query.select(" DISTINCT name,id ").lambda()
                .eq(RaspModule::getStatus, 1) // 有效
                .isNotNull(RaspModule::getVersion);
        List<RaspModule> raspModules = raspModuleMapper.selectList(query);
        List<QueryVo> result = new ArrayList<>();
        for (RaspModule raspModule : raspModules) {
            QueryVo moduleType = new QueryVo(raspModule.getName(), raspModule.getId());
            result.add(moduleType);
        }
        return result;
    }

    @Override
    public List<RaspModule> getRaspModules(List<Integer> ids) {
        return listByIds(ids);
    }
}
