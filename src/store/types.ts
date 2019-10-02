
export interface PaneType {
  command: string,
  width: number,
  height: number,
  x: number,
  y: number,
  z: number,
  path: string,
  id: string,
}

export type PanesType = { [id: string]: PaneType }

export interface RepoType {
  panes: PanesType,
  topZ: number,
  lastAccess: string,
}

export type ReposType = { [path: string]: RepoType }

export interface StoreType {
  repos: ReposType,
  conn: number,
}
