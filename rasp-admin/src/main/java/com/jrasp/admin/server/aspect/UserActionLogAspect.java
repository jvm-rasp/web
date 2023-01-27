package com.jrasp.admin.server.aspect;

import com.jrasp.admin.server.annotation.UserActionLog;
import com.jrasp.admin.server.components.ExpressionEvaluator;
import com.jrasp.admin.server.pojo.RbacUser;
import lombok.extern.slf4j.Slf4j;
import org.aspectj.lang.JoinPoint;
import org.aspectj.lang.annotation.AfterReturning;
import org.aspectj.lang.annotation.Aspect;
import org.aspectj.lang.annotation.Pointcut;
import org.aspectj.lang.reflect.MethodSignature;
import org.springframework.context.expression.AnnotatedElementKey;
import org.springframework.expression.EvaluationContext;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Component;

import java.lang.reflect.Method;

@Aspect
@Component
@Slf4j
public class UserActionLogAspect {
    private ExpressionEvaluator<String> evaluator = new ExpressionEvaluator<>();

    @Pointcut("@annotation(com.jrasp.admin.server.annotation.UserActionLog)")
    public void logPointCut(){
    }

    @AfterReturning("logPointCut()")
    public void userActionLog(JoinPoint joinPoint){
        RbacUser rbacUser = (RbacUser) SecurityContextHolder.getContext().getAuthentication().getPrincipal();
        MethodSignature signature = (MethodSignature) joinPoint.getSignature();
        Method method = signature.getMethod();
        UserActionLog annotation = method.getAnnotation(UserActionLog.class);
        String operation = annotation.value();

        EvaluationContext evaluationContext = evaluator.createEvaluationContext(joinPoint.getTarget(), joinPoint.getTarget().getClass(),
                ((MethodSignature) joinPoint.getSignature()).getMethod(), joinPoint.getArgs());
        AnnotatedElementKey methodKey = new AnnotatedElementKey(((MethodSignature) joinPoint.getSignature()).getMethod(), joinPoint.getTarget().getClass());
        String itemStr = evaluator.condition(annotation.item(), methodKey, evaluationContext, String.class);
        log.info(itemStr);
    }
}
