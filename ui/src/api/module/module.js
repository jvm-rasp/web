import request from '@/utils/request'

export function getModules(params) {
  return request({
    url: '/api/module/list',
    method: 'get',
    params
  })
}

export function createModule(params) {
  return request({
    url: '/api/module/create',
    method: 'post',
    params
  })
}

export function batchDeleteModuleByIds(data) {
  return request({
    url: '/api/module/delete/batch',
    method: 'delete',
    data
  })
}

