export default [
  {
    path: '/user',
    layout: false,
    routes: [
      {
        path: '/user',
        routes: [
          {
            name: 'login',
            path: '/user/login',
            component: './user/Login',
          },
        ],
      },
    ],
  },
  {
    access: 'host',
    path: '/host/index',
    component: './host/index',
  },
  {
    access: 'hostDetail',
    path: '/host/detail/index',
    component: './host/detail/index',
  },
  {
    access: 'attack',
    path: '/attack',
    component: './attack/index',
  },
  {
    access: 'attackDetail',
    path: '/attack/detail/index',
    component: './attack/detail/index',
  },
  {
    access: 'config',
    path: '/config',
    component: './config/index',
  },
  {
    access: 'configDetail',
    path: '/config/detail/index',
    component: './config/detail/index',
  },
  {
    access: 'service',
    path: '/service',
    component: './service/index',
  },
  {
    access: 'module',
    path: '/module',
    component: './module/index',
  },
  {
    access: 'dashboard',
    path: '/dashboard',
    component: './dashboard/index',
  },
  {
    access: 'dependency',
    path: '/dependency',
    component: './dependency/index',
  },
  {
    access: 'Rbac',
    path: '/rbac',
    routes: [
      {
        access: 'rbacRoleIndex',
        path: '/rbac/role/index',
        component: './rbac/role/index',
      },
      {
        access: 'rbacRolePermissionIndex',
        path: '/rbac/role/permission/index',
        component: './rbac/role/permission/index',
      },
      {
        access: 'rbacPermissionIndex',
        path: '/rbac/permission/index',
        component: './rbac/permission/index',
      },
      {
        access: 'rbacUserIndex',
        path: '/rbac/user/index',
        component: './rbac/user/index',
      },
    ],
  },
  {
    access: 'index',
    path: '/',
    component: './index/index',
  },
  {
    access: 'dashboard',
    path: '/',
    component: './dashboard/index',
  },
  {
    component: './404',
  },
];
