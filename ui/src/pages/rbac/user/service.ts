import { request } from 'umi';

import type { TableListItem, TableListItemResponse, TableListParams } from './data';

export async function queryUser(params?: TableListParams) {
  return request<API.Response<TableListItemResponse>>(`/api/rbac/user/index`, {
    method: 'GET',
    params,
  });
}
export async function addUser(payload: TableListItem) {
  return request<API.Response<any>>(`/api/rbac/user/create`, {
    method: 'POST',
    data: {
      ...payload,
    },
  });
}

export async function updateUser(payload: TableListItem) {
  return request(`/api/rbac/user/update`, {
    method: 'POST',
    data: {
      ...payload,
    },
  });
}

export async function statusUser(payload: { id: number; status: number }) {
  return request(`/api/rbac/user/status`, {
    method: 'POST',
    data: {
      ...payload,
    },
  });
}

export async function deleteUser(id: number) {
  return request(`/api/rbac/user/delete`, {
    method: 'GET',
    params: {
      id,
    },
  });
}
