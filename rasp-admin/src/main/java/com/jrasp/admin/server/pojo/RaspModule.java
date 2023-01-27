package com.jrasp.admin.server.pojo;

import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableName;
import com.jrasp.admin.server.mybatishandlers.RaspModuleListTypeHandler;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@TableName(value = "rasp_module", autoResultMap = true)
@Data
@AllArgsConstructor
@NoArgsConstructor
public class RaspModule {

    /**
     * 递增id
     */
    private Integer id;

    /**
     * 模块名称
     */
    private String name;

    /**
     * 模块下载链接
     */
    private String url;

    /**
     * 模块hash
     */
    private String hash;

    /**
     * 模块类型
     */
    private String type;

    /**
     * 模块版本
     */
    private String version;

    /**
     * 模块对应的中间件名称
     */
    private String middlewareName;

    /**
     * 模块对应的中间件版本
     */
    private String middlewareVersion;

    /**
     * 备注信息
     */
    private String message;

    /**
     * 标签
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
     * 参数列表
     */
    @TableField(value = "parameters",
            typeHandler = RaspModuleListTypeHandler.class)
    private List<ParameterItem> parameters;

    @Data
    @NoArgsConstructor
    @AllArgsConstructor
    public static class ParameterItem {
        /**
         * 参数名称
         */
        private String key;

        /**
         * 参数值
         */
        private String value;

        /**
         * 备注信息
         */
        private String info;
    }
}