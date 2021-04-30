import Client from '../client';
import Constants from '../constants';

export function getWopiFilesList() {
    return async (dispatch) => {
        let data = null;
        try {
            data = await Client.getWopiFilesList();
        } catch (error) {
            return {data, error};
        }
        console.log('Wopi files list: ', data, '@@@@@@@@@@@@@@@@@@@@@@@@@@@');
        dispatch({
            type: Constants.ACTION_TYPES.RECEIVED_WOPI_FILES_LIST,
            data,
        });
        return {data, error: null};
    };
}
