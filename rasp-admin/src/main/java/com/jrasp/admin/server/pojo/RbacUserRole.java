package com.jrasp.admin.server.pojo;

import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import java.io.Serializable;
import lombok.Data;

@TableName(value ="rbac_user_role")
@Data
public class RbacUserRole implements Serializable {

    @TableId
    private Integer id;

    private Integer userId;

    private Integer roleId;

    @TableField(exist = false)
    private static final long serialVersionUID = 1L;
}