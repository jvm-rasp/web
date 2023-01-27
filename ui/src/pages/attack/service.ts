import { request } from 'umi';

import type { TableListItemResponse, TableListParams } from './data';

export async function queryAttack(params?: TableListParams) {
  return request<API.Response<TableListItemResponse>>(`/api/attack/index`, {
    method: 'GET',
    params,
  });
}

export async function deleteAttack(id: number) {
  return request(`/api/attack/delete`, {
    method: 'GET',
    params: {
      id,
    },
  });
}

export async function batchDeleteAttack(ids?: (number | string)[]) {
  return request(`/api/attack/batch/delete`, {
    method: 'POST',
    params: {
      ids,
    },
  });
}

export async function markAttack(id: number | undefined, handleStatus: number) {
  return request(`/api/attack/mark`, {
    method: 'GET',
    params: {
      id,
      handleStatus,
    },
  });
}
