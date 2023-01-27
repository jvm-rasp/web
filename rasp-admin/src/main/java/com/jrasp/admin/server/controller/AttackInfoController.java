package com.jrasp.admin.server.controller;

import com.jrasp.admin.common.vo.CommonResult;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.pojo.AttackInfo;
import com.jrasp.admin.server.service.RaspAttackInfoService;
import com.jrasp.admin.server.vo.AttackInfoVo;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/attack")
public class AttackInfoController {

    @Autowired
    private RaspAttackInfoService raspAttackInfoService;

    @GetMapping("/index")
    public CommonResult<PageResult<AttackInfo>> index(
            @RequestParam(value = "host_name", required = false) String hostName,
            @RequestParam(value = "handle_status", required = false) Integer handle_status,
            @RequestParam(value = "local_ip", required = false) String localIp,
            @RequestParam(value = "is_blocked", required = false) Boolean isBlocked,
            @RequestParam(value = "current", required = false, defaultValue = "1") Long pageNum,
            @RequestParam(value = "pageSize", required = false, defaultValue = "10") Long pageSize
    ) {
        PageResult<AttackInfo> result = raspAttackInfoService.getIndex(hostName, handle_status, localIp, isBlocked, pageNum, pageSize);
        return CommonResult.success(result);
    }

    @GetMapping("/delete")
    public CommonResult<Void> delete(@RequestParam(value = "id") Long id) {
        raspAttackInfoService.deleteAttackInfo(id);
        return CommonResult.success("删除成功！");
    }

    @PostMapping("/batch/delete")
    public CommonResult<Void> batchDelete(@RequestParam(value = "ids") Long[] ids) {
        raspAttackInfoService.batchDeleteAttackInfo(ids);
        return CommonResult.success("批量删除成功！");
    }

    @GetMapping("/mark")
    public CommonResult<Void> mark(@RequestParam(value = "id") Long id, @RequestParam(value = "handleStatus") int handleStatus) {
        raspAttackInfoService.mark(id, handleStatus);
        return CommonResult.success("处理状态更新成功！");
    }

    @GetMapping("/data/week")
    public CommonResult<AttackInfoVo> getWeekData() {
        AttackInfoVo todayData = raspAttackInfoService.getWeekData();
        return CommonResult.success("获取数据成功", todayData);
    }

}
