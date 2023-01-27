package com.jrasp.admin.server.kafka;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.ToString;

@Data
@NoArgsConstructor
@AllArgsConstructor
@ToString
public class DaemonMessage{
    private String level;
    private String ts;
    private String caller;
    private String msg;
    private int logId;
    private String ip;
    private String hostName;
    private int pid;
    private String detail;
}
