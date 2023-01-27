package com.jrasp.admin.server.security;

import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.stereotype.Component;
import org.springframework.web.method.HandlerMethod;
import org.springframework.web.servlet.HandlerExecutionChain;
import org.springframework.web.servlet.mvc.method.annotation.RequestMappingHandlerMapping;

import javax.annotation.Resource;
import javax.servlet.http.HttpServletRequest;
import java.util.Collection;

@Component
@Slf4j
public class RbacService {
    @Resource
    private RequestMappingHandlerMapping requestMappingHandlerMapping;

    public boolean hasPermission(HttpServletRequest request, Authentication authentication) throws Exception {
        HandlerExecutionChain handlerExeChain = requestMappingHandlerMapping.getHandler(request);
        if(handlerExeChain!=null){
            HandlerMethod handler = (HandlerMethod) handlerExeChain.getHandler();
            if(handler.hasMethodAnnotation(AuthorizationIgnore.class)){
                return true;
            }
        }
        Object principal = authentication.getPrincipal();
        if(principal instanceof UserDetails){
            UserDetails userDetails= (UserDetails) principal;
            Collection<? extends GrantedAuthority> authorities = userDetails.getAuthorities();
            SimpleGrantedAuthority authority = new SimpleGrantedAuthority(request.getRequestURI());
            if(authorities.contains(authority)){
                return true;
            }
        }
        return false;
    }
}
