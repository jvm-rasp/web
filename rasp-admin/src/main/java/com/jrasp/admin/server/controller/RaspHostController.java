package com.jrasp.admin.server.controller;

import com.jrasp.admin.common.vo.CommonResult;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.pojo.RaspHost;
import com.jrasp.admin.server.service.RaspHostService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/host")
public class RaspHostController {

    @Autowired
    private RaspHostService raspHostService;

    @GetMapping("/index")
    public CommonResult<PageResult<RaspHost>> index(
            @RequestParam(value = "host_name", required = false) String hostName,
            @RequestParam(value = "ip", required = false) String ip,
            @RequestParam(value = "agent_mode", required = false) String agentMode,
            @RequestParam(value = "current", required = false, defaultValue = "1") Long pageNum,
            @RequestParam(value = "pageSize", required = false, defaultValue = "10") Long pageSize
    ) {
        PageResult<RaspHost> result = raspHostService.getIndex(hostName, ip,agentMode, pageNum, pageSize);
        return CommonResult.success(result);
    }

    @GetMapping("/delete")
    public CommonResult<Void> delete(@RequestParam(value = "id") Integer id) {
        raspHostService.deleteHost(id);
        return CommonResult.success("删除成功！");
    }
}
