import type { ReactElement } from 'react';
import { Fragment, useEffect, useState } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';

export type WebSocketConnectProps = {
  url: string;
  sendMessage: boolean;
  message: any;
  onConnected: () => void;
  onMessageSent: () => void;
  onReadyStateChange: (readyState: ReadyState) => void;
  onMessageReceived: (response: any) => void;
};

const WebSocketConnect = (props: WebSocketConnectProps): ReactElement<any> => {
  const [connect] = useState(true);
  const [sendMessage, setSendMessage] = useState(props.sendMessage);
  const websocketService = useWebSocket(
    props.url,
    {
      onOpen: () => props.onConnected(),
      onClose: () => console.log('websocket链接关闭'),
      onError: (event: any) => console.log(event),
      onMessage: (event: any) => {
        console.log('服务器响应: ' + event.data);
        const response = JSON.parse(event.data);
        props.onMessageReceived(response);
      },
      share: true,
      shouldReconnect: () => true,
      reconnectAttempts: 10,
      reconnectInterval: 3000,
    },
    connect,
  );
  useEffect(() => {
    console.log(`websocket状态:  ${ReadyState[websocketService.readyState]}`);
  }, [websocketService.readyState, props.url]);

  useEffect(() => {
    if (props.sendMessage !== sendMessage) {
      setSendMessage(props.sendMessage);
    }
  }, [props.sendMessage, sendMessage]);

  useEffect(() => {
    const intervalId = setInterval(() => {
      console.log('心跳检测');
      websocketService.sendMessage('ping');
    }, 15000);

    return () => clearInterval(intervalId);
  }, [props.url, websocketService]);

  useEffect(() => {
    if (props.sendMessage && props.message) {
      websocketService.sendJsonMessage(props.message);
      props.onMessageSent();
    }
  }, [props, props.message, props.sendMessage, websocketService]);

  return <Fragment />;
};

WebSocketConnect.defaultProps = {
  url: '',
  sendMessage: false,
  message: undefined,
  onMessageSent: () => {},
  onReadyStateChange: () => {},
  onMessageReceived: () => {},
  onConnected: () => () => {},
};

export default WebSocketConnect;
