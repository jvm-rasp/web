import { request } from 'umi';

import type { TableListItem, TableListItemResponse, TableListParams } from './data';
import { ModuleTypeItem } from './data';

export async function queryModule(params?: TableListParams) {
  return request<API.Response<TableListItemResponse>>(`/api/module/index`, {
    method: 'GET',
    params,
  });
}

export async function deleteModule(id: number) {
  return request(`/api/module/delete`, {
    method: 'GET',
    params: {
      id,
    },
  });
}

export async function addModule(payload: TableListItem) {
  return request(`/api/module/create`, {
    method: 'POST',
    data: {
      ...payload,
    },
  });
}

export async function getModule(id: number) {
  return request(`/api/module/detail`, {
    method: 'GET',
    params: {
      id,
    },
  });
}

export async function updateModule(payload: TableListItem) {
  return request(`/api/module/update`, {
    method: 'POST',
    data: {
      ...payload,
    },
  });
}

export async function statusModule(payload: { id: number; status: number }) {
  return request(`/api/module/status`, {
    method: 'POST',
    data: {
      ...payload,
    },
  });
}

export async function getModuleType() {
  return request<API.Response<ModuleTypeItem[]>>(`/api/module/type`, {
    method: 'GET',
    params: {},
  });
}
