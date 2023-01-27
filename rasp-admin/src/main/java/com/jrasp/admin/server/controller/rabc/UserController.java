package com.jrasp.admin.server.controller.rabc;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.jrasp.admin.common.vo.CommonResult;
import com.jrasp.admin.common.vo.PageResult;
import com.jrasp.admin.server.dto.RbacUserDto;
import com.jrasp.admin.server.dto.StatusDto;
import com.jrasp.admin.server.pojo.RbacUser;
import com.jrasp.admin.server.service.RbacUserService;
import com.jrasp.admin.server.vo.MenuDataItem;
import com.jrasp.admin.server.vo.RbacUserVo;
import com.jrasp.admin.server.security.AuthorizationIgnore;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;
import springfox.documentation.annotations.ApiIgnore;

import javax.annotation.Resource;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.util.List;

@RestController("RbacUserController")
@RequestMapping("/rbac/user")
public class UserController {
    @Resource
    private RbacUserService rbacUserService;

    @AuthorizationIgnore
    @GetMapping("current")
    public CommonResult<RbacUser> current(@ApiIgnore @AuthenticationPrincipal RbacUser rbacUser) {
        return CommonResult.success(rbacUser);
    }

    @AuthorizationIgnore
    @GetMapping("menu")
    private void menu(@ApiIgnore @AuthenticationPrincipal RbacUser user, HttpServletResponse response) throws IOException {
        if (user == null || user.getId() == null) {
            return;
        }
        Integer id = user.getId();
        List<MenuDataItem> list = rbacUserService.getMenuByUserId(id);
        CommonResult<List<MenuDataItem>> result = CommonResult.success("获取数据成功！", list);
        ObjectMapper objectMapper = new ObjectMapper();
        String content = objectMapper.writeValueAsString(result);
        response.addHeader("Content-Type", "application/json;charset=utf-8");
        response.getWriter().write(content);
    }

    @PostMapping("create")
    public CommonResult<Void> create(@RequestBody @Validated RbacUserDto dto) {
        rbacUserService.createUser(dto);
        return CommonResult.success("添加成功！");
    }

    @PostMapping("update")
    public CommonResult<Void> update(@RequestBody @Validated RbacUserDto dto) {
        rbacUserService.updateUser(dto);
        return CommonResult.success("编辑成功！");
    }

    @GetMapping("delete")
    public CommonResult<Void> delete(@RequestParam(value = "id") Integer id) {
        rbacUserService.deleteUser(id);
        return CommonResult.success("删除成功！");
    }

    @PostMapping("status")
    public CommonResult<Void> status(@RequestBody @Validated StatusDto dto) {
        rbacUserService.setUserStatus(dto);
        return CommonResult.success("修改成功！");
    }

    @GetMapping("index")
    public CommonResult<PageResult<RbacUserVo>> index(
            @RequestParam(value = "username", required = false) String username,
            @RequestParam(value = "mobile", required = false) String mobile,
            @RequestParam(value = "status", required = false) Integer status,
            @RequestParam(value = "current", required = false, defaultValue = "1") Long pageNum,
            @RequestParam(value = "pageSize", required = false, defaultValue = "20") Long pageSize
    ) {
        PageResult<RbacUserVo> result = rbacUserService.getIndex(username, mobile, status, pageNum, pageSize);
        return CommonResult.success(result);
    }
}
