package com.jrasp.admin.server.vo;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class PublishConfigVo {
    /**
     * 配置名称
     */
    private String configName;

    /**
     * 配置内容
     */
    private String configContent;

    /**
     * nacos groupId
     */
    private String groupId;

    /**
     * nacos dataId
     */
    private List<String> dataId;
}
