import request from '@/utils/request'

export function getAttackLogs(params) {
  return request({
    url: '/attack/list',
    method: 'get',
    params
  })
}

export function deleteAttackLog(params) {
  return request({
    url: '/attack/delete/batch',
    method: 'post',
    data: params
  })
}

export function getAttackDetail(params) {
  return request({
    url: '/attack/detail',
    method: 'get',
    params
  })
}

export function updateStatus(params) {
  return request({
    url: '/attack/update',
    method: 'post',
    data: params
  })
}

export function exportAttackLog(params) {
  return request({
    url: '/attack/export',
    method: 'post',
    data: params,
    responseType: 'blob'
  })
}
