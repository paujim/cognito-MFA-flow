import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import Container from '@material-ui/core/Container';

import { UserProvider  } from "./context/user";
import { useAuthAPI } from "./utils/auth-api";
import LoginDialog from './components/loginDialog'
import Layout from './components/layout'

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

export default function App() {

    const [open, setOpen] = React.useState(false);

    const handleClickOpen = () => {
        if (useAuthAPI.isAuthenticated()) {
            useAuthAPI.logout()
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
        <UserProvider >
            <Container maxWidth="lg">
                <div className={classes.root}>
                    <AppBar position="static">
                        <Toolbar>
                            <Typography variant="h6" className={classes.title}>
                                MFA Sample
                            </Typography>
                            <Button color="inherit" onClick={handleClickOpen}>{!useAuthAPI.isAuthenticated() ? "Login" : "Logout"}</Button>
                            <LoginDialog open={open} handleClose={handleClose} />
                        </Toolbar>
                    </AppBar>
                    <Layout isAuthenticated={useAuthAPI.isAuthenticated()} title={"User:"} />
                </div>
            </Container>
        </UserProvider >
    );
}
