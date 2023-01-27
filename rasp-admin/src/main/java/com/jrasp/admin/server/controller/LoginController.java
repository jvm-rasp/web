package com.jrasp.admin.server.controller;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.jrasp.admin.common.dto.admin.LoginDto;
import com.jrasp.admin.common.exception.BadHttpRequestException;
import com.jrasp.admin.common.vo.CommonResult;
import com.jrasp.admin.server.pojo.RbacUser;
import com.jrasp.admin.server.vo.LoginVo;
import com.jrasp.admin.server.security.JwtAuthService;
import com.jrasp.admin.server.security.AuthorizationIgnore;
import com.nimbusds.jose.JOSEException;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

import javax.annotation.Resource;
import javax.annotation.security.PermitAll;
import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping("/login")
public class LoginController {
    @Resource
    private JwtAuthService jwtAuthService;

    @PostMapping("account")
    public CommonResult<LoginVo> account(@RequestBody @Validated LoginDto dto){
        try {
            LoginVo result = jwtAuthService.login(dto.getUsername(), dto.getPassword());
            return CommonResult.success("登录成功！",result);
        } catch (JsonProcessingException | JOSEException e) {
            e.printStackTrace();
            throw new BadHttpRequestException("登录失败，请重试！");
        }
    }

    @PermitAll
    @AuthorizationIgnore
    @GetMapping("/test")
    public CommonResult<Map<String,Object>> test(@AuthenticationPrincipal RbacUser user){
        HashMap<String, Object> map = new HashMap<>();
        map.put("code",200);
        map.put("user",user);
        return CommonResult.success(map);
    }
}
