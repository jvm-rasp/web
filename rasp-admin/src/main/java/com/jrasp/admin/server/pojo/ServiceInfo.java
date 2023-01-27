package com.jrasp.admin.server.pojo;

import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableId;

import com.baomidou.mybatisplus.annotation.TableName;
import com.jrasp.admin.server.mybatishandlers.ListTypeHandler;
import com.jrasp.admin.server.mybatishandlers.ServiceListTypeHandler;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data
@NoArgsConstructor
@AllArgsConstructor
@TableName(value = "service_info")
public class ServiceInfo {
    /**
     * 递增id
     */
    @TableId
    private Integer id;

    /**
     * 服务名称
     */
    private String serviceName;

    /**
     * 服务描述
     */
    private String description;

    /**
     * 开发语言：java、go、c++、python等
     */
    private String language;

    /**
     * 服务负责人
     */
    private String owners;

    /**
     * 节点数量
     */
    private Integer nodeNumber;

    /**
     * 是否公开：0,公开，所有可见、可操作；1：私有，仅服务负责人可以看到，默认公开
     */
    private Boolean isPublic;

    /**
     * 是否可以被外网访问：0,外网；1：内网，默认外网
     */
    private Boolean isInner;

    /**
     * 服务等级：一般服务,重点服务
     */
    private Boolean level;

    /**
     * 操作低峰期：0～23小时，默认为4，低峰期为4:00～4:59
     */
    private Integer lowPeriod;

    /**
     * 服务的组织：信息安全、交易中心等
     */
    private String organization;

    /**
     * 标签，自定义标签
     */
    private String tag;

    /**
     * 服务信息创建时间
     */
    private String createTime;

    /**
     * 服务信息更新时间
     */
    private String updateTime;

}
