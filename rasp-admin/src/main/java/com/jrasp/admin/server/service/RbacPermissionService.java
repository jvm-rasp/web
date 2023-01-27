package com.jrasp.admin.server.service;

import com.jrasp.admin.server.dto.RbacPermissionDto;
import com.jrasp.admin.server.pojo.RbacPermission;
import com.baomidou.mybatisplus.extension.service.IService;
import com.jrasp.admin.server.vo.MenuDataItem;

import java.util.List;

public interface RbacPermissionService extends IService<RbacPermission> {
    void createPermission(RbacPermissionDto dto);
    void updatePermission(RbacPermissionDto dto);
    void delete(Long id);
    List<RbacPermission> getIndex();
    List<MenuDataItem> getMenuByIds(List<Integer> ids);
}
