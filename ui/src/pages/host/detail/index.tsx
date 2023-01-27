import React from 'react';
import {useHistory} from 'umi';
import {PageHeader, Descriptions, Badge, Tag, Space} from 'antd';
import ProList from '@ant-design/pro-list';
import {GithubIssueItem} from '@/pages/host/detail/data';
import {queryHostJavaInfo} from "@/pages/host/detail/service";

import {BranchesOutlined} from '@ant-design/icons';
import {createFromIconfontCN} from '@ant-design/icons';


const IconFont = createFromIconfontCN({
  scriptUrl: [
    '//at.alicdn.com/t/font_1788044_0dwu4guekcwr.js', // icon-javascript, icon-java, icon-shoppingcart (overrided)
    '//at.alicdn.com/t/font_1788592_a5xf2bdic3u.js', // icon-shoppingcart, icon-python
  ],
});

const HostDetail: React.FC = () => {
  const history = useHistory();
  // @ts-ignore
  const {host} = history.location.state;
  return (
    <div className="site-page-header-ghost-wrapper">
      <PageHeader ghost={false} onBack={() => window.history.back()} title="详情">
        <Descriptions
          title="主机信息"
          bordered
          column={{xxl: 4, xl: 3, lg: 3, md: 3, sm: 2, xs: 1}}
        >
          <Descriptions.Item label="主机名称">{host.host_name}</Descriptions.Item>
          <Descriptions.Item label="IP地址">{host.ip}</Descriptions.Item>
          <Descriptions.Item label="操作系统">{host.os_type}</Descriptions.Item>
          <Descriptions.Item label="系统内存">{host.total_mem} GB</Descriptions.Item>
          <Descriptions.Item label="CPU核数">{host.cpu_counts}</Descriptions.Item>
          <Descriptions.Item label="可用磁盘">{host.free_disk} GB</Descriptions.Item>
        </Descriptions>
        <br/>
        <br/>
        <Descriptions
          title="守护进程信息"
          bordered
          column={{xxl: 4, xl: 3, lg: 3, md: 3, sm: 2, xs: 1}}
        >
          <Descriptions.Item label="运行状态" span={1}>
            <Badge status="processing" text="运行中"/>
          </Descriptions.Item>
          <Descriptions.Item label="版本">{host.version}</Descriptions.Item>
          <Descriptions.Item label="安装路径">{host.install_dir}</Descriptions.Item>
          <Descriptions.Item label="MD5">{host.exe_file_hash}</Descriptions.Item>
          <Descriptions.Item label="编译代码"><BranchesOutlined/> {host.build_git_branch}:{host.build_git_commit}
          </Descriptions.Item>
          <Descriptions.Item label="编译时间">{host.build_date_time}</Descriptions.Item>
        </Descriptions>

        <br/>
        <br/>
        <Descriptions
          title="配置信息"
          bordered
          column={{xxl: 4, xl: 3, lg: 3, md: 3, sm: 2, xs: 1}}
        >
          <Descriptions.Item label="nacos">{host.nacos_info}</Descriptions.Item>
        </Descriptions>

        <br/>
        <br/>
        <ProList<GithubIssueItem>
          bordered
          itemLayout="vertical"
          rowKey="name"
          headerTitle={'进程列表'}
          // @ts-ignore
          request={async (params) => {
            params.hostName = host.host_name
            const response = await queryHostJavaInfo(params);
            return {
              success: response.code === 200,
              data: response!.data!.list,
              page: response!.data!.page_num,
              total: response!.data!.total,
            };
          }}
          pagination={{
            pageSize: 1,
            pageSizeOptions: [1, 2, 5, 10],
            showSizeChanger: true,
          }}
          showActions="hover"
          metas={{
            title: {
              dataIndex: 'pid', // java 进程pid
              render: (dataIndex) => {
                return '进程PID: ' + dataIndex;
              },
            },
            subTitle: {
              dataIndex: 'status',   // 防护状态
              render: (dataIndex) => {
                if ('success inject' === dataIndex) {
                  return (
                    <Space size={0}>
                      <IconFont type="icon-java"/>
                      <Tag color="blue">开启防护</Tag>
                    </Space>
                  );
                }
                if ('not inject' === dataIndex) {
                  return (
                    <Space size={0}>
                      <Tag color="grey">未开启防护</Tag>
                    </Space>
                  );
                } else if ('success uninstall agent' === dataIndex) {
                  return (
                    <Space size={0}>
                      <Tag color="red">关闭防护</Tag>
                    </Space>
                  );
                }else {
                  return (
                    <Space size={0}>
                      <Tag color="red">Agent异常</Tag>
                    </Space>
                  );
                }
                return '进程PID: ' + dataIndex;
              }
            },
            description: {
              dataIndex: 'start_time', // 进程启动时间
              search: false,
              render: (dataIndex) => {
                return '启动时间: ' + dataIndex;
              },
            },
            content: {
              dataIndex: "cmdline_info",  // 进程的cmdlines
            }
          }
          }
        />

      </PageHeader>
    </div>
  );
};

export default HostDetail;
