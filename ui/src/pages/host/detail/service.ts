import {request} from 'umi';

import {GithubIssueItem, ProListItemResponse} from './data';

export async function queryHostJavaInfo(params?: GithubIssueItem) {
  return request<API.Response<ProListItemResponse>>(`/api/java/index`, {
    method: 'GET',
    params
  });
}
