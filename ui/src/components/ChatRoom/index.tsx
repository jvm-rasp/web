import React from 'react';
import ProCard from '@ant-design/pro-card';

const ChatRoom: React.FC = () => {
  return (
    <ProCard split="vertical">
      <ProCard title="左侧详情" colSpan="20%">
        左侧内容
      </ProCard>
      <ProCard title="左右分栏子卡片带标题" headerBordered>
        <div style={{ height: '100%' }}>右侧内容</div>
      </ProCard>
    </ProCard>
  );
};
export default ChatRoom;
