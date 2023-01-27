package com.jrasp.admin.server.vo;

// 前端查询
public class QueryVo {
    private String label;
    private Integer value;

    public QueryVo(String label, Integer value) {
        this.label = label;
        this.value = value;
    }

    public QueryVo() {
    }

    public String getLabel() {
        return label;
    }

    public void setLabel(String label) {
        this.label = label;
    }

    public Integer getValue() {
        return value;
    }

    public void setValue(Integer value) {
        this.value = value;
    }
}
