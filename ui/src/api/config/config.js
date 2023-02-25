import request from '@/utils/request'

// 获取配置列表
export function getConfigs(params) {
  return request({
    url: '/api/config/list',
    method: 'get',
    params
  })
}

// 创建新配置
export function createConfig(params) {
  return request({
    url: '/api/config/create',
    method: 'post',
    params
  })
}

// 批量删除配置
export function batchDeleteConfigByIds(data) {
  return request({
    url: '/api/config/delete/batch',
    method: 'delete',
    data
  })
}

