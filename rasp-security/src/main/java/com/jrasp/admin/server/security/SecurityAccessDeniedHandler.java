package com.jrasp.admin.server.security;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.security.access.AccessDeniedException;
import org.springframework.security.web.access.AccessDeniedHandler;
import org.springframework.stereotype.Component;

import javax.annotation.Resource;
import javax.servlet.ServletException;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.util.HashMap;

@Component
public class SecurityAccessDeniedHandler implements AccessDeniedHandler {
    @Resource
    private ObjectMapper objectMapper;
    @Override
    public void handle(HttpServletRequest request, HttpServletResponse response, AccessDeniedException e) throws IOException, ServletException {
        HashMap<String, Object> result = new HashMap<>();
        result.put("code",403);
        result.put("msg","您没有权限访问该资源！");
        String content = objectMapper.writeValueAsString(result);
        response.setStatus(200);
        response.addHeader("Content-Type","application/json;charset=utf-8");
        response.getWriter().write(content);
    }
}
