package com.jrasp.admin.server.controller;

import com.jrasp.admin.common.exception.BadHttpRequestException;
import com.jrasp.admin.common.vo.CommonResult;
import com.jrasp.admin.server.security.AuthorizationIgnore;
import lombok.extern.slf4j.Slf4j;
import org.springframework.web.bind.annotation.*;

import javax.annotation.security.PermitAll;

@Slf4j
@RestController
@RequestMapping("/license")
public class LicenseController {

    @GetMapping("/get")
    @PermitAll
    @AuthorizationIgnore
    public CommonResult<Void> get() {
        try {
            String result = getLicense();
            return CommonResult.success(result);
        } catch (Exception e) {
            e.printStackTrace();
            throw new BadHttpRequestException("获取密钥Lisence异常！");
        }
    }

    private static final String key="YRK8NKwfNChXDKuwe2yokYl9jHM42SOHw7zIHrcO3BE/Tc8f3ucFw3kQRlCen5OtYHzOYdCtNOWWB6YTCBu6mEIN/N9fOLUsR7Il8kaoSOabWA2MAh/HMSNKhjhkcozAHEl0C3ZRcjev0b+ek9WilNFEFenfB/ltnvK+AuQpTl+3f4A0KloNmlGnyOhDELS2hBcQ0j5lDSnHzPpxR73hmBMDGEvWTNMCd7IFKAG7PEuG8dTmIqfl4ixwZYhFhK3SxT0qCong1x2ig3otLMtDVA7AMjkemHj7M3zxAyNuKdv3bQkCl/Go7GZ/P0QLHMbh3o/hHErrmb2fQfluTnbmdg==";

    private String getLicense() {
        return key;
    }

}

