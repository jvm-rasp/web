import request from '@/utils/request'

export function getHosts(params) {
  return request({
    url: '/api/host/list',
    method: 'get',
    params
  })
}

export function batchDeleteHost(params) {
  return request({
    url: '/api/host/delete/batch',
    method: 'post',
    data: params
  })
}

export function updateHost(params) {
  return request({
    url: '/api/host/update',
    method: 'post',
    data: params
  })
}

export function getDetail(params) {
  return request({
    url: '/api/host/detail',
    method: 'get',
    params
  })
}

export function getProcesss(params) {
  return request({
    url: '/api/process/list',
    method: 'get',
    params
  })
}
