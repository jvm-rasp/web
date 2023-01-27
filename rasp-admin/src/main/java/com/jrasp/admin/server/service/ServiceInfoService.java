package com.jrasp.admin.server.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.pojo.RbacUser;
import com.jrasp.admin.server.pojo.ServiceInfo;

public interface ServiceInfoService extends IService<ServiceInfo> {
    /**
     *
     * @param serviceName
     * @param lowPeriod
     * @param pageNum
     * @param pageSize
     * @return
     */
    PageResult<ServiceInfo> getIndex(String serviceName,Integer lowPeriod, Long pageNum, Long pageSize);

    /**
     *
     * @param serviceInfo
     * @param rbacUser
     */
    void createService(ServiceInfo serviceInfo, RbacUser rbacUser);

    /**
     *
     * @param serviceInfo
     */
    void updateService(ServiceInfo serviceInfo);

    /**
     *
     * @param id
     * @return
     */
    boolean deleteService(Integer id);
}
