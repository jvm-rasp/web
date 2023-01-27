package com.jrasp.admin.server.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.jrasp.admin.server.dto.RbacPermissionRoleAuthDto;
import com.jrasp.admin.server.dto.RbacRoleDto;
import com.jrasp.admin.server.pojo.RbacRole;

import java.util.List;
import java.util.Map;

public interface RbacRoleService extends IService<RbacRole> {
    List<String> getRoleByIds(List<Integer> roleIds);
    Integer createRole(RbacRoleDto dto);
    void updateRole(RbacRoleDto dto);
    void deleteRole(Long id);
    List<RbacRole> getIndex();
    void authPermission(RbacPermissionRoleAuthDto dto);
    List<Integer> getPermissionIdsByRoleId(Integer roleId);

    Map<String,Boolean> getRoleAccess(List<Integer> roleIds);
}
