package com.jrasp.admin.server.service;

import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.dto.RbacUserDto;
import com.jrasp.admin.server.dto.StatusDto;
import com.jrasp.admin.server.pojo.RbacUser;
import com.baomidou.mybatisplus.extension.service.IService;
import com.jrasp.admin.server.vo.MenuDataItem;
import com.jrasp.admin.server.vo.RbacUserVo;

import java.util.List;

public interface RbacUserService extends IService<RbacUser> {
    List<MenuDataItem> getMenuByUserId(Integer userId);
    void createUser(RbacUserDto userDto);
    void updateUser(RbacUserDto userDto);
    void deleteUser(Integer id);
    void setUserStatus(StatusDto dto);
    PageResult<RbacUserVo> getIndex(String username, String mobile, Integer status, Long pageNum, Long pageSize);
    String getHomeUrl(Integer userId);
}
