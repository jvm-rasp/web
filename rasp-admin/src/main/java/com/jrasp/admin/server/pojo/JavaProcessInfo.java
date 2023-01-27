package com.jrasp.admin.server.pojo;

import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
@TableName(value = "rasp_java_info")
public class JavaProcessInfo {
    /**
     * 递增id
     */
    @TableId
    private Long id;

    /**
     * 主机名称
     */
    private String hostName;

    /**
     * java进程命令行信息
     */
    private String cmdlineInfo;

    /**
     * 启动时间
     */
    private String startTime;

    /**
     * 进程号
     */
    private int pid;

    /**
     *
     */
    private String status;

    /**
     * 依赖信息
     */
    private String dependencyInfo;

    /**
     * 依赖信息创建时间
     */
    private String createTime;

    /**
     * 依赖信息更新时间
     */
    private String updateTime;

    public JavaProcessInfo(String hostName, int pid, String dependencyInfo) {
        this.hostName = hostName;
        this.pid = pid;
        this.dependencyInfo = dependencyInfo;
    }

    public JavaProcessInfo(String hostName, String cmdlineInfo, String startTime, int pid, String status) {
        this.hostName = hostName;
        this.cmdlineInfo = cmdlineInfo;
        this.startTime = startTime;
        this.pid = pid;
        this.status = status;
    }

    public JavaProcessInfo(String hostName, int pid, String startTime, String status) {
        this.hostName = hostName;
        this.startTime = startTime;
        this.pid = pid;
        this.status = status;
    }
}