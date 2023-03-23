import request from '@/utils/request'

// 获取角色列表
export function getRoles(params) {
  return request({
    url: '/role/list',
    method: 'get',
    params
  })
}

// 创建角色
export function createRole(data) {
  return request({
    url: '/role/create',
    method: 'post',
    data
  })
}

// 更新角色
export function updateRoleById(roleId, data) {
  return request({
    url: '/role/update/' + roleId,
    method: 'get',
    data
  })
}

// 获取角色的权限菜单
export function getRoleMenusById(roleId) {
  return request({
    url: '/role/menus/get/' + roleId,
    method: 'get'
  })
}

// 更新角色的权限菜单
export function updateRoleMenusById(roleId, data) {
  return request({
    url: '/role/menus/update/' + roleId,
    method: 'get',
    data
  })
}

// 获取角色的权限接口
export function getRoleApisById(roleId) {
  return request({
    url: '/role/apis/get/' + roleId,
    method: 'get'
  })
}

// 更新角色的权限接口
export function updateRoleApisById(roleId, data) {
  return request({
    url: '/role/apis/update/' + roleId,
    method: 'get',
    data
  })
}

// 批量删除角色
export function batchDeleteRoleByIds(data) {
  return request({
    url: '/role/delete/batch',
    method: 'post',
    data
  })
}
