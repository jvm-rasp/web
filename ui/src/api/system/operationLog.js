import request from '@/utils/request'

// 获取操作日志列表
export function getOperationLogs(params) {
  return request({
    url: '/log/operation/list',
    method: 'get',
    params
  })
}

// 批量删除操作日志
export function batchDeleteOperationLogByIds(data) {
  return request({
    url: '/log/operation/delete/batch',
    method: 'post',
    data
  })
}

