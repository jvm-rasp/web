import request from '@/utils/request'

// 获取配置列表
export function getConfigs(params) {
  return request({
    url: '/config/list',
    method: 'get',
    params
  })
}

// 创建新配置
export function createConfig(params) {
  return request({
    url: '/config/create',
    method: 'post',
    data: params
  })
}

// 创建新配置
export function updateConfig(params) {
  return request({
    url: '/config/update',
    method: 'post',
    data: params
  })
}

export function deleteConfig(params) {
  return request({
    url: '/config/delete',
    method: 'post',
    data: params
  })
}

// 批量删除配置
export function batchDeleteConfigByIds(params) {
  return request({
    url: '/config/delete/batch',
    method: 'post',
    data: params
  })
}

export function copyConfigById(params) {
  return request({
    url: '/config/copy',
    method: 'post',
    data: params
  })
}

export function updateStatusById(params) {
  return request({
    url: '/config/update/status',
    method: 'post',
    data: params
  })
}

export function updateDefaultById(params) {
  return request({
    url: '/config/update/default',
    method: 'post',
    data: params
  })
}

export function pushConfigById(params) {
  return request({
    url: '/config/push',
    method: 'post',
    data: params
  })
}

export function exportConfigById(params) {
  return request({
    url: '/config/export',
    method: 'post',
    data: params,
    responseType: 'blob'
  })
}

export function getModules(params) {
  return request({
    url: '/config/module/list',
    method: 'get',
    params
  })
}

export function getConfigHistoryByGuid(params) {
  return request({
    url: '/config/history/list',
    method: 'get',
    params
  })
}

export function getConfigHistoryDataByGuid(params) {
  return request({
    url: '/config/history/data',
    method: 'get',
    params
  })
}

export function syncConfigById(params) {
  return request({
    url: '/config/sync',
    method: 'post',
    data: params
  })
}
