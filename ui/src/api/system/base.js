import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/base/login',
    method: 'post',
    data
  })
}

export function refreshToken() {
  return request({
    url: '/base/refreshToken',
    method: 'post'
  })
}

export function logout() {
  return request({
    url: '/base/logout',
    method: 'post'
  })
}
