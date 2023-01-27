import React from 'react';
import {useHistory} from 'umi';
import {PageHeader} from 'antd';
import ProDescriptions from "@ant-design/pro-descriptions";


const ConfigDetail: React.FC = () => {
  const history = useHistory();
  // @ts-ignore
  const {config} = history.location.state;
  return (
    <div className="site-page-header-ghost-wrapper">
      <PageHeader ghost={false} onBack={() => window.history.back()} title="配置详情">
        <ProDescriptions
          column={2}
          title="配置详情"
          tooltip="包含了配置详情等"
        >
          <ProDescriptions.Item label="配置名称">{config.config_name}</ProDescriptions.Item>
        </ProDescriptions>
        <ProDescriptions>
          <ProDescriptions.Item valueType="jsonCode">
            {config.config_content}
          </ProDescriptions.Item>
        </ProDescriptions>
      </PageHeader>
    </div>
  );
};

export default ConfigDetail;
