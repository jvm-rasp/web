package com.jrasp.admin.server.dto;

import io.swagger.annotations.ApiModelProperty;
import lombok.Data;

import javax.validation.constraints.NotBlank;

@Data
public class RbacPermissionDto {
    private Integer id;
    private Integer parentId;
    @NotBlank(message = "请填写权限名称")
    private String name;
    @ApiModelProperty(value = "图标")
    private String icon;
    @NotBlank(message = "请填写路由")
    private String url;
    @ApiModelProperty(value = "在左侧菜单中隐藏")
    private Integer hideInMenu;
    @ApiModelProperty(value = "隐藏子菜单")
    private Integer hideChildrenInMenu;
    private Long priority;
    private String access;
}
