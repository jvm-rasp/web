package com.jrasp.admin.server.pojo;

import com.baomidou.mybatisplus.annotation.TableField;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class AttackInfo {

    /**
     * 递增id
     */
    private Long id;

    /**
     * 递增id
     */
    private String hostName;

    /**
     * remoteIP 攻击ip
     */
    private String remoteIp;

    /**
     * localIp 本机ip
     */
    private String localIp;

    /**
     * 攻击类型
     */
    private String attackType;

    /**
     * 检查类型
     */
    private String checkType;

    /**
     * 是否阻断
     */
    private Boolean isBlocked;

    /**
     * 危险等级
     *  0-100之间的整数
     */
    private int level;

    /**
     * 危险等级：文字描述
     *  严重、高危、中危、低危
     */
    private String tag;

    /**
     * 处理状态 全部0，未处理1(初始化状态)、处理中2(点击查看详情)、已经处理3
     */
    private int handleStatus;

    /**
     * 处理结果：确认漏洞、误拦截、忽略
     */
    private String handleResult;

    /**
     * 调用栈
     */
    private String stackTrace="";

    /**
     * 请求类型：get、post
     */
    private String httpMethod;

    /**
     * 请求协议
     */
    private String requestProtocol;

    /**
     * 请求路径
     */
    @TableField(value = "request_uri")
    private String requestURI;

    /**
     * 请求参数
     */
    private String requestParameters;

    /**
     * 攻击参数 (payload)
     */
    private String attackParameters;

    /**
     * cookie
     */
    private String cookies;

    /**
     * header
     */
    private String header;

    /**
     * body
     */
    private String body;

    /**
     * 攻击时间,日志上报时间
     */
    private String attackTime;

    /**
     * 创建时间
     */
    private String createTime;

    /**
     * 更新时间
     */
    private String updateTime;

    public AttackInfo(String hostName) {
        this.hostName = hostName;
    }
}
