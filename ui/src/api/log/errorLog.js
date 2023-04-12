import request from '@/utils/request'

export function getLogs(params) {
  return request({
    url: '/rasp-log/list',
    method: 'get',
    params
  })
}

export function batchDeleteLogsByIds(params) {
  return request({
    url: '/rasp-log/delete/batch',
    method: 'post',
    data: params
  })
}
