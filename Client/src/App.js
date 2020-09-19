import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import Container from '@material-ui/core/Container';

import { useAuth } from "./context/auth";
import LoginForm from './components/loginForm'

const useStyles = makeStyles((theme) => ({
    root: {
        flexGrow: 1,
    },
    menuButton: {
        marginRight: theme.spacing(2),
    },
    title: {
        flexGrow: 1,
    },
}));



export default function ButtonAppBar() {

    const [open, setOpen] = React.useState(false);

    const handleClickOpen = () => {
        console.log(useAuth.isAuthenticated())
        if (useAuth.isAuthenticated()) {
            useAuth.logout()
        }
        else {
            setOpen(true);
        }

    };
    const handleClose = () => {
        setOpen(false);
    };

    const classes = useStyles();
    return (
        <Container maxWidth="lg">
            <div className={classes.root}>
                <AppBar position="static">
                    <Toolbar>
                        <Typography variant="h6" className={classes.title}>
                            MFA Sample
                            </Typography>
                        <Button color="inherit" onClick={handleClickOpen}>{!useAuth.isAuthenticated() ? "Login" : "Logout"}</Button>
                        <LoginForm open={open} handleClose={handleClose} />
                    </Toolbar>
                </AppBar>
            </div>
        </Container>
    );
}
