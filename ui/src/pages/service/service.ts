import {request} from 'umi';

import type {TableListItemResponse, TableListParams, AddServiceInfo} from './data';

export async function queryService(params?: TableListParams) {
  return request<API.Response<TableListItemResponse>>(`/api/service/index`, {
    method: 'GET',
    params,
  });
}

export async function deleteService(id: number) {
  return request(`/api/service/delete`, {
    method: 'GET',
    params: {
      id,
    },
  });
}

export async function addService(payload: AddServiceInfo) {
  return request<API.Response<any>>(`/api/service/create`, {
    method: 'POST',
    data: {
      ...payload,
    },
  });
}

export async function updateService(payload: AddServiceInfo) {
  return request<API.Response<any>>(`/api/service/update`, {
    method: 'POST',
    data: {
      ...payload,
    },
  });
}
