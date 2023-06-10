import request from '@/utils/request'

export function updateSetting(name, value) {
  return request({
    url: '/settings/update',
    method: 'post',
    data: { name: name, value: value }
  })
}

export function getSettings() {
  return request({
    url: '/settings/list',
    method: 'get'
  })
}

export function getProjectGuid(params) {
  return request({
    url: '/settings/getProjectInfo',
    method: 'post',
    data: params
  })
}
