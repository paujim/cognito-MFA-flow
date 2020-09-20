import React from 'react';
import { makeStyles } from '@material-ui/core/styles';

import TextField from '@material-ui/core/TextField';
import InputAdornment from '@material-ui/core/InputAdornment';
import LinearProgress from '@material-ui/core/LinearProgress';
import Typography from '@material-ui/core/Typography';

import LockOutlinedIcon from '@material-ui/icons/LockOutlined';


const useStyles = makeStyles((theme) => ({
    form: {
        width: '100%', // Fix IE 11 issue.
        marginTop: theme.spacing(1),
    },
}));


const FormNewPassword = (props) => {
    const classes = useStyles();

    return (
        <form className={classes.form} noValidate>
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
                autoFocus
                onChange={props.handleChangePassword}
            />
            <Typography variant="button" display="block" gutterBottom className={classes.title}>
                New password required
            </Typography>
            {props.isLoading ? <LinearProgress /> : null}
        </form>
    )
}

export default FormNewPassword

