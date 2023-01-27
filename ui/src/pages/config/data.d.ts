export type TableListItem = {
  id?: number;
  config_name: string;
  message: string;
  config_content: string;
  tag: string;
  username: string;
  create_time?: string;
  update_time?: string;
  status: number;
};
// 策略配置
export type ConfigDetail = {
  configName: string;
  modules: number[];
  agentMode: string;
  binFileUrl: string;
  binFileHash: string;
  message: string;
  // v1.1.1 新增
  // 必须显示配置
  nacosAddrs: string;
  // 关闭检测功能
  check_disable: boolean;
  // 阻断配置，可选
  redirect_url: string;
  block_status_code: number;
  json_block_content: string;
  xml_block_content: string;
  html_block_content: string;
};

export type TableListItemResponse = {
  list: TableListItem[];
  page_num: number;
  pageSize: number;
  total: number;
};

export type TableListPagination = {
  total: number;
  pageSize: number;
  current: number;
};

export type TableListParams = {
  configName?: string;
  status?: number;
  pageSize?: number;
  currentPage?: number;
  filter?: Record<string, any[]>;
  sorter?: Record<string, any>;
};

export type SimpleConfigItem = {
  label: string;
  value: number;
};
