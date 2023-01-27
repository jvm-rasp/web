package com.jrasp.admin.server.service.impl;

import cn.hutool.core.bean.BeanUtil;
import cn.hutool.core.util.StrUtil;
import com.baomidou.mybatisplus.extension.conditions.query.LambdaQueryChainWrapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.jrasp.admin.common.exception.BadHttpRequestException;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.dto.RbacUserDto;
import com.jrasp.admin.server.dto.StatusDto;
import com.jrasp.admin.server.mapper.RbacUserMapper;
import com.jrasp.admin.server.pojo.RbacRole;
import com.jrasp.admin.server.pojo.RbacUser;
import com.jrasp.admin.server.pojo.RbacUserRole;
import com.jrasp.admin.server.service.*;
import com.jrasp.admin.server.vo.MenuDataItem;
import com.jrasp.admin.server.vo.RbacUserVo;
import com.jrasp.admin.server.security.RockUserDetailsService;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.AuthorityUtils;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Isolation;
import org.springframework.transaction.annotation.Transactional;

import javax.annotation.Resource;
import java.time.LocalDateTime;
import java.util.*;
import java.util.stream.Collectors;

@Service
public class RbacUserServiceImpl extends ServiceImpl<RbacUserMapper, RbacUser>
        implements RbacUserService, RockUserDetailsService {

    @Resource
    private RbacUserRoleService rbacUserRoleService;

    @Resource
    private RbacRoleService rbacRoleService;

    @Resource
    private RbacRolePermissionService rbacRolePermissionService;

    @Resource
    private RbacPermissionService rbacPermissionService;

    @Resource
    private PasswordEncoder passwordEncoder;

    @Override
    public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
        return getUserDetails(lambdaQuery().eq(RbacUser::getUsername, username));
    }

    @Override
    public List<MenuDataItem> getMenuByUserId(Integer userId) {
        List<Integer> roleIds = rbacUserRoleService.getRoleIdByUserId(userId);
        List<Integer> permissionIds = rbacRolePermissionService.getPermissionIdsByRoleIds(roleIds);
        return rbacPermissionService.getMenuByIds(permissionIds);
    }

    @Override
    @Transactional(isolation = Isolation.READ_COMMITTED)
    public void createUser(RbacUserDto userDto) {
        if(!checkMobile(null,userDto.getMobile())){
            throw new BadHttpRequestException("该手机号已被使用！");
        }
        if(!checkUserName(null,userDto.getUsername())){
            throw new BadHttpRequestException("该用户名已被使用！");
        }
        String password = userDto.getPassword();
        if(StrUtil.isBlank(password)){
            throw new BadHttpRequestException("请设置登录密码！");
        }
        RbacUser user = BeanUtil.copyProperties(userDto, RbacUser.class);
        user.setPassword(passwordEncoder.encode(password));
        user.setCreatedAt(LocalDateTime.now());
        user.setUpdatedAt(LocalDateTime.now());
        save(user);
        rbacUserRoleService.attachRole(user.getId(),userDto.getRoleId());
    }

    @Override
    @Transactional(isolation = Isolation.READ_COMMITTED)
    public void updateUser(RbacUserDto userDto) {
        if(!checkMobile(userDto.getId(),userDto.getMobile())){
            throw new BadHttpRequestException("该手机号已被使用！");
        }
        if(!checkUserName(userDto.getId(),userDto.getUsername())){
            throw new BadHttpRequestException("该用户名已被使用！");
        }
        RbacUser user = BeanUtil.copyProperties(userDto, RbacUser.class);
        String password = userDto.getPassword();
        if(StrUtil.isNotBlank(password)){
            user.setPassword(passwordEncoder.encode(password));
        }
        user.setCreatedAt(LocalDateTime.now());
        user.setUpdatedAt(LocalDateTime.now());
        updateById(user);
        rbacUserRoleService.attachRole(user.getId(),userDto.getRoleId());
    }

    @Override
    @Transactional(isolation = Isolation.READ_COMMITTED)
    public void deleteUser(Integer id) {
        RbacUser rbacUser = getById(id);
        rbacUserRoleService.detachRole(id);
        removeById(rbacUser.getId());
    }

    @Override
    public void setUserStatus(StatusDto dto) {
        lambdaUpdate().eq(RbacUser::getId,dto.getId())
                .set(RbacUser::getStatus,dto.getStatus())
                .set(RbacUser::getUpdatedAt,LocalDateTime.now())
                .update();
    }

    @Override
    public PageResult<RbacUserVo> getIndex(String username, String mobile, Integer status, Long pageNum, Long pageSize) {
        Page<RbacUser> rbacUserPage = new Page<>(pageNum, pageSize);
        Page<RbacUser> page = lambdaQuery().eq(StrUtil.isNotBlank(username), RbacUser::getUsername, username)
                .eq(StrUtil.isNotBlank(mobile), RbacUser::getMobile, mobile)
                .eq(status != null && status != 0, RbacUser::getStatus, status)
                .orderByDesc(RbacUser::getId)
                .page(rbacUserPage);
        List<RbacUserVo> record = page.getRecords().stream().map(rbacUser -> {
            RbacUserVo user = BeanUtil.copyProperties(rbacUser, RbacUserVo.class);
            RbacUserRole userRole = rbacUserRoleService.lambdaQuery()
                    .eq(RbacUserRole::getUserId, user.getId())
                    .last("limit 1")
                    .one();
            RbacRole rbacRole = rbacRoleService.getById(userRole.getRoleId());
            user.setRole(rbacRole);
            user.setRoleId(rbacRole.getId());
            return user;
        }).collect(Collectors.toList());
        return PageResult.page(record,page);
    }

    @Override
    public String getHomeUrl(Integer userId) {
        List<Integer> roleIds = rbacUserRoleService.getRoleIdByUserId(userId);
        List<String> permission = rbacRolePermissionService.getPermissionByRoleIds(roleIds);
        Optional<String> homeUrlOptional = permission.stream().filter(url -> url.contains("index")).findFirst();
        return homeUrlOptional.orElse("/");
    }


    private boolean checkMobile(Integer userId,String mobile){
        return checkParams(userId, lambdaQuery().eq(RbacUser::getMobile, mobile), mobile);
    }

    private boolean checkUserName(Integer userId,String username){
        return checkParams(userId, lambdaQuery().eq(RbacUser::getUsername, username), username);
    }

    private boolean checkParams(Integer userId, LambdaQueryChainWrapper<RbacUser> eq, String mobile) {
        RbacUser rbacUser = eq
                .select(RbacUser::getId)
                .one();
        if(userId==null){
            if(rbacUser!=null){
                return false;
            }
        }else{
            if(rbacUser.getId().equals(userId)){
                return true;
            }
        }
        return true;
    }

    @Override
    public UserDetails getUserDetailsById(Long id) {
        return getUserDetails(lambdaQuery().eq(RbacUser::getId, id));
    }

    private UserDetails getUserDetails(LambdaQueryChainWrapper<RbacUser> eq) {
        RbacUser user = eq.one();
        if(user==null){
            throw new UsernameNotFoundException("用户不存在！");
        }
        List<Integer> roleIds = rbacUserRoleService.getRoleIdByUserId(user.getId());
        List<String> authorities = rbacRolePermissionService.getPermissionByRoleIds(roleIds);
        List<String> roles = rbacRoleService.getRoleByIds(roleIds);
        roles = roles.stream().map(rc -> "ROLE_" + rc).collect(Collectors.toList());

        authorities.addAll(roles);
        List<GrantedAuthority> authorityList = AuthorityUtils.commaSeparatedStringToAuthorityList(String.join(",", authorities));

        user.setAuthorities(authorityList);

        Map<String, Boolean> access = rbacRoleService.getRoleAccess(roleIds);
        user.setAccess(access);
        return user;
    }
}




