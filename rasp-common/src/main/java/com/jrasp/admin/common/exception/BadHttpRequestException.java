package com.jrasp.admin.common.exception;

import lombok.Data;
import lombok.EqualsAndHashCode;

@Data
@EqualsAndHashCode(callSuper = true)
public class BadHttpRequestException extends RuntimeException {
    private Integer code=400;
    public BadHttpRequestException(String message){
        super(message);
    }
}
