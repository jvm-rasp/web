package com.jrasp.admin.server.controller.rabc;

import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.github.xiaoymin.knife4j.annotations.ApiSort;
import com.jrasp.admin.common.vo.CommonResult;
import com.jrasp.admin.server.annotation.UserActionLog;
import com.jrasp.admin.server.dto.RbacPermissionRoleAuthDto;
import com.jrasp.admin.server.dto.RbacRoleDto;
import com.jrasp.admin.server.pojo.RbacRole;
import com.jrasp.admin.server.service.RbacRoleService;
import com.jrasp.admin.server.security.AuthorizationIgnore;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiImplicitParam;
import io.swagger.annotations.ApiOperation;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

import javax.annotation.Resource;
import java.util.List;

@Api(tags = "权限管理--角色管理")
@ApiSort(2)
@RestController
@RequestMapping("/rbac/role")
public class RoleController {
    @Resource
    private RbacRoleService rbacRoleService;

    @UserActionLog(value = "查看角色")
    @ApiOperationSupport(order = 1)
    @ApiOperation("角色列表")
    @GetMapping("index")
    public CommonResult<List<RbacRole>> index(){
        List<RbacRole> list = rbacRoleService.getIndex();
        return CommonResult.success(list);
    }

    @UserActionLog(value = "添加角色")
    @ApiOperationSupport(order = 2)
    @ApiOperation("添加角色")
    @PostMapping("create")
    public CommonResult<Integer> create(@RequestBody @Validated RbacRoleDto dto){
        Integer roleId = rbacRoleService.createRole(dto);
        return CommonResult.success("添加成功！",roleId);
    }

    @UserActionLog(value = "编辑角色",item = "#dto.id")
    @ApiOperationSupport(order = 3)
    @ApiOperation("更新角色")
    @PostMapping("update")
    public CommonResult<Void> update(@RequestBody @Validated RbacRoleDto dto){
        rbacRoleService.updateRole(dto);
        return CommonResult.success("更新成功！");
    }

    @ApiOperationSupport(order = 4)
    @ApiOperation("删除角色")
    @PostMapping("delete")
    @ApiImplicitParam(value = "角色id",name = "id",required = true)
    public CommonResult<Void> delete(@RequestParam(value = "id") Long id){
        rbacRoleService.deleteRole(id);
        return CommonResult.success("删除成功！");
    }

    @ApiOperationSupport(order = 5)
    @ApiOperation("分配权限")
    @PostMapping("auth/permission")
    public CommonResult<Void> authPermission(@RequestBody @Validated RbacPermissionRoleAuthDto dto){
        rbacRoleService.authPermission(dto);
        return CommonResult.success("分配成功！");
    }

    @AuthorizationIgnore
    @ApiOperationSupport(order = 6)
    @ApiOperation("权限id")
    @GetMapping("permission/ids")
    @ApiImplicitParam(value = "角色id",name = "role_id",required = true)
    public CommonResult<List<Integer>> permissionIds(
            @RequestParam(value = "role_id") Integer roleId
    ){
        List<Integer> result = rbacRoleService.getPermissionIdsByRoleId(roleId);
        return CommonResult.success(result);
    }
}
