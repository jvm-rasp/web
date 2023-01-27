package com.jrasp.admin.server.mapper;

import com.jrasp.admin.server.pojo.RbacRolePermission;
import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import org.apache.ibatis.annotations.Param;

import java.util.List;

public interface RbacRolePermissionMapper extends BaseMapper<RbacRolePermission> {
    List<String> getPermissionUrl(@Param("roleIds") List<Integer> roleIds);
}




