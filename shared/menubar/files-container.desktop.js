// @flow
import * as ConfigGen from '../actions/config-gen'
import * as Tabs from '../constants/tabs'
import * as FsTypes from '../constants/types/fs'
import * as FsGen from '../actions/fs-gen'
import {switchTo} from '../actions/route-tree'
import {FilesPreview} from './files.desktop'
import {connect, compose, type Dispatch} from '../util/container'

const mapStateToProps = (state) => ({
  _tlfRows: [
    {path: FsTypes.stringToPath('/keybase/team/zila.test/abc')},
    {path: FsTypes.stringToPath('/keybase/team/zila.test/def')},
  ],
})

const mapDispatchToProps = (dispatch: Dispatch) => ({
  _onSelectPath: (path: FsTypes.Path) => {
    dispatch(ConfigGen.createShowMain())
    dispatch(switchTo([Tabs.fsTab]))
    dispatch(FsGen.createOpenPathItem({
      openDirectly: true,
      path,
    }))
  },
  onViewAll: () => {
    dispatch(ConfigGen.createShowMain())
    dispatch(switchTo([Tabs.fsTab]))
  },
})

const mergeProps = (stateProps, dispatchProps) => ({
  onViewAll: dispatchProps.onViewAll,
  tlfRows: stateProps._tlfRows.map(c => ({
    onSelectPath: () => dispatchProps._onSelectPath(c.path),
    path: FsTypes.pathToString(c.path),
  })),
})

export default compose(
  connect(mapStateToProps, mapDispatchToProps, mergeProps)
)(FilesPreview)