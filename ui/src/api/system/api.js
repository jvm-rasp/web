import request from '@/utils/request'

// 获取接口列表
export function getApis(params) {
  return request({
    url: '/api/list',
    method: 'get',
    params
  })
}

// 获取接口树(按接口Category字段分类)
export function getApiTree(params) {
  return request({
    url: '/api/tree',
    method: 'get',
    params
  })
}

// 创建接口
export function createApi(data) {
  return request({
    url: '/api/create',
    method: 'post',
    data
  })
}

// 更新接口
export function updateApiById(Id, data) {
  return request({
    url: '/api/update/' + Id,
    method: 'post',
    data
  })
}

// 批量删除接口
export function batchDeleteApiByIds(data) {
  return request({
    url: '/api/delete/batch',
    method: 'post',
    data
  })
}
