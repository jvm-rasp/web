package com.jrasp.admin.server.pojo;

import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableName;
import lombok.AllArgsConstructor;
import lombok.Data;

import java.io.Serializable;

@TableName(value = "rasp_host")
@Data
@AllArgsConstructor
public class RaspHost implements Serializable {
    /**
     * 递增id
     */
    private Long id;

    /**
     * 实例名称
     */
    private String hostName;

    /**
     * ip地址
     */
    private String ip;

    /**
     * 配置id
     */
    private Integer configId;

    /**
     * agent接入模式
     * disable、static、dynamic
     */
    private String agentMode;

    /**
     * 最近一次的心跳时间
     */
    private String heartbeatTime;

    /**
     * 是否在线
     */
    @TableField(exist = false)
    private boolean online = true;

    /**
     * java agent信息
     */
    private String agentInfo;

    /**
     * installDir 可执行文件安装的绝对路径
     */
    private String installDir;

    /**
     * version jrasp 版本
     */
    private String version;

    /**
     * exeFileHash 可执行文件的hash
     */
    private String exeFileHash;

    /**
     * osType 实例操作系统类型
     */
    private String osType;

    /**
     * totalMem 实例总内存,单位GB
     */
    private Integer totalMem;

    /**
     * cpu_counts 实例逻辑cpu数量
     */
    private Integer cpuCounts;

    /**
     * free_disk
     */
    private Integer freeDisk;

    /**
     * buildDateTime 编译时间
     */
    private String buildDateTime;

    /**
     * buildGitBranch 分支
     */
    private String buildGitBranch;

    /**
     * buildGitCommit commit
     */
    private String buildGitCommit;

    /**
     * nacos 初始化信息，json字符串
     */
    private String nacosInfo;

    /**
     * agent 配置更新时间
     */
    private String agentConfigUpdateTime;

    /**
     * 创建时间
     */
    private String createTime;

    /**
     * 更新时间
     */
    private String updateTime;


    @TableField(exist = false)
    private static final long serialVersionUID = 1L;

    public RaspHost(Long id) {
        this.id = id;
    }

    public RaspHost(String hostName, String ip, String heartbeatTime) {
        this.hostName = hostName;
        this.ip = ip;
        this.heartbeatTime = heartbeatTime;
    }

    public RaspHost(String hostName, String heartbeatTime) {
        this.hostName = hostName;
        this.heartbeatTime = heartbeatTime;
    }

    public RaspHost(Long id, Integer configId) {
        this.id = id;
        this.configId = configId;
    }

    public RaspHost(Long id, String hostName, String ip, Integer configId, String agentMode, String heartbeatTime, String agentInfo, String installDir, String version, String exeFileHash, String osType, Integer totalMem, Integer cpuCounts, Integer freeDisk, String buildDateTime, String buildGitBranch, String buildGitCommit, String nacosInfo, String agentConfigUpdateTime, String createTime, String updateTime) {
        this.id = id;
        this.hostName = hostName;
        this.ip = ip;
        this.configId = configId;
        this.agentMode = agentMode;
        this.heartbeatTime = heartbeatTime;
        this.agentInfo = agentInfo;
        this.installDir = installDir;
        this.version = version;
        this.exeFileHash = exeFileHash;
        this.osType = osType;
        this.totalMem = totalMem;
        this.cpuCounts = cpuCounts;
        this.freeDisk = freeDisk;
        this.buildDateTime = buildDateTime;
        this.buildGitBranch = buildGitBranch;
        this.buildGitCommit = buildGitCommit;
        this.nacosInfo = nacosInfo;
        this.agentConfigUpdateTime = agentConfigUpdateTime;
        this.createTime = createTime;
        this.updateTime = updateTime;
    }
}
