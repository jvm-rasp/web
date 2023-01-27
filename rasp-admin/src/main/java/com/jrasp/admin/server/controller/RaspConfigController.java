package com.jrasp.admin.server.controller;

import com.jrasp.admin.common.vo.CommonResult;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.pojo.RaspConfig;
import com.jrasp.admin.server.pojo.RbacUser;
import com.jrasp.admin.server.service.RaspConfigService;
import com.jrasp.admin.server.vo.CreateConfigVo;
import com.jrasp.admin.server.vo.QueryVo;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/config")
public class RaspConfigController {

    @Autowired
    private RaspConfigService raspConfigService;

    @GetMapping("/index")
    public CommonResult<PageResult<RaspConfig>> index(
            @RequestParam(value = "configName", required = false) String configName,
            @RequestParam(value = "status", required = false) Integer status,
            @RequestParam(value = "current", required = false, defaultValue = "1") Long pageNum,
            @RequestParam(value = "pageSize", required = false, defaultValue = "10") Long pageSize
    ) {
        PageResult<RaspConfig> result = raspConfigService.getIndex(configName, status, pageNum, pageSize);
        return CommonResult.success(result);
    }

    @GetMapping("/delete")
    public CommonResult<Void> delete(@RequestParam(value = "id") Integer id) {
        raspConfigService.deleteConfig(id);
        return CommonResult.success("删除成功！");
    }

    @GetMapping("/copy")
    public CommonResult<Void> copy(@RequestParam(value = "id") Integer id) {
        raspConfigService.copyConfig(id);
        return CommonResult.success("复制成功！");
    }


    @PostMapping("/create")
    public CommonResult<Void> createConfig(@RequestBody CreateConfigVo raspConfig, @AuthenticationPrincipal RbacUser rbacUser) {
        try {
            raspConfigService.createConfig(raspConfig, rbacUser);
            return CommonResult.success("创建成功！");
        } catch (Exception e) {
            e.printStackTrace();
            return CommonResult.error("创建失败！");
        }
    }

    // 配置列表摘要信息
    @GetMapping("/list/simple")
    public CommonResult<List<QueryVo>> simple() {
        List<QueryVo> queryVos = raspConfigService.listConfig();
        return CommonResult.success(queryVos);
    }

    // 更新配置
    @PostMapping("/update")
    public CommonResult<Boolean> update(@RequestParam(value = "hostName") String hostName, @RequestParam(value = "configId") Integer configId) throws Exception {
        boolean b = raspConfigService.updateConfig(hostName, configId);
        return CommonResult.success(b);
    }

    // 更新配置
    @PostMapping("/batch/update")
    public CommonResult<Map<String, Boolean>> batchUpdate(@RequestParam(value = "hostNames") String hostNames, @RequestParam(value = "configId") Integer configId) throws Exception {
        if (StringUtils.isBlank(hostNames)) {
            return CommonResult.error("hostNames is blank!");
        }
        String[] hostArray = hostNames.split("\n");
        Map<String, Boolean> stringBooleanMap = raspConfigService.batchUpdateConfig(hostArray, configId);
        return CommonResult.success("批量更新配置成功", stringBooleanMap);
    }
}
