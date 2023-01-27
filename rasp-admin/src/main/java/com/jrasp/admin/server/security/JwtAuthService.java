package com.jrasp.admin.server.security;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.jrasp.admin.common.exception.BadHttpRequestException;
import com.jrasp.admin.server.pojo.RbacUser;
import com.jrasp.admin.server.service.RbacUserService;
import com.jrasp.admin.server.vo.LoginVo;
import com.nimbusds.jose.JOSEException;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Component;

import javax.annotation.Resource;

@Component
public class JwtAuthService {
    @Resource
    private RbacUserService rbacUserService;

    @Resource
    private AuthenticationManager authenticationManager;

    @Resource
    private JwtTokenUtil jwtTokenUtil;

    public LoginVo login(String username, String password) throws JsonProcessingException, JOSEException {
        RbacUser rbacUser = rbacUserService.lambdaQuery().eq(RbacUser::getUsername, username)
                .select(RbacUser::getId)
                .last("limit 1")
                .one();
        if(rbacUser==null){
            throw new BadHttpRequestException("账号密码错误！");
        }
        try{
            //使用用户名密码进行登录验证
            UsernamePasswordAuthenticationToken token =
                    new UsernamePasswordAuthenticationToken( username, password );
            Authentication authentication = authenticationManager.authenticate(token);
            SecurityContextHolder.getContext().setAuthentication(authentication);
        }catch(AuthenticationException e){
            throw new BadHttpRequestException("用户名密码错误！");
        }

        //生成JWT
        String token = jwtTokenUtil.generateToken(rbacUser.getId().longValue());
        LoginVo result = new LoginVo();
        result.setToken(token);
        String url = rbacUserService.getHomeUrl(rbacUser.getId());
        result.setUrl(url);
        return result;
    }
}
