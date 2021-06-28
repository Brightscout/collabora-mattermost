import React, {FC, useCallback, useEffect, useState} from 'react';
import {AnyAction} from 'redux';
import {ThunkDispatch} from 'redux-thunk';
import {useDispatch, useSelector} from 'react-redux';

import {FileInfo} from 'mattermost-redux/types/files';
import {getCurrentUser} from 'mattermost-redux/selectors/entities/users';
import {GlobalState} from 'mattermost-redux/types/store';
import {getPost} from 'mattermost-redux/selectors/entities/posts';

import {closeFilePreview} from 'actions/preview';
import {
    collaboraConfig, collaboraFileEditPermissionsEnabled,
    filePreviewModal,
    makeGetCollaboraFilePermissions,
    makeGetIsCurrentUserFileOwner
} from 'selectors';

import FullScreenModal from 'components/full_screen_modal';
import WopiFilePreview from 'components/wopi_file_preview';
import FilePreviewHeader from 'components/file_preview_header';

import {FILE_EDIT_PERMISSIONS} from '../constants';
import {updateFileEditPermission} from '../actions/file';
import {id as pluginId} from '../manifest';

type FilePreviewModalSelector = {
    visible: boolean;
    fileInfo: FileInfo;
}

const FilePreviewModal: FC = () => {
    const dispatch: ThunkDispatch<GlobalState, undefined, AnyAction> = useDispatch();
    const {visible, fileInfo}: FilePreviewModalSelector = useSelector(filePreviewModal);
    const [editable, setEditable] = useState(false);
    const toggleEditing = useCallback(() => {
        setEditable((prevState) => !prevState);
    }, [setEditable]);

    const getIsCurrentUserFileOwner = makeGetIsCurrentUserFileOwner();
    const getCollaboraFilePermissions = makeGetCollaboraFilePermissions();
    const isCurrentUserOwner = useSelector((state: GlobalState) => getIsCurrentUserFileOwner(state, fileInfo));
    const filePermission = useSelector((state: GlobalState) => getCollaboraFilePermissions(state, fileInfo));
    const editPermissionsFeatureEnabled = useSelector(collaboraFileEditPermissionsEnabled);
    const showEditPermissionChangeOption = editPermissionsFeatureEnabled && isCurrentUserOwner;

    const [canChannelEdit, setCanChannelEdit] = useState(false);
    const toggleCanChannelEdit = async () => {
        await setCanChannelEdit((prevState: boolean) => !prevState);
        const updatedPermission = canChannelEdit ? FILE_EDIT_PERMISSIONS.PERMISSION_OWNER : FILE_EDIT_PERMISSIONS.PERMISSION_CHANNEL;
        const response = await dispatch(updateFileEditPermission(fileInfo.id, updatedPermission));
        if ((response as {error: unknown}).error) {
            // TODO handle error
            setCanChannelEdit((prevState: boolean) => !prevState);
        } else if (!canChannelEdit) {
            setEditable(false);
        }
    };

    useEffect(() => {
        let defaultCanChannelEdit = true;
        if (editPermissionsFeatureEnabled) {
            defaultCanChannelEdit = filePermission === FILE_EDIT_PERMISSIONS.PERMISSION_CHANNEL;
        }

        setCanChannelEdit(defaultCanChannelEdit);

        const userHasEditPermission = showEditPermissionChangeOption || canChannelEdit;
        if (!userHasEditPermission) {
            setEditable(false);
        }
    }, [canChannelEdit, editPermissionsFeatureEnabled, fileInfo.id, filePermission, showEditPermissionChangeOption]);

    const handleClose = useCallback((e?: Event): void => {
        if (e && e.preventDefault) {
            e.preventDefault();
        }

        dispatch(closeFilePreview());
        setEditable(false);
    }, [dispatch]);

    return (
        <FullScreenModal
            compact={true}
            show={visible}
        >
            <FilePreviewHeader
                fileInfo={fileInfo}
                onClose={handleClose}
                editable={editable}
                toggleEditing={toggleEditing}
                canChannelEdit={canChannelEdit}
                toggleCanChannelEdit={toggleCanChannelEdit}
                showEditPermissionChangeOption={showEditPermissionChangeOption}
            />
            <WopiFilePreview
                fileInfo={fileInfo}
                editable={editable}
            />
        </FullScreenModal>
    );
};

export default FilePreviewModal;
