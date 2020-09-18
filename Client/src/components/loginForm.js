import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Dialog from '@material-ui/core/Dialog';
import MuiDialogTitle from '@material-ui/core/DialogTitle';
import IconButton from '@material-ui/core/IconButton';
import TextField from '@material-ui/core/TextField';
import Avatar from '@material-ui/core/Avatar';
import Container from '@material-ui/core/Container';
import InputAdornment from '@material-ui/core/InputAdornment';

import ChevronRightIcon from '@material-ui/icons/ChevronRight';
import CloseIcon from '@material-ui/icons/Close';
import LockOutlinedIcon from '@material-ui/icons/LockOutlined';
import AccountCircleIcon from '@material-ui/icons/AccountCircle';
import FingerprintIcon from '@material-ui/icons/Fingerprint';

import Spinner from './spinner';
import { useAuth } from "../context/auth";


const useStyles = makeStyles((theme) => ({
    paper: {
        marginTop: theme.spacing(2),
        marginBottom: theme.spacing(2),
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        height: 380,
    },
    avatar: {
        margin: theme.spacing(4),
        backgroundColor: theme.palette.secondary.main,
        width: theme.spacing(8),
        height: theme.spacing(8),
    },
    form: {
        width: '100%', // Fix IE 11 issue.
        marginTop: theme.spacing(1),
    },
    closeButton: {
        position: 'absolute',
        right: theme.spacing(1),
        top: theme.spacing(1),
        color: theme.palette.grey[500],
    },
    nextButton: {
        position: 'absolute',
        right: theme.spacing(1),
        bottom: theme.spacing(1),
        // color: theme.palette.grey[500],
    },
}));

const { login } = useAuth;


export default function LoginForm(props) {

    const [showSpinner, setshowSpinner] = React.useState(false);

    const [username, setUsername] = React.useState("");
    const handleChangeUsername = event => {
        setUsername(event.target.value);
    };

    const [password, setPassword] = React.useState("");
    const handleChangePassword = event => {
        setPassword(event.target.value);
    };

    const fetchCredentials = () => {
        setshowSpinner(true)
        login(username, password)
            .then(success => {
                setshowSpinner(false)
                console.log("LOGIN_OK")
            })
            .catch(error => {
                setshowSpinner(false)
                console.log("LOGIN_ERROR")
                console.log(error)

            })
            .finally(() => {
                props.handleClose()
            })

    }

    const classes = useStyles();

    return (
        <div>
            <Dialog onClose={props.handleClose} aria-labelledby="customized-dialog-title" open={props.open}>
                <MuiDialogTitle disableTypography>
                    <IconButton aria-label="close" className={classes.closeButton} onClick={props.handleClose}>
                        <CloseIcon />
                    </IconButton>
                </MuiDialogTitle>
                <Container component="main" maxWidth="xs" >
                    <div className={classes.paper}>
                        <Avatar className={classes.avatar}>
                            <FingerprintIcon />
                        </Avatar>
                        <form className={classes.form} noValidate>
                            <TextField
                                InputProps={{
                                    startAdornment: (
                                        <InputAdornment position="start">
                                            <AccountCircleIcon />
                                        </InputAdornment>
                                    ),
                                }}
                                required
                                variant="outlined"
                                fullWidth
                                id="username"
                                label="username"
                                name="username"
                                autoFocus
                                onChange={handleChangeUsername}
                            />
                            <TextField
                                variant="outlined"
                                InputProps={{
                                    startAdornment: (
                                        <InputAdornment position="start">
                                            <LockOutlinedIcon />
                                        </InputAdornment>
                                    ),
                                }}
                                margin="normal"
                                required
                                fullWidth
                                name="password"
                                label="Password"
                                type="password"
                                id="password"
                                onChange={handleChangePassword}
                            />
                        </form>
                        <Spinner show={showSpinner} />
                    </div>
                    <IconButton aria-label="login" className={classes.nextButton} onClick={fetchCredentials}>
                        <ChevronRightIcon />
                    </IconButton>
                </Container>
            </Dialog>
        </div>
    );
}

