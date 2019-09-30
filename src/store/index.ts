import uuid from 'uuid'
import {
  SET_PANE,
  ADD_PANE,
  DEL_PANE,

  INC_TOPZ,

  SET_REPO,
  DEL_REPO,
  SAVE_REPOS,
} from './mutationTypes'
import { StoreType, PaneType } from './types'
import {get, set} from '../persistence'

const REPOS_KEY = 'repos'

export const state = () => {
  const repos = get(REPOS_KEY) || {}
  return {repos}
}
export const mutations = {
  [SET_PANE](state: StoreType, payload: PaneType) {
    const {path, id} = payload
    const repos = state[REPOS_KEY]
    const repo = repos[path]
    const pane = repo.panes[id]
    repo.panes[id] = {... pane, ... payload}
    state[REPOS_KEY] = {... repos}
  },
  [ADD_PANE](state: StoreType, payload: PaneType) {
    const {path} = payload
    const repos = state[REPOS_KEY]
    const repo = repos[path]
    const id = uuid.v1()
    const pane = {id, width: 400, height: 200, x: 100, y: 100, ... payload}
    repo.panes = {... repo.panes, [id]: pane}
    state[REPOS_KEY] = {... repos}
  },
  [DEL_PANE](state: StoreType, payload: {id: string, path: string}) {
    const {path, id} = payload
    const repos = state[REPOS_KEY]
    const repo = repos[path]
    delete repo.panes[id]
    state[REPOS_KEY] = {... repos}
  },
  [INC_TOPZ](state: StoreType, payload: {path: string}) {
    const {path} = payload
    const repos = state[REPOS_KEY]
    const repo = repos[path]
    repo.topZ++
    state[REPOS_KEY] = {... repos}
  },
  [SET_REPO](state: StoreType, payload: {path: string}) {
    const {path} = payload
    const repos = state[REPOS_KEY]
    repos[path] = {... repos[path], ... payload}
    state[REPOS_KEY] = {... repos}
  },
  [DEL_REPO](state: StoreType, payload: {path: string}) {
    const {path} = payload
    const repos = state[REPOS_KEY]
    delete repos[path]
    state[REPOS_KEY] = {... repos}
  },
  [SAVE_REPOS](state: StoreType) {
    set(REPOS_KEY, state[REPOS_KEY])
  },
}
export const actions = {

}
export const getters = {}
