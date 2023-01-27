package com.jrasp.admin.server.service.impl;

import cn.hutool.core.util.StrUtil;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.mapper.JavaProcessInfoMapper;
import com.jrasp.admin.server.pojo.JavaProcessInfo;
import com.jrasp.admin.server.service.JavaProcessInfoService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@Slf4j
public class JavaProcessInfoServiceImpl extends ServiceImpl<JavaProcessInfoMapper, JavaProcessInfo>
        implements JavaProcessInfoService {

    @Override
    public PageResult<JavaProcessInfo> getIndex(String hostName, Long pageNum, Long pageSize, String status) {
        Page<JavaProcessInfo> page = new Page<>(pageNum, pageSize);
        Page<JavaProcessInfo> pageResult = lambdaQuery().eq(StrUtil.isNotBlank(hostName), JavaProcessInfo::getHostName, hostName)
                .eq(StrUtil.isNotBlank(status), JavaProcessInfo::getStatus, status)
                .orderByDesc(JavaProcessInfo::getId)
                .page(page);
        return PageResult.page(pageResult.getRecords(), page);
    }

    @Override
    public void updateOrInsert(JavaProcessInfo javaProcessInfo) {
        JavaProcessInfo javaProcessInfoDb = lambdaQuery()
                .eq(JavaProcessInfo::getHostName, javaProcessInfo.getHostName())
                .eq(JavaProcessInfo::getPid, javaProcessInfo.getPid())
                .select(JavaProcessInfo::getId)
                .last("limit 1")
                .one();
        if (javaProcessInfoDb != null) {
            javaProcessInfo.setId(javaProcessInfoDb.getId());
            updateById(javaProcessInfo);
        } else {
            save(javaProcessInfo);
        }
    }

    @Override
    public void removeByHostNameAndPid(String hostName, int pid) {
        lambdaUpdate().eq(JavaProcessInfo::getHostName, hostName)
                .eq(JavaProcessInfo::getPid, pid).remove();
    }

    @Override
    public void removeByHostName(String hostName) {
        lambdaUpdate().eq(JavaProcessInfo::getHostName, hostName)
                .remove();
    }

    @Override
    public void insertList(List<JavaProcessInfo> list) {
        saveBatch(list);
    }

    @Override
    public void updateStatus(JavaProcessInfo javaProcessInfo) {
        lambdaUpdate().eq(JavaProcessInfo::getHostName, javaProcessInfo.getHostName())
                .eq(JavaProcessInfo::getPid, javaProcessInfo.getPid())
                .update(javaProcessInfo);
    }
}
