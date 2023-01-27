package com.jrasp.admin.server.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.pojo.JavaProcessInfo;

import java.util.List;

public interface JavaProcessInfoService extends IService<JavaProcessInfo> {
    /**
     * @param hostName
     * @param pageNum
     * @param pageSize
     * @param status
     * @return
     */
    PageResult<JavaProcessInfo> getIndex(String hostName, Long pageNum, Long pageSize,String status);

    /**
     * @param javaProcessInfo
     */
    void updateOrInsert(JavaProcessInfo javaProcessInfo);

    /**
     * @param hostName
     * @param pid
     */
    void removeByHostNameAndPid(String hostName, int pid);

    /**
     * @param hostName
     */
    void removeByHostName(String hostName);

    /**
     * @param list
     */
    void insertList(List<JavaProcessInfo> list);

    void updateStatus(JavaProcessInfo javaProcessInfo);

}
