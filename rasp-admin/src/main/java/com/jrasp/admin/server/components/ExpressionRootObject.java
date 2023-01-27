package com.jrasp.admin.server.components;

import lombok.AllArgsConstructor;
import lombok.Data;

@Data
@AllArgsConstructor
public class ExpressionRootObject {
    private final Object object;
    private final Object[] args;
}
