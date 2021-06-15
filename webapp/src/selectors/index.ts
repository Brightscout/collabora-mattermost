import { GlobalState } from 'mattermost-webapp/types/store';

import { id as pluginId } from '../manifest';

//@ts-ignore GlobalState is not complete
const getPluginState = (state: GlobalState) => state['plugins-' + pluginId] || {};

export const wopiFilesList = (state: GlobalState) => getPluginState(state).wopiFilesList;

export const filePreviewModal = (state: GlobalState) => getPluginState(state).filePreviewModal;

//@ts-ignore GlobalState is not complete
const getCorePluginState = (state: GlobalState) => state['plugins-com.meeting.mattermost'] || {};

export const videoMeetingInfo = (state: GlobalState) => getCorePluginState(state).videoMeetingInfo;
