package com.jrasp.admin.server.pojo;

import com.alibaba.fastjson.JSONObject;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class EmailAccount {
    /**
     * 是否启用
     */
    private Boolean enable;
    private String host;
    private Integer port;
    private Boolean auth;
    private String user;
    private String pass;
    private String from;
    private String receiver; // 换行符号分割，每行一个
}
