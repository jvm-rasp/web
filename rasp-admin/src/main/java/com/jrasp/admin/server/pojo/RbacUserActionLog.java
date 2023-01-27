package com.jrasp.admin.server.pojo;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import java.io.Serializable;
import java.time.LocalDateTime;
import lombok.Data;


@TableName(value ="rbac_user_action_log")
@Data
public class RbacUserActionLog implements Serializable {

    @TableId(type = IdType.AUTO)
    private Long id;

    private Integer userId;

    private String operation;

    private String ip;

    private LocalDateTime createdAt;

    @TableField(exist = false)
    private static final long serialVersionUID = 1L;
}