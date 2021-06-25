import React, {FC, useCallback, useState} from 'react';

import {Button} from 'react-bootstrap';

import {FileInfo} from 'mattermost-redux/types/files';

import {useSelector} from 'react-redux';
import {GlobalState} from 'mattermost-redux/types/store';
import {getPost} from 'mattermost-redux/selectors/entities/posts';
import {getCurrentUser} from 'mattermost-redux/selectors/entities/users';

import WopiFilePreview from 'components/wopi_file_preview';
import {collaboraConfig} from '../selectors';
import {id as pluginId} from '../manifest';
import {FILE_EDIT_PERMISSIONS} from "../constants";

type Props = {
    fileInfo: FileInfo;
}

const FilePreviewComponent: FC<Props> = ({fileInfo}: Props) => {
    const [loading, setLoading] = useState(true);
    const [editable, setEditable] = useState(false);
    const enableEditing = useCallback(() => {
        setEditable(true);
    }, []);

    const post = useSelector((state: GlobalState) => getPost(state, fileInfo.post_id || ''));
    const currentUser = useSelector(getCurrentUser);
    const collaboraConf = useSelector(collaboraConfig);

    const editPermissionsFeatureEnabled = collaboraConf.file_edit_permissions;
    const showEditPermissionChangeOption = editPermissionsFeatureEnabled && post?.user_id === currentUser.id;
    const canChannelEdit = post?.props?.[pluginId + '_file_permissions_' + fileInfo.id] === FILE_EDIT_PERMISSIONS.PERMISSION_CHANNEL;
    const canCurrentUserEdit = showEditPermissionChangeOption || canChannelEdit;

    return (
        <>
            <WopiFilePreview
                fileInfo={fileInfo}
                editable={editable}
                setLoading={setLoading}
            />
            {canCurrentUserEdit && !loading && !editable && (
                <Button onClick={enableEditing}>
                    <span className='wopi-switch-to-edit-mode'>
                        <i className='fa fa-pencil-square-o'/>
                        {' Enable Editing'}
                    </span>
                </Button>
            )}
        </>
    );
};

export default FilePreviewComponent;
