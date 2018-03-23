import React from 'react';
import Dialog, {
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
} from 'material-ui/Dialog';
import Button from 'material-ui/Button';
import Dropzone from 'react-dropzone';
import Delete from 'material-ui-icons/Delete';
import CloudUpload from 'material-ui-icons/CloudUpload';
import InsertDriveFile from 'material-ui-icons/InsertDriveFile';
import IconButton from 'material-ui/IconButton';
import Snackbar from 'material-ui/Snackbar';
import { withStyles } from 'material-ui/styles';
import 'app/css/material-ui-dropzone.css';
import {isImage} from 'app/util/dropzone-helpers.js';

const styles = theme => ({
  leftIcon: {
    marginRight: theme.spacing.unit,
  }
});

class MaterialDropZone extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            open: false,
            openSnackBar: false,
            errorMessage: '',
            files: this.props.files || [],
            disabled: true,
            acceptedFiles: this.props.acceptedFiles ||
            ['image/jpeg', 'image/png', 'image/bmp', 'application/vnd.ms-excel',
                'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
                'application/vnd.ms-powerpoint',
                'application/vnd.openxmlformats-officedocument.presentationml.presentation',
                'application/msword', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'],
        };
    }

    componentWillReceiveProps(nextProps) {
        this.setState({
            open: nextProps.open,
            files: nextProps.files,
        });
    }

    handleClose() {
        this.props.closeDialog();
        this.setState({open: false});
    }

    onDrop(files) {

      console.log('onDrop');
      console.log(files);

        let oldFiles = this.state.files;
        const filesLimit = this.props.filesLimit || '3';

        oldFiles = oldFiles.concat(files);
        if (oldFiles.length > filesLimit) {
            this.setState({
                openSnackBar: true,
                errorMessage: 'Cannot upload more then ' + filesLimit + ' items.',
            });
        } else {
            this.setState({
                files: oldFiles,
            }, this.changeButtonDisable);
        }
    }

    handleRemove(file, fileIndex) {
        const files = this.state.files;
        // This is to prevent memory leaks.
        window.URL.revokeObjectURL(file.preview);

        files.splice(fileIndex, 1);
        this.setState(files, this.changeButtonDisable);

        if (file.path) {
            this.props.deleteFile(file);
        }
    }

    changeButtonDisable() {
        if (this.state.files.length !== 0) {
            this.setState({
                disabled: false,
            });
        } else {
            this.setState({
                disabled: true,
            });
        }
    }

    saveFiles() {
        const filesLimit = this.props.filesLimit || '3';

        if (this.state.files.length > filesLimit) {
            this.setState({
                openSnackBar: true,
                errorMessage: 'Cannot upload more then ' + filesLimit + ' items.',
            });
        } else {
            this.props.saveFiles(this.state.files);
        }
    }

    onDropRejected() {
        this.setState({
            openSnackBar: true,
            errorMessage: 'File too big, max size is 3MB',
        });
    }

    handleRequestCloseSnackBar = () => {
        this.setState({
            openSnackBar: false,
        });
    };

    render() {

       const { classes } = this.props;

        let img;
        let previews = '';
        const fileSizeLimit = this.props.maxSize || 1000000000;  // 1GB

        if (this.props.showPreviews === true) {
            previews = this.state.files.map((file, i) => {
                const path = file.preview || '/pic' + file.path;

                if (isImage(file)) {
                    //show image preview.
                    img = <img className="smallPreviewImg" src={path}/>;
                } else {
                    //Show default file image in preview.
                    img = <FileIcon className="smallPreviewImg"/>;
                }

                const divKey = "preview-" + i
                return (<div key={divKey}>
                    <div className={'imageContainer col fileIconImg'} key={i}>
                        {img}
                        <div className="middle">
                            <Button className="removeBtn" size="small" color="inherit">
                              Delete <Delete className={classes.leftIcon} onTouchTap={this.handleRemove(this, file, i)}/>
                            </Button>
                        </div>
                    </div>
                </div>);
            });
        }

        return (
            <div>
                <Dialog
                    title={'Upload File'}
                    open={this.props.open}
                    onClose={this.props.onClose}>
                    <Dropzone
                        accept={this.state.acceptedFiles.join(',')}
                        onDrop={this.onDrop.bind(this)}
                        className={'dropZone'}
                        acceptClassName={'stripes'}
                        rejectClassName={'rejectStripes'}
                        onDropRejected={this.onDropRejected.bind(this)}>
                        <div className={'dropzoneTextStyle'}>
                            <p className={'dropzoneParagraph'}>{'Drag and drop an image file here or click'}</p>
                            <br/>
                            <Button>
                              <CloudUpload className={'uploadIconSize'}/>
                            </Button>
                        </div>
                    </Dropzone>
                    <br/>
                    <div className="row">
                        {this.state.files.length ? <span>Preview:</span> : ''}
                    </div>
                    <div className="row">
                        {previews}
                    </div>
                    <DialogActions>
                      <Button label={'Cancel'} onTouchTap={this.props.onClose}>Cancel</Button>
                      <Button label={'Submit'} disabled={this.state.disabled} onTouchTap={this.saveFiles.bind(this)}>Submit</Button>
                    </DialogActions>
                </Dialog>
                <Snackbar
                    open={this.state.openSnackBar}
                    message={this.state.errorMessage}
                    autoHideDuration={4000}
                    onRequestClose={this.handleRequestCloseSnackBar}/>
            </div>
        );
    }
}

export default withStyles(styles)(MaterialDropZone);
