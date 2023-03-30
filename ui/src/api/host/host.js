import request from '@/utils/request'

export function getHosts(params) {
  return request({
    url: '/host/list',
    method: 'get',
    params
  })
}

export function batchDeleteHost(params) {
  return request({
    url: '/host/delete/batch',
    method: 'post',
    data: params
  })
}

export function updateHost(params) {
  return request({
    url: '/host/update',
    method: 'post',
    data: params
  })
}

export function getDetail(params) {
  return request({
    url: '/host/detail',
    method: 'get',
    params
  })
}

export function getProcesss(params) {
  return request({
    url: '/process/list',
    method: 'get',
    params
  })
}

export function pushConfig(params) {
  return request({
    url: '/host/push/config',
    method: 'post',
    data: params
  })
}
