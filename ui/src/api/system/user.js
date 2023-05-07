import request from '@/utils/request'

// 获取当前登录用户信息
export function getInfo() {
  return request({
    url: '/user/info',
    method: 'post'
  })
}

// 获取用户列表
export function getUsers(params) {
  return request({
    url: '/user/list',
    method: 'get',
    params
  })
}

// 更新用户登录密码
export function changePwd(data) {
  return request({
    url: '/user/changePwd',
    method: 'post',
    data
  })
}

// 创建用户
export function createUser(data) {
  return request({
    url: '/user/create',
    method: 'post',
    data
  })
}

// 更新用户
export function updateUserById(id, data) {
  return request({
    url: '/user/update/' + id,
    method: 'post',
    data
  })
}

// 批量删除用户
export function batchDeleteUserByIds(data) {
  return request({
    url: '/user/delete/batch',
    method: 'post',
    data
  })
}

