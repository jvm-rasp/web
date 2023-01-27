package com.jrasp.admin.server.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.pojo.AttackInfo;
import com.jrasp.admin.server.vo.AttackInfoVo;

public interface RaspAttackInfoService extends IService<AttackInfo> {
    /**
     * @param attackInfo
     */
    void addNewAttackInfo(AttackInfo attackInfo);

    /**
     * @param hostName
     * @param handleStatus
     * @param localIp
     * @param isBlocked
     * @param pageNum
     * @param pageSize
     * @return
     */
    PageResult<AttackInfo> getIndex(String hostName, Integer handleStatus, String localIp, Boolean isBlocked, Long pageNum, Long pageSize);

    /**
     * @param id
     */
    void deleteAttackInfo(Long id);

    /**
     * @param ids
     */
    void batchDeleteAttackInfo(Long[] ids);

    /**
     * 标记状态
     *
     * @param id
     * @param status
     */
    void mark(long id, int status);

    /**
     * 按照时间聚合
     * @return
     */
    AttackInfoVo getWeekData();
}
