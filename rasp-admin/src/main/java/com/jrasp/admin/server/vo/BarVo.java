package com.jrasp.admin.server.vo;

public class BarVo {
    private String name;
    private int value;

    public BarVo() {
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public int getValue() {
        return value;
    }

    public void setValue(int value) {
        this.value = value;
    }

    public BarVo(String name, int value) {
        this.name = name;
        this.value = value;
    }
}
