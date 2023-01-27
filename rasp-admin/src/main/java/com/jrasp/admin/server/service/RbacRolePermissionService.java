package com.jrasp.admin.server.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.jrasp.admin.server.pojo.RbacRolePermission;

import java.util.List;

public interface RbacRolePermissionService extends IService<RbacRolePermission> {
    List<String> getPermissionByRoleIds(List<Integer> roleIds);
    List<Integer> getPermissionIdsByRoleIds(List<Integer> roleIds);
}
