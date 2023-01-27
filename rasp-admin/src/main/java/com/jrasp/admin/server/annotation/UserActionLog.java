package com.jrasp.admin.server.annotation;

import java.lang.annotation.*;

@Target(ElementType.METHOD)
@Retention(RetentionPolicy.RUNTIME)
@Documented
public @interface UserActionLog {
    String value() default "";
    String item() default "0";
}
