package com.jrasp.admin.server.pojo;

import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import java.io.Serializable;
import java.time.LocalDateTime;
import java.util.List;

import lombok.Data;

@TableName(value ="rbac_permission")
@Data
public class RbacPermission implements Serializable {
    @TableId
    private Integer id;

    private Integer parentId;

    private String name;

    private String icon;

    private String url;

    private Integer priority;

    private Integer hideInMenu;

    private Integer hideChildrenInMenu;

    private LocalDateTime createdAt;

    private LocalDateTime updatedAt;

    @TableField(exist = false)
    private static final long serialVersionUID = 1L;
    @TableField(exist = false)
    private List<RbacPermission> children;
    @TableField(exist = false)
    private List<String> parentKeys;
    @TableField(exist = false)
    private Boolean checked;
}