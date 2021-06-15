import React, { FC, useCallback } from 'react';

import { FileInfo } from 'mattermost-redux/types/files';

import { useDispatch, useSelector } from 'react-redux';

import { closeFilePreview } from 'actions/preview';
import { filePreviewModal, videoMeetingInfo } from 'selectors';

import FullScreenModal from 'components/full_screen_modal';
import WopiFilePreview from 'components/wopi_file_preview';
import FilePreviewHeader from 'components/file_preview_header';

type FilePreviewModalSelector = {
    visible: boolean;
    fileInfo: FileInfo;
}

type VideoMeetingSelectror = {
    channelID: string
}

const FilePreviewModal: FC = () => {
    const dispatch = useDispatch();
    const { visible, fileInfo }: FilePreviewModalSelector = useSelector(filePreviewModal);
    const videoMeetingState: VideoMeetingSelectror = useSelector(videoMeetingInfo);

    const videoMeetingVisible = videoMeetingState && 'channelID' in videoMeetingState;

    const handleClose = useCallback((e?: Event): void => {
        if (e && e.preventDefault) {
            e.preventDefault();
        }

        dispatch(closeFilePreview());
    }, [dispatch]);

    return (
        <FullScreenModal
            compact={true}
            show={visible && !videoMeetingVisible}
        >
            <FilePreviewHeader
                fileInfo={fileInfo}
                onClose={handleClose}
            />
            <WopiFilePreview
                fileInfo={fileInfo}
                editable={true}
            />
        </FullScreenModal>
    );
};

export default FilePreviewModal;
