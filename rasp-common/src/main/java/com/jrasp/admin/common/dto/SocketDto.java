package com.jrasp.admin.common.dto;

import lombok.Data;

@Data
public class SocketDto {
    private String event;
    private Long userId;
}
