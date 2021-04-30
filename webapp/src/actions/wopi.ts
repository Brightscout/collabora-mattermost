import {Dispatch} from 'redux';

import Client from '../client';
import Constants from '../constants';

export function getWopiFilesList() {
    return async (dispatch: Dispatch) => {
        let data = null;
        try {
            data = await Client.getWopiFilesList();
        } catch (error) {
            return {data, error};
        }
        dispatch({
            type: Constants.ACTION_TYPES.RECEIVED_WOPI_FILES_LIST,
            data,
        });
        return {data, error: null};
    };
}
