# 常见问题

## 1.npm ERR network Invalid response body while trying to fetch

```shell
npm notice
npm notice New minor version of npm available! 8.5.0 -> 8.7.0
npm notice Changelog: https://github.com/npm/cli/releases/tag/v8.7.0
npm notice Run npm install -g npm@8.7.0 to update!
npm notice
npm ERR! code ERR_SOCKET_TIMEOUT
npm ERR! errno ERR_SOCKET_TIMEOUT
npm ERR! network Invalid response body while trying to fetch https://registry.npmjs.org/depd: Socket timeout
npm ERR! network This is a problem related to network connectivity.
npm ERR! network In most cases you are behind a proxy or have bad network settings.
npm ERR! network
npm ERR! network If you are behind a proxy, please make sure that the
npm ERR! network 'proxy' config is set properly.  See: 'npm help config'
```
通过 config 命令切换源
npm config set registry https://registry.npm.taobao.org