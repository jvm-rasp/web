import request from '@/utils/request'

export function getAttackLogs(params) {
  return request({
    url: '/api/attack/list',
    method: 'get',
    params
  })
}

export function deleteAttackLog(params) {
  return request({
    url: '/api/attack/delete/batch',
    method: 'post',
    data: params
  })
}

export function getAttackDetail(params) {
  return request({
    url: '/api/attack/detail',
    method: 'get',
    params
  })
}
