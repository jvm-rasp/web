import { request } from 'umi';

import type { TableListItemResponse, TableListParams } from './data';
import { UpdateConfigDetail } from './data';

export async function queryHost(params?: TableListParams) {
  return request<API.Response<TableListItemResponse>>(`/api/host/index`, {
    method: 'GET',
    params,
  });
}

export async function deleteHost(id: number) {
  return request(`/api/host/delete`, {
    method: 'GET',
    params: {
      id,
    },
  });
}

export async function batchUpdateConfig(payload?: UpdateConfigDetail) {
  return request(`/api/config/batch/update`, {
    method: 'POST',
    params: {
      ...payload
    },
  });
}

export async function updateConfig(payload?: UpdateConfigDetail) {
  return request(`/api/config/update`, {
    method: 'POST',
    params: {
      ...payload,
    },
  });
}
