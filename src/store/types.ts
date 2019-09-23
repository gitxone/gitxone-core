export interface StyleType {[s: string]: any}


export interface PanePropsType {
  command: string,
  x: number,
  y: number,
  z: number,
  width: number,
  height: number,
  style: StyleType,
}

export type PanePropsListType = PanePropsType[]

export interface MainState {
  paneList: PanePropsListType,
  serial: number,
}
