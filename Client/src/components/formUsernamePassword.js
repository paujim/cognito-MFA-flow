import React from 'react';
import { makeStyles } from '@material-ui/core/styles';

import TextField from '@material-ui/core/TextField';
import InputAdornment from '@material-ui/core/InputAdornment';
import LinearProgress from '@material-ui/core/LinearProgress';
import LockOutlinedIcon from '@material-ui/icons/LockOutlined';
import AccountCircleIcon from '@material-ui/icons/AccountCircle';


const useStyles = makeStyles((theme) => ({
    form: {
        width: '100%', // Fix IE 11 issue.
        marginTop: theme.spacing(1),
    },
}));

const FormUsernamePassword = (props) => {
    const classes = useStyles();

    return (
        <form className={classes.form} noValidate>
            <TextField
                disabled={props.isLoading}
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
                onChange={props.handleChangeUsername}
            />
            <TextField
                disabled={props.isLoading}
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
                onChange={props.handleChangePassword}
            />
            {props.isLoading ? <LinearProgress /> : null}
        </form>
    )
}

export default FormUsernamePassword
