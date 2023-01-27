package com.jrasp.admin.server.controller;

import com.jrasp.admin.common.vo.CommonResult;
import com.jrasp.admin.server.pojo.EmailAccount;
import com.jrasp.admin.server.service.EmailService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.security.ProtectionDomain;


@RestController
@RequestMapping("/email")
@Slf4j
public class EmailController {


    @Autowired
    private EmailService emailService;

    /**
     * 获取邮件配置信息
     *
     * @return
     * @throws Exception
     */
    // 权限设置
    @GetMapping("/get")
    public CommonResult<EmailAccount> get() {
        try {
            ProtectionDomain protectionDomain = this.getClass().getProtectionDomain();
            EmailAccount emailAccount = emailService.get();
            return CommonResult.success("获取邮件配置成功！", emailAccount);
        } catch (Exception e) {
            log.error("获取邮件配置失败", e);
        }
        return CommonResult.error("获取邮件配置失败！");
    }

    // 更新配置
    @PostMapping("/update")
    public CommonResult<Boolean> update(@RequestBody EmailAccount emailAccount) {
        try {
            boolean b = emailService.update(emailAccount);
            return CommonResult.success("更新邮件配置成功！", b);
        } catch (Exception e) {
            log.error("更新邮件配置失败！", e);
        }
        return CommonResult.error("更新邮件配置失败！");
    }

}
