import {id as pluginId} from '../manifest';

const getPluginState = (state) => state['plugins-' + pluginId] || {};

export const wopiFilesList = (state) => getPluginState(state).wopiFilesList;
