package common

import (
	"errors"
	"server/config"
	"server/model"
	"server/util"

	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

// 初始化mysql数据
// TODO 用户权限系统默认配置
func InitData() {
	// 是否初始化数据
	if !config.Conf.System.InitData {
		return
	}

	// 1.写入角色数据
	newRoles := make([]*model.Role, 0)
	roles := []*model.Role{
		{
			Model:   gorm.Model{ID: 1},
			Name:    "管理员",
			Keyword: "admin",
			Desc:    new(string),
			Sort:    1,
			Status:  1,
			Creator: "系统",
		},
		{
			Model:   gorm.Model{ID: 2},
			Name:    "普通用户",
			Keyword: "user",
			Desc:    new(string),
			Sort:    3,
			Status:  1,
			Creator: "系统",
		},
		{
			Model:   gorm.Model{ID: 3},
			Name:    "访客",
			Keyword: "guest",
			Desc:    new(string),
			Sort:    5,
			Status:  1,
			Creator: "系统",
		},
	}

	for _, role := range roles {
		err := DB.First(&role, role.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newRoles = append(newRoles, role)
		}
	}

	if len(newRoles) > 0 {
		err := DB.Create(&newRoles).Error
		if err != nil {
			Log.Errorf("写入系统角色数据失败：%v", err)
		}
	}

	// 2写入菜单
	newMenus := make([]model.Menu, 0)

	var rootId uint = 0
	var hostId uint = 1000
	var attackId uint = 2000
	var moduleId uint = 3000
	var configId uint = 4000
	var logId uint = 5000
	var systemId uint = 6000

	systemUserStr := "/system/user"
	attackLogStr := "/attack/list"
	// direct
	configDirectStr := "/config/list"
	hostDirectStr := "/host/list"
	// TODO 其他菜单页面的重定向需要修改

	// icon 配置
	hostIcon := "host"
	moduleIcon := "module"
	attackLogIcon := "attackLog"
	jarIcon := "jar"
	listIcon := "list2"
	jsonIcon := "json"
	runLogIcon := "runLog"
	systemIcon := "system"
	apiIcon := "api"
	apiLogIcon := "apiLog"
	menuIcon := "menu"
	roleIcon := "role"
	userIcon := "user"

	menus := []model.Menu{
		{
			Model:     gorm.Model{ID: hostId},
			Name:      "Host",
			Title:     "实例管理",
			Icon:      &hostIcon,
			Path:      "/host",
			Component: "Layout",
			Redirect:  &hostDirectStr,
			Sort:      hostId,
			ParentId:  &rootId,
			Roles:     roles,
			Creator:   "系统",
		},
		{
			Model:     gorm.Model{ID: hostId + 1},
			Name:      "List",
			Title:     "实例列表",
			Icon:      &listIcon,
			Path:      "list",
			Component: "/host/list/index",
			Sort:      hostId + 1,
			ParentId:  &hostId,
			Roles:     roles,
			Creator:   "系统",
		},
		{
			Model:     gorm.Model{ID: hostId + 2},
			Name:      "Detail",
			Title:     "实例详情",
			Hidden:    1,
			Icon:      &listIcon,
			Path:      "/detail",
			Component: "/host/detail/index",
			Sort:      hostId + 2,
			ParentId:  &hostId,
			Roles:     roles,
			Creator:   "系统",
		},
		{
			Model:     gorm.Model{ID: attackId},
			Name:      "Attack",
			Title:     "攻击告警",
			Icon:      &attackLogIcon,
			Path:      "/attack",
			Component: "Layout",
			Redirect:  &attackLogStr,
			Sort:      attackId,
			ParentId:  &rootId,
			Roles:     roles,
			Creator:   "系统",
		},
		{
			Model:     gorm.Model{ID: attackId + 1},
			Name:      "List",
			Title:     "日志列表",
			Icon:      &listIcon,
			Path:      "list",
			Component: "/attack/list/index",
			Sort:      attackId + 1,
			ParentId:  &attackId,
			Roles:     roles,
			Creator:   "系统",
		},
		{
			Model:     gorm.Model{ID: moduleId},
			Name:      "Module",
			Title:     "模块配置",
			Icon:      &moduleIcon,
			Path:      "/module",
			Component: "Layout",
			Redirect:  &systemUserStr,
			Sort:      moduleId,
			ParentId:  &rootId,
			Roles:     roles,
			Creator:   "系统",
		},
		{
			Model:     gorm.Model{ID: moduleId + 1},
			Name:      "List",
			Title:     "模块列表",
			Icon:      &listIcon,
			Path:      "list",
			Component: "/module/list/index",
			Sort:      moduleId + 1,
			ParentId:  &moduleId,
			Roles:     roles,
			Creator:   "系统",
		},
		{
			Model:     gorm.Model{ID: moduleId + 2},
			Name:      "Jar",
			Title:     "jar包管理",
			Icon:      &jarIcon,
			Path:      "jar",
			Component: "/module/jar/index",
			Sort:      moduleId + 2,
			ParentId:  &moduleId,
			Roles:     roles,
			Creator:   "系统",
		},
		{
			Model:     gorm.Model{ID: configId},
			Name:      "Config",
			Title:     "策略配置",
			Icon:      &jsonIcon,
			Path:      "/config",
			Component: "Layout",
			Redirect:  &configDirectStr,
			Sort:      configId,
			ParentId:  &rootId,
			Roles:     roles,
			Creator:   "系统",
		},
		{
			Model:     gorm.Model{ID: configId + 1},
			Name:      "List",
			Title:     "策略列表",
			Icon:      &listIcon,
			Path:      "list",
			Component: "/config/list/index",
			Sort:      configId + 1,
			ParentId:  &configId,
			Roles:     roles,
			Creator:   "系统",
		},
		{
			Model:     gorm.Model{ID: logId},
			Name:      "Log",
			Title:     "日志管理",
			Icon:      &runLogIcon,
			Path:      "/log",
			Component: "Layout",
			Redirect:  &systemUserStr,
			Sort:      logId,
			ParentId:  &rootId,
			Roles:     roles,
			Creator:   "系统",
		},
		{
			Model:     gorm.Model{ID: logId + 1},
			Name:      "List",
			Title:     "错误日志",
			Icon:      &listIcon,
			Path:      "list",
			Component: "/log/list/index",
			Sort:      logId + 1,
			ParentId:  &logId,
			Roles:     roles,
			Creator:   "系统",
		},
		{
			Model:     gorm.Model{ID: systemId},
			Name:      "System",
			Title:     "系统管理",
			Icon:      &systemIcon,
			Path:      "/system",
			Component: "Layout",
			Redirect:  &systemUserStr,
			Sort:      systemId,
			ParentId:  &rootId,
			Roles:     roles[:1],
			Creator:   "系统",
		},
		{
			Model:     gorm.Model{ID: systemId + 1},
			Name:      "User",
			Title:     "用户管理",
			Icon:      &userIcon,
			Path:      "user",
			Component: "/system/user/index",
			Sort:      systemId + 1,
			ParentId:  &systemId,
			Roles:     roles[:1],
			Creator:   "系统",
		},
		{
			Model:     gorm.Model{ID: systemId + 2},
			Name:      "Role",
			Title:     "角色管理",
			Icon:      &roleIcon,
			Path:      "role",
			Component: "/system/role/index",
			Sort:      systemId + 2,
			ParentId:  &systemId,
			Roles:     roles[:1],
			Creator:   "系统",
		},
		{
			Model:     gorm.Model{ID: systemId + 3},
			Name:      "Menu",
			Title:     "菜单管理",
			Icon:      &menuIcon,
			Path:      "menu",
			Component: "/system/menu/index",
			Sort:      systemId + 3,
			ParentId:  &systemId,
			Roles:     roles[:1],
			Creator:   "系统",
		},
		{
			Model:     gorm.Model{ID: systemId + 4},
			Name:      "Api",
			Title:     "接口管理",
			Icon:      &apiIcon,
			Path:      "api",
			Component: "/system/api/index",
			Sort:      systemId + 4,
			ParentId:  &systemId,
			Roles:     roles[:1],
			Creator:   "系统",
		},
		{
			Model:     gorm.Model{ID: systemId + 5},
			Name:      "OperationLog",
			Title:     "操作日志",
			Icon:      &apiLogIcon,
			Path:      "operation-log",
			Component: "/system/operation-log/index",
			Sort:      systemId + 5,
			ParentId:  &systemId,
			Roles:     roles[:1],
			Creator:   "系统",
		},
	}
	for _, menu := range menus {
		err := DB.First(&menu, menu.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newMenus = append(newMenus, menu)
		}
	}
	if len(newMenus) > 0 {
		err := DB.Create(&newMenus).Error
		if err != nil {
			Log.Errorf("写入系统菜单数据失败：%v", err)
		}
	}

	// 3.写入用户
	newUsers := make([]model.User, 0)
	users := []model.User{
		{
			Model:        gorm.Model{ID: 1},
			Username:     "root",
			Password:     util.GenPasswd("123456"),
			Mobile:       "13188887777",
			Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
			Introduction: new(string),
			Status:       1,
			Creator:      "系统",
			Roles:        roles[:1],
		},
		{
			Model:        gorm.Model{ID: 2},
			Username:     "admin",
			Password:     util.GenPasswd("123456"),
			Mobile:       "13188888888",
			Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
			Introduction: new(string),
			Status:       1,
			Creator:      "系统",
			Roles:        roles[:1],
		},
	}

	for _, user := range users {
		err := DB.First(&user, user.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUsers = append(newUsers, user)
		}
	}

	if len(newUsers) > 0 {
		err := DB.Create(&newUsers).Error
		if err != nil {
			Log.Errorf("写入用户数据失败：%v", err)
		}
	}

	// 4.写入api
	apis := []model.Api{
		{
			Method:   "POST",
			Path:     "/base/login",
			Category: "base",
			Desc:     "用户登录",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/base/logout",
			Category: "base",
			Desc:     "用户登出",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/base/refreshToken",
			Category: "base",
			Desc:     "刷新JWT令牌",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/base/report",
			Category: "base",
			Desc:     "收集日志",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/base/remote/config",
			Category: "base",
			Desc:     "获取配置",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/base/file/download",
			Category: "file",
			Desc:     "文件下载",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/user/info",
			Category: "user",
			Desc:     "获取当前登录用户信息",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/user/list",
			Category: "user",
			Desc:     "获取用户列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/user/changePwd",
			Category: "user",
			Desc:     "更新用户登录密码",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/user/create",
			Category: "user",
			Desc:     "创建用户",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/user/update/:userId",
			Category: "user",
			Desc:     "更新用户",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/user/delete/batch",
			Category: "user",
			Desc:     "批量删除用户",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/role/list",
			Category: "role",
			Desc:     "获取角色列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/role/create",
			Category: "role",
			Desc:     "创建角色",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/role/update/:roleId",
			Category: "role",
			Desc:     "更新角色",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/role/menus/get/:roleId",
			Category: "role",
			Desc:     "获取角色的权限菜单",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/role/menus/update/:roleId",
			Category: "role",
			Desc:     "更新角色的权限菜单",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/role/apis/get/:roleId",
			Category: "role",
			Desc:     "获取角色的权限接口",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/role/apis/update/:roleId",
			Category: "role",
			Desc:     "更新角色的权限接口",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/role/delete/batch",
			Category: "role",
			Desc:     "批量删除角色",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/menu/list",
			Category: "menu",
			Desc:     "获取菜单列表",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/menu/tree",
			Category: "menu",
			Desc:     "获取菜单树",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/menu/create",
			Category: "menu",
			Desc:     "创建菜单",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/menu/update/:menuId",
			Category: "menu",
			Desc:     "更新菜单",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/menu/delete/batch",
			Category: "menu",
			Desc:     "批量删除菜单",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/menu/access/list/:userId",
			Category: "menu",
			Desc:     "获取用户的可访问菜单列表",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/menu/access/tree/:userId",
			Category: "menu",
			Desc:     "获取用户的可访问菜单树",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/api/list",
			Category: "api",
			Desc:     "获取接口列表",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/api/tree",
			Category: "api",
			Desc:     "获取接口树",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/api/create",
			Category: "api",
			Desc:     "创建接口",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/api/update/:roleId",
			Category: "api",
			Desc:     "更新接口",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/api/delete/batch",
			Category: "api",
			Desc:     "批量删除接口",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/log/operation/list",
			Category: "log",
			Desc:     "获取操作日志列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/log/operation/delete/batch",
			Category: "log",
			Desc:     "批量删除操作日志",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/config/create",
			Category: "config",
			Desc:     "创建配置",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/config/list",
			Category: "config",
			Desc:     "获取配置列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/config/update",
			Category: "config",
			Desc:     "更新配置",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/config/update/status",
			Category: "config",
			Desc:     "更新配置状态",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/config/update/default",
			Category: "config",
			Desc:     "更新配置默认状态",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/config/delete/batch",
			Category: "config",
			Desc:     "批量删除配置",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/config/delete",
			Category: "config",
			Desc:     "删除配置",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/config/push",
			Category: "config",
			Desc:     "批量推送配置",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/config/copy",
			Category: "config",
			Desc:     "复制一条配置",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/config/module/list",
			Category: "config",
			Desc:     "获取模块列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/config/export",
			Category: "config",
			Desc:     "导出配置",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/config/import",
			Category: "config",
			Desc:     "导入配置",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/config/history/list",
			Category: "config",
			Desc:     "获取配置历史记录",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/config/history/data",
			Category: "config",
			Desc:     "获取配置历史信息",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/config/sync",
			Category: "config",
			Desc:     "同步配置策略",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/module/create",
			Category: "module",
			Desc:     "创建模块",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/module/list",
			Category: "module",
			Desc:     "获取模块列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/module/update",
			Category: "module",
			Desc:     "更新模块",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/module/upgrade",
			Category: "module",
			Desc:     "升级模块版本",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/module/update/status",
			Category: "module",
			Desc:     "更新模块状态",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/module/delete/batch",
			Category: "module",
			Desc:     "批量删除模块",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/module/delete",
			Category: "module",
			Desc:     "删除指定模块",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/module/component/list",
			Category: "module",
			Desc:     "获取组件列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/module/component/delete",
			Category: "module",
			Desc:     "删除组件列表",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/host/list",
			Category: "host",
			Desc:     "获取实例列表",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/host/detail",
			Category: "host",
			Desc:     "获取实例详情",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/host/delete/batch",
			Category: "host",
			Desc:     "批量删除实例列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/host/push/config",
			Category: "host",
			Desc:     "批量推送配置",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/host/update",
			Category: "host",
			Desc:     "更新主机配置",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/host/add",
			Category: "host",
			Desc:     "新增主机配置",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/attack/list",
			Category: "attack",
			Desc:     "获取攻击日志列表",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/attack/detail",
			Category: "attack",
			Desc:     "获取攻击日志详情",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/attack/delete/batch",
			Category: "attack",
			Desc:     "批量删除日志列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/attack/update",
			Category: "attack",
			Desc:     "更新攻击处理状态",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/process/list",
			Category: "process",
			Desc:     "获取进程日志列表",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/dashboard/attackData",
			Category: "dashboard",
			Desc:     "获取攻击数据",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/dashboard/attackTrends",
			Category: "dashboard",
			Desc:     "获取攻击趋势",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/dashboard/attackTypes",
			Category: "dashboard",
			Desc:     "获取攻击类型",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/ws/:hostName",
			Category: "ws",
			Desc:     "socket获取配置",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/file/upload",
			Category: "file",
			Desc:     "文件上传",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/file/list",
			Category: "file",
			Desc:     "文件列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/file/delete/batch",
			Category: "file",
			Desc:     "文件删除",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/file/getFileInfo/module",
			Category: "file",
			Desc:     "获取模块基本信息",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/rasp-log/list",
			Category: "logs",
			Desc:     "获取rasp日志信息",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/rasp-log/delete/batch",
			Category: "logs",
			Desc:     "批量删除rasp日志",
			Creator:  "系统",
		},
	}
	newApi := make([]model.Api, 0)
	newRoleCasbin := make([]model.RoleCasbin, 0)
	for i, api := range apis {
		api.ID = uint(i + 1)
		err := DB.First(&api, api.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newApi = append(newApi, api)

			// 管理员拥有所有API权限
			newRoleCasbin = append(newRoleCasbin, model.RoleCasbin{
				Keyword: roles[0].Keyword,
				Path:    api.Path,
				Method:  api.Method,
			})

			// 非管理员拥有基础权限
			basePaths := []string{
				"/base/login",
				"/base/logout",
				"/base/refreshToken",
				"/base/report",
				"/base/remote/config",
				"/user/info",
				"/menu/access/tree/:userId",
				"/ws/:hostName",
			}

			if funk.ContainsString(basePaths, api.Path) {
				newRoleCasbin = append(newRoleCasbin, model.RoleCasbin{
					Keyword: roles[1].Keyword,
					Path:    api.Path,
					Method:  api.Method,
				})
				newRoleCasbin = append(newRoleCasbin, model.RoleCasbin{
					Keyword: roles[2].Keyword,
					Path:    api.Path,
					Method:  api.Method,
				})
			}
		}
	}

	if len(newApi) > 0 {
		if err := DB.Create(&newApi).Error; err != nil {
			Log.Errorf("写入api数据失败：%v", err)
		}
	}

	if len(newRoleCasbin) > 0 {
		rules := make([][]string, 0)
		for _, c := range newRoleCasbin {
			rules = append(rules, []string{
				c.Keyword, c.Path, c.Method,
			})
		}
		isAdd, err := CasbinEnforcer.AddPolicies(rules)
		if !isAdd {
			Log.Errorf("写入casbin数据失败：%v", err)
		}
	}
}
