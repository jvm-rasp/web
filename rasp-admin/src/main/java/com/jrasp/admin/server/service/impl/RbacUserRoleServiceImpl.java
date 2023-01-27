package com.jrasp.admin.server.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.jrasp.admin.server.pojo.RbacUserRole;
import com.jrasp.admin.server.service.RbacUserRoleService;
import com.jrasp.admin.server.mapper.RbacUserRoleMapper;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.stream.Collectors;

@Service
public class RbacUserRoleServiceImpl extends ServiceImpl<RbacUserRoleMapper, RbacUserRole>
        implements RbacUserRoleService{
    @Override
    public List<Integer> getRoleIdByUserId(Integer userId) {
        return lambdaQuery().eq(RbacUserRole::getUserId, userId)
                .select(RbacUserRole::getRoleId)
                .list()
                .stream().mapToInt(RbacUserRole::getRoleId).boxed().collect(Collectors.toList());
    }

    @Override
    @Transactional
    public void attachRole(Integer userId, Integer roleId) {
        lambdaUpdate().eq(RbacUserRole::getUserId,userId)
                .remove();
        RbacUserRole userRole = new RbacUserRole();
        userRole.setRoleId(roleId);
        userRole.setUserId(userId);
        save(userRole);
    }

    @Override
    public void detachRole(Integer userId) {
        lambdaUpdate().eq(RbacUserRole::getUserId,userId)
                .remove();
    }
}




