import { PanePropsListType, MainState } from './types'

export const SET_PANE_LIST = 'setPaneList'
export const SET_SERIAL = 'setSerial'
export const INCREMENT_SERIAL = 'incrementSerial'

export const state = () => ({
  serial: 1,
  paneList: [
    {command: 'log --graph --pretty=format:"%h %ad %s %d" --all --date=format:"%y%m%d %H:%M"', x: 100, y: 0, width: 500, height: 200, style: {}},
    {command: "status --short", x: 0, y: 200, width: 200, height: 100, style: {}},
    {command: "branch", x: 0, y: 400, width: 200, height: 100, style: {}},
    {command: "tag", x: 0, y: 500, width: 200, height: 100, style: {}},
    {command: "remote", x: 0, y: 600, width: 200, height: 100, style: {}},

  ]
})
export const actions = {}
export const mutations = {
  [SET_PANE_LIST](state: MainState, paneList: PanePropsListType) {
    state.paneList = paneList
  },
  [SET_SERIAL](state: MainState, serial: number) {
    state.serial = serial
  },
  [INCREMENT_SERIAL](state: MainState) {
    console.log('state:', state);
    state.serial = state.serial + 1
    return state.serial
  },
}
export const getters = {}
