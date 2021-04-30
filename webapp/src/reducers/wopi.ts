import Constants from '../constants';

export const wopiFilesList = (state = {}, action) => {
    switch (action.type) {
    case Constants.ACTION_TYPES.RECEIVED_WOPI_FILES_LIST:
        return action.data;
    default:
        return state;
    }
};
