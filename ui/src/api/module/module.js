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
    data: params
  })
}

export function updateModule(params) {
  return request({
    url: '/api/module/update',
    method: 'post',
    data: params
  })
}

export function deleteModule(params) {
  return request({
    url: '/api/module/delete',
    method: 'post',
    data: params
  })
}

export function batchDeleteModuleByIds(params) {
  return request({
    url: '/api/module/delete/batch',
    method: 'get',
    data: params
  })
}

