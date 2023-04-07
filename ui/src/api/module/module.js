import request from '@/utils/request'

export function getModules(params) {
  return request({
    url: '/module/list',
    method: 'get',
    params
  })
}

export function createModule(params) {
  return request({
    url: '/module/create',
    method: 'post',
    data: params
  })
}

export function updateModule(params) {
  return request({
    url: '/module/update',
    method: 'post',
    data: params
  })
}

export function deleteModule(params) {
  return request({
    url: '/module/delete',
    method: 'post',
    data: params
  })
}

export function batchDeleteModuleByIds(params) {
  return request({
    url: '/module/delete/batch',
    method: 'delete',
    data: params
  })
}

export function updateStatusById(params) {
  return request({
    url: '/module/update/status',
    method: 'post',
    data: params
  })
}

export function getUploadFiles(params) {
  return request({
    url: '/file/list',
    method: 'get',
    params
  })
}

export function batchDeleteFilesByIds(params) {
  return request({
    url: '/file/delete/batch',
    method: 'post',
    data: params
  })
}

export function getModuleInfoById(params) {
  return request({
    url: '/file/getFileInfo/module',
    method: 'get',
    params
  })
}
