package com.jrasp.admin.server.service;

import com.jrasp.admin.server.pojo.RbacUserRole;
import com.baomidou.mybatisplus.extension.service.IService;

import java.util.List;

public interface RbacUserRoleService extends IService<RbacUserRole> {
    List<Integer> getRoleIdByUserId(Integer userId);
    void attachRole(Integer userId,Integer roleId);
    void detachRole(Integer userId);
}
