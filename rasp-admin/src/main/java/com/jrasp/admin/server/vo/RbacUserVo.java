package com.jrasp.admin.server.vo;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.jrasp.admin.server.pojo.RbacRole;
import lombok.Data;

import java.time.LocalDateTime;

@Data
public class RbacUserVo {
    private Integer id;
    private String username;
    private String avatar;
    private String mobile;
    private RbacRole role;
    private Integer roleId;
    private Integer status;
    @JsonFormat(pattern = "yyyy-MM-dd HH:mm:ss")
    private LocalDateTime createdAt;
    @JsonFormat(pattern = "yyyy-MM-dd HH:mm:ss")
    private LocalDateTime updatedAt;
}
