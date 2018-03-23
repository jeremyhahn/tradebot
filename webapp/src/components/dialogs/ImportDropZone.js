import DropZoneModal from 'app/components/MaterialDropZone';
import React, {Component} from 'react';
import Button from 'material-ui/Button';
import {MuiThemeProvider} from 'material-ui';

class ImportDialog extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            openUploadModal: false,
            files: [],
        };
    }

    closeDialog() {
        this.setState({openUploadModal: false});
    }

    saveFiles(files) {
        //Saving files to state for further use and closing Modal.
        console.log(files);
        this.setState({files: files, openUploadModal: false});
    }

    handleOpenUpload() {
        this.setState({
            openUploadModal: true,
        });
    }

    deleteFile(fileName) {
        this.props.deleteFile(fileName);
    }

    render() {
        //If we already saved files they will be shown again in modal preview.
        let files = this.state.files;
        let style = {
            addFileBtn: {
                'marginTop': '15px',
            },
        };

        return (
                <div>
                    <DropZoneModal
                        open={this.props.open}
                        onClose={this.props.onClose}
                        saveFiles={this.saveFiles.bind(this)}
                        deleteFile={this.deleteFile.bind(this)}
                        acceptedFiles={['image/jpeg', 'image/png', 'image/bmp']}
                        files={files}
                        showPreviews={true}
                        maxSize={5000000}
                        closeDialog={this.props.close}/>
                </div>
        );
    }
}

export default ImportDialog;
