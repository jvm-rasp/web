import { request } from 'umi';

import type { TableListItem, TableListItemResponse, TableListParams } from './data';
import { ConfigDetail, SimpleConfigItem } from './data';

export async function queryConfig(params?: TableListParams) {
  return request<API.Response<TableListItemResponse>>(`/api/config/index`, {
    method: 'GET',
    params,
  });
}

export async function deleteConfig(id: number) {
  return request(`/api/config/delete`, {
    method: 'GET',
    params: {
      id,
    },
  });
}

export async function copyConfig(id: number) {
  return request(`/api/config/copy`, {
    method: 'GET',
    params: {
      id,
    },
  });
}

export async function addConfig(payload: ConfigDetail) {
  return request<API.Response<any>>(`/api/config/create`, {
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

// 获取配置的信息
export async function getSimpleConfigList() {
  return request<API.Response<SimpleConfigItem[]>>(`/api/config/list/simple`, {
    method: 'GET',
    params: {},
  });
}
