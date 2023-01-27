package com.jrasp.admin.server.configruation;

import com.github.xiaoymin.knife4j.spring.annotations.EnableKnife4j;
import com.google.common.collect.Lists;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import springfox.documentation.builders.ApiInfoBuilder;
import springfox.documentation.builders.PathSelectors;
import springfox.documentation.builders.RequestHandlerSelectors;
import springfox.documentation.service.*;
import springfox.documentation.spi.DocumentationType;
import springfox.documentation.spi.service.contexts.SecurityContext;
import springfox.documentation.spring.web.plugins.Docket;
import springfox.documentation.swagger2.annotations.EnableSwagger2;

import java.util.List;

@Configuration
@EnableSwagger2
@EnableKnife4j
public class Knife4jConfigure {

    @Bean(value = "AdminApi")
    public Docket userApi() {
        return new Docket(DocumentationType.SWAGGER_2)
                .apiInfo(adminApiInfo())
                .groupName("管理后台接口")
                .select()
                .apis(RequestHandlerSelectors.basePackage("com.icodeview.rock.admin.controller"))
                .paths(PathSelectors.any())
                .build().securityContexts(Lists.newArrayList(
                        securityContext()
                )).securitySchemes(Lists.<SecurityScheme>newArrayList(apiKey()));
    }
    private ApiInfo adminApiInfo() {
        return new ApiInfoBuilder()
                .title("后台接口")
                .description("后台接口")
                .termsOfServiceUrl("https://www.icodeview.com")
                .version("1.0")
                .build();
    }

    private ApiKey apiKey(){
        return new ApiKey("BearerToken","Authorization","header");
    }
    private List<SecurityReference> defaultAuth() {
        AuthorizationScope authorizationScope = new AuthorizationScope("global", "accessEverything");
        AuthorizationScope[] authorizationScopes = new AuthorizationScope[1];
        authorizationScopes[0] = authorizationScope;
        return Lists.newArrayList(new SecurityReference("BearerToken", authorizationScopes));
    }
    private SecurityContext securityContext() {
        return SecurityContext.builder()
                .securityReferences(defaultAuth())
                .forPaths(PathSelectors.regex("/.*"))
                .build();
    }
}
