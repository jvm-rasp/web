package com.jrasp.admin.server.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.jrasp.admin.server.pojo.RbacRolePermission;
import com.jrasp.admin.server.service.RbacRolePermissionService;
import com.jrasp.admin.server.mapper.RbacRolePermissionMapper;
import org.springframework.stereotype.Service;

import javax.annotation.Resource;
import java.util.ArrayList;
import java.util.List;
import java.util.stream.Collectors;

@Service
public class RbacRolePermissionServiceImpl extends ServiceImpl<RbacRolePermissionMapper, RbacRolePermission>
        implements RbacRolePermissionService{
    @Resource
    private RbacRolePermissionMapper rbacRolePermissionMapper;
    @Override
    public List<String> getPermissionByRoleIds(List<Integer> roleIds) {
        return rbacRolePermissionMapper.getPermissionUrl(roleIds);
    }

    @Override
    public List<Integer> getPermissionIdsByRoleIds(List<Integer> roleIds) {
        if(roleIds.isEmpty()){
            return new ArrayList<>();
        }
        return lambdaQuery().in(RbacRolePermission::getRoleId, roleIds)
                .select(RbacRolePermission::getPermissionId)
                .list().stream().mapToInt(RbacRolePermission::getPermissionId).boxed().collect(Collectors.toList());
    }
}




