import {AnyAction, Store} from 'redux';
import {ThunkDispatch} from 'redux-thunk';

//@ts-ignore PluginRegistry doesn't have types yet
import {PluginRegistry} from 'mattermost-webapp/plugins/registry';

import {GlobalState} from 'mattermost-webapp/types/store';
import {FileInfo} from 'mattermost-redux/types/files';

import {getWopiFilesList} from 'actions/wopi';
import {wopiFilesList} from 'selectors';
import Reducer from 'reducers';

import FilePreviewOverride from 'components/file_preview_override';
import FilePreviewModal from 'components/file_preview/file_preview_modal';

import {id as pluginId} from './manifest';

export default class Plugin {
    public initialize(registry: PluginRegistry, store: Store<GlobalState>): void {
        registry.registerReducer(Reducer);
        registry.registerRootComponent(FilePreviewModal);
        registry.registerFilePreviewComponent(
            (fileInfo: FileInfo) => {
                const state = store.getState();
                const wopiFiles = wopiFilesList(state);
                return Boolean(wopiFiles?.[fileInfo.extension]);
            },
            FilePreviewOverride,
        );

        (store.dispatch as ThunkDispatch<GlobalState, undefined, AnyAction>)(getWopiFilesList());
    }
}

// @ts-ignore
window.registerPlugin(pluginId, new Plugin());
