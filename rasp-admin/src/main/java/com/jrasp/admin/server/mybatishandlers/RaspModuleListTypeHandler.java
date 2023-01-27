package com.jrasp.admin.server.mybatishandlers;

import com.alibaba.fastjson.TypeReference;
import com.jrasp.admin.server.pojo.RaspModule;

import java.util.List;

public class RaspModuleListTypeHandler extends ListTypeHandler<RaspModule.ParameterItem> {

    @Override
    protected TypeReference<List<RaspModule.ParameterItem>> specificType() {
        return new TypeReference<List<RaspModule.ParameterItem>>() {
        };
    }

}

