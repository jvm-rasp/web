package com.jrasp.admin.server.controller.rabc;

import com.jrasp.admin.common.vo.CommonResult;
import com.jrasp.admin.server.dto.RbacPermissionDto;
import com.jrasp.admin.server.pojo.RbacPermission;
import com.jrasp.admin.server.service.RbacPermissionService;
import com.jrasp.admin.server.security.AuthorizationIgnore;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

import javax.annotation.Resource;
import java.util.List;

@RestController("RbacPermission")
@RequestMapping("/rbac/permission")
public class PermissionController {
    @Resource
    private RbacPermissionService rbacPermissionService;

    @AuthorizationIgnore
    @GetMapping("index")
    public CommonResult<List<RbacPermission>> index(){
        List<RbacPermission> list = rbacPermissionService.getIndex();
        return CommonResult.success(list);
    }

    @PostMapping("create")
    public CommonResult<Void> create(@RequestBody @Validated RbacPermissionDto dto){
        rbacPermissionService.createPermission(dto);
        return CommonResult.success("添加成功");
    }

    @PostMapping("update")
    public CommonResult<Void> update(@RequestBody @Validated RbacPermissionDto dto){
        rbacPermissionService.updatePermission(dto);
        return CommonResult.success("编辑成功");
    }

    @PostMapping("delete")
    public CommonResult<Void> delete(@RequestParam(value = "id") Long id){
        rbacPermissionService.delete(id);
        return CommonResult.success("删除成功！");
    }
}
