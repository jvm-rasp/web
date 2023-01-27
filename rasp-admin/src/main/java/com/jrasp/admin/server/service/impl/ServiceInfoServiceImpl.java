package com.jrasp.admin.server.service.impl;

import cn.hutool.core.util.StrUtil;
import com.baomidou.mybatisplus.extension.conditions.query.LambdaQueryChainWrapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.jrasp.admin.common.exception.BadHttpRequestException;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.mapper.ServiceInfoMapper;
import com.jrasp.admin.server.pojo.RbacUser;
import com.jrasp.admin.server.pojo.ServiceInfo;
import com.jrasp.admin.server.service.ServiceInfoService;
import org.apache.commons.lang3.StringUtils;
import org.springframework.stereotype.Service;

@Service
public class ServiceInfoServiceImpl extends ServiceImpl<ServiceInfoMapper, ServiceInfo>
        implements ServiceInfoService {

    @Override
    public PageResult<ServiceInfo> getIndex(String serviceName, Integer lowPeriod, Long pageNum, Long pageSize) {
        Page<ServiceInfo> page = new Page<>(pageNum, pageSize);
        Page<ServiceInfo> pageResult = lambdaQuery().
                like(StrUtil.isNotBlank(serviceName), ServiceInfo::getServiceName, serviceName)
                .eq(lowPeriod != null && lowPeriod >= 0, ServiceInfo::getLowPeriod, lowPeriod)
                .orderByDesc(ServiceInfo::getId)
                .page(page);
        return PageResult.page(pageResult.getRecords(), page);
    }

    @Override
    public void createService(ServiceInfo serviceInfo, RbacUser rbacUser) {
        if (StringUtils.isBlank(serviceInfo.getServiceName()) || !checkServiceName(serviceInfo.getServiceName())) {
            throw new BadHttpRequestException("服务名称[" + serviceInfo.getServiceName() + "]已经存在！");
        }
        serviceInfo.setOwners(rbacUser.getUsername());
        save(serviceInfo);
    }

    @Override
    public void updateService(ServiceInfo serviceInfo) {
        updateById(serviceInfo);
    }

    @Override
    public boolean deleteService(Integer id) {
        return removeById(id);
    }

    private boolean checkServiceName(String name) {
        LambdaQueryChainWrapper<ServiceInfo> eq = lambdaQuery().eq(ServiceInfo::getServiceName, name);
        ServiceInfo serviceInfo = eq
                .select(ServiceInfo::getServiceName)
                .one();
        if (serviceInfo != null) {
            return false;
        }
        return true;
    }

}
