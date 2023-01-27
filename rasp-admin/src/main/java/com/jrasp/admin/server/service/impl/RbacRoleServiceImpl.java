package com.jrasp.admin.server.service.impl;

import cn.hutool.core.bean.BeanUtil;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.jrasp.admin.common.exception.BadHttpRequestException;
import com.jrasp.admin.server.dto.RbacPermissionRoleAuthDto;
import com.jrasp.admin.server.dto.RbacRoleDto;
import com.jrasp.admin.server.mapper.RbacRoleMapper;
import com.jrasp.admin.server.pojo.RbacPermission;
import com.jrasp.admin.server.pojo.RbacRole;
import com.jrasp.admin.server.pojo.RbacRolePermission;
import com.jrasp.admin.server.pojo.RbacUserRole;
import com.jrasp.admin.server.service.RbacPermissionService;
import com.jrasp.admin.server.service.RbacRolePermissionService;
import com.jrasp.admin.server.service.RbacRoleService;
import com.jrasp.admin.server.service.RbacUserRoleService;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Isolation;
import org.springframework.transaction.annotation.Transactional;

import javax.annotation.Resource;
import java.util.*;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.stream.Collectors;

@Service
public class RbacRoleServiceImpl extends ServiceImpl<RbacRoleMapper, RbacRole>
        implements RbacRoleService {
    @Resource
    private RbacUserRoleService rbacUserRoleService;
    @Resource
    private RbacRolePermissionService rbacRolePermissionService;
    @Resource
    private RbacPermissionService rbacPermissionService;

    @Override
    public List<String> getRoleByIds(List<Integer> roleIds) {
        return lambdaQuery().in(roleIds != null && !roleIds.isEmpty(), RbacRole::getId, roleIds)
                .select(RbacRole::getCode)
                .list()
                .stream().map(RbacRole::getCode).collect(Collectors.toList());
    }

    @Override
    public Integer createRole(RbacRoleDto dto) {
        RbacRole rbacRole = BeanUtil.copyProperties(dto, RbacRole.class);
        save(rbacRole);
        return rbacRole.getId();
    }

    @Override
    public void updateRole(RbacRoleDto dto) {
        if (dto.getId() == null || dto.getId() == 0) {
            throw new BadHttpRequestException("请填写角色id！");
        }
        RbacRole rbacRole = BeanUtil.copyProperties(dto, RbacRole.class);
        updateById(rbacRole);
    }

    @Override
    @Transactional(isolation = Isolation.READ_COMMITTED)
    public void deleteRole(Long id) {
        RbacRole rbacRole = getById(id);
        Integer count = rbacUserRoleService.lambdaQuery().eq(RbacUserRole::getRoleId, id).count();
        if (count > 0) {
            throw new BadHttpRequestException("请移除所有管理员的" + rbacRole.getName() + "角色后，再删除该角色！");
        }
        removeById(id);
    }

    @Override
    public List<RbacRole> getIndex() {
        return lambdaQuery().orderByAsc(RbacRole::getId).list();
    }

    @Override
    @Transactional(isolation = Isolation.READ_COMMITTED)
    public void authPermission(RbacPermissionRoleAuthDto dto) {
        Integer roleId = dto.getRoleId();
        List<Integer> permissionIds = dto.getPermissionIds();
        if (permissionIds.isEmpty()) {
            throw new BadHttpRequestException("请选择要分配的权限！");
        }
        rbacRolePermissionService.lambdaUpdate().eq(RbacRolePermission::getRoleId, dto.getRoleId())
                .remove();
        ArrayList<RbacRolePermission> rbacRolePermissions = new ArrayList<>();
        permissionIds.forEach(permissionId -> {
            RbacRolePermission rbacRolePermission = new RbacRolePermission();
            rbacRolePermission.setRoleId(roleId);
            rbacRolePermission.setPermissionId(permissionId);
            rbacRolePermissions.add(rbacRolePermission);
        });
        rbacRolePermissionService.saveBatch(rbacRolePermissions);
    }

    @Override
    public List<Integer> getPermissionIdsByRoleId(Integer roleId) {
        return rbacRolePermissionService.getPermissionIdsByRoleIds(Collections.singletonList(roleId));
    }

    @Override
    public Map<String, Boolean> getRoleAccess(List<Integer> roleIds) {
        HashMap<String, Boolean> result = new HashMap<>();
        List<Integer> permissionIdsByRoleIds = rbacRolePermissionService.getPermissionIdsByRoleIds(roleIds);
        HashSet<Integer> permissionIds = new HashSet<>(permissionIdsByRoleIds);
        List<RbacPermission> list = rbacPermissionService.lambdaQuery().select(RbacPermission::getId, RbacPermission::getUrl).list();

        list.forEach(permission -> {
            String url = permission.getUrl().substring(1);
            result.put(lineToHump(url), false);
            if (permissionIds.contains(permission.getId())) {
                result.put(lineToHump(url), true);
            }
        });
        return result;
    }

    private static final Pattern linePattern = Pattern.compile("/(\\w)");

    private static String lineToHump(String str) {
        str = str.toLowerCase();
        Matcher matcher = linePattern.matcher(str);
        StringBuffer sb = new StringBuffer();
        while (matcher.find()) {
            matcher.appendReplacement(sb, matcher.group(1).toUpperCase());
        }
        matcher.appendTail(sb);
        return sb.toString();
    }
}




