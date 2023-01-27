/**
 * @see https://umijs.org/zh-CN/plugins/plugin-access
 * */
export default function access(initialState: { currentUser?: API.CurrentUser | undefined }) {
  const { currentUser } = initialState || {};
  let accessData = {};
  if (currentUser) {
    accessData = currentUser!.access || {};
  }
  return accessData;
}
