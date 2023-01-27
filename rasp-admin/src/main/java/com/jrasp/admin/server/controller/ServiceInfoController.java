package com.jrasp.admin.server.controller;

import com.jrasp.admin.common.vo.CommonResult;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.pojo.RbacUser;
import com.jrasp.admin.server.pojo.ServiceInfo;
import com.jrasp.admin.server.service.ServiceInfoService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.*;

@Slf4j
@RestController
@RequestMapping("/service")
public class ServiceInfoController {

    @Autowired
    private ServiceInfoService serviceInfoService;

    @GetMapping("/index")
    public CommonResult<PageResult<ServiceInfo>> index(
            @RequestParam(value = "service_name", required = false) String serviceName,
            @RequestParam(value = "low_period", required = false) Integer lowPeriod,
            @RequestParam(value = "current", required = false, defaultValue = "1") Long pageNum,
            @RequestParam(value = "pageSize", required = false, defaultValue = "10") Long pageSize
    ) {
        PageResult<ServiceInfo> result = serviceInfoService.getIndex(serviceName,lowPeriod, pageNum, pageSize);
        return CommonResult.success(result);
    }

    @PostMapping("/create")
    public CommonResult<Void> create(@RequestBody ServiceInfo serviceInfo, @AuthenticationPrincipal RbacUser rbacUser) {
        try {
            serviceInfoService.createService(serviceInfo, rbacUser);
            return CommonResult.success("创建服务成功！");
        } catch (Exception e) {
            log.error("创建服务失败", e);
            return CommonResult.error("创建服务失败！");
        }
    }

    @GetMapping("/delete")
    public CommonResult<Void> delete(@RequestParam(value = "id") Integer id) {
        serviceInfoService.deleteService(id);
        return CommonResult.success("删除成功！");
    }

    @PostMapping("/update")
    public CommonResult<Void> update(@RequestBody ServiceInfo serviceInfo) {
        serviceInfoService.updateService(serviceInfo);
        return CommonResult.success("更新成功！");
    }

}
