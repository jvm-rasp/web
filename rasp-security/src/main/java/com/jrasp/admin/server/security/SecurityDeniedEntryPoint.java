package com.jrasp.admin.server.security;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.web.AuthenticationEntryPoint;
import org.springframework.stereotype.Component;

import javax.annotation.Resource;
import javax.servlet.ServletException;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.util.HashMap;

@Component
public class SecurityDeniedEntryPoint implements AuthenticationEntryPoint {
    @Resource
    private ObjectMapper objectMapper;
    @Override
    public void commence(HttpServletRequest request, HttpServletResponse response, AuthenticationException e) throws IOException, ServletException {
        HashMap<String, Object> result = new HashMap<>();
        result.put("code",401);
        result.put("msg","请登录后继续操作!");
        String content = objectMapper.writeValueAsString(result);
        response.setStatus(200);
        response.addHeader("Content-Type","application/json;charset=utf-8");
        response.getWriter().write(content);
    }
}
