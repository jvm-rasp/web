package com.jrasp.admin.server.service.impl;

import cn.hutool.core.date.DateTime;
import cn.hutool.core.date.DateUtil;
import cn.hutool.core.util.StrUtil;
import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.mapper.RaspAttackInfoMapper;
import com.jrasp.admin.server.pojo.AttackInfo;
import com.jrasp.admin.server.service.RaspAttackInfoService;
import com.jrasp.admin.server.vo.AttackInfoVo;
import com.jrasp.admin.server.vo.BarVo;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Arrays;
import java.util.Date;

@Service
@Slf4j
public class RaspAttackInfoServiceImpl extends ServiceImpl<RaspAttackInfoMapper, AttackInfo> implements RaspAttackInfoService {

    @Autowired
    private RaspAttackInfoMapper raspAttackInfoMapper;

    @Override
    public void addNewAttackInfo(AttackInfo attackInfo) {
        save(attackInfo);
    }

    @Override
    public PageResult<AttackInfo> getIndex(String hostName, Integer handleStatus, String localIp, Boolean isBlocked, Long pageNum, Long pageSize) {
        Page<AttackInfo> page = new Page<>(pageNum, pageSize);
        Page<AttackInfo> pageResult = lambdaQuery().eq(StrUtil.isNotBlank(hostName), AttackInfo::getHostName, hostName)
                // 处理状态不为空并且小于Integer.MAX_VALUE
                .eq(handleStatus != null && handleStatus < Integer.MAX_VALUE, AttackInfo::getHandleStatus, handleStatus)
                .eq(StrUtil.isNotBlank(localIp), AttackInfo::getLocalIp, localIp)
                .eq(isBlocked != null, AttackInfo::getIsBlocked, isBlocked)
                .orderByDesc(AttackInfo::getId)
                .page(page);
        return PageResult.page(pageResult.getRecords(), page);
    }

    @Override
    public void deleteAttackInfo(Long id) {
        removeById(id);
    }

    @Override
    public void batchDeleteAttackInfo(Long[] ids) {
        removeByIds(Arrays.asList(ids));
    }

    @Override
    public void mark(long id, int status) {
        lambdaUpdate().eq(AttackInfo::getId, id).set(AttackInfo::getHandleStatus, status).update();
    }

    @Override
    public AttackInfoVo getWeekData() {
        Date now = new Date();
        AttackInfoVo attackInfoVo = new AttackInfoVo();
        attackInfoVo.setTime(DateUtil.formatChineseDate(now, false, true));
        // 最近7天每日漏洞数量
        DateTime start = DateUtil.beginOfDay(now);
        DateTime end = DateUtil.endOfDay(now);
        BarVo[] list = new BarVo[7];
        for (int i = 0; i < 7; i++) {
            QueryWrapper<AttackInfo> query0 = new QueryWrapper<>();
            query0.lambda().between(AttackInfo::getAttackTime, DateUtil.offsetDay(start, -i), DateUtil.offsetDay(end, -i));
            Integer allAttackCnt = raspAttackInfoMapper.selectCount(query0);
            String time = DateUtil.formatDate(DateUtil.offsetDay(start, -i));
            BarVo barVo = new BarVo(time, allAttackCnt);
            list[i] = (barVo);
        }
        attackInfoVo.setItemList(list);
        // 全部漏洞
        QueryWrapper<AttackInfo> query0 = new QueryWrapper<>();
        query0.lambda().ge(AttackInfo::getLevel, 0).between(AttackInfo::getAttackTime, DateUtil.offsetDay(now, -7), now);
        Integer allAttackCnt1 = raspAttackInfoMapper.selectCount(query0);
        query0.lambda().ge(AttackInfo::getLevel, 0).between(AttackInfo::getAttackTime, DateUtil.offsetDay(now, -14), DateUtil.offsetDay(now, -7));
        Integer allAttackCnt2 = raspAttackInfoMapper.selectCount(query0);
        if (allAttackCnt2 != null && allAttackCnt1 != null) {
            if (allAttackCnt1 >= allAttackCnt2) {
                attackInfoVo.setAllAttackTrend("up");
            } else {
                attackInfoVo.setAllAttackTrend("down");
            }
            attackInfoVo.setAllAttackDiffValue(Math.abs(allAttackCnt1 - allAttackCnt2));
        }
        attackInfoVo.setAllAttackCnt(allAttackCnt1);

        // 高危漏洞 level>=90
        QueryWrapper<AttackInfo> query1 = new QueryWrapper<>();
        query1.lambda().ge(AttackInfo::getLevel, 90).between(AttackInfo::getAttackTime, DateUtil.offsetDay(now, -7), now);
        Integer highLevelAttackCnt1 = raspAttackInfoMapper.selectCount(query1);
        query1.lambda().ge(AttackInfo::getLevel, 90).between(AttackInfo::getAttackTime, DateUtil.offsetDay(now, -14), DateUtil.offsetDay(now, -7));
        Integer highLevelAttackCnt2 = raspAttackInfoMapper.selectCount(query1);
        if (highLevelAttackCnt2 != null && highLevelAttackCnt1 != null) {
            if (highLevelAttackCnt1 >= highLevelAttackCnt2) {
                attackInfoVo.setWeekHighAttackTrend("up");
            } else {
                attackInfoVo.setWeekHighAttackTrend("down");
            }
            attackInfoVo.setWeekHighAttackDiffValue(Math.abs(highLevelAttackCnt1 - highLevelAttackCnt2));
        }
        attackInfoVo.setWeekHighAttackCnt(highLevelAttackCnt1);

        //阻断漏洞
        QueryWrapper<AttackInfo> query2 = new QueryWrapper<>();
        query2.lambda().eq(AttackInfo::getIsBlocked, true).between(AttackInfo::getAttackTime, DateUtil.offsetDay(now, -7), now);
        Integer blockAttackCnt1 = raspAttackInfoMapper.selectCount(query2);
        query2.lambda().eq(AttackInfo::getIsBlocked, true).between(AttackInfo::getAttackTime, DateUtil.offsetDay(now, -14), DateUtil.offsetDay(now, -7));
        Integer blockAttackCnt2 = raspAttackInfoMapper.selectCount(query2);
        if (blockAttackCnt1 != null && blockAttackCnt2 != null) {
            if (blockAttackCnt1 >= blockAttackCnt2) {
                attackInfoVo.setWeekAttackBlockTrend("up");
            } else {
                attackInfoVo.setWeekAttackBlockTrend("down");
            }
            attackInfoVo.setWeekAttackBlockCnt(Math.abs(blockAttackCnt1 - blockAttackCnt2));
        }
        attackInfoVo.setWeekAttackBlockCnt(blockAttackCnt1);
        return attackInfoVo;
    }

}
