package com.jrasp.admin.server.dto;

import io.swagger.annotations.ApiModelProperty;
import lombok.Data;

import javax.validation.constraints.NotNull;
import javax.validation.constraints.Pattern;

@Data
public class RbacUserDto {
    private Integer id;
    @ApiModelProperty(value = "头像",name = "avatar")
    private String avatar;
    @ApiModelProperty(value = "用户名",name = "username",required = true)
    private String username;
    @ApiModelProperty(value = "手机号",name = "mobile",required = true)
    @Pattern(regexp ="^1[3-9][0-9]\\d{8}$",message = "请填写正确的手机号")
    private String mobile;
    @ApiModelProperty(value = "用户状态",name = "status",required = true)
    private Integer status;

    @ApiModelProperty(value = "角色id",name = "role_id",required = true)
    @NotNull(message = "请选择所属角色！")
    private Integer roleId;

    @ApiModelProperty(value = "登录密码",name = "password")
    private String password;
}
