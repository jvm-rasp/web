package com.jrasp.admin.server.dto;

import io.swagger.annotations.ApiModelProperty;
import lombok.Data;

import javax.validation.constraints.NotBlank;

@Data
public class RbacRoleDto {
    private Integer id;

    @ApiModelProperty(name = "name",value = "角色名称",required = true)
    @NotBlank(message = "请填写角色名称！")
    private String name;

    @ApiModelProperty(name = "code",value = "角色编码",required = true)
    @NotBlank(message = "请填写角色编码")
    private String code;

    @ApiModelProperty(name = "remark",value = "角色描述",required = false)
    @NotBlank(message = "请填写角色描述")
    private String remark;
}
