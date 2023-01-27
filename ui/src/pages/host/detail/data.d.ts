export type RaspJavaInfoItem = {
  id: number;
  pid: number;
  host_name: number;
  cmdline_info: string;
  updated_time: string;
};

// 发送给后端的数据
export type QueryParams = {
  hostName?: string;
  pageSize?: number;
  currentPage?: number;
};

// 请求返回端数据结构
export type ProListItemResponse = {
  list: RaspJavaInfoItem[];
  page_num: number;
  page_size: number;
  total: number;
};


type GithubIssueItem = {
  pageSize?: number;
  current?: number;
  hostName?: string;
};
