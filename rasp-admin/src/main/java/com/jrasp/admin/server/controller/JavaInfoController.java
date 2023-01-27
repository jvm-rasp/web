package com.jrasp.admin.server.controller;

import com.jrasp.admin.common.vo.CommonResult;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.pojo.JavaProcessInfo;
import com.jrasp.admin.server.service.JavaProcessInfoService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import javax.annotation.security.PermitAll;

@RestController
@RequestMapping("/java")
public class JavaInfoController {

    @Autowired
    private JavaProcessInfoService javaProcessInfoService;

    @GetMapping("/index")
    @PermitAll
    public CommonResult<PageResult<JavaProcessInfo>> index(
            @RequestParam(value = "current", required = false, defaultValue = "1") Long pageNum,
            @RequestParam(value = "pageSize", required = false, defaultValue = "10") Long pageSize,
            @RequestParam(value = "status", required = false) String status,
            @RequestParam(value = "hostName", required = false) String hostName
    ) {
        PageResult<JavaProcessInfo> result = javaProcessInfoService.getIndex(hostName, pageNum, pageSize, status);
        return CommonResult.success(result);
    }
}
