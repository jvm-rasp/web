package com.jrasp.admin.server.vo;

public class AttackInfoVo {

    // 全部攻击数量
    private int allAttackCnt;

    // 同比值
    private int allAttackDiffValue;

    // 趋势
    private String  allAttackTrend;

    // 最近7天攻击高危总数
    private int weekHighAttackCnt;

    // 同比值
    private int weekHighAttackDiffValue;

    // 趋势
    private String weekHighAttackTrend;

    // 最近7天攻击阻断总数
    private int weekAttackBlockCnt;

    // 同比值
    private int weekAttackBlockDiffValue;

    // 趋势
    private String  weekAttackBlockTrend;

    private String time;

    // 最近7天，每日攻击数量
    BarVo[] itemList;

    public AttackInfoVo() {
    }

    public int getWeekHighAttackCnt() {
        return weekHighAttackCnt;
    }

    public void setWeekHighAttackCnt(int weekHighAttackCnt) {
        this.weekHighAttackCnt = weekHighAttackCnt;
    }

    public int getWeekAttackBlockCnt() {
        return weekAttackBlockCnt;
    }

    public void setWeekAttackBlockCnt(int weekAttackBlockCnt) {
        this.weekAttackBlockCnt = weekAttackBlockCnt;
    }

    public String getTime() {
        return time;
    }

    public void setTime(String time) {
        this.time = time;
    }

    public String getWeekHighAttackTrend() {
        return weekHighAttackTrend;
    }

    public void setWeekHighAttackTrend(String weekHighAttackTrend) {
        this.weekHighAttackTrend = weekHighAttackTrend;
    }

    public String getWeekAttackBlockTrend() {
        return weekAttackBlockTrend;
    }

    public void setWeekAttackBlockTrend(String weekAttackBlockTrend) {
        this.weekAttackBlockTrend = weekAttackBlockTrend;
    }

    public int getWeekHighAttackDiffValue() {
        return weekHighAttackDiffValue;
    }

    public void setWeekHighAttackDiffValue(int weekHighAttackDiffValue) {
        this.weekHighAttackDiffValue = weekHighAttackDiffValue;
    }

    public int getWeekAttackBlockDiffValue() {
        return weekAttackBlockDiffValue;
    }

    public void setWeekAttackBlockDiffValue(int weekAttackBlockDiffValue) {
        this.weekAttackBlockDiffValue = weekAttackBlockDiffValue;
    }

    public int getAllAttackCnt() {
        return allAttackCnt;
    }

    public void setAllAttackCnt(int allAttackCnt) {
        this.allAttackCnt = allAttackCnt;
    }

    public int getAllAttackDiffValue() {
        return allAttackDiffValue;
    }

    public void setAllAttackDiffValue(int allAttackDiffValue) {
        this.allAttackDiffValue = allAttackDiffValue;
    }

    public String getAllAttackTrend() {
        return allAttackTrend;
    }

    public void setAllAttackTrend(String allAttackTrend) {
        this.allAttackTrend = allAttackTrend;
    }

    public BarVo[] getItemList() {
        return itemList;
    }

    public void setItemList(BarVo[] itemList) {
        this.itemList = itemList;
    }
}
