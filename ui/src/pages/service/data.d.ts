export type TableListItem = {
  id?: number;
  service_name?: string;
  low_period?: number;
  tag?: string;
  level?: boolean;
  owners?: string;
  language?: string;
  create_time?: string;
  update_time?: string;
};

// 查询/返回的服务信息
export type ServiceInfo = {
  serviceName: string;
  description: string;
  language: string;
  nodeNumber: number;
  lowPeriod: number;
  organization: string;
};

// 提交的服务信息
export type AddServiceInfo = {
  service_name: string;
  description: string;
  language: string;
  organization: string;
};


export type TableListItemResponse = {
  list: TableListItem[];
  page_num: number;
  pageSize: number;
  total: number;
};

export type TableListParams = {
  serviceName?: string;
  pageSize?: number;
  currentPage?: number;
  filter?: Record<string, any[]>;
  sorter?: Record<string, any>;
};

export type SimpleConfigItem = {
  label: string;
  value: number;
};
