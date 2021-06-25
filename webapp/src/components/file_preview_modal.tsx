import React, {FC, useCallback, useState} from 'react';

import {FileInfo} from 'mattermost-redux/types/files';

import {useDispatch, useSelector} from 'react-redux';

import {closeFilePreview} from 'actions/preview';
import {filePreviewModal, riffMeetingInfo} from 'selectors';

import FullScreenModal from 'components/full_screen_modal';
import WopiFilePreview from 'components/wopi_file_preview';
import FilePreviewHeader from 'components/file_preview_header';

type FilePreviewModalSelector = {
    visible: boolean;
    fileInfo: FileInfo;
}

type RiffMeetingSelectror = {
    channelID: string
}

const FilePreviewModal: FC = () => {
    const dispatch = useDispatch();
    const {visible, fileInfo}: FilePreviewModalSelector = useSelector(filePreviewModal);
    const riffMeetingState: RiffMeetingSelectror = useSelector(riffMeetingInfo);

    const riffMeetingVisible = riffMeetingState && 'channelID' in riffMeetingState;
    
    const [editable, setEditable] = useState(false);
    const toggleEditing = useCallback(() => {
        setEditable((prevState) => !prevState);
    }, [setEditable]);

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
            show={visible && !riffMeetingVisible}
        >
            <FilePreviewHeader
                fileInfo={fileInfo}
                onClose={handleClose}
                editable={editable}
                toggleEditing={toggleEditing}
            />
            <WopiFilePreview
                fileInfo={fileInfo}
                editable={editable}
            />
        </FullScreenModal>
    );
};

export default FilePreviewModal;
