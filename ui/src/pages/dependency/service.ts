import { request } from 'umi';

import type { TableListItemResponse, TableListParams } from './data';

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
