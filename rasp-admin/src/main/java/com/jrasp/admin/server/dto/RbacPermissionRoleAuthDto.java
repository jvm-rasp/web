package com.jrasp.admin.server.dto;

import lombok.Data;

import javax.validation.constraints.NotNull;
import java.util.List;

@Data
public class RbacPermissionRoleAuthDto {
    @NotNull(message = "请选择角色id")
    private Integer roleId;
    @NotNull(message = "请选择要分配的权限！")
    private List<Integer> permissionIds;
}
