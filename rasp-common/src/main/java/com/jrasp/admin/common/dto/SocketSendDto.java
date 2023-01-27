package com.jrasp.admin.common.dto;

import lombok.Data;

@Data
public class SocketSendDto {
    private Long userId;
    private String message;
}
