import {createSelector} from 'reselect';

import {GlobalState} from 'mattermost-redux/types/store';
import {getPost} from 'mattermost-redux/selectors/entities/posts';
import {getCurrentUser} from 'mattermost-redux/selectors/entities/users';
import {FileInfo} from 'mattermost-redux/types/files';

import {FILE_EDIT_PERMISSIONS} from '../constants';
import {id as pluginId} from '../manifest';

//@ts-ignore GlobalState is not complete
const getPluginState = (state: GlobalState) => state['plugins-' + pluginId] || {};

export const wopiFilesList = (state: GlobalState) => getPluginState(state).wopiFilesList;

export const filePreviewModal = (state: GlobalState) => getPluginState(state).filePreviewModal;

//@ts-ignore GlobalState is not complete
const getCorePluginState = (state: GlobalState) => state['plugins-ai.riffanalytics.core'] || {};

export const riffMeetingInfo = (state: GlobalState) => getCorePluginState(state).riffMeetingInfo;
export const createFileModal = (state: GlobalState) => getPluginState(state).createFileModal;

export const collaboraConfig = (state: GlobalState) => getPluginState(state).config;

export const collaboraFileEditPermissionsEnabled = (state: GlobalState) => Boolean(collaboraConfig(state)?.file_edit_permissions);

export function makeGetIsCurrentUserFileOwner(): (state: GlobalState, fileInfo: FileInfo) => boolean {
    return createSelector(
        (state: GlobalState, fileInfo: FileInfo) => getPost(state, fileInfo.post_id || ''),
        (state: GlobalState) => getCurrentUser(state),
        (post, currentUser) => {
            return Boolean(post?.user_id === currentUser.id);
        },
    );
}

export function makeGetCollaboraFilePermissions(): (state: GlobalState, fileInfo: FileInfo) => FILE_EDIT_PERMISSIONS {
    return createSelector(
        (state: GlobalState) => collaboraFileEditPermissionsEnabled(state),
        (state: GlobalState, fileInfo: FileInfo) => getPost(state, fileInfo.post_id || ''),
        (state: GlobalState, fileInfo: FileInfo) => fileInfo.id,
        (featureEnabled, post, fileID) => {
            if (!featureEnabled) {
                // if the feature id disabled, then everyone in the channel can edit
                return FILE_EDIT_PERMISSIONS.PERMISSION_CHANNEL;
            }

            return post?.props?.[pluginId + '_file_permissions_' + fileID];
        },
    );
}
