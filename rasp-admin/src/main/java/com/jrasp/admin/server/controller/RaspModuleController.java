package com.jrasp.admin.server.controller;

import com.jrasp.admin.common.vo.CommonResult;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.dto.StatusDto;
import com.jrasp.admin.server.pojo.RaspModule;
import com.jrasp.admin.server.pojo.RbacUser;
import com.jrasp.admin.server.service.RaspModuleService;
import com.jrasp.admin.server.vo.QueryVo;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/module")
public class RaspModuleController {

    @Autowired
    private RaspModuleService raspModuleService;

    @GetMapping("/index")
    public CommonResult<PageResult<RaspModule>> index(
            @RequestParam(value = "name", required = false) String name,
            @RequestParam(value = "status", required = false) Integer status,
            @RequestParam(value = "current", required = false, defaultValue = "1") Long pageNum,
            @RequestParam(value = "pageSize", required = false, defaultValue = "10") Long pageSize
    ) {
        PageResult<RaspModule> result = raspModuleService.getIndex(name, status, pageNum, pageSize);
        return CommonResult.success(result);
    }

    @GetMapping("/delete")
    public CommonResult<Void> delete(@RequestParam(value = "id") Integer id) {
        raspModuleService.deleteModule(id);
        return CommonResult.success("删除成功！");
    }

    @PostMapping("/create")
    public CommonResult<Void> createConfig(@RequestBody RaspModule raspModule, @AuthenticationPrincipal RbacUser rbacUser) {
        try {
            raspModuleService.createModule(raspModule, rbacUser);
            return CommonResult.success("更新成功！");
        } catch (Exception e) {
            e.printStackTrace();
            return CommonResult.error("更新失败！");
        }
    }

    @GetMapping("/detail")
    public CommonResult<RaspModule> detail(@RequestParam(value = "id") Integer id) {
        RaspModule raspModule = raspModuleService.getModule(id);
        return CommonResult.success(raspModule);
    }

    @PostMapping("/update")
    public CommonResult<Void> update(@RequestBody RaspModule raspModule) {
        try {
            raspModuleService.updateModule(raspModule);
            return CommonResult.success("更新成功！");
        } catch (Exception e) {
            e.printStackTrace();
            return CommonResult.error("更新失败！");
        }
    }

    @PostMapping("/status")
    public CommonResult<Void> status(@RequestBody @Validated StatusDto dto) {
        raspModuleService.setModuleStatus(dto);
        return CommonResult.success("修改成功！");
    }

    @GetMapping("/type")
    public CommonResult<List<QueryVo>> type() {
        List<QueryVo> moduleType = raspModuleService.getModuleType();
        return CommonResult.success("success",moduleType);
    }
}