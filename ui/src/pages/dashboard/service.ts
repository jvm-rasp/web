import { request } from 'umi';

export async function fetchWeekData() {
  return request<WeekDataResult>(`/api/attack/data/week`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });
}
