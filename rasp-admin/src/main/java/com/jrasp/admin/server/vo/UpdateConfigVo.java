package com.jrasp.admin.server.vo;

public class UpdateConfigVo {

    private String hostNames;

    private Integer configId;

    public UpdateConfigVo(String hostNames, Integer configId) {
        this.hostNames = hostNames;
        this.configId = configId;
    }

    public UpdateConfigVo() {
    }

    public String getHostNames() {
        return hostNames;
    }

    public void setHostNames(String hostNames) {
        this.hostNames = hostNames;
    }

    public Integer getConfigId() {
        return configId;
    }

    public void setConfigId(Integer configId) {
        this.configId = configId;
    }
}
