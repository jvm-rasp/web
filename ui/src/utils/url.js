
export function getRootPath() {
  const curWwwPath = window.document.location.href
  const pathName = window.document.location.pathname
  const pos = curWwwPath.indexOf(pathName)
  const localhostPath = curWwwPath.substring(0, pos)
  const projectName = pathName.substring(0, pathName.substr(1).indexOf('/') + 1)
  return (localhostPath + projectName)
}
