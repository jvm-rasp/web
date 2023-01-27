package com.jrasp.admin.common.dto.admin;

import io.swagger.annotations.ApiModelProperty;
import lombok.Data;

import javax.validation.constraints.NotBlank;

@Data
public class LoginDto {
    @ApiModelProperty(value = "用户名",name = "username",required = true)
    @NotBlank(message = "请填写登录名！")
    private String username;

    @ApiModelProperty(value = "登录密码",name = "password",required = true)
    @NotBlank(message = "请输入登录密码！")
    private String password;
}
