import React, {FC, useMemo, useState} from 'react';
import {useSelector} from 'react-redux';
import clsx from 'clsx';
import {Button} from 'react-bootstrap';

import {FileInfo} from 'mattermost-redux/types/files';
import {GlobalState} from 'mattermost-redux/types/store';
import {getPost} from 'mattermost-redux/selectors/entities/posts';
import {getChannel} from 'mattermost-redux/selectors/entities/channels';
import {getCurrentUser} from 'mattermost-redux/selectors/entities/users';

import Client from 'client';
import {collaboraConfig} from 'selectors';

import {CHANNEL_TYPES} from '../constants';

import CloseIcon from './close_icon';

type Props = {
    fileInfo: FileInfo;
    onClose: () => void;
    editable: boolean;
    toggleEditing: () => void;
}

export const FilePreviewHeader: FC<Props> = ({fileInfo, onClose, editable, toggleEditing}: Props) => {
    const post = useSelector((state: GlobalState) => getPost(state, fileInfo.post_id || ''));
    const channel = useSelector((state: GlobalState) => getChannel(state, post?.channel_id));
    const channelName: React.ReactNode = useMemo(() => {
        if (!channel) {
            return '';
        }

        switch (channel.type) {
        case CHANNEL_TYPES.CHANNEL_DIRECT:
            return 'Direct Message';

        case CHANNEL_TYPES.CHANNEL_GROUP:
            return 'Group Message';

        default:
            return channel.display_name;
        }
    }, [channel]);

    const currentUser = useSelector(getCurrentUser);
    const collaboraConf = useSelector(collaboraConfig);

    const editPermissionsFeatureEnabled = collaboraConf.file_edit_permissions;
    const showEditPermissionChangeOption = editPermissionsFeatureEnabled && post.user_id === currentUser.id;
    const [canChannelEdit, setCanChannelEdit] = useState(true);
    const toggleCanChannelEdit = () => {
        setCanChannelEdit((prevState) => !prevState);
    };

    return (
        <>
            <div
                id='header'
                style={{
                    fontSize: 15,
                    lineHeight: 1.46668,
                    fontWeight: 400,
                    borderBottom: '1px solid #e1e1e1',
                    boxShadow: 'inset 0 1px 0 rgb(0 0 0 / 20%)',
                    height: 64,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    flex: '0 0 auto',
                }}
            >
                <div
                    id='headerMeta'
                    style={{
                        display: 'flex',
                        alignItems: 'center',
                        padding: '12px 16px',
                        minWidth: 0,
                    }}
                >
                    <div
                        style={{
                            maxHeight: 40,
                            minWidth: 0,
                        }}
                    >
                        <div
                            style={{
                                display: 'block',
                                textOverflow: 'ellipsis',
                                overflow: 'hidden',
                                whiteSpace: 'nowrap',
                                fontWeight: 700,
                            }}
                        >
                            {fileInfo.name}
                        </div>
                        <div
                            style={{
                                display: 'flex',
                                fontSize: 13,
                                lineHeight: 1.38463,
                                fontWeight: 400,
                            }}
                        >
                            <span
                                style={{
                                    color: '#606060',
                                    fontWeight: 700,
                                    paddingRight: 4,
                                    whiteSpace: 'nowrap',
                                    overflow: 'hidden',
                                    textOverflow: 'ellipsis',
                                }}
                            >
                                {channelName}
                            </span>
                        </div>
                    </div>
                </div>
                <div className='collabora-header-actions'>
                    <Button
                        bsSize='large'
                        bsStyle='large'
                        title='Download'
                        aria-label='Download'
                        className='collabora-header-action-button'
                        href={Client.getFileUrl(fileInfo.id)}
                        target='_blank'
                        rel='noopener noreferrer'
                        download={true}
                    >
                        <i className='fa fa-cloud-download'/>
                    </Button>
                    <Button
                        bsSize='large'
                        bsStyle='large'
                        onClick={toggleEditing}
                        className='collabora-header-action-button'
                        title={`${editable ? 'Lock' : 'Unlock'} Editing`}
                        aria-label={`${editable ? 'Lock' : 'Unlock'} Editing`}
                    >
                        <i
                            className={clsx(
                                'fa',
                                {
                                    'fa-lock': !editable,
                                    'fa-unlock': editable,
                                },
                            )}
                        />
                    </Button>
                    {
                        showEditPermissionChangeOption && (
                            <Button
                                bsStyle='large'
                                onClick={toggleCanChannelEdit}
                                className='collabora-header-action-button'
                                title={canChannelEdit ? 'Everyone in the channel can edit.' : 'Only you can edit.'}
                                aria-label={canChannelEdit ? 'Everyone in the channel can edit.' : 'Only you can edit.'}
                            >
                                <i
                                    className={clsx(
                                        'fa',
                                        {
                                            'fa-users': canChannelEdit,
                                            'fa-user': !canChannelEdit,
                                        },
                                    )}
                                />
                            </Button>
                        )
                    }
                    <div className='collabora-header-actions-separator'/>
                    <CloseIcon
                        id='closeIcon'
                        title='Close'
                        aria-label='Close'
                        className='close-x collabora-header-action-button'
                        onClick={onClose}
                    />
                </div>
            </div>
        </>
    );
};

export default FilePreviewHeader;
