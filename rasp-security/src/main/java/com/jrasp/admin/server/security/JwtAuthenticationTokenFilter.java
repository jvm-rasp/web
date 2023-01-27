package com.jrasp.admin.server.security;

import cn.hutool.core.util.StrUtil;
import com.jrasp.admin.common.exception.UnAuthenticationException;
import com.nimbusds.jose.JOSEException;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import javax.annotation.Resource;
import javax.servlet.FilterChain;
import javax.servlet.ServletException;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.text.ParseException;

@Component
@Slf4j
public class JwtAuthenticationTokenFilter extends OncePerRequestFilter {

    @Resource
    private JwtTokenUtil jwtTokenUtil;
    @Resource
    private RockUserDetailsService rockUserDetailsService;

    @Override
    protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response, FilterChain chain) throws ServletException, IOException {
        String jwtToken = request.getHeader(jwtTokenUtil.getHeader());
        if(StrUtil.isNotBlank(jwtToken) && !request.getRequestURI().startsWith("/login")){
            Long userId = null;
            try {
                userId = jwtTokenUtil.getUserIdFromToken(jwtToken);
            } catch (ParseException | JOSEException e) {
                throw new UnAuthenticationException("请登录后继续操作！");
            }
            if(userId!=null && SecurityContextHolder.getContext().getAuthentication()==null){
                UserDetails userDetails = rockUserDetailsService.getUserDetailsById(userId);
                if(jwtTokenUtil.validateToken(jwtToken,userDetails)){
                    UsernamePasswordAuthenticationToken authenticationToken = new UsernamePasswordAuthenticationToken(userDetails, null, userDetails.getAuthorities());
                    SecurityContextHolder.getContext().setAuthentication(authenticationToken);
                }
            }
        }
        chain.doFilter(request,response);
    }
}
