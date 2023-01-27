package com.jrasp.admin.server.service.impl;

import com.alibaba.fastjson.JSONObject;
import com.alibaba.nacos.api.annotation.NacosInjected;
import com.alibaba.nacos.api.config.ConfigService;
import com.alibaba.nacos.api.exception.NacosException;
import com.jrasp.admin.server.pojo.EmailAccount;
import com.jrasp.admin.server.service.EmailService;
import org.springframework.stereotype.Service;

@Service
public class EmailServiceImpl implements EmailService {

    private final String EMAIL_ACCOUNT_CONFIG_ID = "EMAIL_ACCOUNT_CONFIG";

    private final String GROUP_ID = "DEFAULT_GROUP";

    @NacosInjected
    private ConfigService nacosConfigService;

    @Override
    public EmailAccount get() throws NacosException {
        String config = nacosConfigService.getConfig(EMAIL_ACCOUNT_CONFIG_ID, GROUP_ID, 10000);
        return JSONObject.parseObject(config, EmailAccount.class);
    }

    @Override
    public boolean update(EmailAccount emailAccount) throws NacosException {
        String config = JSONObject.toJSONString(emailAccount);
        return nacosConfigService.publishConfig(EMAIL_ACCOUNT_CONFIG_ID, GROUP_ID, config);
    }
}
