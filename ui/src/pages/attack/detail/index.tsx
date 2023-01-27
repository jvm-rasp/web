import React from 'react';
import { useHistory } from 'umi';
import { PageHeader, Button, message } from 'antd';
import ProDescriptions from '@ant-design/pro-descriptions';
import { markAttack } from '@/pages/attack/service';

const HostDetail: React.FC = () => {
  const history = useHistory();
  const handleMarkAttack = async (id: number, status: number) => {
    // 如果是处理状态，则标记为处理完成的状态
    const response = await markAttack(id, status);
    const { code, msg } = response;
    if (code !== 200) {
      message.error(msg);
      return;
    }
    message.success(msg);
  };
  // @ts-ignore
  const { attackInfo } = history.location.state;
  return (
    <div className="site-page-header-ghost-wrapper">
      <PageHeader ghost={false} onBack={() => window.history.back()} title="详情">
        <ProDescriptions
          column={2}
          title="攻击详情"
          tooltip="包含了服务器信息、攻击日志、检测算法等"
        >
          <ProDescriptions.Item label="文本" valueType="option">
            <Button
              type="primary"
              onClick={async () => {
                // 如果是处理状态，则标记为处理完成的状态
                handleMarkAttack(attackInfo.id, 2);
              }}
              key="confirm"
            >
              确认
            </Button>
            <Button
              key="misreport"
              onClick={async () => {
                // 如果是已经处理或者处理中的状态，则标记为误报
                handleMarkAttack(attackInfo.id, 3);
              }}
            >
              误报
            </Button>
            <Button
              key="ignore"
              onClick={async () => {
                // 如果是已经处理或者处理中的状态，则标记为忽略
                handleMarkAttack(attackInfo.id, 4);
              }}
            >
              忽略
            </Button>
          </ProDescriptions.Item>
          <ProDescriptions.Item label="主机名称">{attackInfo.host_name}</ProDescriptions.Item>
          <ProDescriptions.Item label="受攻击IP">{attackInfo.local_ip}</ProDescriptions.Item>
          <ProDescriptions.Item label="攻击IP">{attackInfo.remote_ip}</ProDescriptions.Item>
          <ProDescriptions.Item label="攻击时间">{attackInfo.attack_time}</ProDescriptions.Item>
          <ProDescriptions.Item label="危险等级">{attackInfo.level}</ProDescriptions.Item>
          <ProDescriptions.Item
            label="阻断状态"
            valueEnum={{
              true: {
                text: '阻断',
                status: 'Error',
              },
              false: {
                text: '放行',
                status: 'Success',
              },
            }}
          >
            {attackInfo.is_blocked}
          </ProDescriptions.Item>
          <ProDescriptions.Item label="请求协议">
            {attackInfo.request_protocol}
          </ProDescriptions.Item>
          <ProDescriptions.Item label="请求方法">{attackInfo.http_method}</ProDescriptions.Item>
          <ProDescriptions.Item label="请求uri">{attackInfo.request_uri}</ProDescriptions.Item>
          <ProDescriptions.Item label="请求参数">
            {attackInfo.request_parameters}
          </ProDescriptions.Item>
          <ProDescriptions.Item label="攻击参数" span={2}>
            {attackInfo.attack_parameters}
          </ProDescriptions.Item>
          <ProDescriptions.Item label="检测结果" span={2}>
            {attackInfo.check_type}
          </ProDescriptions.Item>
          <ProDescriptions.Item
            label="攻击调用栈"
            span={2}
            valueType="jsonCode"
            tooltip={'栈深度仅显示50层'}
          >
            {attackInfo.stack_trace}
          </ProDescriptions.Item>
          <ProDescriptions.Item label="body" span={2}  >
            {attackInfo.body}
          </ProDescriptions.Item>
          <ProDescriptions.Item label="header" span={2} valueType="jsonCode">
            {attackInfo.header}
          </ProDescriptions.Item>
        </ProDescriptions>
      </PageHeader>
    </div>
  );
};

export default HostDetail;
