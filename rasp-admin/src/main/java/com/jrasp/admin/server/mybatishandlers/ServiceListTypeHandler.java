package com.jrasp.admin.server.mybatishandlers;

import com.alibaba.fastjson.TypeReference;
import com.jrasp.admin.server.pojo.RbacUser;

import java.util.List;

public class ServiceListTypeHandler extends ListTypeHandler<RbacUser> {
    @Override
    protected TypeReference<List<RbacUser>> specificType() {
        return new TypeReference<List<RbacUser>>() {

        };
    }
}
