package com.jrasp.admin.server.dto;

import lombok.Data;

import javax.validation.constraints.NotNull;

@Data
public class StatusDto {

    @NotNull(message = "请填写id！")
    private Integer id;

    @NotNull(message = "请填写状态！")
    private Integer status;
}
