package com.jrasp.admin.server.kafka;

import com.alibaba.fastjson.annotation.JSONField;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.ToString;

@Data
@NoArgsConstructor
@AllArgsConstructor
@ToString
/**
 * agent/module 日志结构一样
 */
public class AgentMessage {

    private String level;

    @JSONField(name = "@timestamp")
    private String ts;

    @JSONField(name = "logger_name")
    private String loggerName;

    private String message;

    @JSONField(name = "log_id")
    private int logId;

    @JSONField(name = "host_name")
    private String hostName;

    @JSONField(name = "thread_name")
    private String threadName;

    @JSONField(name = "level_value")
    private int levelValue;

    @JSONField(name = "pid")
    private String pid;
}
