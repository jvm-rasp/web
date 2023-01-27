package com.jrasp.admin.server.service;

import com.alibaba.nacos.api.exception.NacosException;
import com.jrasp.admin.server.pojo.EmailAccount;

public interface EmailService {
    EmailAccount get() throws NacosException;

    boolean update(EmailAccount emailAccount) throws NacosException;
}
