package com.jrasp.admin.server.vo;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class CreateConfigVo {

    /**
     * 配置名称
     */
    private String configName;

    /**
     * daemon url
     */
    private String binFileUrl;

    /**
     * daemon file hash
     */
    private String binFileHash;

    /**
     * agentMode
     */
    private String agentMode;

    /**
     * 备注消息
     */
    private String message;

    /**
     * 模块信息集合
     */
    private List<Integer> modules;

    private Boolean checkDisable;
    /**
     * 阻断状态码
     */
    private Integer blockStatusCode;

    /**
     * 阻断反馈链接
     */
    private String redirectUrl;

    /**
     * html格式的阻断文本
     */
    private String htmlBlockContent;

    /**
     * json格式的阻断文本
     */
    private String jsonBlockContent;

    /**
     * xml格式的阻断文本
     */
    private String xmlBlockContent;

}
