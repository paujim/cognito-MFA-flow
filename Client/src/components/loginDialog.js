import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Dialog from '@material-ui/core/Dialog';
import MuiDialogTitle from '@material-ui/core/DialogTitle';
import IconButton from '@material-ui/core/IconButton';
import Avatar from '@material-ui/core/Avatar';
import Container from '@material-ui/core/Container';
import Snackbar from '@material-ui/core/Snackbar';
import MobileStepper from '@material-ui/core/MobileStepper';
import Button from '@material-ui/core/Button';

import KeyboardArrowLeft from '@material-ui/icons/KeyboardArrowLeft';
import KeyboardArrowRight from '@material-ui/icons/KeyboardArrowRight';
import CloseIcon from '@material-ui/icons/Close';
import FingerprintIcon from '@material-ui/icons/Fingerprint';

import FormUsernamePassword from './formUsernamePassword'
import FormNewPassword from './formNewPassword'

import { useAuthAPI } from "../utils/auth-api";
import { UserContext } from "../context/user";

const useStyles = makeStyles((theme) => ({
    title: {
        padding: theme.spacing(1),
    },
    stepper: {
        maxWidth: 400,
        flexGrow: 1,
    },
    paper: {
        marginTop: theme.spacing(2),
        marginBottom: theme.spacing(2),
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        height: 380,
        width: 390,
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

const FormSwitch = (props) => {
    switch (props.step) {
        case 0:
            return (<FormUsernamePassword {...props} />)
        case 1:
            return (<FormNewPassword {...props} />)
        default:
            return (<form noValidate></form>)
    }
}

const SnackbarError = (props) => {
    return (
        <Snackbar
            anchorOrigin={{
                vertical: 'bottom',
                horizontal: 'left',
            }}
            open={props.open}
            autoHideDuration={3000}
            message={props.message}
            action={
                <React.Fragment>
                    <IconButton size="small" aria-label="close" color="inherit">
                        <CloseIcon fontSize="small" />
                    </IconButton>
                </React.Fragment>
            }
        />
    )
}

export default function LoginDialog(props) {

    const { setUser } = React.useContext(UserContext)

    const [isLoading, setIsLoading] = React.useState(false);

    const [activeStep, setActiveStep] = React.useState(0);
    const handleNext = () => {

        switch (activeStep) {
            case 0:
                CallGetCredentials()
                break;
            case 1:
                CallChangePassword()
                break;
            default:
                break;
        }
    };
    const handleBack = () => {
        setActiveStep(0);
    };

    const [errorMessage, setErrorMessage] = React.useState({ message: "", show: false });

    const [username, setUsername] = React.useState("");
    const handleChangeUsername = event => {
        setUsername(event.target.value);
    };

    const handleError = (error) => {
        console.log(error.message)
        setErrorMessage({ message: error.message, show: true })
    }

    const handleAccessToken = (accessToken) => {
        useAuthAPI.setToken(accessToken)
        let user = useAuthAPI.decodeUser(accessToken)
        setUser(user)
        props.handleClose()
    }

    const [password, setPassword] = React.useState("");
    const handleChangePassword = event => {
        setPassword(event.target.value);
    };

    const CallGetCredentials = () => {
        setIsLoading(true)
        useAuthAPI.login(username, password)
            .then(data => {
                setIsLoading(false)
                console.log(data)
                if (data && data.message === "New password required") {
                    console.log(data.session)
                    useAuthAPI.setSession(data.session)
                    setActiveStep((prevActiveStep) => prevActiveStep + 1);
                }
                if (data && data.accessToken) {
                    console.log("Logged in Successfully")
                    handleAccessToken(data.accessToken)
                }

            })
            .catch(error => {
                setIsLoading(false)
                console.log("Login Falied")
                handleError(error)
            })
    }

    const CallChangePassword = () => {
        setIsLoading(true)
        useAuthAPI.changePassword(username, password, useAuthAPI.getSession())
            .then(data => {
                setIsLoading(false)
                if (data && data.accessToken) {
                    console.log("Password Changed Successfully")
                    handleAccessToken(data.accessToken)
                }
            })
            .catch(error => {
                setIsLoading(false)
                console.log("Password Change Falied")
                handleError(error)

            })
    }

    const classes = useStyles();
    const maxSteps = 3

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
                        <FormSwitch handleChangeUsername={handleChangeUsername} handleChangePassword={handleChangePassword} isLoading={isLoading} step={activeStep} />
                        <SnackbarError open={errorMessage.show} message={errorMessage.message} />
                    </div>
                    <MobileStepper
                        variant="dots"
                        steps={maxSteps}
                        position="static"
                        activeStep={activeStep}
                        className={classes.stepper}
                        nextButton={
                            <Button size="small" onClick={handleNext} disabled={activeStep === maxSteps - 1}> Next
                                <KeyboardArrowRight />
                            </Button>
                        }
                        backButton={
                            <Button size="small" onClick={handleBack} disabled={activeStep === 0}> Back
                                 <KeyboardArrowLeft />
                            </Button>
                        }
                    />
                </Container>
            </Dialog>
        </div>
    );
}
