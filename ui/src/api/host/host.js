import request from '@/utils/request'

export function getHosts(params) {
  return request({
    url: '/api/host/list',
    method: 'get',
    params
  })
}

export function deleteHost(params) {
  return request({
    url: '/api/host/delete',
    method: 'post',
    params
  })
}
