
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

export interface RepoState {
  panes: PanesType,
  topZ: number,
}

export type ReposState = { [path: string]: RepoState }

export interface StoreState {
  repos: ReposState,
}
