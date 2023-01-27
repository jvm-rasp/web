package com.jrasp.admin.server.pojo;

import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import java.io.Serializable;
import lombok.Data;

@TableName(value ="rbac_role_permission")
@Data
public class RbacRolePermission implements Serializable {

    @TableId
    private Integer id;

    private Integer roleId;

    private Integer permissionId;

    @TableField(exist = false)
    private static final long serialVersionUID = 1L;
}