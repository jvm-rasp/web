package com.jrasp.admin.server.pojo;

import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableName;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@TableName(value = "rasp_config")
@Data
@AllArgsConstructor
@NoArgsConstructor
public class RaspConfig {

    /**
     * 递增id
     */
    private Integer id;
    /**
     * 策略名
     */
    private String configName;

    /**
     * 策略内容
     */
    private String configContent;

    /**
     * 备注信息
     */
    private String message;

    /**
     * 策略标签：通用策略，1.0版本策略，2.0版本策略
     */
    private String tag;

    /**
     * 是否有效,1 有效；2 失效
     */
    private Integer status;

    /**
     * 创建策略的用户
     */
    private String username;

    /**
     * 策略创建时间
     */
    private String createTime;
    /**
     * 策略更新时间
     */
    private String updateTime;

    /**
     * nacos groupId
     */
    @TableField(exist = false)
    private String groupId;

    /**
     * nacos dataId
     */
    @TableField(exist = false)
    private List<String> dataId;
}
